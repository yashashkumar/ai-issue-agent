package repository

import (
	"context"
	"database/sql"

	"github.com/yashashkumar/ai-issue-agent/internal/database"
	"github.com/yashashkumar/ai-issue-agent/internal/database/queries"
	"github.com/yashashkumar/ai-issue-agent/internal/errors"
	"github.com/yashashkumar/ai-issue-agent/internal/models"
)

type ProjectRepo struct {
	db *database.DB
}

func NewProjectRepo(db *database.DB) *ProjectRepo {
	return &ProjectRepo{db: db}
}

func (r *ProjectRepo) Create(ctx context.Context, p *models.Project) error {
	_, err := r.db.ExecContext(ctx, queries.CreateProject,
		p.ID, p.Name, p.Description, p.RootFolder, p.AllowedEmailsJSON(),
		p.GitHubOwner, p.GitHubRepo, p.GitHubWebhookSecret,
		p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create project")
	}
	return nil
}

func (r *ProjectRepo) GetByID(ctx context.Context, id string) (*models.Project, error) {
	var p models.Project
	var emailsStr string
	err := r.db.QueryRowContext(ctx, queries.GetProjectByID, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.RootFolder, &emailsStr,
		&p.GitHubOwner, &p.GitHubRepo, &p.GitHubWebhookSecret,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get project")
	}
	p.AllowedEmails = models.ParseAllowedEmails(emailsStr)
	return &p, nil
}

func (r *ProjectRepo) GetByGitHubRepo(ctx context.Context, owner, repo string) (*models.Project, error) {
	var p models.Project
	var emailsStr string
	err := r.db.QueryRowContext(ctx, queries.GetProjectByGitHubRepo, owner, repo).Scan(
		&p.ID, &p.Name, &p.Description, &p.RootFolder, &emailsStr,
		&p.GitHubOwner, &p.GitHubRepo, &p.GitHubWebhookSecret,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get project by repo")
	}
	p.AllowedEmails = models.ParseAllowedEmails(emailsStr)
	return &p, nil
}

func (r *ProjectRepo) Update(ctx context.Context, p *models.Project) error {
	res, err := r.db.ExecContext(ctx, queries.UpdateProject,
		p.Name, p.Description, p.RootFolder, p.AllowedEmailsJSON(),
		p.GitHubOwner, p.GitHubRepo, p.GitHubWebhookSecret, p.UpdatedAt,
		p.ID,
	)
	if err != nil {
		return errors.Wrap(err, "failed to update project")
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *ProjectRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, queries.DeleteProject, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete project")
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *ProjectRepo) List(ctx context.Context) ([]models.Project, error) {
	rows, err := r.db.QueryContext(ctx, queries.ListProjects)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list projects")
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		var emailsStr string
		if err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.RootFolder, &emailsStr,
			&p.GitHubOwner, &p.GitHubRepo, &p.GitHubWebhookSecret,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, errors.Wrap(err, "failed to scan project")
		}
		p.AllowedEmails = models.ParseAllowedEmails(emailsStr)
		projects = append(projects, p)
	}
	return projects, nil
}

func (r *ProjectRepo) Search(ctx context.Context, query string) ([]models.Project, error) {
	likeQuery := "%" + query + "%"
	rows, err := r.db.QueryContext(ctx, queries.SearchProjects, likeQuery, likeQuery)
	if err != nil {
		return nil, errors.Wrap(err, "failed to search projects")
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		var emailsStr string
		if err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.RootFolder, &emailsStr,
			&p.GitHubOwner, &p.GitHubRepo, &p.GitHubWebhookSecret,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, errors.Wrap(err, "failed to scan project")
		}
		p.AllowedEmails = models.ParseAllowedEmails(emailsStr)
		projects = append(projects, p)
	}
	return projects, nil
}
