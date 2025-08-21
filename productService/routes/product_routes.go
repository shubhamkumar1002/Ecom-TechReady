package routes

import (
	"github.com/kataras/iris/v12"
	"productService/controller"
	"productService/middleware"
)

func RegisterProductRoutes(app *iris.Application, productHandler *controller.ProductController) {

	app.Get("/payment", func(ctx iris.Context) {
		ctx.StatusCode(iris.StatusOK)
		ctx.WriteString("OK")
	})

	protectedRoutes := app.Party("/product")

	protectedRoutes.Use(middleware.ValidateTokenMiddleware)
	{
		protectedRoutes.Get("/getall", productHandler.GetProducts)
		protectedRoutes.Get("/getproductbyid/:id", productHandler.GetProductByID)
		protectedRoutes.Post("/create", productHandler.CreateProduct)
		protectedRoutes.Put("/updateproduct/:id", productHandler.UpdateProduct)
		protectedRoutes.Post("/details", productHandler.GetProductDetails)
		protectedRoutes.Delete("/deleteproduct/:id", productHandler.DeleteProduct)
	}

}
