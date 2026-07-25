package dto

import "errors"

// ===== Request =====

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// ===== Response =====

type AuthResponse struct {
	Token string `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// ===== Error =====
// Error sentinel memungkinkan lapisan service menandakan *apa* yang salah
// dan lapisan handler menentukan status HTTP, tanpa mencocokkan string pesan.
var (
	ErrEmailAlreadyExists = errors.New("email sudah terdaftar")
	ErrInvalidCredential  = errors.New("email atau password salah")
	ErrUserNotFound       = errors.New("user tidak ditemukan")
	ErrGenerateToken      = errors.New("gagal generate token")
	ErrTokenInvalid       = errors.New("token tidak valid atau kadaluwarsa")
	ErrUnauthorized       = errors.New("token tidak ditemukan atau format salah")
)
