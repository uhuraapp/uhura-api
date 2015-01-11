package services

import (
	"encoding/json"
	"strconv"
	"time"

	"bitbucket.org/dukex/uhura-api/entities"
	"bitbucket.org/dukex/uhura-api/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type UserSubscriptionService struct {
	DB gorm.DB
}

func NewUserSubscriptionService(db gorm.DB) UserSubscriptionService {
	return UserSubscriptionService{DB: db}
}

func (s UserSubscriptionService) Index(c *gin.Context) {
	var ids []int

	subscriptions := make([]entities.Subscription, 0)
	_userId, _ := c.Get("user_id")
	userId := _userId.(string)

	s.DB.Table(models.Subscription{}.TableName()).Where("user_id = ?", userId).
		Pluck("channel_id", &ids)

	if len(ids) > 0 {
		s.DB.Table(models.Channel{}.TableName()).Where("id in (?)", ids).Find(&subscriptions)
	}

	for i, _ := range subscriptions {
		//subscriptions[i].Uri = channel.FixUri()
		//go subscriptions[i].SetSubscribed(userId)
		//subscriptions[i].SetEpisodesIds()
		subscriptions[i].ToView = subscriptions[i].GetToView(s.DB, userId)
		// subscriptions[i].Subscribed = true
	}

	c.JSON(200, gin.H{"subscriptions": subscriptions})
}

func (s UserSubscriptionService) Show(c *gin.Context) {
	var channel entities.Channel

	channelURI := c.Params.ByName("uri")
	_userId, _ := c.Get("user_id")
	userId := _userId.(string)

	s.DB.Table(models.Channel{}.TableName()).Where("uri = ?", channelURI).First(&channel)

	if userId != "" {
		channel.Subscribed = true
	}

	c.JSON(200, gin.H{"subscription": channel})
}

func (s UserSubscriptionService) Delete(c *gin.Context) {
	var channel entities.Channel

	channelURI := c.Params.ByName("uri")
	_userId, _ := c.Get("user_id")
	userId := _userId.(string)

	if s.DB.Table(models.Channel{}.TableName()).Where("uri = ?", channelURI).First(&channel).Error != gorm.RecordNotFound {
		s.DB.Table(models.Subscription{}.TableName()).Where("channel_id = ? AND user_id = ?", channel.Id, userId).
			Delete(models.Subscription{})

		channel.Subscribed = false

		c.JSON(200, gin.H{"subscription": channel})
	}

	c.Abort(404)
}

func (s UserSubscriptionService) Create(c *gin.Context) {
	var (
		channel entities.Channel
		sp      struct {
			Subscription struct {
				ChannelId string `json:"channel_id"`
			}
		}
	)

	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&sp)
	if err != nil {
		c.Abort(500)
	}

	_userId, _ := c.Get("user_id")
	userId, _ := strconv.Atoi(_userId.(string))

	if s.DB.Table(models.Channel{}.TableName()).Where("uri = ?", sp.Subscription.ChannelId).First(&channel).Error != gorm.RecordNotFound {
		subscription := models.Subscription{
			UserId:    int64(userId),
			ChannelId: channel.Id,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		s.DB.Table(models.Subscription{}.TableName()).Where("user_id = ?", userId).
			Where("channel_id = ?", channel.Id).FirstOrCreate(&subscription)
		channel.Subscribed = true
		c.JSON(200, gin.H{"subscription": channel})
	} else {
		c.Abort(404)
	}
}
