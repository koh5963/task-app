package task

import "context"

type Repository interface {
	ListByUser(ctx context.Context, userID string) ([]Task, error)
	Create(ctx context.Context, t Task) error
	UpdateStatus(ctx context.Context, id string, userID string, status string) error
	Delete(ctx context.Context, id string, userID string) error
}
