package reap

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// ClientV1 is an implementation of Client interface
type ClientV1 struct {
	sandboxURL string
	apiKey     string
	httpClient *http.Client
}

var _ Client = (*ClientV1)(nil)

type ClientConfig struct {
	SandboxURL string
	APIKey     string
	HTTPClient *http.Client
}

func NewClient(cfg ClientConfig) *ClientV1 {
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = newHTTPClient()
	}

	// TODO: Validate apikey is non-empty?
	return &ClientV1{
		sandboxURL: cfg.SandboxURL,
		apiKey:     cfg.APIKey,
		httpClient: httpClient,
	}
}

func (c *ClientV1) GetAccountBalance(ctx context.Context) (*GetAccountBalanceResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "account/balance", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	var body struct {
		GetAccountBalanceResponse
		Error
	}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		e := body.Error
		return nil, fmt.Errorf("unexpected status %d with error code %q and message %q", resp.StatusCode, e.Code, e.Message)
	}

	return &body.GetAccountBalanceResponse, nil
}

// CreateCard implements Client.
func (c *ClientV1) CreateCard(ctx context.Context, params CreateCardParams) (*CreateCardResponse, error) {
	// TODO: validate required params

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(params)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "cards", nil, buf)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	var body struct {
		CreateCardResponse
		Error
	}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		e := body.Error
		return nil, fmt.Errorf("unexpected status %d with error code %q and message %q", resp.StatusCode, e.Code, e.Message)
	}

	return &body.CreateCardResponse, nil
}

// GetCard implements Client.
func (c *ClientV1) GetCard(ctx context.Context, params GetCardParams) (*GetCardResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("cards/%s", params.CardID), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	var body struct {
		GetCardResponse
		Error
	}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		e := body.Error
		return nil, fmt.Errorf("unexpected status %d with error code %q and message %q", resp.StatusCode, e.Code, e.Message)
	}

	return &body.GetCardResponse, nil
}

// GetCards implements Client.
func (c *ClientV1) GetCards(ctx context.Context, params GetCardsParams) (*GetCardsResponse, error) {
	cardStatus := params.Status
	if cardStatus == "" {
		cardStatus = "ACTIVE"
	}

	q := url.Values{}
	q.Add("status", cardStatus)
	if len(params.MetadataIDs) > 0 {
		q.Add("metadataId", strings.Join(params.MetadataIDs, ","))
	}

	req, err := c.newRequest(ctx, http.MethodGet, "cards", q, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	var body struct {
		GetCardsResponse
		Error
	}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		e := body.Error
		return nil, fmt.Errorf("unexpected status %d with error code %q and message %q", resp.StatusCode, e.Code, e.Message)
	}

	return &body.GetCardsResponse, nil
}

// GetCardTransactions implements Client.
func (c *ClientV1) GetCardTransactions(ctx context.Context, params GetCardTransactionsParams) (*GetCardTransactionsResponse, error) {
	requestBody := struct {
		FromDate string `json:"fromDate"`
		Limit    int    `json:"limit"`
	}{
		FromDate: params.FromDate,
		Limit:    params.Limit,
	}

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(requestBody)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	path := fmt.Sprintf("cards/%s/transactions", params.CardID)
	req, err := c.newRequest(ctx, http.MethodGet, path, nil, buf)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	var body struct {
		GetCardTransactionsResponse
		Error
	}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		e := body.Error
		return nil, fmt.Errorf("unexpected status %d with error code %q and message %q", resp.StatusCode, e.Code, e.Message)
	}

	return &body.GetCardTransactionsResponse, nil
}

func (c *ClientV1) GetAllTransactions(ctx context.Context, params GetAllTransactionsParams) (*GetAllTransactionsResponse, error) {
	// TODO: Exclude transactions from deleted cards??

	requestBody := struct {
		FromDate string `json:"fromDate"`
		Limit    int    `json:"limit"`
	}{
		FromDate: params.FromDate,
		Limit:    params.Limit,
	}

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(requestBody)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodGet, "transactions", nil, buf)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	var body struct {
		GetAllTransactionsResponse
		Error
	}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if err := mustHaveStatusOrError(http.StatusOK, resp.StatusCode, body.Error); err != nil {
		return nil, err
	}

	return &body.GetAllTransactionsResponse, nil
}

func mustHaveStatusOrError(expectStatus int, gotStatus int, e Error) error {
	if expectStatus != gotStatus {
		return fmt.Errorf("expecting status %d, but got %d: error code %q, message %q", expectStatus, gotStatus, e.Code, e.Message)
	}

	return nil
}

// GetCardBalanceHistory implements Client.
func (c *ClientV1) GetCardBalanceHistory(ctx context.Context, params GetCardBalanceHistoryParams) (*GetCardBalanceHistoryResponse, error) {
	requestBody := struct {
		FromDate string `json:"fromDate"`
		Limit    string `json:"limit"`
	}{
		FromDate: params.FromDate,
		Limit:    strconv.Itoa(params.Limit),
	}

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(requestBody)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	path := fmt.Sprintf("cards/%s/balance-history", params.CardID)
	req, err := c.newRequest(ctx, http.MethodGet, path, nil, buf)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	var body struct {
		GetCardBalanceHistoryResponse
		Error
	}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		e := body.Error
		return nil, fmt.Errorf("unexpected status %d with error code %q and message %q", resp.StatusCode, e.Code, e.Message)
	}

	return &body.GetCardBalanceHistoryResponse, nil
}

// AdjustCardBalance implements Client.
func (c *ClientV1) AdjustCardBalance(context.Context, AdjustCardBalanceParams) (*AdjustCardBalanceResponse, error) {
	panic("unimplemented")
}

func (c *ClientV1) newRequest(
	ctx context.Context,
	method string,
	path string,
	q url.Values,
	body io.Reader,
) (*http.Request, error) {
	urlStr, err := url.JoinPath(c.sandboxURL, path)
	if err != nil {
		return nil, fmt.Errorf("url join: %w", err)
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	u.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("new request with context: %w", err)
	}

	req.Header.Add("accept-version", "v1.0")
	req.Header.Add("x-reap-api-key", c.apiKey)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	return req, nil
}
