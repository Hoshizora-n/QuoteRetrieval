package service

import (
	"context"
	"errors"
	"integration_test/mocks"
	"integration_test/model/delivery"
	"integration_test/model/entity"
	"integration_test/model/platform"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

type test struct {
	name         string
	req          delivery.Request
	platformMock func() (platformMock *mocks.QuoteClientInterface)
	repoMock     func() (repoMock *mocks.QuoteRepoInterface)
	want         want
}

type want struct {
	res delivery.Response
	err error
}

func TestQuoteService(t *testing.T) {
	tests := []test{
		{
			name: "error in platform",
			req: delivery.Request{
				Category: "freedom",
			},
			platformMock: func() (platformMock *mocks.QuoteClientInterface) {
				platformMock = mocks.NewQuoteClientInterface(t)
				platformMock.On("Get", mock.Anything, "?category=freedom").Return(nil, nil, http.ErrHandlerTimeout)
				return
			},
			repoMock: func() (repoMock *mocks.QuoteRepoInterface) {
				repoMock = mocks.NewQuoteRepoInterface(t)
				return
			},
			want: want{
				res: delivery.Response{
					Error: http.ErrHandlerTimeout.Error(),
				},
				err: http.ErrHandlerTimeout,
			},
		},
		{
			name: "platform has error in response",
			req: delivery.Request{
				Category: "freedom",
			},
			platformMock: func() (platformMock *mocks.QuoteClientInterface) {
				platformMock = mocks.NewQuoteClientInterface(t)
				platformMock.On("Get", mock.Anything, "?category=freedom").Return(nil, &platform.QuoteError{
					Error: "missing api key",
				}, nil)
				return
			},
			repoMock: func() (repoMock *mocks.QuoteRepoInterface) {
				repoMock = mocks.NewQuoteRepoInterface(t)
				return
			},
			want: want{
				res: delivery.Response{
					Error: "missing api key",
				},
				err: errors.New("get response error from platform"),
			},
		},
		{
			name: "error in insert",
			req: delivery.Request{
				Category: "freedom",
			},
			platformMock: func() (platformMock *mocks.QuoteClientInterface) {
				platformMock = mocks.NewQuoteClientInterface(t)
				platformMock.On("Get", mock.Anything, "?category=freedom").Return([]platform.QuoteResponse{{
					Quote:    "if you don't fight, you won't win",
					Author:   "eren yeager",
					Category: "freedom",
				}}, nil, nil)
				return
			},
			repoMock: func() (repoMock *mocks.QuoteRepoInterface) {
				repoMock = mocks.NewQuoteRepoInterface(t)
				repoMock.On("Insert", mock.Anything, entity.Quote{
					Quote:    "if you don't fight, you won't win",
					Author:   "eren yeager",
					Category: "freedom",
				}).Return(mongo.ErrClientDisconnected)
				return
			},
			want: want{
				res: delivery.Response{
					Error: mongo.ErrClientDisconnected.Error(),
				},
				err: mongo.ErrClientDisconnected,
			},
		},
		{
			name: "success",
			req: delivery.Request{
				Category: "freedom",
			},
			platformMock: func() (platformMock *mocks.QuoteClientInterface) {
				platformMock = mocks.NewQuoteClientInterface(t)
				platformMock.On("Get", mock.Anything, "?category=freedom").Return([]platform.QuoteResponse{{
					Quote:    "if you don't fight, you won't win",
					Author:   "eren yeager",
					Category: "freedom",
				}}, nil, nil)
				return
			},
			repoMock: func() (repoMock *mocks.QuoteRepoInterface) {
				repoMock = mocks.NewQuoteRepoInterface(t)
				repoMock.On("Insert", mock.Anything, entity.Quote{
					Quote:    "if you don't fight, you won't win",
					Author:   "eren yeager",
					Category: "freedom",
				}).Return(nil)
				return
			},
			want: want{
				res: delivery.Response{
					Quote:    "if you don't fight, you won't win",
					Author:   "eren yeager",
					Category: "freedom",
				},
				err: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			platformMock := test.platformMock()
			repoMock := test.repoMock()
			service := NewQuoteService(platformMock, repoMock)

			res, err := service.GetQuote(context.Background(), test.req)
			assert.Equal(t, test.want.res, res)
			assert.Equal(t, test.want.err, err)

			platformMock.AssertExpectations(t)
			repoMock.AssertExpectations(t)
		})
	}

}
