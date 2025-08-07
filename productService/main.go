package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/swaggo/http-swagger"
	"productService/config"
	"productService/controller"
	"productService/middleware"
	"productService/pubsub"
	//_ "productService/docs"
	"productService/repository"
	"productService/service"
)

func main() {
	app := iris.New()
	db, error := config.ConnectToDB()
	if error != nil {
		fmt.Printf("Connection Lost")
		return
	}

	repo := &repository.ProductRepository{DB: db}
	service := &service.ProductService{Repo: repo}
	productHandler := &controller.ProductController{Service: *service}

	pubsub.CheckForCreateOrder(repo)

	app.Get("/product/swagger/{any:path}", iris.FromStd(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	)))

	app.Get("/product", func(ctx iris.Context) {
		ctx.StatusCode(iris.StatusOK)
		ctx.WriteString("OK")
	})
	app.Get("/product/health", func(ctx iris.Context) {
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

	app.Listen(":8080")
}
