package users_service

import (
	"context"
	"fmt"
)

func (s *UsersService) DeleteUserById(
	ctx context.Context,
	id int,
) error {
	if err := s.usersRepository.DeleteUserById(ctx, id); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}
