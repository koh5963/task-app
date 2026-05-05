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
