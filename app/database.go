package app

import (
	"belajar-go-restful-api/helper"
	models "belajar-go-restful-api/model/domain"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewDB() *gorm.DB{
	dsn := "host=localhost user=bagaskara password=@Bagas123 dbname=ecommerce port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	
	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)

	gormConfig := &gorm.Config{
		// enhance performance config
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger:                 dbLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	helper.PanicIfError(err)

	RunMigration(db)

	return db
}

func RunMigration(db *gorm.DB){
	db.AutoMigrate(
		&models.Product{},
	)
}