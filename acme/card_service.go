package acme

import (
	"context"
	"time"
)

type CardService interface {
	GetAccountBalance(context.Context) (*AccountBalance, error)

	// CreateCard creates a virtual card for consumer use
	CreateCard(context.Context, CreateCardParams) (*CreateCardResponse, error)

	GetCard(ctx context.Context, cardID string) (*Card, error)

	ListCards(
		ctx context.Context,
		params ListCardsParams,
	) (*ListCardsResponse, error)

	ListCardTransactions(
		ctx context.Context,
		cardID string,
	) (*ListCardTransactionsResponse, error)

	ListCardBalanceHistory(
		ctx context.Context,
		cardID string,
	) (*ListCardBalanceHistoryResponse, error)

	ListTransactions(
		ctx context.Context,
	) (*ListTransactionsResponse, error)
}

type AccountBalance struct {
	AvailableBalance    string
	AvailableToAllocate string
}

type CreateCardParams struct {
	FirstName   string
	LastName    string
	DOB         string
	Address     Address
	ContactInfo ContactInfo
	IDDocument  IDDocument
}

type CreateCardResponse struct {
	CardID string
}

type ListCardsParams struct{}
type ListCardsResponse struct {
	Cards []Card
}

type ListCardTransactionsResponse struct {
	Transactions []Transaction
}

type ListTransactionsResponse struct {
	Transactions []Transaction
}

type ListCardBalanceHistoryParams struct{}
type ListCardBalanceHistoryResponse struct {
	BalanceChanges []BalanceChange
}

type Address struct {
	Line1       string
	Line2       string
	City        string
	CountryCode string
}

type ContactInfo struct {
	Email       string
	DialCode    int
	PhoneNumber string
}

type IDDocument struct {
	Type   string
	Number string
}

type Card struct {
	ID              string
	Name            string
	Last4           string
	AvailableCredit string
	ContactInfo     ContactInfo
}

type Transaction struct {
	ID       string
	CardID   string
	Category string
	Status   string
	Channel  string

	Amount   string
	Currency string

	Fees     FeeDetails
	Merchant MerchantDetails

	CreatedAt time.Time
}

type MerchantDetails struct {
	ID      string
	Name    string
	City    string
	Country string
}

type FeeDetails struct {
	ATMFees string
	FXFees  string
}

type BalanceChange struct {
	ID       string
	Date     time.Time
	Type     string
	Status   string
	Amount   string
	Currency string
}
