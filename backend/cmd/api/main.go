package main

import (
	"log"
	"net/http"
	"os"

	"github.com/koh5963/task-app/internal/auth"
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

	authenticator, err := auth.NewAuthenticator(
		os.Getenv("SUPABASE_JWKS_URL"),
		os.Getenv("SUPABASE_ISSUER"),
		os.Getenv("SUPABASE_AUDIENCE"),
	)
	if err != nil {
		log.Fatal(err)
	}

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

	mux.Handle("GET /tasks", auth.Middleware(authenticator, http.HandlerFunc(handler.ListByUser)))
	mux.Handle("POST /tasks", auth.Middleware(authenticator, http.HandlerFunc(handler.Create)))
	mux.Handle("PATCH /tasks/{id}/{status}", auth.Middleware(authenticator, http.HandlerFunc(handler.Update)))
	mux.Handle("DELETE /tasks/{id}", auth.Middleware(authenticator, http.HandlerFunc(handler.Delete)))

	addr := ":8080"
	log.Printf("Starting server on %s\n", addr)

	if err := http.ListenAndServe(addr, cors(mux)); err != nil {
		log.Fatal(err)
	}
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := os.Getenv("FRONTEND_ORIGIN")
		// Vite の開発サーバーを許可
		w.Header().Set("Access-Control-Allow-Origin", origin)

		// 許可するHTTPメソッド
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")

		// Authorization を許可しないと Bearer token 付き fetch で落ちる
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// プリフライトリクエストはここで返す
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
