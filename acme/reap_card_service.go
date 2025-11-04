package acme

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/stevenferrer/acme-cards-api/reap"
)

// ReapCardService implements CardService
type ReapCardService struct {
	reapClient reap.Client
	cardRepo   CardRepository
}

var _ CardService = (*ReapCardService)(nil)

func NewReapCardService(
	reapClient reap.Client,
	cardRepo CardRepository,
) *ReapCardService {
	return &ReapCardService{
		cardRepo:   cardRepo,
		reapClient: reapClient,
	}
}

func (s *ReapCardService) GetAccountBalance(ctx context.Context) (*AccountBalance, error) {
	resp, err := s.reapClient.GetAccountBalance(ctx)
	if err != nil {
		return nil, fmt.Errorf("get reap account balance: %w", err)
	}

	return &AccountBalance{
		AvailableBalance:    fmt.Sprintf("%.2f", resp.AvailableBalance),
		AvailableToAllocate: fmt.Sprintf("%.2f", resp.AvailableToAllocate),
	}, nil
}

func (s *ReapCardService) CreateCard(ctx context.Context, params CreateCardParams) (*CreateCardResponse, error) {
	// TODO: Validate required params

	// generate internal ID
	cardID := uuid.New().String()
	cardID = strings.ReplaceAll(cardID, "-", "")

	// create card on Reap side
	resp, err := s.reapClient.CreateCard(ctx, toReapCreateCardParams(cardID, params))
	if err != nil {
		return nil, fmt.Errorf("create reap card: %w", err)
	}

	// save id mapping to database
	err = s.cardRepo.SaveCardID(ctx, cardID, resp.CardID)
	if err != nil {
		return nil, fmt.Errorf("save card ID %q external(%q): %w", cardID, resp.CardID, err)
	}

	return &CreateCardResponse{CardID: cardID}, nil
}

func toReapCreateCardParams(cardID string, params CreateCardParams) reap.CreateCardParams {
	return reap.CreateCardParams{
		CardType:          "Virtual",
		CustomerType:      "Consumer",
		PreferredCardName: fmt.Sprintf("%s %s", params.FirstName, params.LastName),
		// The current default balance is 0, admin can update balance later
		SpendLimit: 0,
		KYC: reap.KYC{
			ConsumerInfo: reap.ConsumerInfo{
				FirstName:          params.FirstName,
				LastName:           params.LastName,
				DOB:                params.DOB,
				ResidentialAddress: toReapAddress(params.Address),
				IDDocumentType:     params.IDDocument.Type,
				IDDocumentNumber:   params.IDDocument.Number,
			},
		},
		Meta: reap.Meta{
			ID:    cardID,
			Email: params.ContactInfo.Email,
			OTPPhoneNumber: reap.PhoneNumberInfo{
				DialCode:    params.ContactInfo.DialCode,
				PhoneNumber: params.ContactInfo.PhoneNumber,
			},
		},
	}
}

func toReapAddress(addr Address) reap.AddressInfo {
	return reap.AddressInfo{
		Line1:   addr.Line1,
		Line2:   addr.Line2,
		City:    addr.City,
		Country: addr.CountryCode,
	}
}

func (s *ReapCardService) GetCard(ctx context.Context, cardID string) (*Card, error) {
	// get card id mapping
	reapCardID, err := s.cardRepo.GetExternalID(ctx, cardID)
	if err != nil {
		return nil, fmt.Errorf("get reap card ID: %w", err)
	}

	// fetch card from reap
	resp, err := s.reapClient.GetCard(ctx, reap.GetCardParams{
		CardID: reapCardID,
	})
	if err != nil {
		return nil, fmt.Errorf("get reap card: %w", err)
	}

	card := toAcmeCard(resp.Card)
	return &card, nil
}

func toAcmeCard(params reap.Card) Card {
	return Card{
		ID:              params.Meta.ID,
		Last4:           params.Last4,
		Name:            params.CardName,
		AvailableCredit: params.AvailableCredit,
		ContactInfo: ContactInfo{
			Email:       params.Meta.Email,
			DialCode:    params.Meta.OTPPhoneNumber.DialCode,
			PhoneNumber: params.Meta.OTPPhoneNumber.PhoneNumber,
		},
	}
}

