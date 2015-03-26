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
		Find(&episodes)

	for _, e := range episodes {
		ids = append(ids, e.Id)
	}

	var listeneds []*models.Listened

	if len(ids) > 0 {
		s.DB.Table(models.Listened{}.TableName()).
			Where("item_id IN (?)", ids).
			Where("user_id = ?", userId).
			Find(&listeneds)
	} else {
		listeneds = make([]*models.Listened, 0)
	}

	mapListened := make(map[int64]*models.Listened, 0)
	for _, listen := range listeneds {
		mapListened[listen.ItemId] = listen
	}

	for _, episode := range episodes {
		episode.ChannelUri = channelUri

		if mapListened[episode.Id] != nil {
			episode.Listened = mapListened[episode.Id].Viewed
			episode.StoppedAt = mapListened[episode.Id].StoppedAt
		}
	}

	return
}
