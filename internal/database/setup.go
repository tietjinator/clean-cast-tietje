package setup

import (
	"ikoyhn/podcast-sponsorblock/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func ConnectDatabase() (*gorm.DB, error) {

	database, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	err = database.AutoMigrate(&models.PodcastEpisode{})
	if err != nil {
		return nil, err
	}

	return database, nil
}
