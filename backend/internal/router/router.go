package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/dhruvsolanki0811/webgen/internal/handler"
	"github.com/dhruvsolanki0811/webgen/internal/middleware"
)

type Config struct {
	FrontendURL string
}

func New(cfg Config, authMW *middleware.AuthMiddleware, authHandler *handler.AuthHandler) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimw.RequestID)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.FrontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check — no auth
	r.Get("/health", healthHandler)

	// Auth routes — no auth
	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/signup", authHandler.Signup)
		r.Post("/login", authHandler.Login)
	})

	// Protected routes — auth required
	r.Route("/api", func(r chi.Router) {
		r.Use(authMW.Required)
		// Phase 2 adds: /projects, /projects/{id}, /projects/{id}/chat, /projects/{id}/generate
		// Phase 3 adds: /projects/{id}/deploy, /projects/{id}/stream
	})

	return r
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
