package handler

import (
	"net/http"

	"github.com/Zelviarani16/taskflow-api/constants"
	"github.com/Zelviarani16/taskflow-api/dto"
	"github.com/Zelviarani16/taskflow-api/service"
	"github.com/Zelviarani16/taskflow-api/utils"
	"github.com/gin-gonic/gin"
)

type IAuthHandler interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type AuthHandler struct {
	authService service.IAuthService
}

func NewAuthHandler(authService service.IAuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register menangani POST /api/v1/auth/register
func (h *AuthHandler) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BuildError(constants.MsgFailedBadRequest+": "+err.Error()))
		return
	}

	result, err := h.authService.Register(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(utils.StatusFromError(err), utils.BuildError(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.BuildSuccess(constants.MsgSuccessRegister, result))
}

// Login menangani POST /api/v1/auth/login
func (h *AuthHandler) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BuildError(constants.MsgFailedBadRequest+": "+err.Error()))
		return
	}

	result, err := h.authService.Login(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(utils.StatusFromError(err), utils.BuildError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess(constants.MsgSuccessLogin, result))
}
