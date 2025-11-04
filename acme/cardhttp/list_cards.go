package cardhttp

import (
	"fmt"
	"net/http"

	"github.com/stevenferrer/acme-cards-api/acme"
	"github.com/stevenferrer/acme-cards-api/x/xhttp"
)

func makeListCardsHandler(cardSvc acme.CardService) http.Handler {
	return xhttp.WrapXHTTP(xhttp.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		resp, err := cardSvc.ListCards(r.Context(), acme.ListCardsParams{})
		if err != nil {
			return fmt.Errorf("list cards: %w", err)
		}

		cards := make([]card, 0, len(resp.Cards))
		for _, c := range resp.Cards {
			cards = append(cards, toAcmeCardResponse(c))
		}

		err = renderResponse(http.StatusOK, w, listCardsResponse{
			Cards: cards,
		})
		if err != nil {
			return fmt.Errorf("render response: %w", err)
		}

		return nil
	}))
}
