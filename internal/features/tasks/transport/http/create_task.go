package tasks_transport_http

import (
	"net/http"

	"github.com/NikitaKissa/golang-todo-app/internal/core/domain"
	core_http_request "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/request"
	_ "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/response"
)

type CreateTaskRequest struct {
	Title        string  `json:"title" validate:"required,min=1,max=100"`
	Description  *string `json:"description" validate:"omitempty,min=1,max=1000"`
	AuthorUserID int     `json:"author_user_id" validate:"required"`
}

type CreateTaskResponse TaskDTOResponse

// CreateTask 	godoc
// @Summary 	Create task
// @Description Create new task in system
// @Tags 		tasks
// @Accept 		json
// @Produce 	json
// @Param 		request body CreateTaskRequest true "CreateTaskRequest request body"
// @Success 	201 {object} CreateTaskResponse "Successfully created task"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 	404 {object} core_http_response.ErrorResponse "Author user not found"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 		/tasks [post]
func (h *TasksHttpHandler) CreateTask(rw http.ResponseWriter, r *http.Request) {
	ctx, responseHandler := core_http_request.NewContext(rw, r)

	var request CreateTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate http request",
		)
		return
	}

	taskDomain := domain.NewTaskUninitialized(
		request.Title,
		request.Description,
		request.AuthorUserID,
	)

	taskDomain, err := h.tasksService.CreateTask(ctx, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create task",
		)
		return
	}

	response := CreateTaskResponse(taskDTOFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusCreated)

}
