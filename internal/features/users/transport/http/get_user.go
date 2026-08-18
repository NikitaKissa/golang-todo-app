package users_transport_http

import (
	"net/http"

	core_logger "github.com/NikitaKissa/golang-todo-app/internal/core/logger"
	core_http_request "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/response"
)

type GetUserResponse UserDTOResponse

// GetUser 	godoc
// @Summary 	Get user
// @Description Get user by id
// @Tags 		users
// @Produce 	json
// @Param 		user_id path int true "User id"
// @Success 	200 {object} GetUserResponse "Successfully got user"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 	404 {object} core_http_response.ErrorResponse "Not found"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 		/users/{user_id} [get]
func (h *UsersHttpHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(log, w)

	userId, err := core_http_request.GetIntPathValue(r, "id")
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
