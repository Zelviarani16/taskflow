package repository

import (
	"context"
	"errors"

	"github.com/Zelviarani16/taskflow-api/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IProjectRepository interface {
	Create(ctx context.Context, project entity.Project) (entity.Project, error)
	FindByID(ctx context.Context, id uuid.UUID) (entity.Project, bool, error)
	FindByOwner(ctx context.Context, ownerID uuid.UUID, offset, limit int) ([]entity.Project, int64, error)
	TaskCount(ctx context.Context, projectID uuid.UUID) (int64, error)
	Update(ctx context.Context, project entity.Project) (entity.Project, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(ctx context.Context, project entity.Project) (entity.Project, error) {
	if err := r.db.WithContext(ctx).Create(&project).Error; err != nil {
		return entity.Project{}, err
	}
	return project, nil
}

func (r *ProjectRepository) FindByID(ctx context.Context, id uuid.UUID) (entity.Project, bool, error) {
	var project entity.Project
	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&project).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Project{}, false, nil
		}
		return entity.Project{}, false, err
	}
	return project, true, nil
}

// FindByOwner mengembalikan satu halaman project milik owner tertentu beserta
// total jumlah (untuk metadata paginasi), diurutkan dari yang terbaru.
func (r *ProjectRepository) FindByOwner(ctx context.Context, ownerID uuid.UUID, offset, limit int) ([]entity.Project, int64, error) {
	var projects []entity.Project
	var total int64

	base := r.db.WithContext(ctx).Model(&entity.Project{}).Where("owner_id = ?", ownerID)

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := base.Order("created_at DESC").Offset(offset).Limit(limit).Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

func (r *ProjectRepository) TaskCount(ctx context.Context, projectID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Task{}).Where("project_id = ?", projectID).Count(&count).Error
	return count, err
}

func (r *ProjectRepository) Update(ctx context.Context, project entity.Project) (entity.Project, error) {
	if err := r.db.WithContext(ctx).Save(&project).Error; err != nil {
		return entity.Project{}, err
	}
	return project, nil
}

// Delete adalah soft delete - GORM mengisi deleted_at alih-alih menghapus
// baris, sehingga task di bawah project tetap menyimpan riwayatnya.
func (r *ProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Project{}, "id = ?", id).Error
}
