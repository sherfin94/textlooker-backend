package models

import (
	"textlooker-backend/database"
)

func ApplyMigrations(databaseName string) error {
	database.ConnectDatabase(databaseName)

	return database.Database.AutoMigrate(
		&User{},
		&UserRegistration{},
		&Source{},
	)
}
