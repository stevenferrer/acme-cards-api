package reap_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stevenferrer/acme-cards-api/reap"
)

type M = map[any]any

func TestClientV1(t *testing.T) {
	httpmock.Activate(t)
	defer httpmock.DeactivateAndReset()

	sandboxURL := "https://sandbox.example.xyz"
	apiKey := "asdfqwerty"

	client := reap.NewClient(reap.ClientConfig{
		SandboxURL: sandboxURL,
		APIKey:     apiKey,
		// use default client so httpmock can intercept
		HTTPClient: http.DefaultClient,
	})

	createCardParams := reap.CreateCardParams{
		CardType:          "Virtual",
		CustomerType:      "Consumer",
		PreferredCardName: "Hua Liang",
		SpendLimit:        5000,
		KYC: reap.KYC{
			reap.ConsumerInfo{
				FirstName: "Hua",
				LastName:  "Liang",
				DOB:       "1990-08-08",
				ResidentialAddress: reap.AddressInfo{
					Line1:   "Tung Ning Bldg",
					Line2:   "Western District",
					City:    "Hong Kong",
					Country: "HKG",
				},
				IDDocumentType:   "Passport",
				IDDocumentNumber: "1000000",
			},
		},
		Meta: reap.Meta{
			ID:    "1000000",
			Email: "hualiang@myspace.xyz",
			OTPPhoneNumber: reap.PhoneNumberInfo{
				DialCode:    852,
				PhoneNumber: "25441194",
			},
		},
	}

	t.Run("Create card ok", func(t *testing.T) {
		httpmock.RegisterResponder(
			http.MethodPost,
			fmt.Sprintf("%s/cards", sandboxURL),
			newResponderWithStatus(
				t, http.StatusCreated,
				"create_card_request.json",
				"create_card_response.json",
			),
		)

		resp, err := client.CreateCard(context.TODO(), createCardParams)
		require.NoError(t, err)

		assert.Equal(t, &reap.CreateCardResponse{CardID: "1234"}, resp)
	})

	t.Run("Create card error", func(t *testing.T) {
		httpmock.RegisterResponder(
			http.MethodPost,
			fmt.Sprintf("%s/cards", sandboxURL),
			newResponderWithStatus(
				t, http.StatusBadRequest,
				"create_card_request.json",
				"create_card_response_error.json",
			),
		)

		_, err := client.CreateCard(context.TODO(), createCardParams)
		assert.ErrorContains(t, err, "unexpected status 400")
	})

	chinZengCard := reap.Card{
		CardName:           "Chin Zeng",
		Last4:              "2112",
		AvailableCredit:    "5000.00",
		Status:             "ACTIVE",
		CardType:           "virtual",
		PhysicalCardStatus: "NOT_PHYSICAL_CARD",
		SpendControl: reap.SpendControl{
			SpendControlAmount: reap.SpendControlAmount{
				DailySpent:   "0.00",
				WeeklySpent:  "0.00",
				MonthlySpent: "0.00",
				YearlySpent:  "0.00",
				AllTimeSpent: "0.00",
			},
			SpendControlCap: reap.SpendControlCap{
				TransactionLimit: "0.00",
				DailyLimit:       "0.00",
				WeeklyLimit:      "0.00",
				MonthlyLimit:     "0.00",
				YearlyLimit:      "0.00",
				AllTimeLimit:     "0.00",
			},
			ATMControl: reap.ATMControl{
				DailyFrequency:    "0.00",
				MonthlyFrequency:  "0.00",
				DailyWithdrawal:   "0.00",
				MonthlyWithdrawal: "0.00",
			},
		},
		CardDesign: "a8fb9fb3-2a0e-468c-b98e-f47955db20c7",
		Meta: reap.Meta{
			ID:    "3000000",
			Email: "chinzeng@myspace.xyz",
			OTPPhoneNumber: reap.PhoneNumberInfo{
				DialCode: 852, PhoneNumber: "31070474",
			},
		},
	}

	t.Run("Get card ok", func(t *testing.T) {
		cardID := "1234"
		httpmock.RegisterResponder(
			http.MethodGet,
			fmt.Sprintf("%s/cards/%s", sandboxURL, cardID),
			newResponderWithStatus(
				t, http.StatusOK, "",
				"get_card_response.json",
			),
		)

		resp, err := client.GetCard(context.TODO(), reap.GetCardParams{
			CardID: cardID,
		})
		require.NoError(t, err)

		expect := &reap.GetCardResponse{
			Card: chinZengCard,
		}
		assert.Equal(t, expect, resp)
	})

	t.Run("Get cards ok", func(t *testing.T) {
		q := url.Values{}
		q.Add("metadataId", "1000000")
		q.Add("status", "ACTIVE")

		httpmock.RegisterResponder(
			http.MethodGet,
			fmt.Sprintf("%s/cards?%s", sandboxURL, q.Encode()),
			newResponderWithStatus(
				t, http.StatusOK, "",
				"get_cards_response.json",
			),
		)

		resp, err := client.GetCards(context.TODO(), reap.GetCardsParams{
			Status:      "ACTIVE",
			MetadataIDs: []string{"1000000"},
		})
		require.NoError(t, err)

		expect := &reap.GetCardsResponse{
			Items: []reap.Card{chinZengCard},
			Meta: reap.Pagination{
				TotalItems:   1,
				ItemCount:    1,
				ItemsPerPage: 10,
				TotalPages:   1,
				CurrentPage:  1,
			},
		}
		assert.Equal(t, expect, resp)
	})
}

func newResponderWithStatus(
	t *testing.T,
	status int,
	requestFilename string,
	responseFilename string,
) httpmock.Responder {
	return func(r *http.Request) (*http.Response, error) {
		// assert request body if not empty
		if requestFilename != "" {
			expectBody, err := decodeExpectRequestBody(requestFilename)
			if err != nil {
				return nil, fmt.Errorf("decode expect request body: %w", err)
			}

			var actualBody any
			err = json.NewDecoder(r.Body).Decode(&actualBody)
			if err != nil {
				return nil, fmt.Errorf("json decode: %w", err)
			}

			if !assert.Equal(t, expectBody, actualBody, "unexpected request body") {
				return nil, fmt.Errorf("unexpected request body")
			}
		}

		return httpmock.NewJsonResponse(status, httpmock.File(path.Join("fixtures", responseFilename)))
	}
}

func decodeExpectRequestBody(requestFilename string) (any, error) {
	f, err := os.Open(path.Join("fixtures", requestFilename))
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	var expectBody any
	err = json.NewDecoder(f).Decode(&expectBody)
	if err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	}

	return expectBody, nil
}
