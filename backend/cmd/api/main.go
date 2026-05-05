package main

import (
	"log"
	"net/http"

	taskHandler "github.com/koh5963/task-app/internal/handler/task"
	db "github.com/koh5963/task-app/internal/infra/db"
	taskRepo "github.com/koh5963/task-app/internal/infra/repository/task"
	taskUc "github.com/koh5963/task-app/internal/usecase/task"
)

func main() {
	mux := http.NewServeMux()

	database, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	repo := taskRepo.New(database)
	usecase := taskUc.NewUseCase(repo)
	handler := taskHandler.NewHandler(usecase)

	mux.HandleFunc("/health/db", func(w http.ResponseWriter, r *http.Request) {
		if err := db.HealthCheck(database); err != nil {
			http.Error(w, "DB NG", 500)
			return
		}
		w.Write([]byte("DB OK"))
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World!:P"))
	})

	mux.HandleFunc("/tasks", handler.ListByUser)

	addr := ":8080"
	log.Printf("Starting server on %s\n", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
