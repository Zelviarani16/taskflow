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

type ITaskHandler interface {
	Create(ctx *gin.Context)
	List(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type TaskHandler struct {
	taskService service.ITaskService
}

func NewTaskHandler(taskService service.ITaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

// Create menangani POST /api/v1/projects/:id/tasks
func (h *TaskHandler) Create(ctx *gin.Context) {
	ownerID, projectID, ok := h.identity(ctx)
	if !ok {
		return
	}

	var req dto.CreateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BuildError(constants.MsgFailedBadRequest+": "+err.Error()))
		return
	}

	result, err := h.taskService.Create(ctx.Request.Context(), ownerID, projectID, req)
	if err != nil {
		ctx.JSON(utils.StatusFromError(err), utils.BuildError(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.BuildSuccess(constants.MsgSuccessCreated, result))
}

// List menangani GET /api/v1/projects/:id/tasks?status=&priority=&page=&limit=
func (h *TaskHandler) List(ctx *gin.Context) {
	ownerID, projectID, ok := h.identity(ctx)
	if !ok {
		return
	}

	var filter dto.TaskFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BuildError(constants.MsgFailedBadRequest+": "+err.Error()))
		return
	}
	filter.Normalize()

	results, total, err := h.taskService.List(ctx.Request.Context(), ownerID, projectID, filter)
	if err != nil {
		ctx.JSON(utils.StatusFromError(err), utils.BuildError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess(constants.MsgSuccessFetched, dto.BuildPaginatedResponse(results, filter.PaginationQuery, total)))
}

// Update menangani PUT /api/v1/projects/:id/tasks/:taskId
func (h *TaskHandler) Update(ctx *gin.Context) {
	ownerID, projectID, ok := h.identity(ctx)
	if !ok {
		return
	}

	taskID, err := uuid.Parse(ctx.Param("taskId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BuildError(constants.MsgFailedBadRequest+": id task tidak valid"))
		return
	}

	var req dto.UpdateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BuildError(constants.MsgFailedBadRequest+": "+err.Error()))
		return
	}

	result, err := h.taskService.Update(ctx.Request.Context(), ownerID, projectID, taskID, req)
	if err != nil {
		ctx.JSON(utils.StatusFromError(err), utils.BuildError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess(constants.MsgSuccessUpdated, result))
}

// Delete menangani DELETE /api/v1/projects/:id/tasks/:taskId
func (h *TaskHandler) Delete(ctx *gin.Context) {
	ownerID, projectID, ok := h.identity(ctx)
	if !ok {
		return
	}

	taskID, err := uuid.Parse(ctx.Param("taskId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BuildError(constants.MsgFailedBadRequest+": id task tidak valid"))
		return
	}

	if err := h.taskService.Delete(ctx.Request.Context(), ownerID, projectID, taskID); err != nil {
		ctx.JSON(utils.StatusFromError(err), utils.BuildError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess(constants.MsgSuccessDeleted, nil))
}

// identity mengambil ID user yang sudah terautentikasi dan parameter path :id
// (project id) yang menjadi induk setiap route task. Menulis response error
// sendiri dan mengembalikan ok=false jika salah satunya tidak ada/tidak valid.
func (h *TaskHandler) identity(ctx *gin.Context) (uuid.UUID, uuid.UUID, bool) {
	ownerID, err := utils.CurrentUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.BuildError(constants.MsgFailedUnauthorized))
		return uuid.Nil, uuid.Nil, false
	}

	projectID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BuildError(constants.MsgFailedBadRequest+": id project tidak valid"))
		return uuid.Nil, uuid.Nil, false
	}

	return ownerID, projectID, true
}
