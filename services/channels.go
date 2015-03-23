package services

import (
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
	var userId string

	channelURI := c.Params.ByName("uri")

	_userId, err := c.Get("user_id")
	if err == nil {
		userId = _userId.(string)
	}

	s.DB.Table(models.Channel{}.TableName()).Where("uri = ?", channelURI).First(&channel)
	channel.Episodes, episodes = s.getEpisodes(channel.Id, channelURI, userId)

	if userId != "" {
		channel.Subscribed = s.DB.Table(models.Subscription{}.TableName()).Where("user_id = ?", userId).
			Where("channel_id = ?", channel.Id).
			Find(&models.Subscription{}).Error != gorm.RecordNotFound
	}

	c.JSON(200, gin.H{"channel": channel, "episodes": episodes})
}

func (s ChannelsService) getEpisodes(channelID int64, channelUri string, userId string) (ids []int64, episodes []*entities.Episode) {
	s.DB.Table(models.Episode{}.TableName()).
		Where("channel_id = ?", channelID).
		Order("published_at DESC").
		Limit(20).
		Find(&episodes).
		Pluck("id", &ids)

	var listeneds []int64

	if len(ids) > 0 {
		s.DB.Table(models.Listened{}.TableName()).
			Where("item_id IN (?)", ids).
			Where("viewed = true").
			Where("user_id = ?", userId).
			Pluck("item_id", &listeneds)
	} else {
		listeneds = make([]int64, 0)
	}

	for _, episode := range episodes {
		episode.ChannelUri = channelUri
		episode.Listened = helpers.Contains(listeneds, episode.Id)
	}

	return
}
