package models

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

func ConnectDatabase(databaseName string) *gorm.DB {
	dsn := fmt.Sprintln("host=localhost user=gorm password=gorm dbname=", databaseName, " port=5432 sslmode=disable TimeZone=Asia/Kolkata")
	var err error
	Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Could not connect to database")
	}

	return Database
}

func ApplyMigrations(databaseName string) error {
	ConnectDatabase(databaseName)
	return Database.AutoMigrate(&User{})
}
