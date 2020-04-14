package migrations

import (
	"sms-aiforesee-be/database"
	"sms-aiforesee-be/models"
)

func Migrate() error {
	db, err := database.ConnectToORM()

	// unable to connect to database
	if err != nil {
		return err
	}

	// ping to database
	err = db.DB().Ping()

	// error ping to database
	if err != nil {
		return err
	}
	// migration
	db.AutoMigrate(&models.User{}, &models.Event{})
	return nil
}
