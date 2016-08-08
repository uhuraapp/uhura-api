package services

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/uhuraapp/uhura-api/channels"
	"github.com/uhuraapp/uhura-api/entities"
	"github.com/uhuraapp/uhura-api/helpers"
	"github.com/uhuraapp/uhura-api/models"
)

type ChannelsService struct {
	DB *gorm.DB
}

func NewChannelsService(db *gorm.DB) ChannelsService {
	return ChannelsService{DB: db}
}

func (s ChannelsService) Top(c *gin.Context) {
	var channels []entities.Channel
	s.DB.Table(models.Subscription{}.TableName()).
		Select("channels.title, channels.image_url, channels.description, channels.uri, channels.id, COUNT(*) AS subscribers_count").
		Joins("INNER JOIN channels ON user_channels.channel_id = channels.id").
		Group("1,2,3,4,5 ORDER BY subscribers_count DESC").
		Limit(5).
		Find(&channels)

	c.JSON(200, gin.H{"channels": channels})
}

func (s ChannelsService) Get(c *gin.Context) {
	channel, episodes, _, found := channels.Find(s.DB, c.Params.ByName("uri"))

	if !found {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	userId, _ := helpers.GetUser(c)
	episodes = s.getPlayed(userId, episodes)

	if userId != 0 {
		channel.Subscribed = s.DB.Table(models.Subscription{}.TableName()).Where("user_id = ?", userId).
			Where("channel_id = ?", channel.Id).
			Find(&models.Subscription{}).Error != gorm.ErrRecordNotFound
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

func (s ChannelsService) getPlayed(userId int, episodes entities.Episodes) entities.Episodes {
	var played []*models.Listened

	ids := episodes.IDs()

	if len(ids) > 0 {
		s.DB.Table(models.Listened{}.TableName()).
			Where("item_uid IN (?)", ids).
			Where("user_id = ?", userId).
			Find(&played)
	}

	mapPlayed := make(map[string]*models.Listened, 0)
	for _, play := range played {
		mapPlayed[play.ItemUID] = play
	}

	for _, episode := range episodes {
		if mapPlayed[episode.Id] != nil {
			episode.Listened = mapPlayed[episode.Id].Viewed
			episode.StoppedAt = &mapPlayed[episode.Id].StoppedAt
			if episode.Listened {
				z := int64(0)
				episode.StoppedAt = &z
			}
		}
	}

	return episodes
}
