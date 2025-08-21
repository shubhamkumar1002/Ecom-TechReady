package common

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ApplicationContext struct {
	DbConnector *gorm.DB
	Logger      *zap.Logger
}

var App *ApplicationContext

func Init(db *gorm.DB, logger *zap.Logger) {
	App = &ApplicationContext{
		DbConnector: db,
		Logger:      logger,
	}
}
