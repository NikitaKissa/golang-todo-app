package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_http_request "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/request"
)

type GetTasksResponse []TaskDTOResponse

func (h *TasksHttpHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx, responseHandler := core_http_request.NewContext(rw, r)

	userId, limit, offset, err := getUserIdLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'limit' or 'offset' or 'user_id' query parameter",
		)
		return
	}

	tasksDomains, err := h.tasksService.GetTasks(ctx, userId, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get tasks",
		)
		return
	}

	response := GetTasksResponse(tasksDTOFromDomains(tasksDomains))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func getUserIdLimitOffsetQueryParams(r *http.Request) (*int, *int, *int, error) {
	const (
		userIdKey      = "user_id"
		limitParamKey  = "limit"
		offsetParamKey = "offset"
	)

	userId, err := core_http_request.GetIntQueryParams(r, userIdKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	limit, err := core_http_request.GetIntQueryParams(r, limitParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParams(r, offsetParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return userId, limit, offset, nil
}
