package services

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	// "github.com/uhuraapp/uhura-api/cache"
	"github.com/uhuraapp/uhura-api/entities"
	// "github.com/uhuraapp/uhura-api/helpers"
	"github.com/uhuraapp/uhura-api/models"
)

type ChannelsService struct {
	DB gorm.DB
}

func NewChannelsService(db gorm.DB) ChannelsService {
	return ChannelsService{DB: db}
}

func (s ChannelsService) Get(c *gin.Context) {
	var channel entities.Channel

	channelURI := c.Params.ByName("uri")

	s.DB.Table(models.Channel{}.TableName()).Where("uri = ?", channelURI).First(&channel)

	c.JSON(200, gin.H{"channel": channel})
}
