package api

import (
	"integration_test/model/delivery"
	"integration_test/service"
	"integration_test/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type App struct {
	quoteSvc service.QuoteInterface
}

func NewApp(quoteSvc service.QuoteInterface) *App {
	return &App{
		quoteSvc: quoteSvc,
	}
}

func (a *App) CreateServer(address string) (*http.Server, error) {
	gin.SetMode(util.Configuration.Server.Mode)

	r := gin.Default()
	r.Use(gin.Recovery())

	r.GET("/ping", a.checkConnectivity)
	r.GET(util.Configuration.Server.Endpoint, a.getQuote)

	server := &http.Server{
		Addr:    address,
		Handler: r,
	}

	return server, nil
}

func (a *App) checkConnectivity(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (a *App) getQuote(ctx *gin.Context) {
	var req delivery.Request

	if ctx.ShouldBindQuery(&req) != nil {
		ctx.JSON(500, delivery.Response{
			Error: "error binding request",
		})
		return
	}

	res, err := a.quoteSvc.GetQuote(ctx, req)
	if err != nil {
		ctx.JSON(500, res)
		return
	}

	ctx.JSON(200, res)
}
