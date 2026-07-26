package initializers

import "gochatroom/models"

func SyncDb() {
	DB.AutoMigrate(&models.User{})
}
