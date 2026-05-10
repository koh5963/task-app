package task

import (
	"context"
	sql "database/sql"
	"fmt"
	"time"

	tasks "github.com/koh5963/task-app/internal/domain/task"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListByUser(ctx context.Context, userID string) ([]tasks.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, `
		SELECT
		  id
		  , user_id
		  , title
		  , description
		  , status
		  , created_at
		  , updated_at
		FROM
		  tasks
		WHERE
		  user_id = $1
	`, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	taskList := []tasks.Task{}

	for rows.Next() {
		var t tasks.Task

		if err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		taskList = append(taskList, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate task rows: %w", err)
	}

	return taskList, nil
}

func (r *Repository) Create(ctx context.Context, userID string, t tasks.Task) ([]tasks.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	tx, errT := r.db.BeginTx(ctx, nil)
	if errT != nil {
		return nil, fmt.Errorf("failed to open transaction: %w", errT)
	}
	defer tx.Rollback()

	_, err := tx.ExecContext(ctx, `
		INSERT INTO
		tasks (
		  user_id
		  , title
		  , description
		  , status
		  , created_at
		  , updated_at
		) VALUES (
		  $1
			, $2
			, $3
			, $4
			, NOW()
			, NOW()
		)
	`, userID, t.Title, t.Description, t.Status)

	if err != nil {
		return nil, fmt.Errorf("failed to insert tasks table: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return r.ListByUser(ctx, userID)
}

func (r *Repository) UpdateStatus(ctx context.Context, id string, userID string, status string) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	tx, errT := r.db.BeginTx(ctx, nil)
	if errT != nil {
		return fmt.Errorf("failed to open transaction: %w", errT)
	}
	defer tx.Rollback()

	_, err := tx.ExecContext(ctx, `
		UPDATE tasks
		SET
		  status = $1
			, updated_at = NOW()
		WHERE
			id = $2
		  AND user_id = $3
	`, status, id, userID)

	if err != nil {
		return fmt.Errorf("failed to update tasks table: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id string, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	tx, errT := r.db.BeginTx(ctx, nil)
	if errT != nil {
		return fmt.Errorf("failed to open transaction: %w", errT)
	}
	defer tx.Rollback()

	_, err := tx.ExecContext(ctx, `
		DELETE
		FROM
		  tasks
		WHERE
			id = $1
		  AND user_id = $2
	`, id, userID)

	if err != nil {
		return fmt.Errorf("failed to delete tasks table: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
