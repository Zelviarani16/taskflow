package service

import (
	"context"

	"github.com/Zelviarani16/taskflow-api/dto"
	"github.com/Zelviarani16/taskflow-api/entity"
	"github.com/Zelviarani16/taskflow-api/repository"
	"github.com/google/uuid"
)

type ITaskService interface {
	Create(ctx context.Context, ownerID, projectID uuid.UUID, req dto.CreateTaskRequest) (dto.TaskResponse, error)
	List(ctx context.Context, ownerID, projectID uuid.UUID, filter dto.TaskFilter) ([]dto.TaskResponse, int64, error)
	Update(ctx context.Context, ownerID, projectID, taskID uuid.UUID, req dto.UpdateTaskRequest) (dto.TaskResponse, error)
	Delete(ctx context.Context, ownerID, projectID, taskID uuid.UUID) error
}

type TaskService struct {
	taskRepo    repository.ITaskRepository
	projectRepo repository.IProjectRepository
	userRepo    repository.IUserRepository
}

func NewTaskService(taskRepo repository.ITaskRepository, projectRepo repository.IProjectRepository, userRepo repository.IUserRepository) *TaskService {
	return &TaskService{taskRepo: taskRepo, projectRepo: projectRepo, userRepo: userRepo}
}

func (s *TaskService) Create(ctx context.Context, ownerID, projectID uuid.UUID, req dto.CreateTaskRequest) (dto.TaskResponse, error) {
	if err := s.mustOwnProject(ctx, ownerID, projectID); err != nil {
		return dto.TaskResponse{}, err
	}

	assignee, err := s.resolveAssignee(ctx, req.AssignedTo)
	if err != nil {
		return dto.TaskResponse{}, err
	}

	task := entity.Task{
		ProjectID:   projectID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    entity.TaskPriority(defaultString(req.Priority, string(entity.PriorityMedium))),
		DueDate:     req.DueDate,
		AssignedTo:  assignee,
	}

	created, err := s.taskRepo.Create(ctx, task)
	if err != nil {
		return dto.TaskResponse{}, err
	}

	return toTaskResponse(created), nil
}

func (s *TaskService) List(ctx context.Context, ownerID, projectID uuid.UUID, filter dto.TaskFilter) ([]dto.TaskResponse, int64, error) {
	if err := s.mustOwnProject(ctx, ownerID, projectID); err != nil {
		return nil, 0, err
	}

	tasks, total, err := s.taskRepo.FindByProject(ctx, projectID, filter.Status, filter.Priority, filter.Offset(), filter.Limit)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.TaskResponse, 0, len(tasks))
	for _, t := range tasks {
		responses = append(responses, toTaskResponse(t))
	}

	return responses, total, nil
}

func (s *TaskService) Update(ctx context.Context, ownerID, projectID, taskID uuid.UUID, req dto.UpdateTaskRequest) (dto.TaskResponse, error) {
	if err := s.mustOwnProject(ctx, ownerID, projectID); err != nil {
		return dto.TaskResponse{}, err
	}

	task, found, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return dto.TaskResponse{}, err
	}
	if !found || task.ProjectID != projectID {
		return dto.TaskResponse{}, dto.ErrTaskNotFound
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Status != nil {
		task.Status = entity.TaskStatus(*req.Status)
	}
	if req.Priority != nil {
		task.Priority = entity.TaskPriority(*req.Priority)
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}
	if req.AssignedTo != nil {
		assignee, err := s.resolveAssignee(ctx, req.AssignedTo)
		if err != nil {
			return dto.TaskResponse{}, err
		}
		task.AssignedTo = assignee
	}

	updated, err := s.taskRepo.Update(ctx, task)
	if err != nil {
		return dto.TaskResponse{}, err
	}

	return toTaskResponse(updated), nil
}

func (s *TaskService) Delete(ctx context.Context, ownerID, projectID, taskID uuid.UUID) error {
	if err := s.mustOwnProject(ctx, ownerID, projectID); err != nil {
		return err
	}

	task, found, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return err
	}
	if !found || task.ProjectID != projectID {
		return dto.ErrTaskNotFound
	}

	return s.taskRepo.Delete(ctx, taskID)
}

func (s *TaskService) mustOwnProject(ctx context.Context, ownerID, projectID uuid.UUID) error {
	project, found, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return err
	}
	if !found {
		return dto.ErrProjectNotFound
	}
	if project.OwnerID != ownerID {
		return dto.ErrProjectForbidden
	}
	return nil
}

// resolveAssignee mengubah string user-id opsional menjadi *uuid.UUID yang
// sudah divalidasi, memastikan user tersebut benar-benar ada sebelum task
// merujuk ke mereka.
func (s *TaskService) resolveAssignee(ctx context.Context, rawID *string) (*uuid.UUID, error) {
	if rawID == nil || *rawID == "" {
		return nil, nil
	}

	parsed, err := uuid.Parse(*rawID)
	if err != nil {
		return nil, dto.ErrAssigneeNotFound
	}

	_, found, err := s.userRepo.FindByID(ctx, parsed)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, dto.ErrAssigneeNotFound
	}

	return &parsed, nil
}

func defaultString(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}

func toTaskResponse(t entity.Task) dto.TaskResponse {
	var assignedTo *string
	if t.AssignedTo != nil {
		s := t.AssignedTo.String()
		assignedTo = &s
	}

	return dto.TaskResponse{
		ID:          t.ID.String(),
		ProjectID:   t.ProjectID.String(),
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		Priority:    string(t.Priority),
		DueDate:     t.DueDate,
		AssignedTo:  assignedTo,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
