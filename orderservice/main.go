package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
	"log"
	"orderService/common"
	"orderService/config"
	"orderService/controller"
	"orderService/pubsub"
	"orderService/repository"
	"orderService/routes"
	"orderService/service"
)

var logger *zap.Logger

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using environment variables")
	}
	var err error
	logger, err = zap.NewDevelopment()

	if err != nil {
		panic(err)
	}
	defer logger.Sync()
}

func main() {
	app := iris.New()
	dbConnector, error := config.ConnectToDB()
	if error != nil {
		fmt.Printf("Connection Lost")
		return
	}

	common.Init(dbConnector, logger)

	repo := &repository.OrderRepository{DB: dbConnector}
	service := &service.OrderService{Repo: repo}
	orderHandler := &controller.OrderController{Service: *service}

	pubsub.CheckForFailedOrderEvent(repo)
	app.Get("/order", func(ctx iris.Context) {
		ctx.StatusCode(iris.StatusOK)
		ctx.WriteString("OK")
	})
	app.Get("/order/health", func(ctx iris.Context) {
		ctx.StatusCode(iris.StatusOK)
		ctx.WriteString("OK")
	})

	routes.RegisterOrderRoutes(app, orderHandler)

	app.Listen(":8080")
}
