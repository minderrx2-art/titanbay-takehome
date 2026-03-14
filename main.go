package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"titanbay/internal/handler"
	"titanbay/internal/service"
	"titanbay/internal/store/postgres"
)

func main() {
	migrateDb := flag.Bool("migrate", false, "Initialize and migrate the database")
	flag.Parse()

	dbURL := os.Getenv("DATABASE_URL")
	// Exit Status: 0 = success, 1 = failure
	if dbURL == "" {
		slog.Error("DATABASE_URL environment variable is not set")
		os.Exit(1)
	}

	database, err := postgres.NewPostgresStore(dbURL)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	if *migrateDb {
		slog.Info("Running migrations...")
		if err := postgres.Migrate(database); err != nil {
			slog.Error("Migration failed", "error", err)
		}
		os.Exit(0)
	}

	store := postgres.NewStoreModel(database)
	app := service.NewService(store)
	h := handler.NewHandler(app)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler.LoggingMiddleware(mux),
	}

	srvErrors := make(chan error, 1)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	slog.Info("Starting api server on :8080")
	go func() {
		srvErrors <- srv.ListenAndServe()
	}()

	select {
	case err := <-srvErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server listener failure", "error", err)
			os.Exit(1)
		}
	case sig := <-signalChan:
		slog.Info("shutdown signal received", "signal", sig.String())
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		store.Close()
		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("server shutdown failed", "error", err)
		}
	}
}
