package acme

import "context"

type CardRepository interface {
	SaveCardID(ctx context.Context, cardID string, externalCardID string) error
	FindCardIDs(context.Context) ([]string, error)
	GetExternalID(ctx context.Context, cardID string) (externalCardID string, err error)
	GetExternalIDMapping(ctx context.Context, externalIDs ...string) (map[string]string, error)
}
