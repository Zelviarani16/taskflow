package service

import (
	"context"

	"github.com/Zelviarani16/taskflow-api/dto"
	"github.com/Zelviarani16/taskflow-api/entity"
	"github.com/Zelviarani16/taskflow-api/repository"
	"github.com/google/uuid"
)

type IProjectService interface {
	Create(ctx context.Context, ownerID uuid.UUID, req dto.CreateProjectRequest) (dto.ProjectResponse, error)
	List(ctx context.Context, ownerID uuid.UUID, page dto.PaginationQuery) ([]dto.ProjectResponse, int64, error)
	GetByID(ctx context.Context, ownerID, projectID uuid.UUID) (dto.ProjectResponse, error)
	Update(ctx context.Context, ownerID, projectID uuid.UUID, req dto.UpdateProjectRequest) (dto.ProjectResponse, error)
	Delete(ctx context.Context, ownerID, projectID uuid.UUID) error
}

type ProjectService struct {
	projectRepo repository.IProjectRepository
}

func NewProjectService(projectRepo repository.IProjectRepository) *ProjectService {
	return &ProjectService{projectRepo: projectRepo}
}

func (s *ProjectService) Create(ctx context.Context, ownerID uuid.UUID, req dto.CreateProjectRequest) (dto.ProjectResponse, error) {
	project := entity.Project{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     ownerID,
	}

	created, err := s.projectRepo.Create(ctx, project)
	if err != nil {
		return dto.ProjectResponse{}, err
	}

	return toProjectResponse(created, 0), nil
}

func (s *ProjectService) List(ctx context.Context, ownerID uuid.UUID, page dto.PaginationQuery) ([]dto.ProjectResponse, int64, error) {
	projects, total, err := s.projectRepo.FindByOwner(ctx, ownerID, page.Offset(), page.Limit)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.ProjectResponse, 0, len(projects))
	for _, p := range projects {
		count, _ := s.projectRepo.TaskCount(ctx, p.ID)
		responses = append(responses, toProjectResponse(p, count))
	}

	return responses, total, nil
}

func (s *ProjectService) GetByID(ctx context.Context, ownerID, projectID uuid.UUID) (dto.ProjectResponse, error) {
	project, err := s.mustOwnProject(ctx, ownerID, projectID)
	if err != nil {
		return dto.ProjectResponse{}, err
	}

	count, err := s.projectRepo.TaskCount(ctx, project.ID)
	if err != nil {
		return dto.ProjectResponse{}, err
	}

	return toProjectResponse(project, count), nil
}

func (s *ProjectService) Update(ctx context.Context, ownerID, projectID uuid.UUID, req dto.UpdateProjectRequest) (dto.ProjectResponse, error) {
	project, err := s.mustOwnProject(ctx, ownerID, projectID)
	if err != nil {
		return dto.ProjectResponse{}, err
	}

	if req.Name != nil {
		project.Name = *req.Name
	}
	if req.Description != nil {
		project.Description = *req.Description
	}

	updated, err := s.projectRepo.Update(ctx, project)
	if err != nil {
		return dto.ProjectResponse{}, err
	}

	return toProjectResponse(updated, 0), nil
}

func (s *ProjectService) Delete(ctx context.Context, ownerID, projectID uuid.UUID) error {
	if _, err := s.mustOwnProject(ctx, ownerID, projectID); err != nil {
		return err
	}
	return s.projectRepo.Delete(ctx, projectID)
}

// mustOwnProject memusatkan pengecekan "apakah project ini ada DAN milik
// user ini" agar setiap method menerapkan kepemilikan dengan cara yang sama.
func (s *ProjectService) mustOwnProject(ctx context.Context, ownerID, projectID uuid.UUID) (entity.Project, error) {
	project, found, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return entity.Project{}, err
	}
	if !found {
		return entity.Project{}, dto.ErrProjectNotFound
	}
	if project.OwnerID != ownerID {
		return entity.Project{}, dto.ErrProjectForbidden
	}
	return project, nil
}

func toProjectResponse(p entity.Project, taskCount int64) dto.ProjectResponse {
	return dto.ProjectResponse{
		ID:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		OwnerID:     p.OwnerID.String(),
		TaskCount:   taskCount,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
