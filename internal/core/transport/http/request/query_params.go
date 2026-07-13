package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/NikitaKissa/golang-todo-app/internal/core/errors"
)

func GetIntQueryParams(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' is not valid int: %v: %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return &val, nil
}
