package services

import (
	"bitbucket.org/dukex/uhura-api/cache"
	"bitbucket.org/dukex/uhura-api/entities"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type SuggestionsService struct {
	DB gorm.DB
}

func NewSuggestionsService(db gorm.DB) SuggestionsService {
	return SuggestionsService{DB: db}
}

func (s SuggestionsService) Get(c *gin.Context) {
	var (
		channelsIds []int
	)

	channels := make([]entities.Channel, 0)
	episodes := make([]entities.Episode, 0)
	_userId, _ := c.Get("user_id")
	userId := _userId.(string)

	subscriptionsCached, err := cache.Get("s:ids:"+userId, channelsIds)
	if err == nil {
		channelsIds, _ = subscriptionsCached.([]int)

		if len(channelsIds) > 0 {
			s.DB.Table("channels").Where("channels.id in (?)", channelsIds).Find(&channels)
			s.DB.Raw("SELECT * FROM (SELECT items.*,row_number() OVER (PARTITION BY items.channel_id ORDER BY items.title) AS number_rows FROM items LEFT JOIN user_items ON user_items.item_id = items.id AND user_items.user_id = ? AND user_items.viewed = FALSE WHERE items.channel_id IN (?)) AS itemS WHERE number_rows <= 5 ORDER BY title", userId, channelsIds).Scan(&episodes)
		}
	}

	c.JSON(200, gin.H{"suggestions": channels})
}
