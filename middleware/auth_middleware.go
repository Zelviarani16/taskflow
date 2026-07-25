package middleware

import (
	"net/http"
	"strings"

	"github.com/Zelviarani16/taskflow-api/constants"
	"github.com/Zelviarani16/taskflow-api/dto"
	"github.com/Zelviarani16/taskflow-api/service"
	"github.com/Zelviarani16/taskflow-api/utils"
	"github.com/gin-gonic/gin"
)

// Authentication memvalidasi header "Authorization: Bearer <token>" dan,
// jika berhasil, menyimpan user id di context agar handler bisa membacanya
// melalui utils.CurrentUserID. Setiap grup endpoint yang butuh user yang
// sudah login mendapatkan middleware ini di routes.
func Authentication(jwtService service.IJWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			abortUnauthorized(ctx, dto.ErrUnauthorized.Error())
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			abortUnauthorized(ctx, dto.ErrUnauthorized.Error())
			return
		}

		tokenString := parts[1]

		token, err := jwtService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			abortUnauthorized(ctx, dto.ErrTokenInvalid.Error())
			return
		}

		userID, err := jwtService.GetUserIDByToken(tokenString)
		if err != nil {
			abortUnauthorized(ctx, dto.ErrTokenInvalid.Error())
			return
		}

		ctx.Set("user_id", userID)
		ctx.Next()
	}
}

func abortUnauthorized(ctx *gin.Context, reason string) {
	ctx.AbortWithStatusJSON(http.StatusUnauthorized,
		utils.BuildError(constants.MsgFailedUnauthorized+": "+reason),
	)
}
