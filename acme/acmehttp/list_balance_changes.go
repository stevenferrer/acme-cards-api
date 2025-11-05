package acmehttp

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/stevenferrer/acme-cards-api/acme"
	"github.com/stevenferrer/acme-cards-api/x/xhttp"
)

func makeListBalanceHistoryHandler(cardSvc acme.CardService) http.Handler {
	return xhttp.WrapXHTTP(xhttp.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		cardID := chi.URLParam(r, "cardID")

		resp, err := cardSvc.ListCardBalanceHistory(r.Context(), cardID)
		if err != nil {
			return fmt.Errorf("list card transactions: %w", err)
		}

		bcs := make([]balanceChange, 0, len(resp.BalanceChanges))
		for _, bc := range resp.BalanceChanges {
			bcs = append(bcs, toAcmeBalanceChangeResponse(bc))
		}

		err = renderResponse(http.StatusOK, w, listBalanceChangesResponse{
			BalanceChanges: bcs,
		})
		if err != nil {
			return fmt.Errorf("render response: %w", err)
		}

		return nil
	}))
}

func toAcmeBalanceChangeResponse(bc acme.BalanceChange) balanceChange {
	return balanceChange{
		ID:       bc.ID,
		Date:     bc.Date.Format("2006-01-02 15:04:05"),
		Type:     bc.Type,
		Status:   bc.Status,
		Amount:   bc.Amount,
		Currency: bc.Currency,
	}
}
