package service

import (
	"context"
	"integration_test/model/delivery"
)

type QuoteInterface interface {
	GetQuote(ctx context.Context, request delivery.Request) (delivery.Response, error)
}
