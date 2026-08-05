// This the control panel
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"automatomic/internal/auth"
	"automatomic/internal/config"
	"automatomic/internal/http/handler"
	"automatomic/internal/database"
	"automatomic/internal/http/middleware"
	"automatomic/internal/model"
	"automatomic/internal/repository"

)

func main() {
	// Load our config
	cfg := config.Load()

	// Run Database Migrations Automatically
	if err := database.RunMigrations(cfg.DatabaseURL, "internal/database/migrations"); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	//  Connect to PostgreSQL
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbPool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer dbPool.Close()
	

	// Initialize jwt Manager, and the repo pool
	userRepo := repository.NewPostgresRepo(dbPool)
	jwtMgr := auth.NewJWTManager(cfg.JWTSecret, 24*time.Hour)

	redirectURL := fmt.Sprintf("http://localhost:%s/api/v1/auth/github/callback", cfg.Port)
	ghOAuth := auth.NewGitHubOAuth(cfg.GitHubClientID, cfg.GitHubSecret, redirectURL)

	authHandler := handler.NewAuthHandler(ghOAuth, jwtMgr, userRepo)

	// Set Up HTTP Router & Middleware
	mux := http.NewServeMux()

	// Public Routes
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
	})


	// OAuth Endpoints
	mux.HandleFunc("GET /api/v1/auth/github/login", authHandler.HandleGitHubLogin)
	mux.HandleFunc("GET /api/v1/auth/github/callback", authHandler.HandleGitHubCallback)

	// Protected Pipeline Routes (Requires valid JWT & pipeline:read scope)
	protectedPipelineMux := http.NewServeMux()
	protectedPipelineMux.HandleFunc("GET /api/v1/pipelines", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"message":"Successfully accessed protected pipeline resource"}`))
	})

	// Authentication & Scope Guards
	authMiddleware := middleware.JWTMiddleware(jwtMgr)
	scopeMiddleware := middleware.RequireScope(model.ScopePipelineRead)

	mux.Handle("/api/v1/pipelines", authMiddleware(scopeMiddleware(protectedPipelineMux)))

	// Configure HTTP Server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      corsMiddleware(mux, cfg.FrontendURL),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Shutdown Handler
	serverCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("Control plane API server running on port %s...", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v\n", err)
		}
	}()

	<-serverCtx.Done()
	log.Println("Shutting down server gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced shutdown: %v", err)
	}

	log.Println("Server exited cleanly")
}

// Basic CORS Middleware for Local Frontend Development
func corsMiddleware(next http.Handler, frontendURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", frontendURL)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
