package server

import (
	"log/slog"
	"net/http"

	"github.com/PumpkinSeed/frost-playground/internal/core"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Run() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Post("/private-key", core.CreatePrivateKeyHandler)
	r.Post("/split-private-key", core.SplitPrivateKeyHandler)
	r.Post("/commit-key-share", core.CommitKeyShare)
	r.Post("/sign-key-share", core.SignKeyShareHandler)

	slog.Info("starting server", slog.String("address", ":3000"))
	if err := http.ListenAndServe(":3000", r); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
