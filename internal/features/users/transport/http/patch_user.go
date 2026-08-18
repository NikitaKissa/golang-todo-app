package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/NikitaKissa/golang-todo-app/internal/core/domain"
	core_http_request "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/request"
	_ "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/response"
	core_http_types "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"    swaggertype:"string" example:"Jack Yard"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" swaggertype:"string" example:"+48123123456"`
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("'full_name' can't be NULL")
		}

		fullNameLen := len([]rune(*r.FullName.Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("'full_name' length must be between 3 and 100 symbols")
		}
	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phoneNumberLen := len([]rune(*r.PhoneNumber.Value))
			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf("'phone_number' length must be between 10 and 15 symbols")
			}

			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf("'phone_number' must start from '+' symbol")
			}
		}
	}

	return nil
}

type PatchUserResponse UserDTOResponse

// CreateUser 	godoc
// @Summary 	Patch user
// @Description Patch user in system
// @Description ### Three-state logic of fields updating
// @Description 1. **Field is unsigned**: `phone_number` is ignored
// @Description 2. **Field is provided**: `phone_number` is updated
// @Description 3. **Field is null**: `phone_number` is deleted
// @Description **Invalid state:** `full_name` set to null
// @Tags 		users
// @Accept 		json
// @Produce 	json
// @Param 		user_id path int true "User id"
// @Param 		request body PatchUserRequest true "CreateUserRequest request body"
// @Success 	200 {object} PatchUserResponse "Successfully created user"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 	404 {object} core_http_response.ErrorResponse "Not found"
// @Failure 	409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 		/users/{user_id} [patch]
func (h *UsersHttpHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx, responseHandler := core_http_request.NewContext(rw, r)

	userId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"Missing user id",
		)
		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userId, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)
		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.FullName.ToDomain(),
		request.PhoneNumber.ToDomain(),
	)
}
