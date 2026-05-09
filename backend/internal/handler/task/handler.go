package task

import (
	"encoding/json"
	"fmt"
	"net/http"

	auth "github.com/koh5963/task-app/internal/auth"
	taskUsecase "github.com/koh5963/task-app/internal/usecase/task"
)

type Handler struct {
	usecase *taskUsecase.UseCase
}

func NewHandler(usecase *taskUsecase.UseCase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) ListByUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")

	// 実際にはURLパラメータからユーザID取得？
	token, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := token

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
