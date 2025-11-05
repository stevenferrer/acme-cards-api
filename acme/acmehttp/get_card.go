package acmehttp

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/stevenferrer/acme-cards-api/acme"
	"github.com/stevenferrer/acme-cards-api/x/xhttp"
)

func makeGetCardHandler(cardSvc acme.CardService) http.Handler {
	return xhttp.WrapXHTTP(xhttp.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		cardID := chi.URLParam(r, "cardID")

		card, err := cardSvc.GetCard(r.Context(), cardID)
		if err != nil {
			return fmt.Errorf("get card: %w", err)
		}

		err = renderResponse(http.StatusOK, w, toAcmeCardResponse(*card))
		if err != nil {
			return fmt.Errorf("render response: %w", err)
		}

		return nil
	}))
}

func toAcmeCardResponse(c acme.Card) card {
	return card{
		ID:              c.ID,
		Name:            c.Name,
		Last4:           c.Last4,
		AvailableCredit: c.AvailableCredit,
		ContactInfo:     toAcmeContactDetailsResponse(c.ContactInfo),
	}
}

func toAcmeContactDetailsResponse(c acme.ContactInfo) contactDetails {
	return contactDetails{
		Email:       c.Email,
		DialCode:    c.DialCode,
		PhoneNumber: c.PhoneNumber,
	}
}
