package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CurrentUserID membaca key "user_id" yang disimpan auth middleware di
// context dan mengubahnya menjadi uuid.UUID. Handler memanggil ini alih-alih
// langsung menggunakan ctx.MustGet, sehingga nilai yang hilang/rusak menjadi
// error biasa bukan panic.
func CurrentUserID(ctx *gin.Context) (uuid.UUID, error) {
	raw, exists := ctx.Get("user_id")
	if !exists {
		return uuid.Nil, errors.New("user_id tidak ditemukan di context")
	}

	str, ok := raw.(string)
	if !ok {
		return uuid.Nil, errors.New("user_id di context bukan string")
	}

	return uuid.Parse(str)
}
