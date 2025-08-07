package main

import (
	"authservice/config"
	"authservice/jwt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var logger *zap.Logger
var jwtManager *jwt.JWTManager
var dbConnector *gorm.DB
var err error

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
	dbConnector, err = config.ConnectToDB()
	if err != nil {
		logger.Fatal("Failed to connect to the database", zap.Error(err))
	}

	// * Create a new jwt manager`
	jwtManager = jwt.NewJWTManager("SECRET_KEY", 15*time.Minute, 7*24*time.Hour)

	app := iris.New()
	app.Get("/auth/health", func(ctx iris.Context) {
		ctx.StatusCode(iris.StatusOK)
		ctx.WriteString("OK")
	})

	app.Post("/auth/register", AddUser)
	app.Post("/auth/login", AuthenticateUser)
	app.Get("/auth/users", GetAllUsers)
	app.Post("/auth/validate", ValidateToken)
	app.Post("/auth/refresh", RefreshTokenHandler)

	app.Get("/auth", func(ctx iris.Context) {
		ctx.StatusCode(iris.StatusOK)
		ctx.WriteString("OK")
	})

	app.Listen(":8080")
}
