package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/NikitaKissa/golang-todo-app/internal/core/logger"
	core_http_request "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/response"
)

type GetUsersResponse []UserDTOResponse

func (h *UsersHttpHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(log, w)
	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'limit' or 'offset' query parameter",
		)
		return
	}

	usersDomains, err := h.usersService.GetUsers(r.Context(), limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get users",
		)
		return
	}

	response := GetUsersResponse(usersDTOFromDomains(usersDomains))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	const (
		limitParamKey  = "limit"
		offsetParamKey = "offset"
	)

	limit, err := core_http_request.GetIntQueryParams(r, limitParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param", err)
	}
	offset, err := core_http_request.GetIntQueryParams(r, offsetParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param", err)
	}

	return limit, offset, nil
}
