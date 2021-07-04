package models

import (
	"textlooker-backend/database"
)

func ApplyMigrations(databaseName string, logMode database.LogMode) error {
	database.ConnectDatabase(databaseName, logMode)

	return database.Database.AutoMigrate(
		&User{},
		&UserRegistration{},
		&Source{},
		&Insight{},
		&Dashboard{},
		&DashboardInsight{},
	)
}
