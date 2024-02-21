package service

import (
	"context"
	"integration_test/model/delivery"
)

type QuoteServiceInterface interface {
	GetQuote(ctx context.Context, request delivery.Request) (delivery.Response, error)
}
