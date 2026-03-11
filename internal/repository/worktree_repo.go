package repository

import (
	"context"
	"database/sql"

	"github.com/yashashkumar/ai-issue-agent/internal/database"
	"github.com/yashashkumar/ai-issue-agent/internal/database/queries"
	"github.com/yashashkumar/ai-issue-agent/internal/errors"
	"github.com/yashashkumar/ai-issue-agent/internal/models"
)

type WorktreeRepo struct {
	db *database.DB
}

func NewWorktreeRepo(db *database.DB) *WorktreeRepo {
	return &WorktreeRepo{db: db}
}

func (r *WorktreeRepo) Create(ctx context.Context, w *models.Worktree) error {
	_, err := r.db.ExecContext(ctx, queries.CreateWorktree,
		w.ID, w.ProjectID, w.AgentID, w.BranchName, w.WorktreePath,
		w.GitHubIssueNumber, w.IsActive, w.CreatedAt, w.CleanedAt,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create worktree")
	}
	return nil
}

func (r *WorktreeRepo) GetActiveForIssue(ctx context.Context, projectID string, issueNumber int) (*models.Worktree, error) {
	var w models.Worktree
	err := r.db.QueryRowContext(ctx, queries.GetActiveWorktreeForIssue, projectID, issueNumber).Scan(
		&w.ID, &w.ProjectID, &w.AgentID, &w.BranchName, &w.WorktreePath,
		&w.GitHubIssueNumber, &w.IsActive, &w.CreatedAt, &w.CleanedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get active worktree")
	}
	return &w, nil
}

func (r *WorktreeRepo) MarkCleaned(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, queries.MarkWorktreeCleaned, id)
	return errors.Wrap(err, "failed to mark worktree cleaned")
}
