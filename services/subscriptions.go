package services

import (
	"bitbucket.org/dukex/uhura-api/entities"
	"bitbucket.org/dukex/uhura-api/gpodder"
	"bitbucket.org/dukex/uhura-api/helpers"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type SubscriptionService struct {
	DB gorm.DB
}

func NewSubscriptionService(db gorm.DB) SubscriptionService {
	return SubscriptionService{DB: db}
}

func (s SubscriptionService) Top(c *gin.Context) {
	data, err := gpodder.Top("20")
	if err != nil {
		c.AbortWithStatus(500)
	}

	channels := make([]entities.Channel, len(data))

	for i, _ := range data {
		channels[i] = entities.Channel{
			Title:       data[i].Title,
			Description: data[i].Description,
			ImageUrl:    data[i].LogoUrl,
			Uri:         new(helpers.Uriable).MakeUri(data[i].Title),
		}
	}

	c.JSON(200, map[string]interface{}{"channels": channels})
}
