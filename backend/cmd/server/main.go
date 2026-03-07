package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dhruvsolanki0811/webgen/internal/config"
	"github.com/dhruvsolanki0811/webgen/internal/handler"
	"github.com/dhruvsolanki0811/webgen/internal/middleware"
	"github.com/dhruvsolanki0811/webgen/internal/repository"
	"github.com/dhruvsolanki0811/webgen/internal/router"
	"github.com/dhruvsolanki0811/webgen/internal/service"
)

const shutdownTimeout = 10 * time.Second

func main() {
	if err := run(); err != nil {
		log.Fatalf("fatal: %v", err)
	}
}

func run() error {
	// 1. Load config
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// 2. Connect database
	db, err := repository.NewMongoDB(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		return err
	}
	defer db.Close(context.Background())

	// 3. Build repositories
	userRepo := repository.NewUserRepo(db)
	_ = repository.NewProjectRepo(db) // Used in Phase 2

	// 4. Build services
	tokenSvc := service.NewTokenService(cfg.JWTSecret)
	authSvc := service.NewAuthService(userRepo, tokenSvc)

	// 5. Build middleware
	authMW := middleware.NewAuthMiddleware(tokenSvc)

	// 6. Build handlers
	authHandler := handler.NewAuthHandler(authSvc)

	// 7. Build router
	r := router.New(router.Config{
		FrontendURL: cfg.FrontendURL,
	}, authMW, authHandler)

	// 8. Start server with graceful shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("server starting on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	<-done
	log.Println("server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	return srv.Shutdown(ctx)
}
