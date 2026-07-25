package handler

import (
	"net/http"

	"github.com/Zelviarani16/taskflow-api/constants"
	"github.com/Zelviarani16/taskflow-api/dto"
	"github.com/Zelviarani16/taskflow-api/service"
	"github.com/Zelviarani16/taskflow-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IProjectHandler interface {
	Create(ctx *gin.Context)
	List(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type ProjectHandler struct {
	projectService service.IProjectService
}

func NewProjectHandler(projectService service.IProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}

// Create menangani POST /api/v1/projects
func (h *ProjectHandler) Create(ctx *gin.Context) {
	ownerID, err := utils.CurrentUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.BuildError(constants.MsgFailedUnauthorized))
		return
	}

	var req dto.CreateProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BuildError(constants.MsgFailedBadRequest+": "+err.Error()))
		return
	}

	result, err := h.projectService.Create(ctx.Request.Context(), ownerID, req)
	if err != nil {
		ctx.JSON(utils.StatusFromError(err), utils.BuildError(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.BuildSuccess(constants.MsgSuccessCreated, result))
}

// List menangani GET /api/v1/projects
func (h *ProjectHandler) List(ctx *gin.Context) {
	ownerID, err := utils.CurrentUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.BuildError(constants.MsgFailedUnauthorized))
		return
	}

	var page dto.PaginationQuery
	_ = ctx.ShouldBindQuery(&page)
	page.Normalize()

	results, total, err := h.projectService.List(ctx.Request.Context(), ownerID, page)
	if err != nil {
		ctx.JSON(utils.StatusFromError(err), utils.BuildError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess(constants.MsgSuccessFetched, dto.BuildPaginatedResponse(results, page, total)))
}

// GetByID menangani GET /api/v1/projects/:id
func (h *ProjectHandler) GetByID(ctx *gin.Context) {
	ownerID, err := utils.CurrentUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.BuildError(constants.MsgFailedUnauthorized))
		return
	}

	projectID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BuildError(constants.MsgFailedBadRequest+": id project tidak valid"))
		return
	}

	result, err := h.projectService.GetByID(ctx.Request.Context(), ownerID, projectID)
	if err != nil {
		ctx.JSON(utils.StatusFromError(err), utils.BuildError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess(constants.MsgSuccessFetched, result))
}

// Update menangani PUT /api/v1/projects/:id
func (h *ProjectHandler) Update(ctx *gin.Context) {
	ownerID, err := utils.CurrentUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.BuildError(constants.MsgFailedUnauthorized))
		return
	}

	projectID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BuildError(constants.MsgFailedBadRequest+": id project tidak valid"))
		return
	}

	var req dto.UpdateProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BuildError(constants.MsgFailedBadRequest+": "+err.Error()))
		return
	}

	result, err := h.projectService.Update(ctx.Request.Context(), ownerID, projectID, req)
	if err != nil {
		ctx.JSON(utils.StatusFromError(err), utils.BuildError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess(constants.MsgSuccessUpdated, result))
}

// Delete menangani DELETE /api/v1/projects/:id
func (h *ProjectHandler) Delete(ctx *gin.Context) {
	ownerID, err := utils.CurrentUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.BuildError(constants.MsgFailedUnauthorized))
		return
	}

	projectID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BuildError(constants.MsgFailedBadRequest+": id project tidak valid"))
		return
	}

	if err := h.projectService.Delete(ctx.Request.Context(), ownerID, projectID); err != nil {
		ctx.JSON(utils.StatusFromError(err), utils.BuildError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess(constants.MsgSuccessDeleted, nil))
}
