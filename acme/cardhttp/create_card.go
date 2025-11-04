package cardhttp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stevenferrer/acme-cards-api/acme"
	"github.com/stevenferrer/acme-cards-api/x/xhttp"
)

func makeCreateCardHandler(cardSvc acme.CardService) http.Handler {
	return xhttp.WrapXHTTP(xhttp.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		var c createCardRequest
		err := json.NewDecoder(r.Body).Decode(&c)
		if err != nil {
			return fmt.Errorf("decode request: %w", err)
		}

		createCardResp, err := cardSvc.CreateCard(r.Context(), toCreateCardParams(c))
		if err != nil {
			return fmt.Errorf("create card: %w", err)
		}

		err = renderResponse(http.StatusCreated, w, createCardResponse{
			CardID: createCardResp.CardID,
		})
		if err != nil {
			return fmt.Errorf("render response: %w", err)
		}

		return nil
	}))
}

func toCreateCardParams(p createCardRequest) acme.CreateCardParams {
	addr := p.Address
	otp := p.OTP
	idDoc := p.IDDocument
	return acme.CreateCardParams{
		FirstName: p.FirstName,
		LastName:  p.LastName,
		DOB:       p.DOB,

		Address: acme.Address{
			Line1:       addr.Line1,
			Line2:       addr.Line2,
			City:        addr.City,
			CountryCode: addr.Country,
		},
		ContactInfo: acme.ContactInfo{
			Email:       otp.Email,
			DialCode:    otp.DialCode,
			PhoneNumber: otp.PhoneNumber,
		},
		IDDocument: acme.IDDocument{
			Type:   idDoc.IDType,
			Number: idDoc.IDNumber,
		},
	}
}
