package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/NikitaKissa/golang-todo-app/internal/core/domain"
	core_http_request "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/request"
)

type GetStatisticsResponse struct {
	TasksCreated               int      `json:"tasks_created"`
	TasksCompleted             int      `json:"tasks_completed"`
	TasksCompletedRate         *float64 `json:"tasks_completed_rate"`
	TasksAverageCompletionTime *string  `json:"tasks_average_completion_time"`
}

func toDTOFromDomain(statistics domain.Statistics) GetStatisticsResponse {
	var avgTime *string
	if statistics.TasksAverageCompletionTime != nil {
		duration := statistics.TasksAverageCompletionTime.String()
		avgTime = &duration
	}

	return GetStatisticsResponse{
		TasksCreated:               statistics.TasksCreated,
		TasksCompleted:             statistics.TasksCompleted,
		TasksCompletedRate:         statistics.TasksCompletedRate,
		TasksAverageCompletionTime: avgTime,
	}
}

func (h *StatisticsHttpHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx, responseHandler := core_http_request.NewContext(rw, r)

	userId, from, to, err := getUserIdFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'from'/'to'/'user_id' query parameter",
		)
		return
	}

	statisticsDomain, err := h.statisticsService.GetStatistics(ctx, userId, from, to)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get statistics",
		)
		return
	}

	response := toDTOFromDomain(statisticsDomain)
	responseHandler.JSONResponse(response, http.StatusOK)

}

func getUserIdFromToQueryParams(r *http.Request) (*int, *time.Time, *time.Time, error) {
	const (
		userIdKey    = "user_id"
		fromParamKey = "from"
		toParamKey   = "to"
	)

	userId, err := core_http_request.GetIntQueryParams(r, userIdKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	from, err := core_http_request.GetDateQueryParams(r, fromParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'from' query param: %w", err)
	}

	to, err := core_http_request.GetDateQueryParams(r, toParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'to' query param: %w", err)
	}

	return userId, from, to, nil
}
