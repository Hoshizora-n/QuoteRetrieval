package main

import (
	"context"
	"fmt"
	"integration_test/http/api"
	"integration_test/platform"
	"integration_test/provider"
	"integration_test/repository"
	"integration_test/service"
	"integration_test/util"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	if err := util.LoadConfig("."); err != nil {
		log.Fatal(err)
	}
}

func main() {
	mongoClient, err := provider.NewMongoDBClient()
	if err != nil {
		log.Fatal(err)
	}

	db := mongoClient.Database(util.Configuration.MongoDB.Database)
	quoteRepo := repository.NewQuoteRepo(db)
	quoteClient := platform.NewQuoteClient()
	quoteService := service.NewQuoteService(quoteClient, quoteRepo)

	app := api.NewApp(quoteService)
	addr := fmt.Sprintf(":%v", util.Configuration.Server.Port)
	server, err := app.CreateServer(addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting server...")

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}

	}()

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)

	sig := <-shutdownCh
	log.Printf("Receiving signal: %s", sig)

	ctx, cancel := context.WithTimeout(context.Background(), util.Configuration.Server.ShutdownTimeout)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	func(db *mongo.Client) {
		db.Disconnect(ctx)
		log.Println("Successfully disconnected from DB..")
		log.Println("Successfully shutdown server..")
	}(mongoClient)
}
