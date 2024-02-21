package repository

import (
	"context"
	"integration_test/model/entity"
)

type QuoteInterface interface {
	Insert(ctx context.Context, quote entity.Quote) error
}
