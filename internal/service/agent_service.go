package service

import (
	"context"

	"github.com/yashashkumar/ai-issue-agent/internal/models"
	"github.com/yashashkumar/ai-issue-agent/internal/repository"
)

type AgentService struct {
	repo *repository.AgentRepo
}

func NewAgentService(repo *repository.AgentRepo) *AgentService {
	return &AgentService{repo: repo}
}

func (s *AgentService) GetAgent(ctx context.Context, id string) (*models.Agent, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *AgentService) ListProjectAgents(ctx context.Context, projectID string) ([]models.Agent, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *AgentService) ListAllAgents(ctx context.Context) ([]models.Agent, error) {
	return s.repo.ListAll(ctx)
}

func (s *AgentService) GetAgentLogs(ctx context.Context, agentID, stream string) ([]models.AgentLog, error) {
	if stream != "stdout" && stream != "stderr" && stream != "all" {
		stream = "all"
	}
	return s.repo.GetLogs(ctx, agentID, stream)
}