// GetAllCards implements CardService.
func (s *ReapCardService) ListCards(ctx context.Context, params ListCardsParams) (*ListCardsResponse, error) {
	cardIDs, err := s.cardRepo.FindCardIDs(ctx)
	if err != nil {
		return nil, fmt.Errorf("get card ids: %w", err)
	}

	resp, err := s.reapClient.GetCards(ctx, reap.GetCardsParams{
		MetadataIDs: cardIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("get reap cards: %w", err)
	}

	cards := make([]Card, 0, len(resp.Items))
	for _, reapCard := range resp.Items {
		cards = append(cards, toAcmeCard(reapCard))
	}

	return &ListCardsResponse{
		Cards: cards,
	}, nil
}

func (s *ReapCardService) ListCardTransactions(ctx context.Context, cardID string) (*ListCardTransactionsResponse, error) {
	reapCardID, err := s.cardRepo.GetExternalID(ctx, cardID)
	if err != nil {
		return nil, fmt.Errorf("get reap card id: %w", err)
	}

	fromDate := time.Now().Format("2006-01-02")
	resp, err := s.reapClient.GetCardTransactions(ctx, reap.GetCardTransactionsParams{
		CardID:   reapCardID,
		FromDate: fromDate,
		Limit:    10,
	})
	if err != nil {
		return nil, fmt.Errorf("get reap card transactions: %w", err)
	}

	transactions := make([]Transaction, 0, len(resp.Transactions))
	for _, t := range resp.Transactions {
		transactions = append(transactions, toAcmeTransaction(cardID, t))
	}

	return &ListCardTransactionsResponse{
		Transactions: transactions,
	}, nil
}

func (s *ReapCardService) ListTransactions(ctx context.Context) (*ListTransactionsResponse, error) {

	fromDate := time.Now().Format("2006-01-02")
	resp, err := s.reapClient.GetAllTransactions(ctx, reap.GetAllTransactionsParams{
		FromDate: fromDate,
		Limit:    10,
	})
	if err != nil {
		return nil, fmt.Errorf("get reap card transactions: %w", err)
	}

	// collect unique reap card ids
	externalIDsSet := make(map[string]struct{})
	externalIDs := make([]string, 0, len(externalIDsSet))
	for _, t := range resp.Transactions {
		if _, ok := externalIDsSet[t.CardID]; ok {
			continue
		}
		externalIDsSet[t.CardID] = struct{}{}
		externalIDs = append(externalIDs, t.CardID)
	}

	cardIDMapping, err := s.cardRepo.GetExternalIDMapping(ctx, externalIDs...)
	if err != nil {
		return nil, fmt.Errorf("get external id mapping: %w", err)
	}

	transactions := make([]Transaction, 0, len(resp.Transactions))
	for _, t := range resp.Transactions {
		// skip transactions with no mapping from database
		cardID, ok := cardIDMapping[t.CardID]
		if !ok {
			continue
		}
		transactions = append(transactions, toAcmeTransaction(cardID, t))
	}

	return &ListTransactionsResponse{
		Transactions: transactions,
	}, nil
}

func toAcmeTransaction(cardID string, reapTx reap.Transaction) Transaction {
	return Transaction{
		ID:        reapTx.ID,
		CardID:    cardID,
		Category:  reapTx.Category,
		Status:    reapTx.Status,
		Channel:   reapTx.Channel,
		Amount:    reapTx.BillAmount,
		Currency:  reapTx.BillCurrency,
		CreatedAt: reapTx.CreatedAt,
		Fees: FeeDetails{
			ATMFees: reapTx.Fees.ATMFees,
			FXFees:  reapTx.Fees.FXFees,
		},
		Merchant: MerchantDetails{
			ID:      reapTx.Merchant.ID,
			Name:    reapTx.Merchant.Name,
			City:    reapTx.Merchant.City,
			Country: reapTx.Merchant.Country,
		},
	}
}

// ListCardBalanceHistory implements CardService.
func (s *ReapCardService) ListCardBalanceHistory(ctx context.Context, cardID string) (*ListCardBalanceHistoryResponse, error) {
	reapCardID, err := s.cardRepo.GetExternalID(ctx, cardID)
	if err != nil {
		return nil, fmt.Errorf("get reap card id: %w", err)
	}

	fromDate := time.Now().Format("2006-01-02")
	resp, err := s.reapClient.GetCardBalanceHistory(ctx, reap.GetCardBalanceHistoryParams{
		CardID:   reapCardID,
		FromDate: fromDate,
		Limit:    10,
	})
	if err != nil {
		return nil, fmt.Errorf("get reap card balance history: %w", err)
	}

	balanceChanges := make([]BalanceChange, 0, len(resp.BalanceChanges))
	for _, b := range resp.BalanceChanges {
		balanceChanges = append(balanceChanges, toAcmeBalanceChange(b))
	}

	return &ListCardBalanceHistoryResponse{
		BalanceChanges: balanceChanges,
	}, nil
}

func toAcmeBalanceChange(bc reap.BalanceChange) BalanceChange {
	return BalanceChange{
		ID:       bc.ID,
		Date:     bc.Date,
		Type:     bc.Type,
		Status:   bc.Status,
		Amount:   bc.Amount,
		Currency: bc.Currency,
	}
}
