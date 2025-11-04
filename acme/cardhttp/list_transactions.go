package cardhttp

import (
	"fmt"
	"net/http"

	"github.com/stevenferrer/acme-cards-api/acme"
	"github.com/stevenferrer/acme-cards-api/x/xhttp"
)

func makeListTransactionsHandler(cardSvc acme.CardService) http.Handler {
	return xhttp.WrapXHTTP(xhttp.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		resp, err := cardSvc.ListTransactions(r.Context())
		if err != nil {
			return fmt.Errorf("list transactions: %w", err)
		}

		transactions := make([]transaction, 0, len(resp.Transactions))
		for _, t := range resp.Transactions {
			transactions = append(transactions, toAcmeTransactionResponse(t))
		}

		err = renderResponse(http.StatusOK, w, listTransactionsResponse{
			Transactions: transactions,
		})
		if err != nil {
			return fmt.Errorf("render response: %w", err)
		}

		return nil
	}))
}
