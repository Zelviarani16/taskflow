package service

import (
	"context"

	"github.com/Zelviarani16/taskflow-api/dto"
	"github.com/Zelviarani16/taskflow-api/entity"
	"github.com/Zelviarani16/taskflow-api/helpers"
	"github.com/Zelviarani16/taskflow-api/repository"
)

type IAuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (dto.AuthResponse, error)
}

type AuthService struct {
	userRepo   repository.IUserRepository
	jwtService IJWTService
}

func NewAuthService(userRepo repository.IUserRepository, jwtService IJWTService) *AuthService {
	return &AuthService{userRepo: userRepo, jwtService: jwtService}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (dto.AuthResponse, error) {
	_, found, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	if found {
		return dto.AuthResponse{}, dto.ErrEmailAlreadyExists
	}

	hashed, err := helpers.HashPassword(req.Password)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	newUser := entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashed,
		Role:     entity.RoleMember,
	}

	created, err := s.userRepo.Create(ctx, newUser)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return s.buildAuthResponse(created)
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (dto.AuthResponse, error) {
	user, found, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	// Sengaja mengembalikan error yang sama baik email tidak ada
	// maupun password salah - jangan bocorkan kasus mana yang terjadi.
	if !found {
		return dto.AuthResponse{}, dto.ErrInvalidCredential
	}
	if !helpers.CheckPassword(user.Password, req.Password) {
		return dto.AuthResponse{}, dto.ErrInvalidCredential
	}

	return s.buildAuthResponse(user)
}

func (s *AuthService) buildAuthResponse(user entity.User) (dto.AuthResponse, error) {
	token, err := s.jwtService.GenerateToken(user.ID.String(), string(user.Role))
	if err != nil {
		return dto.AuthResponse{}, dto.ErrGenerateToken
	}

	return dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
			Role:  string(user.Role),
		},
	}, nil
}
