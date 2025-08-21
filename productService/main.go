package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
	"log"
	"productService/routes"

	"productService/config"
	"productService/controller"
	"productService/pubsub"
	"productService/repository"
	"productService/service"
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

	repo := &repository.ProductRepository{DB: db}
	service := &service.ProductService{Repo: repo}
	productHandler := &controller.ProductController{Service: *service}

	pubsub.CheckForCreateOrder(repo)
	
	app.Get("/product", func(ctx iris.Context) {
		ctx.StatusCode(iris.StatusOK)
		ctx.WriteString("OK")
	})
	app.Get("/product/health", func(ctx iris.Context) {
		ctx.StatusCode(iris.StatusOK)
		ctx.WriteString("OK")
	})

	routes.RegisterProductRoutes(app, productHandler)

	app.Listen(":8080")
}
