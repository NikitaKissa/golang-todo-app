package users_transport_http

import (
	"context"
	"net/http"

	"github.com/NikitaKissa/golang-todo-app/internal/core/domain"
	core_http_server "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/server"
)

type UsersHttpHandler struct {
	usersService UsersService
}

type UsersService interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)

	GetUserById(
		ctx context.Context,
		id int,
	) (domain.User, error)

	DeleteUserById(
		ctx context.Context,
		id int,
	) error

	PatchUser(
		ctx context.Context,
		id int,
		patch domain.UserPatch,
	) (domain.User, error)
}

func NewUsersHttpHandler(usersService UsersService) *UsersHttpHandler {
	return &UsersHttpHandler{
		usersService: usersService,
	}
}

func (h *UsersHttpHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: h.GetUsers,
			//Middleware: []core_http_middleware.Middleware{
			//	core_http_middleware.Dummy("get users middleware"),
			//},
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/{id}",
			Handler: h.GetUser,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/users/{id}",
			Handler: h.DeleteUser,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/users/{id}",
			Handler: h.PatchUser,
		},
	}
}
