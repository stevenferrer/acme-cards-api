package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/stevenferrer/acme-cards-api/acme"
)

type CardRepository struct {
	db *sql.DB
}

var _ acme.CardRepository = (*CardRepository)(nil)

func NewCardRepository(db *sql.DB) *CardRepository {
	return &CardRepository{db: db}
}

// SaveCardID implements acme.CardRepository.
func (r *CardRepository) SaveCardID(ctx context.Context, cardID string, externalCardID string) error {
	stmnt := `insert into cards (id, external_id) values ($1, $2)`
	_, err := r.db.ExecContext(ctx, stmnt, cardID, externalCardID)
	if err != nil {
		return fmt.Errorf("exec context: %w", err)
	}

	return nil
}

// GetExternalCardID implements acme.CardRepository.
func (r *CardRepository) GetExternalID(ctx context.Context, cardID string) (string, error) {
	stmnt := `select external_id from cards where id = $1`

	var externalID string
	err := r.db.QueryRowContext(ctx, stmnt, cardID).Scan(&externalID)
	if err != nil {
		return "", fmt.Errorf("query row context: %w", err)
	}

	return externalID, nil
}

// FindCardIDs implements acme.CardRepository.
func (r *CardRepository) FindCardIDs(ctx context.Context) ([]string, error) {
	// TODO: Filter by card status
	stmnt := `select id from cards order by created_at desc`

	rows, err := r.db.QueryContext(ctx, stmnt)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("query context: %w", err)
	}
	defer rows.Close()

	cardIDs := make([]string, 0)
	for rows.Next() {
		var cardID string
		err = rows.Scan(&cardID)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}
		cardIDs = append(cardIDs, cardID)
	}

	return cardIDs, nil
}

func (r *CardRepository) GetExternalIDMapping(ctx context.Context, externalIDs ...string) (map[string]string, error) {
	args := make([]any, 0, len(externalIDs))
	dollars := make([]string, 0, len(externalIDs))
	for i, externalID := range externalIDs {
		dollars = append(dollars, fmt.Sprintf("$%d", i+1))
		args = append(args, externalID)
		fmt.Printf("%s\n", externalID)
	}
	stmnt := fmt.Sprintf(`select "id", external_id from cards where external_id in (%s)`, strings.Join(dollars, ","))

	rows, err := r.db.QueryContext(ctx, stmnt, args...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("query context: %w", err)
	}
	defer rows.Close()

	cardIDMap := make(map[string]string, len(externalIDs))
	for rows.Next() {
		var cardID, externalID string
		err = rows.Scan(&cardID, &externalID)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}
		cardIDMap[externalID] = cardID
	}

	return cardIDMap, nil
}
