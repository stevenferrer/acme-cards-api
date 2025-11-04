package cardhttp

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/stevenferrer/acme-cards-api/acme"
	"github.com/stevenferrer/acme-cards-api/x/xhttp"
)

func makeListCardTransactionsHandler(cardSvc acme.CardService) http.Handler {
	return xhttp.WrapXHTTP(xhttp.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		cardID := chi.URLParam(r, "cardID")

		resp, err := cardSvc.ListCardTransactions(r.Context(), cardID)
		if err != nil {
			return fmt.Errorf("list card transactions: %w", err)
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

func toAcmeTransactionResponse(t acme.Transaction) transaction {
	return transaction{
		ID:       t.ID,
		CardID:   t.CardID,
		Category: t.Category,
		Status:   t.Status,
		Channel:  t.Channel,
		Amount:   t.Amount,
		Currency: t.Currency,

		Merchant: merchantDetails{
			Name:    t.Merchant.Name,
			City:    t.Merchant.City,
			Country: t.Merchant.Country,
		},

		Fees: feeDetails{
			ATMFees: t.Fees.ATMFees,
			FXFees:  t.Fees.FXFees,
		},

		Date: t.CreatedAt.Format("2006-01-02"),
	}
}
