package core_http_request

import (
	"context"
	"net/http"

	core_logger "github.com/NikitaKissa/golang-todo-app/internal/core/logger"
	core_http_response "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/response"
)

func NewContext(
	rw http.ResponseWriter,
	r *http.Request,
) (context.Context, *core_http_response.HttpResponseHandler) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	rh := core_http_response.NewHttpResponseHandler(log, rw)

	return ctx, rh
}
