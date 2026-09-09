package main

import (
	"log"
	"net/http"
	"os"

	"dj/backend/internal/app"
	"dj/backend/internal/server"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	addr := envOrDefault("DJ_SERVER_ADDR", ":8080")
	handler := server.NewHandler(application)

	log.Printf("dev-journal server listening on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
