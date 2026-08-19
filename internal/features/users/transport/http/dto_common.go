package users_transport_http

import "github.com/NikitaKissa/golang-todo-app/internal/core/domain"

type UserDTOResponse struct {
	ID          int     `json:"id"           example:"11"`
	Version     int     `json:"version"      example:"1"`
	FullName    string  `json:"full_name"    example:"Jack Daniels"`
	PhoneNumber *string `json:"phone_number" example:"+48123123123"`
}

func userDTOFromDomain(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:          user.ID,
		Version:     user.Version,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}
}

func usersDTOFromDomains(users []domain.User) []UserDTOResponse {
	usersDTO := make([]UserDTOResponse, len(users))

	for i, user := range users {
		usersDTO[i] = userDTOFromDomain(user)
	}

	return usersDTO
}
