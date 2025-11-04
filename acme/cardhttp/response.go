package cardhttp

type card struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	Last4           string         `json:"last4"`
	AvailableCredit string         `json:"availableCredit"`
	ContactInfo     contactDetails `json:"contactInfo"`
}

type contactDetails struct {
	Email       string `json:"email"`
	DialCode    int    `json:"dialCode"`
	PhoneNumber string `json:"phoneNumber"`
}

type listCardsResponse struct {
	Cards []card `json:"cards"`
}

type createCardResponse struct {
	CardID string `json:"cardId"`
}

type accountBalance struct {
	Balance   string `json:"balance"`
	Available string `json:"available"`
}

type transaction struct {
	ID       string `json:"id"`
	CardID   string `json:"cardId"`
	Category string `json:"category"`
	Status   string `json:"status"`
	Channel  string `json:"channel"`
	Amount   string `json:"amount"`
	Currency string `json:"currency"`

	Fees     feeDetails      `json:"fees"`
	Merchant merchantDetails `json:"merchant"`

	Date string `json:"date"`
}

type feeDetails struct {
	ATMFees string `json:"atmFees"`
	FXFees  string `json:"fxFees"`
}

type merchantDetails struct {
	Name    string `json:"name"`
	City    string `json:"city"`
	Country string `json:"country"`
}

type listTransactionsResponse struct {
	Transactions []transaction `json:"transactions"`
}

type balanceChange struct {
	ID       string `json:"id"`
	Date     string `json:"date"`
	Type     string `json:"type"`
	Status   string `json:"status"`
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type listBalanceChangesResponse struct {
	BalanceChanges []balanceChange `json:"balanceChanges"`
}
