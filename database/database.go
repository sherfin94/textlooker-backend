package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Database *gorm.DB

func ConnectDatabase(databaseName string) *gorm.DB {
	dsn := fmt.Sprintln("host=localhost user=gorm password=gorm dbname=", databaseName, " port=5432 sslmode=disable TimeZone=Asia/Kolkata")
	var err error
	Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic("Could not connect to database")
	}

	return Database
}
