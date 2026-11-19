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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"precedent/internal/db"
	"precedent/internal/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "precedent.db"
	}

	templateDir := os.Getenv("TEMPLATE_DIR")
	if templateDir == "" {
		templateDir = "web/templates"
	}

	log.Printf("[INIT] Opening SQLite database at %s...", dbPath)
	store, err := db.NewStore(dbPath)
	if err != nil {
		log.Fatalf("[FATAL] Database initialization failed: %v", err)
	}
	defer store.Close()

	log.Printf("[INIT] Compiling HTML templates from %s...", templateDir)
	gameHandler, err := handlers.NewGameHandler(store, templateDir)
	if err != nil {
		log.Fatalf("[FATAL] Handler initialization failed: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Health check endpoint for Cloud Run
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Game routes
	gameHandler.RegisterRoutes(r)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("[READY] PRECEDENT Customs Desk listening on http://localhost:%s", port)
		serverErrors <- server.ListenAndServe()
	}()

	// Graceful shutdown handling
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("[FATAL] Server error: %v", err)
		}
	case sig := <-shutdown:
		log.Printf("[SHUTDOWN] Signal %v received. Shutting down gracefully...", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("[WARN] Graceful shutdown error: %v", err)
			_ = server.Close()
		}
	}

	fmt.Println("[EXIT] Precedent server exited.")
}
