package service

import (
	"fmt"
	"os"
	"time"

	"github.com/Zelviarani16/taskflow-api/dto"
	"github.com/golang-jwt/jwt/v5"
)

type IJWTService interface {
	GenerateToken(userID string, role string) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	GetUserIDByToken(tokenString string) (string, error)
}

type jwtCustomClaim struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secretKey []byte
	issuer    string
}

func NewJWTService() *JWTService {
	return &JWTService{
		secretKey: []byte(os.Getenv("JWT_SECRET")),
		issuer:    "taskflow-api",
	}
}

func (j *JWTService) GenerateToken(userID string, role string) (string, error) {
	claims := jwtCustomClaim{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    j.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", dto.ErrGenerateToken
	}
	return signed, nil
}

func (j *JWTService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &jwtCustomClaim{}, func(t *jwt.Token) (any, error) {
		// Tolak apapun yang tidak ditandatangani dengan HMAC untuk mencegah
		// serangan "algorithm confusion" di mana token dipalsukan menggunakan
		// metode penandatanganan yang berbeda dari yang diharapkan server.
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return j.secretKey, nil
	})
}

func (j *JWTService) GetUserIDByToken(tokenString string) (string, error) {
	token, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", dto.ErrTokenInvalid
	}

	claims, ok := token.Claims.(*jwtCustomClaim)
	if !ok || !token.Valid {
		return "", dto.ErrTokenInvalid
	}

	return claims.UserID, nil
}
