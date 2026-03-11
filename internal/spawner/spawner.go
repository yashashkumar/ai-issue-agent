package spawner

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"syscall"
	"time"

	"github.com/yashashkumar/ai-issue-agent/internal/models"
	"github.com/yashashkumar/ai-issue-agent/internal/repository"
)

type RunningAgent struct {
	ID        string
	PID       int
	Cancel    context.CancelFunc
	StartedAt time.Time
}

type Spawner struct {
	requests      chan SpawnRequest
	agentRepo     *repository.AgentRepo
	worktreeRepo  *repository.WorktreeRepo
	logger        *slog.Logger
	wg            sync.WaitGroup
	maxConcurrent int
	sem           chan struct{}
	ctx           context.Context
	cancel        context.CancelFunc
	mu            sync.RWMutex
	running       map[string]*RunningAgent
}

func NewSpawner(agentRepo *repository.AgentRepo, worktreeRepo *repository.WorktreeRepo, logger *slog.Logger, maxConcurrent int, buffer int) *Spawner {
	ctx, cancel := context.WithCancel(context.Background())
	return &Spawner{
		requests:      make(chan SpawnRequest, buffer),
		agentRepo:     agentRepo,
		worktreeRepo:  worktreeRepo,
		logger:        logger,
		maxConcurrent: maxConcurrent,
		sem:           make(chan struct{}, maxConcurrent),
		ctx:           ctx,
		cancel:        cancel,
		running:       make(map[string]*RunningAgent),
	}
}

func (s *Spawner) Spawn(req SpawnRequest) {
	// Create agent record in pending state
	agent := &models.Agent{
		ID:                req.ID,
		ProjectID:         req.ProjectID,
		Name:              req.Name,
		Description:       req.Description,
		Status:            models.AgentStatusPending,
		Source:            req.Source,
		SystemPrompt:      req.SystemPrompt,
		WorkDir:           req.WorkDir,
		GitHubIssueNumber: req.GitHubIssueNumber,
	}

	// Normalize empty project_id to nil to avoid FK constraint failure
	if agent.ProjectID != nil && *agent.ProjectID == "" {
		agent.ProjectID = nil
	}

	if err := s.agentRepo.Create(context.Background(), agent); err != nil {
		s.logger.Error("failed to create agent record", "id", agent.ID, "err", err)
		return
	}

	s.requests <- req
	s.logger.Info("agent spawn request queued", "id", agent.ID)
}

func (s *Spawner) Start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case <-s.ctx.Done():
				return
			case req, ok := <-s.requests:
				if !ok {
					return
				}
				s.sem <- struct{}{} // Acquire semaphore slot

				s.wg.Add(1)
				go func(r SpawnRequest) {
					defer s.wg.Done()
					defer func() { <-s.sem }() // Release semaphore
					s.executeAgent(r)
				}(req)
			}
		}
	}()
}

func (s *Spawner) CancelAgent(agentID string) error {
	s.mu.RLock()
	agent, ok := s.running[agentID]
	s.mu.RUnlock()

	if !ok {
		return fmt.Errorf("agent not running or not found")
	}

	agent.Cancel()
	return nil
}

func (s *Spawner) Stop() {
	s.cancel() // Stop accepting new and cancel existing

	s.mu.RLock()
	for _, agent := range s.running {
		agent.Cancel()
	}
	s.mu.RUnlock()

	// Wait with a small timeout natively or let caller do it
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		s.logger.Info("spawner stopped gracefully")
	case <-time.After(30 * time.Second):
		s.logger.Warn("spawner stopped with timeout")
	}
}

func (s *Spawner) executeAgent(req SpawnRequest) {
	startedAt := time.Now()
	s.logger.Info("executing agent process", "id", req.ID)

	err := s.agentRepo.UpdateStatus(context.Background(), req.ID, models.AgentStatusRunning, nil, nil, &startedAt, nil)
	if err != nil {
		s.logger.Error("failed to update agent to running", "id", req.ID, "err", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)

	// Create isolated process group
	proc := NewAgentProcess(s.logger, s.agentRepo)

	s.mu.Lock()
	s.running[req.ID] = &RunningAgent{
		ID:        req.ID,
		Cancel:    cancel,
		StartedAt: startedAt,
	}
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.running, req.ID)
		s.mu.Unlock()
		cancel()
	}()

	err = proc.Run(ctx, req)

	finishedAt := time.Now()
	var exitCode *int
	var errMsg *string
	finalStatus := models.AgentStatusFinished

	if err != nil {
		s.logger.Error("agent process failed", "id", req.ID, "err", err)
		finalStatus = models.AgentStatusError
		em := err.Error()
		if len(em) > 1000 {
			em = em[:1000] // Truncate
		}
		errMsg = &em

		// If it's an exit error, grab the code
		if exitErr, ok := err.(*processExitError); ok {
			code := exitErr.Code()
			exitCode = &code
			// If cancelled via context
			if ctx.Err() != nil {
				em = "cancelled by user or timeout"
				errMsg = &em
				// Kill pgid
				syscall.Kill(-exitErr.Pid(), syscall.SIGKILL)
			}
		}
	} else {
		code := 0
		exitCode = &code
	}

	// Persist final status
	if err := s.agentRepo.UpdateStatus(context.Background(), req.ID, finalStatus, errMsg, exitCode, nil, &finishedAt); err != nil {
		s.logger.Error("failed to update final agent status", "id", req.ID, "err", err)
	}
}
