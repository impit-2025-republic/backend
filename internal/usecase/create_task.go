package usecase

import (
	"b8boost/backend/internal/entities"
	"context"
)

type (
	CreateTaskUseCase interface {
		Execute(ctx context.Context, input CreateTaskInput) (CreateTaskOutput, error)
	}

	CreateTaskInput struct {
		Title string `json:"title"`
	}

	CreateTaskOutput struct {
		ID uint `json:"id"`
	}

	createTaskInteractor struct {
		eventRepo entities.EventRepo
	}
)
