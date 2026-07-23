package tasks_transport_http

import (
	"net/http"

	core_http_server "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/server"
)

type TasksHttpHandler struct {
	tasksService TasksService
}

type TasksService interface {
}

func NewTasksHttpHandler(tasksService TasksService) *TasksHttpHandler {
	return &TasksHttpHandler{tasksService}
}

func (h *TasksHttpHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/tasks",
			Handler: h.CreateTask,
		},
	}
}
