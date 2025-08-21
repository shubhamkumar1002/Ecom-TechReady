package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
	"log"
	"paymentService/common"
	"paymentService/config"
	"paymentService/controller"
	"paymentService/pubsub"
	"paymentService/repository"
	"paymentService/routes"
	"paymentService/service"
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
	db, error := config.ConnectToDB()
	if error != nil {
		fmt.Printf("Connection Lost")
		return
	}

	common.Init(db, logger)
	repo := &repository.PaymentRepository{DB: db}
	service := &service.PaymentService{Repo: repo}
	paymentHandler := &controller.PaymentController{Service: *service}

	pubsub.CheckForPublishedPayments(repo)
	pubsub.CheckForSuccessOrderEvent(repo)

	app.Get("/payment/health", func(ctx iris.Context) {
		ctx.StatusCode(iris.StatusOK)
		ctx.WriteString("OK")
	})
	routes.RegisterPaymentRoutes(app, paymentHandler)
	app.Listen(":8080")
}
