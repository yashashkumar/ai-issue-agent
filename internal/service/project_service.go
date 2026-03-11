package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/yashashkumar/ai-issue-agent/internal/models"
	"github.com/yashashkumar/ai-issue-agent/internal/repository"
)

type ProjectService struct {
	repo *repository.ProjectRepo
}

func NewProjectService(repo *repository.ProjectRepo) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) CreateProject(ctx context.Context, p *models.Project) (*models.Project, error) {
	p.ID = uuid.NewString()
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProjectService) GetProject(ctx context.Context, id string) (*models.Project, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProjectService) UpdateProject(ctx context.Context, p *models.Project) error {
	return s.repo.Update(ctx, p)
}

func (s *ProjectService) DeleteProject(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *ProjectService) ListProjects(ctx context.Context) ([]models.Project, error) {
	return s.repo.List(ctx)
}
