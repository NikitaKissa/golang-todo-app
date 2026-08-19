package tasks_transport_http

import (
	"net/http"

	core_http_request "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/request"
	_ "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/response"
)

// DeleteTask 	godoc
// @Summary 	Delete task
// @Description Delete task in system
// @Tags 		tasks
// @Param 		task_id path int true "Task id"
// @Success 	204 "Successfully deleted task"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 	404 {object} core_http_response.ErrorResponse "Task not found"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 		/tasks/{task_id} [delete]
func (h *TasksHttpHandler) DeleteTask(rw http.ResponseWriter, r *http.Request) {
	ctx, responseHandler := core_http_request.NewContext(rw, r)

	taskId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"Missing task id",
		)
		return
	}

	if err := h.tasksService.DeleteTask(ctx, taskId); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete task",
		)
		return
	}

	responseHandler.DeleteResponse()
}
