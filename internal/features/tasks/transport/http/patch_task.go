package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/NikitaKissa/golang-todo-app/internal/core/domain"
	core_http_request "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/request"
	_ "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/response"
	core_http_types "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title"       swaggertype:"string" example:"Bake a cake"`
	Description core_http_types.Nullable[string] `json:"description" swaggertype:"string" example:"With wildberries"`
	Completed   core_http_types.Nullable[bool]   `json:"completed"   swaggertype:"boolean" example:"true"`
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

// PatchTask 	godoc
// @Summary 	Patch task
// @Description Patch task in system
// @Description ### Three-state logic of fields updating
// @Description 1. **Field is unsigned**: `description` is ignored
// @Description 2. **Field is provided**: `description` is updated
// @Description 3. **Field is null**: `description` is deleted
// @Description **Invalid state:** `completed`/`title` set to null
// @Tags 		tasks
// @Accept 		json
// @Produce 	json
// @Param 		task_id path int true "Task id"
// @Param 		request body PatchTaskRequest true "PatchTaskRequest request body"
// @Success 	200 {object} PatchTaskResponse "Successfully patched task"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 	404 {object} core_http_response.ErrorResponse "Task not found"
// @Failure 	409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 		/task/{task_id} [patch]
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
