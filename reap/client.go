package reap

import (
	"context"
	"time"
)

type Client interface {
	GetAccountBalance(
		context.Context,
	) (*GetAccountBalanceResponse, error)

	CreateCard(
		context.Context,
		CreateCardParams,
	) (*CreateCardResponse, error)

	GetCards(
		context.Context,
		GetCardsParams,
	) (*GetCardsResponse, error)

	GetCard(
		context.Context,
		GetCardParams,
	) (*GetCardResponse, error)

	AdjustCardBalance(
		context.Context,
		AdjustCardBalanceParams,
	) (*AdjustCardBalanceResponse, error)

	GetCardBalanceHistory(
		context.Context,
		GetCardBalanceHistoryParams,
	) (*GetCardBalanceHistoryResponse, error)

	GetCardTransactions(
		context.Context,
		GetCardTransactionsParams,
	) (*GetCardTransactionsResponse, error)

	GetAllTransactions(
		context.Context,
		GetAllTransactionsParams,
	) (*GetAllTransactionsResponse, error)
}

// TODO: Define separate request/response for http calls

type GetAccountBalanceResponse struct {
	AvailableBalance    float64 `json:"availableBalance"`
	AvailableToAllocate float64 `json:"availableToAllocate"`
}

type CreateCardParams struct {
	CardType          string  `json:"cardType"`
	SpendLimit        float64 `json:"spendLimit"`
	CustomerType      string  `json:"customerType"`
	KYC               KYC     `json:"kyc"`
	PreferredCardName string  `json:"preferredCardName"`
	Meta              Meta    `json:"meta"`
}
type CreateCardResponse struct {
	CardID string `json:"id"`
}

type GetCardParams struct {
	CardID string
}
type GetCardResponse struct {
	Card
}

type GetCardsParams struct {
	Status      string
	MetadataIDs []string
}
type GetCardsResponse struct {
	Items []Card     `json:"items"`
	Meta  Pagination `json:"meta"`
}

type AdjustCardBalanceParams struct{}
type AdjustCardBalanceResponse struct{}

type GetCardBalanceHistoryParams struct {
	CardID   string
	FromDate string
	Limit    int
}
type GetCardBalanceHistoryResponse struct {
	BalanceChanges []BalanceChange `json:"items"`
	Meta           Pagination      `json:"meta"`
}

type GetCardTransactionsParams struct {
	CardID   string
	FromDate string
	Limit    int
}
type GetCardTransactionsResponse struct {
	Transactions []Transaction `json:"items"`
	Meta         Pagination    `json:"meta"`
}

type GetAllTransactionsParams struct {
	FromDate string
	Limit    int
}
type GetAllTransactionsResponse struct {
	Transactions []Transaction `json:"items"`
	Meta         Pagination    `json:"meta"`
}

type Error struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

type Card struct {
	CardName           string       `json:"cardName"`
	SecondaryCardName  string       `json:"secondaryCardName"`
	Last4              string       `json:"last4"`
	AvailableCredit    string       `json:"availableCredit"`
	Status             string       `json:"status"`
	CardType           string       `json:"cardType"`
	PhysicalCardStatus string       `json:"physicalCardStatus"`
	ShippingAddress    *AddressInfo `json:"shippingAddress"`
	SpendControl       SpendControl `json:"spendControl"`
	CardDesign         string       `json:"cardDesign"`
	ShippingInfo       ShippingInfo `json:"shippingInfo"`
	Meta               Meta         `json:"meta"`
}

// TODO: Support business card type
type KYC struct {
	ConsumerInfo
}

type ConsumerInfo struct {
	FirstName          string      `json:"firstName"`
	LastName           string      `json:"lastName"`
	DOB                string      `json:"dob"`
	ResidentialAddress AddressInfo `json:"residentialAddress"`
	IDDocumentType     string      `json:"idDocumentType"`
	IDDocumentNumber   string      `json:"idDocumentNumber"`
}

// type BizInfo struct {
// 	FullName          string
// 	EntityType        string
// 	RegisteredAddress AddressInfo
// }

type AddressInfo struct {
	Line1   string `json:"line1"`
	Line2   string `json:"line2"`
	City    string `json:"city"`
	Country string `json:"country"`

	// PostalCode string `json:"postalCode"`
}

type Meta struct {
	ID             string          `json:"id"`
	Email          string          `json:"email"`
	OTPPhoneNumber PhoneNumberInfo `json:"otpPhoneNumber"`
}

type PhoneNumberInfo struct {
	DialCode    int    `json:"dialCode"`
	PhoneNumber string `json:"phoneNumber"`
}

type SpendControl struct {
	SpendControlAmount SpendControlAmount `json:"spendControlAmount"`
	SpendControlCap    SpendControlCap    `json:"spendControlCap"`
	ATMControl         ATMControl         `json:"atmControl"`
}

type SpendControlAmount struct {
	DailySpent   string `json:"dailySpent"`
	WeeklySpent  string `json:"weeklySpent"`
	MonthlySpent string `json:"monthlySpent"`
	YearlySpent  string `json:"yearlySpent"`
	AllTimeSpent string `json:"AllTimeSpent"`
}

type SpendControlCap struct {
	TransactionLimit string `json:"transactionLimit"`
	DailyLimit       string `json:"dailyLimit"`
	WeeklyLimit      string `json:"weeklyLimit"`
	MonthlyLimit     string `json:"monthlyLimit"`
	YearlyLimit      string `json:"yearlyLimit"`
	AllTimeLimit     string `json:"AllTimeLimit"`
}

type ATMControl struct {
	DailyFrequency    string `json:"dailyFrequency"`
	MonthlyFrequency  string `json:"monthlyFrequency"`
	DailyWithdrawal   string `json:"dailyWithdrawal"`
	MonthlyWithdrawal string `json:"monthlyWithdrawal"`
}

type ShippingInfo struct {
	BulkShippingID string `json:"bulkShippingID"`
	SKU            string `json:"sku"`
}

type Pagination struct {
	TotalItems   int `json:"totalItems"`
	ItemCount    int `json:"itemCount"`
	ItemsPerPage int `json:"itemsPerPage"`
	TotalPages   int `json:"totalPages"`
	CurrentPage  int `json:"currentPage"`
}

type Transaction struct {
	ID                  string    `json:"id"`
	CardID              string    `json:"card_id"`
	Merchant            Merchant  `json:"merchant_data"`
	Category            string    `json:"category"`
	Fees                Fees      `json:"fees"`
	BillAmount          string    `json:"bill_amount"`
	BillCurrency        string    `json:"bill_currency"`
	TransactionAmount   string    `json:"transaction_amount"`
	TransactionCurrency string    `json:"transaction_currency"`
	ConversionRate      string    `json:"conversion_rate"`
	Status              string    `json:"status"`
	Channel             string    `json:"channel"`
	CreatedAt           time.Time `json:"created_at"`
}

type Merchant struct {
	ID          string `json:"merchant_id"`
	Name        string `json:"merchant_name"`
	City        string `json:"merchant_city"`
	PostCode    string `json:"merchant_post_code"`
	State       string `json:"merchant_state"`
	Country     string `json:"merchant_country"`
	MCCCategory string `json:"mcc_category"`
	MCCCode     string `json:"mcc_code"`
}
type Fees struct {
	ATMFees string `json:"atm_fees"`
	FXFees  string `json:"fx_fees"`
}

type BalanceChange struct {
	ID          string    `json:"id"`
	Date        time.Time `json:"date"`
	Source      string    `json:"source"`
	Destination string    `json:"destination"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	Amount      string    `json:"amount"`
	Currency    string    `json:"currency"`
}
