package database

import (
	"log"
	"os"

	"bitbucket.org/dukex/uhura-api/models"
	"github.com/jinzhu/gorm"
	pq "github.com/lib/pq"
)

func New() gorm.DB {
	var database gorm.DB
	var err error

	databaseUrl, _ := pq.ParseURL(os.Getenv("DATABASE_URL"))
	database, err = gorm.Open("postgres", databaseUrl)

	if err != nil {
		log.Fatalln(err.Error())
	}

	err = database.DB().Ping()
	if err != nil {
		log.Fatalln(err.Error())
	}

	database.LogMode(os.Getenv("DEBUG") == "true")

	return database
}

func Migrations(database gorm.DB) {
	database.AutoMigrate(models.Episode{})
	database.AutoMigrate(models.Listened{})
	database.AutoMigrate(models.Channel{})
	database.AutoMigrate(models.Subscription{})
}
