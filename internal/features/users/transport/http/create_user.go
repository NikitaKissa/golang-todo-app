package users_transport_http

import (
	"net/http"

	"github.com/NikitaKissa/golang-todo-app/internal/core/domain"
	core_http_request "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/request"
	_ "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name"    validate:"required,min=3,max=100"               example:"Jack Daniels"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+" example:"+48123123123"`
}

type CreateUserResponse UserDTOResponse

// CreateUser 	godoc
// @Summary 	Create user
// @Description Create new user in system
// @Tags 		users
// @Accept 		json
// @Produce 	json
// @Param 		request body CreateUserRequest true "CreateUserRequest request body"
// @Success 	201 {object} CreateUserResponse "Successfully created user"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 		/users [post]
func (h *UsersHttpHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx, responseHandler := core_http_request.NewContext(rw, r)

	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userDomain := domainFromDTO(request)

	userDomain, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")
		return
	}

	response := CreateUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUninitialized(dto.FullName, dto.PhoneNumber)
}
