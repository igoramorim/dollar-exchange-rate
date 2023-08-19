package repository

import "context"

type Repository interface {
	Save(ctx context.Context, exchangeRate float64) error
}
