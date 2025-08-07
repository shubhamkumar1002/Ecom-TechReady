package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"paymentService/config"
	"paymentService/controller"
	"paymentService/middleware"
	"paymentService/pubsub"
	"paymentService/repository"
	"paymentService/service"
)

// @title Payment Service API
// @version 1.0
// @description This is a Simple application for checking payments status
// @BasePath /
func main() {
	app := iris.New()
	db, error := config.ConnectToDB()
	if error != nil {
		fmt.Printf("Connection Lost")
		return
	}

	repo := &repository.PaymentRepository{DB: db}
	service := &service.PaymentService{Repo: repo}
	paymentHandler := &controller.PaymentController{Service: *service}

	pubsub.CheckForPublishedPayments(repo)
	pubsub.CheckForSuccessOrderEvent(repo)
	
	app.Get("/payment/health", func(ctx iris.Context) {
		ctx.StatusCode(iris.StatusOK)
		ctx.WriteString("OK")
	})
	app.Get("/payment", func(ctx iris.Context) {
		ctx.StatusCode(iris.StatusOK)
		ctx.WriteString("OK")
	})

	protectedRoutes := app.Party("/payment")

	protectedRoutes.Use(middleware.ValidateTokenMiddleware)
	{
		protectedRoutes.Get("/payments", paymentHandler.GetPayments)
		protectedRoutes.Get("/paymentbyorderid/{id}", paymentHandler.GetPaymentByOrderID)
	}
	app.Listen(":8080")
}
