package services

import (
	"bitbucket.org/dukex/uhura-api/cache"
	"bitbucket.org/dukex/uhura-api/entities"
	"bitbucket.org/dukex/uhura-api/helpers"
	"bitbucket.org/dukex/uhura-api/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type SubscriptionService struct {
	DB gorm.DB
}

func NewSubscriptionService(db gorm.DB) SubscriptionService {
	return SubscriptionService{DB: db}
}

func (s SubscriptionService) Get(c *gin.Context) {
	var ids []int

	subscriptions := make([]entities.Subscription, 0)
	_userId, _ := c.Get("user_id")
	userId := _userId.(string)

	if !helpers.IsABotUser(userId) {
		subscriptionsCached, err := cache.Get("s:ids:"+userId, ids)

		if err == nil {
			var ok bool
			ids, ok = subscriptionsCached.([]int)
			if !ok {
				s.DB.Table(models.Subscription{}.TableName()).Where("user_id = ?", userId).
					Pluck("channel_id", &ids)
				go cache.Set("s:ids:"+userId, ids)
			}
		} else {
			s.DB.Table(models.Subscription{}.TableName()).Where("user_id = ?", userId).
				Pluck("channel_id", &ids)
			go cache.Set("s:ids:"+userId, ids)
		}

		if len(ids) > 0 {
			s.DB.Table(models.Channel{}.TableName()).Where("id in (?)", ids).Find(&subscriptions)
		}

		for i, _ := range subscriptions {
			//subscriptions[i].Uri = channel.FixUri()
			//go subscriptions[i].SetSubscribed(userId)
			//subscriptions[i].SetEpisodesIds()
			subscriptions[i].ToView = subscriptions[i].GetToView(s.DB, userId)
		}
	}

	c.JSON(200, gin.H{"subscriptions": subscriptions})
}
