package dto

import "github.com/EvansTrein/Workmate/internal/entity"

type TaskCreateRequest struct {
	Description string `json:"description" validate:"required"`
}

type TaskCreateResponce struct {
	Task *entity.Task `json:"task"`
}

type AllTasksResponce struct {
	Tasks []*entity.Task `json:"tasks"`
}
