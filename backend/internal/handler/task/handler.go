package task

import (
	"encoding/json"
	"fmt"
	"net/http"

	auth "github.com/koh5963/task-app/internal/auth"
	taskDomain "github.com/koh5963/task-app/internal/domain/task"
	taskUsecase "github.com/koh5963/task-app/internal/usecase/task"
)

type Handler struct {
	usecase *taskUsecase.UseCase
}

func NewHandler(usecase *taskUsecase.UseCase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) ListByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tasks, err := h.usecase.ListByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to list tasks: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var task taskDomain.Task
	parseErr := json.NewDecoder(r.Body).Decode(&task)
	if parseErr != nil {
		http.Error(w, "BadRequest", http.StatusBadRequest)
		return
	}

	newTasks, err := h.usecase.Create(r.Context(), userID, task)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to list tasks: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newTasks); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	id := r.PathValue("id")

	delErr := h.usecase.Delete(r.Context(), id, userID)
	if delErr != nil {
		http.Error(w, fmt.Errorf("failed to delete task: %w", delErr).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	id := r.PathValue("id")
	status := r.PathValue("status")

	updErr := h.usecase.UpdateStatus(r.Context(), id, userID, status)
	if updErr != nil {
		http.Error(w, fmt.Errorf("failed to delete task: %w", updErr).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
