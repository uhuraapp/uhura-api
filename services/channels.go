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
	var episodes []entities.Episode

	channelURI := c.Params.ByName("uri")

	s.DB.Table(models.Channel{}.TableName()).Where("uri = ?", channelURI).First(&channel)
	channel.Episodes, episodes = s.getEpisodesIDs(channel.Id)

	c.JSON(200, gin.H{"channel": channel, "episodes": episodes})
}

func (s ChannelsService) getEpisodesIDs(channelID int64) (ids []int64, episodes []entities.Episode) {
	s.DB.Table(models.Episode{}.TableName()).
		Where("channel_id = ?", channelID).
		Order("published_at DESC").
		Find(&episodes).
		Pluck("id", &ids)

	return
}
