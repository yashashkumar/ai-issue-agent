package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/yashashkumar/ai-issue-agent/internal/database"
	"github.com/yashashkumar/ai-issue-agent/internal/database/queries"
	"github.com/yashashkumar/ai-issue-agent/internal/errors"
	"github.com/yashashkumar/ai-issue-agent/internal/models"
)

type AgentRepo struct {
	db *database.DB
}

func NewAgentRepo(db *database.DB) *AgentRepo {
	return &AgentRepo{db: db}
}

func (r *AgentRepo) Create(ctx context.Context, a *models.Agent) error {
	_, err := r.db.ExecContext(ctx, queries.CreateAgent,
		a.ID, a.ProjectID, a.Name, a.Description, a.Status, a.Source,
		a.SystemPrompt, a.WorkDir, a.ExitCode, a.ErrorMessage,
		a.GitHubIssueNumber, a.GitHubPRNumber, a.PID,
		a.StartedAt, a.FinishedAt, a.CreatedAt, a.UpdatedAt,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create agent")
	}
	return nil
}

func (r *AgentRepo) UpdateStatus(ctx context.Context, id string, status models.AgentStatus, errMsg *string, exitCode *int, startedAt, finishedAt *time.Time) error {
	_, err := r.db.ExecContext(ctx, queries.UpdateAgentStatus,
		status, errMsg, exitCode, startedAt, finishedAt, time.Now(), id,
	)
	if err != nil {
		return errors.Wrap(err, "failed to update agent status")
	}
	return nil
}

func (r *AgentRepo) UpdatePID(ctx context.Context, id string, pid int) error {
	_, err := r.db.ExecContext(ctx, queries.UpdateAgentPID, pid, time.Now(), id)
	return err
}

func (r *AgentRepo) UpdatePRNumber(ctx context.Context, id string, prNum int) error {
	_, err := r.db.ExecContext(ctx, queries.UpdateAgentPRNumber, prNum, time.Now(), id)
	return err
}

func (r *AgentRepo) GetByID(ctx context.Context, id string) (*models.Agent, error) {
	var a models.Agent
	err := r.db.QueryRowContext(ctx, queries.GetAgentByID, id).Scan(
		&a.ID, &a.ProjectID, &a.Name, &a.Description, &a.Status, &a.Source,
		&a.SystemPrompt, &a.WorkDir, &a.ExitCode, &a.ErrorMessage,
		&a.GitHubIssueNumber, &a.GitHubPRNumber, &a.PID,
		&a.StartedAt, &a.FinishedAt, &a.CreatedAt, &a.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get agent")
	}
	return &a, nil
}

func (r *AgentRepo) ListByProject(ctx context.Context, projectID string) ([]models.Agent, error) {
	rows, err := r.db.QueryContext(ctx, queries.ListProjectAgents, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list project agents")
	}
	defer rows.Close()

	var agents []models.Agent
	for rows.Next() {
		var a models.Agent
		if err := rows.Scan(
			&a.ID, &a.ProjectID, &a.Name, &a.Description, &a.Status, &a.Source,
			&a.SystemPrompt, &a.WorkDir, &a.ExitCode, &a.ErrorMessage,
			&a.GitHubIssueNumber, &a.GitHubPRNumber, &a.PID,
			&a.StartedAt, &a.FinishedAt, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, errors.Wrap(err, "failed to scan agent")
		}
		agents = append(agents, a)
	}
	return agents, nil
}

func (r *AgentRepo) ListAll(ctx context.Context) ([]models.Agent, error) {
	rows, err := r.db.QueryContext(ctx, queries.ListAllAgents)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list all agents")
	}
	defer rows.Close()

	var agents []models.Agent
	for rows.Next() {
		var a models.Agent
		if err := rows.Scan(
			&a.ID, &a.ProjectID, &a.Name, &a.Description, &a.Status, &a.Source,
			&a.SystemPrompt, &a.WorkDir, &a.ExitCode, &a.ErrorMessage,
			&a.GitHubIssueNumber, &a.GitHubPRNumber, &a.PID,
			&a.StartedAt, &a.FinishedAt, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, errors.Wrap(err, "failed to scan agent")
		}
		agents = append(agents, a)
	}
	return agents, nil
}

func (r *AgentRepo) InsertLog(ctx context.Context, agentID string, stream models.LogStream, content string) error {
	_, err := r.db.ExecContext(ctx, queries.InsertAgentLog, agentID, stream, content, time.Now())
	if err != nil {
		return errors.Wrap(err, "failed to insert agent log")
	}
	return nil
}

func (r *AgentRepo) GetLogs(ctx context.Context, agentID string, streamOption string) ([]models.AgentLog, error) {
	rows, err := r.db.QueryContext(ctx, queries.GetAgentLogs, agentID, streamOption, streamOption)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get agent logs")
	}
	defer rows.Close()

	var logs []models.AgentLog
	for rows.Next() {
		var l models.AgentLog
		if err := rows.Scan(&l.ID, &l.AgentID, &l.Stream, &l.Content, &l.CreatedAt); err != nil {
			return nil, errors.Wrap(err, "failed to scan agent log")
		}
		logs = append(logs, l)
	}
	return logs, nil
}
