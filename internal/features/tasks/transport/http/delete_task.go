package tasks_transport_http

import (
	"net/http"

	core_http_request "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/request"
)

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
