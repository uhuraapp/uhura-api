package services

import (
	"net/http"
	"time"

	"bitbucket.org/dukex/uhura-api/database"
	"bitbucket.org/dukex/uhura-api/entities"
	"bitbucket.org/dukex/uhura-api/helpers"
	"bitbucket.org/dukex/uhura-api/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type SubscriptionsService struct {
	DB gorm.DB
}

func NewSubscriptionsService(db gorm.DB) SubscriptionsService {
	return SubscriptionsService{DB: db}
}

func (s SubscriptionsService) Get(c *gin.Context) {
	var ids []int
	subscriptions := make([]entities.Subscription, 0)

	t := time.Now()
	keyOrUsername := c.Params.ByName("key")

	cacheRes, err := database.CACHE.Value("u:" + keyOrUsername)
	if err == nil {
		t = *cacheRes.Data().(*time.Time)
	}

	if helpers.CacheHeader(c, t) {
		c.AbortWithStatus(http.StatusNotModified)
		return
	}

	var u models.User

	err = s.DB.Table(models.Profile{}.TableName()).Where("key = ?", keyOrUsername).First(&u).Error
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	s.DB.Table(models.Subscription{}.TableName()).Where("user_id = ?", u.Id).
		Order("channel_id").
		Pluck("channel_id", &ids)

	if len(ids) > 0 {
		s.DB.Table(models.Channel{}.TableName()).Where("id in (?)", ids).Order("title ASC").Find(&subscriptions)
	}

	database.CACHE.Add("u:"+keyOrUsername, ((7 * 24) * time.Hour), &t)
	c.JSON(200, gin.H{"subscriptions": subscriptions})
}
