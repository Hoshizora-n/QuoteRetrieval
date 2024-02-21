package testing

import (
	"context"
	"fmt"
	"integration_test/http/api"
	"integration_test/platform"
	"integration_test/repository"
	"integration_test/service"
	"integration_test/util"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testSvc struct {
	suite.Suite
	ctx            context.Context
	mongoClient    *mongo.Client
	app            *api.App
	platformServer *httptest.Server
}

func (suite *testSvc) SetupSuite() {
	suite.ctx = context.Background()

	suite.loadConfig()

	suite.initMongoTest()
	suite.initQuoteServer()
	suite.initApp()
}

func (suite *testSvc) TearDownSuite() {
	t := suite.T()

	t.Log("hi im teardown")
}

func (suite *testSvc) loadConfig() {
	t := suite.T()

	viper.AddConfigPath("./")
	viper.SetConfigName("config.test")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	assert.NoError(t, err)

	conf := util.Config{}
	err = viper.Unmarshal(&conf)
	assert.NoError(t, err)
	util.Configuration = conf
}

func (suite *testSvc) initMongoTest() {
	t := suite.T()

	mongoC, err := mongodb.RunContainer(
		suite.ctx,
		testcontainers.WithImage("mongo:4.4.23"),
	)
	assert.NoError(t, err)

	uri, err := mongoC.ConnectionString(suite.ctx)
	assert.NoError(t, err)

	mongoClient, err := mongo.Connect(suite.ctx, options.Client().ApplyURI(uri))
	assert.NoError(t, err)
	suite.mongoClient = mongoClient
}

func (suite *testSvc) initQuoteServer() {
	suite.platformServer = PlatformServer()
	util.Configuration.Platform.Host = suite.platformServer.URL
}

func (suite *testSvc) initApp() {
	db := suite.mongoClient.Database(util.Configuration.MongoDB.Database)
	quoteRepo := repository.NewQuoteRepo(db)
	quoteClient := platform.NewQuoteClient()
	quoteService := service.NewQuoteService(quoteClient, quoteRepo)

	suite.app = api.NewApp(quoteService)
	addr := fmt.Sprintf(":%v", util.Configuration.Server.Port)
	server, err := suite.app.CreateServer(addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting server...")

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	time.Sleep(2 * time.Second)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(testSvc))
}
