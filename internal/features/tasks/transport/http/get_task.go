package tasks_transport_http

import (
	"net/http"

	core_http_request "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/request"
	_ "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/response"
)

type GetTaskResponse TaskDTOResponse

// GetUser 	godoc
// @Summary 	Get task
// @Description Get task by id
// @Tags 		tasks
// @Produce 	json
// @Param 		task_id path int true "Task id"
// @Success 	200 {object} GetTaskResponse "Successfully got task"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 	404 {object} core_http_response.ErrorResponse "Task not found"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 		/tasks/{task_id} [get]
func (h *TasksHttpHandler) GetTask(rw http.ResponseWriter, r *http.Request) {
	ctx, responseHandler := core_http_request.NewContext(rw, r)

	taskId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"Missing task id",
		)
		return
	}

	taskDomain, err := h.tasksService.GetTask(ctx, taskId)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get task",
		)
		return
	}

	response := GetTaskResponse(taskDTOFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}
