package services

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/uhuraapp/uhura-api/channels"
	"github.com/uhuraapp/uhura-api/database"
	"github.com/uhuraapp/uhura-api/entities"
	"github.com/uhuraapp/uhura-api/helpers"
	"github.com/uhuraapp/uhura-api/models"
	"github.com/uhuraapp/uhura-worker/sync"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// UserSubscriptionService TODO
type UserSubscriptionService struct {
	DB *gorm.DB
}

// NewUserSubscriptionService TODO
func NewUserSubscriptionService(db *gorm.DB) UserSubscriptionService {
	return UserSubscriptionService{DB: db}
}

// Index TODO
func (s UserSubscriptionService) Index(c *gin.Context) {
	_userID, _ := c.Get("user_id")
	userID := _userID.(string)

	origin := c.Request.Header.Get("Origin")
	referer, _ := url.Parse(origin)
	host := referer.Host

	if helpers.UseHTTPCache("user-subscriptions-index:"+userID+":"+host, database.CACHE, c) {
		return
	}

	subscriptions, _ := helpers.UserSubscriptions(userID, s.DB, models.Subscription{}.TableName(), models.Channel{}.TableName(), "")

	c.JSON(200, gin.H{"subscriptions": subscriptions})
}

// Show TODO
func (s UserSubscriptionService) Show(c *gin.Context) {
	var channel entities.Channel

	channelURI := c.Params.ByName("uri")
	_userID, _ := c.Get("user_id")
	userID := _userID.(string)

	s.DB.Table(models.Channel{}.TableName()).Where("uri = ?", channelURI).First(&channel)

	if userID != "" {
		channel.Subscribed = true
	}

	channel.Episodes = make([]int64, 0)

	c.JSON(200, gin.H{"subscription": channel})
}

// Delete TODO
func (s UserSubscriptionService) Delete(c *gin.Context) {
	var channel entities.Channel

	channelURI := c.Params.ByName("uri")
	_userID, _ := c.Get("user_id")
	userID := _userID.(string)

	database.CACHE.Delete("s:" + userID)

	if s.DB.Table(models.Channel{}.TableName()).Where("uri = ?", channelURI).First(&channel).Error != gorm.ErrRecordNotFound {
		s.DB.Table(models.Subscription{}.TableName()).Where("channel_id = ? AND user_id = ?", channel.Id, userID).
			Delete(models.Subscription{})

		channel.Subscribed = false

		channel.Episodes = make([]int64, 0)

		c.JSON(200, gin.H{"subscription": channel})
	}

	c.AbortWithStatus(404)
}

// Create TODO
func (s UserSubscriptionService) Create(c *gin.Context) {
	var (
		channel *entities.Channel
		ok      bool
		params  struct {
			Subscription struct {
				ChannelID  string `json:"channel_id"`
				ChannelURL string `json:"channel_url"`
			}
		}
	)

	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&params)
	if err != nil {
		c.AbortWithStatus(500)
	}
	_userID, _ := c.Get("user_id")
	userID := _userID.(string)

	database.CACHE.Delete("s:" + userID)

	if params.Subscription.ChannelURL != "" {
		_channel, ok := channels.Create(s.DB, params.Subscription.ChannelURL)

		if !ok {
			c.AbortWithStatus(500)
			return
		}

		sync.Sync(_channel.Id, s.DB)
		params.Subscription.ChannelID = _channel.Uri
	}

	userIDInt, _ := helpers.GetUser(c)
	channel, ok = s.createByChannelID(userIDInt, params.Subscription.ChannelID)

	if ok {
		go channels.Ping(s.DB, channel.Id)

		channel.Episodes = make([]int64, 0)

		c.JSON(200, gin.H{"subscription": channel})
	} else {
		c.AbortWithStatus(404)
	}
}

func (s UserSubscriptionService) createByChannelID(userID int, channelID string) (*entities.Channel, bool) {
	var channel entities.Channel

	if s.DB.Table(models.Channel{}.TableName()).Where("uri = ?", channelID).First(&channel).Error != gorm.ErrRecordNotFound {
		subscription := models.Subscription{
			UserId:    int64(userID),
			ChannelId: channel.Id,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := s.DB.Table(models.Subscription{}.TableName()).Where("user_id = ?", userID).Where("channel_id = ?", channel.Id).First(&models.Subscription{}).Error; err == gorm.ErrRecordNotFound {
			s.DB.Table(models.Subscription{}.TableName()).Save(&subscription)
		}

		channel.Subscribed = true

		return &channel, true
	}

	return nil, false
}
