package acmehttp

import (
	"fmt"
	"net/http"

	"github.com/stevenferrer/acme-cards-api/acme"
	"github.com/stevenferrer/acme-cards-api/x/xhttp"
)

func makeGetAccountBalHandler(cardSvc acme.CardService) http.Handler {
	return xhttp.WrapXHTTP(xhttp.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {

		bal, err := cardSvc.GetAccountBalance(r.Context())
		if err != nil {
			return fmt.Errorf("get account balance: %w", err)
		}

		err = renderResponse(http.StatusOK, w, accountBalance{
			Balance:   bal.AvailableBalance,
			Available: bal.AvailableToAllocate,
		})
		if err != nil {
			return fmt.Errorf("render response: %w", err)
		}

		return nil
	}))
}
