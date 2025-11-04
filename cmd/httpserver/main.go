package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/stevenferrer/acme-cards-api/httpserver"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_DSN"))
	if err != nil {
		fatalError(logger, "sql open", err)
	}
	defer db.Close()

	srvr := httpserver.New(httpserver.Config{
		ReapAPIKey:    os.Getenv("REAP_API_KEY"),
		ReapSandoxURL: os.Getenv("REAP_SANDBOX_URL"),
		DB:            db,
		Logger:        logger,
	})

	// start server
	go func() {
		logger.Info("listening...", "addr", srvr.Addr)

		if err := srvr.ListenAndServe(); err != nil {
			fatalError(logger, "listen", err)
		}
	}()

	// setup signal capturing
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// wait for SIGINT (pkill -2)
	<-sigChan

	logger.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srvr.Shutdown(ctx); err != nil {
		fatalError(logger, "shutdown", err)
	}
}

func fatalError(logger *slog.Logger, msg string, err error) {
	logger.Error(msg, "err", err)
	os.Exit(1)
}
