package acmehttp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/stevenferrer/acme-cards-api/acme"
)

func NewAccountHTTPHandler(cardSvc acme.CardService) http.Handler {
	mux := chi.NewMux()

	mux.Method(http.MethodGet, "/balance", makeGetAccountBalHandler(cardSvc))
	mux.Method(http.MethodGet, "/transactions", makeListTransactionsHandler(cardSvc))

	return mux
}

func NewHTTPHandler(cardSvc acme.CardService) http.Handler {
	mux := chi.NewMux()

	mux.Method(http.MethodPost, "/", makeCreateCardHandler(cardSvc))
	mux.Method(http.MethodGet, "/", makeListCardsHandler(cardSvc))
	mux.Method(http.MethodGet, "/{cardID}", makeGetCardHandler(cardSvc))
	mux.Method(http.MethodGet, "/{cardID}/transactions", makeListCardTransactionsHandler(cardSvc))
	mux.Method(http.MethodGet, "/{cardID}/balance-history", makeListBalanceHistoryHandler(cardSvc))

	return mux
}

func renderResponse(status int, w http.ResponseWriter, body any) error {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)

	if body == nil {
		return nil
	}

	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		return fmt.Errorf("json encode: %w", err)
	}

	return nil
}
