package routes

import (
	"github.com/kataras/iris/v12"
	"orderService/controller"
	"orderService/middleware"
)

func RegisterOrderRoutes(app *iris.Application, orderHandler *controller.OrderController) {
	
	protectedRoutes := app.Party("/order")

	protectedRoutes.Use(middleware.ValidateTokenMiddleware)
	{
		protectedRoutes.Post("/create", orderHandler.CreateOrder)
		protectedRoutes.Get("/getall", orderHandler.GetOrders)
		protectedRoutes.Get("/{id}", orderHandler.GetOrderByID)
		protectedRoutes.Patch("/{id}", orderHandler.UpdateOrderStatus)
	}

}
