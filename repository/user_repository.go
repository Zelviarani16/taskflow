package repository

import (
	"context"
	"errors"

	"github.com/Zelviarani16/taskflow-api/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// IUserRepository adalah kontrak yang dipakai oleh lapisan service. Service
// hanya tahu "ada FindByEmail dan Create" - mereka tidak pernah melihat SQL.
type IUserRepository interface {
	FindByEmail(ctx context.Context, email string) (entity.User, bool, error)
	FindByID(ctx context.Context, id uuid.UUID) (entity.User, bool, error)
	Create(ctx context.Context, user entity.User) (entity.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByEmail mengembalikan (user, ditemukan, error). "tidak ditemukan" dan
// "error DB" sengaja dibedakan - baris yang tidak ada bukan sebuah kegagalan.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (entity.User, bool, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("email = ?", email).Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, false, nil
		}
		return entity.User{}, false, err
	}
	return user, true, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (entity.User, bool, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, false, nil
		}
		return entity.User{}, false, err
	}
	return user, true, nil
}

func (r *UserRepository) Create(ctx context.Context, user entity.User) (entity.User, error) {
	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		return entity.User{}, err
	}
	return user, nil
}
