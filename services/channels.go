package services

import (
	"strings"
	"time"

	"bitbucket.org/dukex/uhura-api/channels"
	"bitbucket.org/dukex/uhura-api/entities"
	"bitbucket.org/dukex/uhura-api/helpers"
	"bitbucket.org/dukex/uhura-api/models"
	"bitbucket.org/dukex/uhura-api/parser"
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

	uri := c.Params.ByName("uri")
	channelURI := strings.Replace(uri, "/", "", 1)

	err := s.DB.Table(models.Channel{}.TableName()).Where("uri = ?", channelURI).First(&channel).Error

	if err == gorm.RecordNotFound {
		url, err := helpers.ParseURL(uri)

		if err != nil {
			c.AbortWithStatus(404)
			return
		}

		feed, err := parser.URL(url)

		if err != nil {
			c.AbortWithStatus(404)
			return
		}

		channel = channels.TranslateFromFeedToEntity(channel, feed[0])

		episodes, ids := channels.TranslateEpisodesFromFeedToEntity(feed[0])
		channel.Episodes = ids

		c.JSON(200, gin.H{"channel": channel, "episodes": episodes})
		return
	}

	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	if helpers.CacheHeader(c, channel.UpdatedAt) {
		c.AbortWithStatus(304)
		return
	}

	userId, _ := helpers.GetUser(c)
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

	models.SetListenAttributesToEpisode(s.DB, userId, episodes, channelUri)

	return
}
