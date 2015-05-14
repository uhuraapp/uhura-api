package services

import (
	"time"

	"bitbucket.org/dukex/uhura-api/entities"
	"bitbucket.org/dukex/uhura-api/helpers"
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
	var userId int

	channelURI := c.Params.ByName("uri")

	userId, _ = helpers.GetUser(c)

	err := s.DB.Table(models.Channel{}.TableName()).Where("uri = ?", channelURI).First(&channel).Error

	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	channel.Episodes, episodes = s.getEpisodes(channel.Id, channelURI, userId)

	if userId != 0 {
		channel.Subscribed = s.DB.Table(models.Subscription{}.TableName()).Where("user_id = ?", userId).
			Where("channel_id = ?", channel.Id).
			Find(&models.Subscription{}).Error != gorm.RecordNotFound
	}

	c.JSON(200, gin.H{"channel": channel, "episodes": episodes})
}

func (s ChannelsService) Open(c *gin.Context) {
	var channel entities.Channel
	channelURI := c.Params.ByName("uri")
	s.DB.Table(models.Channel{}.TableName()).Update(&models.Channel{
		VisitedAt: time.Now().UTC(),
	}).Where("uri = ?", channelURI).First(&channel)

	c.JSON(200, gin.H{})
}

func (s ChannelsService) getEpisodes(channelID int64, channelUri string, userId int) (ids []int64, episodes []*entities.Episode) {
	s.DB.Table(models.Episode{}.TableName()).
		Where("channel_id = ?", channelID).
		Order("published_at DESC").
		Limit(20).
		Find(&episodes)

	for _, e := range episodes {
		ids = append(ids, e.Id)
	}

	entities.SetListenAttributesToEpisode(s.DB, userId, episodes, channelUri)

	return
}
