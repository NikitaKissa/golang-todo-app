package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/NikitaKissa/golang-todo-app/internal/core/domain"
	core_http_request "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/request"
	core_http_types "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title"`
	Description core_http_types.Nullable[string] `json:"description"`
	Completed   core_http_types.Nullable[bool]   `json:"completed"`
}

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("`Title` can't be null")
		}

		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("`Title` must be between 1 and 100 chars")
		}
	}

	if r.Description.Set && r.Description.Value != nil {
		descriptionLen := len([]rune(*r.Description.Value))
		if descriptionLen < 1 || descriptionLen > 1000 {
			return fmt.Errorf("`Description` must be between 1 and 1000 chars")
		}
	}

	if r.Completed.Set && r.Completed.Value == nil {
		return fmt.Errorf("`Completed` cant't be null")
	}

	return nil
}

type PatchTaskResponse TaskDTOResponse

func (h *TasksHttpHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx, responseHandler := core_http_request.NewContext(rw, r)

	taskId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"Missing task id",
		)
		return
	}

	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"Failed to decode and validate Http request",
		)
		return
	}

	taskPatch := taskPatchFromRequest(request)

	taskDomain, err := h.tasksService.PatchTask(ctx, taskId, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch task",
		)
		return
	}

	response := PatchTaskResponse(taskDTOFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusOK)

}

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
