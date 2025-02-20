package services

import (
	"context"

	"goal-tracker/api/pkg/todoist"
)

type TodoistService struct {
	client    todoist.Client
	projectID string
}

func (service *TodoistService) GetSections(
	ctx context.Context,
) ([]todoist.Section, error) {
	sections, err := service.client.GetAllSections(ctx, service.projectID)
	if err != nil {
		return nil, err
	}

	return sections, nil
}

func (service *TodoistService) GetTasks(ctx context.Context) ([]todoist.Task, error) {
	return service.client.GetActiveTasks(ctx, service.projectID)
}

func (service *TodoistService) GetTaskByID(
	ctx context.Context,
	id string,
) (*todoist.Task, error) {
	return service.client.GetActiveTask(ctx, id)
}
