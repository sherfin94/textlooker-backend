package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type LogMode uint8

var Database *gorm.DB

const Silent, OnlyErrors, Loud = 1, 2, 3

func ConnectDatabase(databaseName string, logMode LogMode) *gorm.DB {
	dsn := fmt.Sprintln("host=localhost user=gorm password=gorm dbname=", databaseName, " port=5432 sslmode=disable TimeZone=Asia/Kolkata")
	var err error
	var logLevel logger.LogLevel

	switch logMode {
	case Silent:
		logLevel = logger.Silent
	case Loud:
		logLevel = logger.Info
	case OnlyErrors:
		logLevel = logger.Error
	}

	Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		panic("Could not connect to database")
	}

	return Database
}
