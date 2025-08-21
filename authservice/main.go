package main

import (
	"authservice/common"
	"authservice/config"
	"authservice/jwt"
	"authservice/routes"
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
	//ctx := context.Background()
	dbConnector, err = config.ConnectToDB()
	if err != nil {
		logger.Fatal("Failed to connect to the database", zap.Error(err))
	}

	secretKey := config.GetSecretKey()
	jwtManager = jwt.NewJWTManager(secretKey, 15*time.Minute, 7*24*time.Hour)

	common.Init(dbConnector, logger, jwtManager)

	app := iris.New()
	routes.RegisterAuthRoutes(app)
	app.Listen(":8080")
}
