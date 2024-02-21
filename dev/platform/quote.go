package platform

import (
	"context"
	"integration_test/model/platform"
)

type QuoteInterface interface {
	Get(ctx context.Context, category string) ([]platform.QuoteResponse, *platform.QuoteError, error)
}
