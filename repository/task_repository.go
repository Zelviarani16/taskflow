package repository

import (
	"context"
	"errors"

	"github.com/Zelviarani16/taskflow-api/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ITaskRepository interface {
	Create(ctx context.Context, task entity.Task) (entity.Task, error)
	FindByID(ctx context.Context, id uuid.UUID) (entity.Task, bool, error)
	FindByProject(ctx context.Context, projectID uuid.UUID, status, priority string, offset, limit int) ([]entity.Task, int64, error)
	Update(ctx context.Context, task entity.Task) (entity.Task, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task entity.Task) (entity.Task, error) {
	if err := r.db.WithContext(ctx).Create(&task).Error; err != nil {
		return entity.Task{}, err
	}
	return task, nil
}

func (r *TaskRepository) FindByID(ctx context.Context, id uuid.UUID) (entity.Task, bool, error) {
	var task entity.Task
	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Task{}, false, nil
		}
		return entity.Task{}, false, err
	}
	return task, true, nil
}

// FindByProject mendukung filter status/prioritas opsional - string kosong
// berarti "jangan filter berdasarkan field ini".
func (r *TaskRepository) FindByProject(ctx context.Context, projectID uuid.UUID, status, priority string, offset, limit int) ([]entity.Task, int64, error) {
	var tasks []entity.Task
	var total int64

	base := r.db.WithContext(ctx).Model(&entity.Task{}).Where("project_id = ?", projectID)
	if status != "" {
		base = base.Where("status = ?", status)
	}
	if priority != "" {
		base = base.Where("priority = ?", priority)
	}

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := base.Order("created_at DESC").Offset(offset).Limit(limit).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *TaskRepository) Update(ctx context.Context, task entity.Task) (entity.Task, error) {
	if err := r.db.WithContext(ctx).Save(&task).Error; err != nil {
		return entity.Task{}, err
	}
	return task, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Task{}, "id = ?", id).Error
}
