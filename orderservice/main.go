package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"orderService/config"
	"orderService/controller"
	"orderService/middleware"
	"orderService/pubsub"
	"orderService/repository"
	"orderService/service"
)

func main() {
	app := iris.New()
	db, error := config.ConnectToDB()
	if error != nil {
		fmt.Printf("Connection Lost")
		return
	}

	repo := &repository.OrderRepository{DB: db}
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
	protectedRoutes := app.Party("/order")

	protectedRoutes.Use(middleware.ValidateTokenMiddleware)
	{
		protectedRoutes.Post("/create", orderHandler.CreateOrder)
		protectedRoutes.Get("/getall", orderHandler.GetOrders)
		protectedRoutes.Get("/{id}", orderHandler.GetOrderByID)
		protectedRoutes.Patch("/{id}", orderHandler.UpdateOrderStatus)
	}
	app.Listen(":8080")
}
