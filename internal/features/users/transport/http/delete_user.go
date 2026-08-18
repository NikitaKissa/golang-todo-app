package users_transport_http

import (
	"net/http"

	core_logger "github.com/NikitaKissa/golang-todo-app/internal/core/logger"
	core_http_request "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/response"
)

// DeleteUser 	godoc
// @Summary 	Delete user
// @Description Delete user in system
// @Tags 		users
// @Param 		user_id path int true "User id"
// @Success 	204 "Successfully deleted user"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 	404 {object} core_http_response.ErrorResponse "User not found"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 		/users/{user_id} [delete]
func (h *UsersHttpHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
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

	if err := h.usersService.DeleteUserById(ctx, userId); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete user",
		)
		return
	}

	responseHandler.DeleteResponse()
}
