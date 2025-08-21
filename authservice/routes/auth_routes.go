package routes

import (
	"authservice/controller"
	"github.com/kataras/iris/v12"
)

func RegisterAuthRoutes(app *iris.Application) {
	app.Get("/auth/health", func(ctx iris.Context) {
		ctx.StatusCode(iris.StatusOK)
		ctx.WriteString("OK")
	})

	app.Post("/auth/register", controller.AddUser)
	app.Post("/auth/login", controller.AuthenticateUser)
	app.Get("/auth/users", controller.GetAllUsers)
	app.Post("/auth/validate", controller.ValidateToken)
	app.Post("/auth/refresh", controller.RefreshTokenHandler)

	app.Get("/auth", func(ctx iris.Context) {
		ctx.StatusCode(iris.StatusOK)
		ctx.WriteString("OK")
	})

}
