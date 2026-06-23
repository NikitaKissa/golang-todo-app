package users_transport_http

import (
	"net/http"

	core_logger "github.com/NikitaKissa/golang-todo-app/internal/core/logger"
	core_http_response "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/response"
	core_http_utils "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/utils"
)

type GetUserResponse UserDTOResponse

func (h *UsersHttpHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(log, w)

	userId, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"Missing user id",
		)
		return
	}

	userDomain, err := h.usersService.GetUserById(ctx, userId)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get user",
		)
		return
	}

	response := GetUserResponse(userDTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}
