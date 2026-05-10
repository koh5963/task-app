package task

import (
	"context"

	tasks "github.com/koh5963/task-app/internal/domain/task"
)

type UseCase struct {
	repo tasks.Repository
}

func NewUseCase(repo tasks.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) ListByUser(ctx context.Context, userID string) ([]tasks.Task, error) {
	return uc.repo.ListByUser(ctx, userID)
}

func (uc *UseCase) Create(ctx context.Context, userID string, t tasks.Task) ([]tasks.Task, error) {
	return uc.repo.Create(ctx, userID, t)
}

func (uc *UseCase) UpdateStatus(ctx context.Context, id string, userID string, status string) error {
	return uc.repo.UpdateStatus(ctx, id, userID, status)
}

func (uc *UseCase) Delete(ctx context.Context, id string, userID string) error {
	return uc.repo.Delete(ctx, id, userID)
}
