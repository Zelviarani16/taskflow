package utils

import (
	"errors"
	"net/http"

	"github.com/Zelviarani16/taskflow-api/dto"
)

// StatusFromError memeriksa error dari lapisan service dan memilih status
// HTTP yang sesuai. Apapun yang tidak dikenali akan fallback ke 500, karena
// error yang tidak dikenali berarti sesuatu yang tidak terduga terjadi di server.
func StatusFromError(err error) int {
	switch {
	case errors.Is(err, dto.ErrEmailAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, dto.ErrInvalidCredential):
		return http.StatusUnauthorized
	case errors.Is(err, dto.ErrUserNotFound),
		errors.Is(err, dto.ErrProjectNotFound),
		errors.Is(err, dto.ErrTaskNotFound),
		errors.Is(err, dto.ErrAssigneeNotFound):
		return http.StatusNotFound
	case errors.Is(err, dto.ErrProjectForbidden):
		return http.StatusForbidden
	case errors.Is(err, dto.ErrTokenInvalid), errors.Is(err, dto.ErrUnauthorized):
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
