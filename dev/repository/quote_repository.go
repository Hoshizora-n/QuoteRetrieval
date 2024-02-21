package repository

import (
	"context"
	"integration_test/model/entity"
	"integration_test/util"

	"go.mongodb.org/mongo-driver/mongo"
)

type quoteRepo struct {
	db *mongo.Database
}

func NewQuoteRepo(db *mongo.Database) QuoteRepoInterface {
	return &quoteRepo{
		db: db,
	}
}

func (q *quoteRepo) Insert(ctx context.Context, quote entity.Quote) error {
	_, err := q.db.Collection(util.Configuration.MongoDB.Collection).InsertOne(ctx, quote)

	return err
}
