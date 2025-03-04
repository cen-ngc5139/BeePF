package database

import (
	"fmt"
	"log"

	"github.com/cen-ngc5139/BeePF/server/conf"

	"gorm.io/gorm/logger"

	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

// Setup initializes the database instance
func Setup() {
	config := conf.Config()

	var err error
	db, err := gorm.Open(mysql.Open(
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local",
			config.Database.User,
			config.Database.Password,
			config.Database.Host,
			config.Database.Name)), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	if config.Database.LogMode == "debug" {
		db.Logger = logger.Default.LogMode(logger.Info)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	sqlDB.SetMaxIdleConns(config.Database.MaxIdle)
	sqlDB.SetMaxOpenConns(config.Database.MaxOpen)

	DB = db
}
