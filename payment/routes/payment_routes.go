package routes

import (
	"github.com/kataras/iris/v12"
	"paymentService/controller"
	"paymentService/middleware"
)

func RegisterPaymentRoutes(app *iris.Application, paymentHandler *controller.PaymentController) {

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
}
