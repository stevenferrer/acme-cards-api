package acmehttp

type createCardRequest struct {
	FirstName  string         `json:"firstName"`
	LastName   string         `json:"lastName"`
	DOB        string         `json:"dob"`
	Address    addressInfo    `json:"address"`
	IDDocument idDocument     `json:"idDocument"`
	OTP        contactDetails `json:"otp"`
}

type addressInfo struct {
	Line1   string `json:"line1"`
	Line2   string `json:"line2"`
	City    string `json:"city"`
	Country string `json:"country"`
}

type idDocument struct {
	IDType   string `json:"idType"`
	IDNumber string `json:"idNumber"`
}
