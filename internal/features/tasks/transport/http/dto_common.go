package tasks_transport_http

import (
	"time"

	"github.com/NikitaKissa/golang-todo-app/internal/core/domain"
)

type TaskDTOResponse struct {
	ID      int `json:"id"     example:"1"`
	Version int `json:"version" example:"1"`

	Title        string     `json:"title"          example:"Make homework"`
	Description  *string    `json:"description"    example:"Exercise 15 p. 12"`
	Completed    bool       `json:"completed"      example:"false"`
	CreatedAt    time.Time  `json:"created_at"     example:"2026-08-06T21:56:12.633263Z"`
	CompletedAt  *time.Time `json:"completed_at"   example:"2026-08-16T09:59:20.778079Z"`
	AuthorUserID int        `json:"author_user_id" example:"123"`
}

func taskDTOFromDomain(task domain.Task) TaskDTOResponse {
	return TaskDTOResponse{
		ID:           task.ID,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		CreatedAt:    task.CreatedAt,
		CompletedAt:  task.CompletedAt,
		AuthorUserID: task.AuthorUserID,
	}
}

func tasksDTOFromDomains(tasks []domain.Task) []TaskDTOResponse {
	tasksDTO := make([]TaskDTOResponse, len(tasks))

	for i, task := range tasks {
		tasksDTO[i] = taskDTOFromDomain(task)
	}

	return tasksDTO
}
