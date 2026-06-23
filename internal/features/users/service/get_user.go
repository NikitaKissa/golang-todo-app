package users_service

import (
	"context"
	"fmt"

	"github.com/NikitaKissa/golang-todo-app/internal/core/domain"
)

func (s *UsersService) GetUserById(
	ctx context.Context,
	id int,
) (domain.User, error) {
	user, err := s.usersRepository.GetUserById(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user from repository: %w", err)
	}

	return user, nil
}
