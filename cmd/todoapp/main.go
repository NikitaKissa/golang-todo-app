package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/NikitaKissa/golang-todo-app/internal/core/logger"
	core_postgres_pool "github.com/NikitaKissa/golang-todo-app/internal/core/repository/postgres/pool"
	core_http_middleware "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/middleware"
	core_http_server "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/server"
	users_postgres_repository "github.com/NikitaKissa/golang-todo-app/internal/features/users/repository/postgres"
	users_service "github.com/NikitaKissa/golang-todo-app/internal/features/users/service"
	users_transport_http "github.com/NikitaKissa/golang-todo-app/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("Failed to init app logger: ", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("initializing postgres connection pool")
	pool, err := core_postgres_pool.NewConnectionPool(ctx, core_postgres_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHttp := users_transport_http.NewUsersHttpHandler(usersService)

	logger.Debug("initializing http server")
	httpServer := core_http_server.NewHttpServer(
		core_http_server.NewConfigMust(),
		logger,

		// middleware declarations
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)

	apiVersionRouter := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHttp.Routes()...)
	httpServer.RegisterApiRoutes(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("Http server run error:", zap.Error(err))
	}
}
