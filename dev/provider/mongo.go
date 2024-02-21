package provider

import (
	"context"
	"fmt"
	"integration_test/util"
	"net/url"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// NewMongoDBClient creates a new MongoDB client
func NewMongoDBClient() (*mongo.Client, error) {
	cfg := util.Configuration.MongoDB

	url := (&url.URL{
		Scheme:   "mongodb",
		User:     url.UserPassword(cfg.Username, cfg.Password),
		Host:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Path:     cfg.Database,
		RawQuery: strings.Join(cfg.Options, "&"),
	}).String()

	// Set the MongoDB client options
	clientOptions := options.Client().ApplyURI(url)

	// Connect to the MongoDB server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.PrimaryPreferred())
	if err != nil {
		return nil, err
	}

	return client, nil
}
