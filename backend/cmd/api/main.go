package main

import (
	"log"
	"net/http"

	db "github.com/koh5963/task-app/internal/infra/db"
)

func main() {
	mux := http.NewServeMux()

	database, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	mux.HandleFunc("/health/db", func(w http.ResponseWriter, r *http.Request) {
		if err := db.HealthCheck(database); err != nil {
			http.Error(w, "DB NG", 500)
			return
		}
		defer database.Close()
		w.Write([]byte("DB OK"))
		log.Println("DB OK")
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World!:P"))
	})

	addr := ":8080"
	log.Printf("Starting server on %s\n", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
