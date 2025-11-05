package httpserver

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	sloghttp "github.com/samber/slog-http"

	"github.com/stevenferrer/acme-cards-api/acme"
	"github.com/stevenferrer/acme-cards-api/acme/acmehttp"
	"github.com/stevenferrer/acme-cards-api/acme/postgres"
	"github.com/stevenferrer/acme-cards-api/reap"
)

type Config struct {
	ReapAPIKey    string
	ReapSandoxURL string
	DB            *sql.DB
	Logger        *slog.Logger
}

func New(cfg Config) *http.Server {
	logger := cfg.Logger
	if logger == nil {
		logger = slog.Default()
	}

	var cardHTTPHandler, accountHTTPHandler http.Handler
	{
		cardRepo := postgres.NewCardRepository(cfg.DB)

		reapClient := reap.NewClient(reap.ClientConfig{
			APIKey:     cfg.ReapAPIKey,
			SandboxURL: cfg.ReapSandoxURL,
		})

		cardSvc := acme.NewReapCardService(reapClient, cardRepo)
		cardHTTPHandler = acmehttp.NewHTTPHandler(cardSvc)
		accountHTTPHandler = acmehttp.NewAccountHTTPHandler(cardSvc)
	}

	mux := chi.NewMux()
	mux.Use(
		sloghttp.Recovery,
		sloghttp.New(logger),
		cors.Default().Handler,
	)

	mux.Mount("/account", accountHTTPHandler)
	mux.Mount("/cards", cardHTTPHandler)

	return &http.Server{
		Addr:           ":9000",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
