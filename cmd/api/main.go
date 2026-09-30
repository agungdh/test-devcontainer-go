package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"todoapp/internal/config"
	"todoapp/internal/database"
	"todoapp/internal/todo"
	"todoapp/pkg/response"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	db, err := database.New(cfg.DSN())
	if err != nil {
		slog.Error("gagal konek database", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	repo := todo.NewRepository(db)
	svc := todo.NewService(repo)
	handler := todo.NewHandler(svc)

	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.RealIP, middleware.Logger, middleware.Recoverer)
	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	}).Handler)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	r.Route("/api/v1", func(r chi.Router) {
		todo.RegisterRoutes(r, handler)
	})

	srv := &http.Server{Addr: cfg.Addr(), Handler: r}

	go func() {
		slog.Info("server jalan", "addr", cfg.Addr())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	slog.Info("server berhenti")
}
