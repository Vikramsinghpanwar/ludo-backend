package player

import "context"

type Repository interface {
	GetProfile(ctx context.Context, playerID int64) (*string, error)
}
