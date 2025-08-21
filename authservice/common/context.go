package common

import (
	"authservice/jwt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ApplicationContext struct {
	DbConnector *gorm.DB
	Logger      *zap.Logger
	JWTManager  *jwt.JWTManager
}

var App *ApplicationContext

func Init(db *gorm.DB, logger *zap.Logger, jwtManager *jwt.JWTManager) {
	App = &ApplicationContext{
		DbConnector: db,
		Logger:      logger,
		JWTManager:  jwtManager,
	}
}
