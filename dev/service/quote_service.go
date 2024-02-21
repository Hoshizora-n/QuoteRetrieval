package service

import (
	"context"
	"errors"
	"fmt"
	"integration_test/model/delivery"
	"integration_test/model/entity"
	"integration_test/platform"
	"integration_test/repository"
)

type quoteService struct {
	quoteClient platform.QuoteInterface
	quoteRepo   repository.QuoteInterface
}

func NewQuoteService(quoteClient platform.QuoteInterface, quoteRepo repository.QuoteInterface) QuoteInterface {
	return &quoteService{
		quoteClient: quoteClient,
		quoteRepo:   quoteRepo,
	}
}

func (q *quoteService) GetQuote(ctx context.Context, req delivery.Request) (res delivery.Response, err error) {
	category := fmt.Sprintf("?category=%s", req.Category)

	quoteData, quoteErr, err := q.quoteClient.Get(ctx, category)
	if err != nil {
		return
	}

	if err != nil {
		res.Error = err.Error()
		return
	}

	if quoteErr != nil {
		res.Error = quoteErr.Error
		err = errors.New("get response error from platform")
		return
	}

	err = q.quoteRepo.Insert(ctx, entity.Quote{
		Author:   quoteData[0].Author,
		Quote:    quoteData[0].Quote,
		Category: quoteData[0].Category,
	})
	if err != nil {
		res.Error = err.Error()
		return
	}

	res.Author = quoteData[0].Author
	res.Quote = quoteData[0].Quote
	res.Category = quoteData[0].Category

	return res, nil
}
