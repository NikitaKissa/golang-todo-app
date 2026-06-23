package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/NikitaKissa/golang-todo-app/internal/core/domain"
	core_errors "github.com/NikitaKissa/golang-todo-app/internal/core/errors"
	"github.com/jackc/pgx/v5"
)

func (r *UsersRepository) PatchUser(
	ctx context.Context,
	id int,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE todoapp.users
	SET
	    full_name=$1,
	    phone_number=$2,
		version=version+1
	WHERE id=$3 AND version=$4
	RETURNING 
		id,
	    full_name,
	    phone_number,
	    version
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		user.FullName,
		user.PhoneNumber,
		id,
		user.Version,
	)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.FullName,
		&userModel.PhoneNumber,
		&userModel.Version,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with id='%d' concurrently accessed: %w",
				id,
				core_errors.ErrConflict,
			)
		}

		return domain.User{}, fmt.Errorf(
			"scan error: %w",
			err,
		)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
	)
	return userDomain, nil
}
