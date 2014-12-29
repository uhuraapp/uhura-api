package services

import (
	"bitbucket.org/dukex/uhura-api/entities"
	"bitbucket.org/dukex/uhura-api/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type ChannelsService struct {
	DB gorm.DB
}

func NewChannelsService(db gorm.DB) ChannelsService {
	return ChannelsService{DB: db}
}

func (s ChannelsService) Get(c *gin.Context) {
	var channel entities.Channel
	var episodes []*entities.Episode

	channelURI := c.Params.ByName("uri")

	s.DB.Table(models.Channel{}.TableName()).Where("uri = ?", channelURI).First(&channel)
	channel.Episodes, episodes = s.getEpisodes(channel.Id, channelURI)

	c.JSON(200, gin.H{"channel": channel, "episodes": episodes})
}

func (s ChannelsService) getEpisodes(channelID int64, channelUri string) (ids []int64, episodes []*entities.Episode) {
	s.DB.Table(models.Episode{}.TableName()).
		Where("channel_id = ?", channelID).
		Order("published_at DESC").
		Limit(20).
		Find(&episodes).
		Pluck("id", &ids)

	for _, episode := range episodes {
		episode.ChannelUri = channelUri
	}

	return
}
