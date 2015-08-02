package services

import (
	"encoding/json"
	"strconv"
	"time"

	"bitbucket.org/dukex/uhura-api/channels"
	"bitbucket.org/dukex/uhura-api/database"
	"bitbucket.org/dukex/uhura-api/entities"
	"bitbucket.org/dukex/uhura-api/helpers"
	"bitbucket.org/dukex/uhura-api/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// UserSubscriptionService TODO
type UserSubscriptionService struct {
	DB gorm.DB
}

// NewUserSubscriptionService TODO
func NewUserSubscriptionService(db gorm.DB) UserSubscriptionService {
	return UserSubscriptionService{DB: db}
}

// Index TODO
func (s UserSubscriptionService) Index(c *gin.Context) {
	var ids []int
	subscriptions := make([]entities.Subscription, 0)

	_userID, _ := c.Get("user_id")
	userID := _userID.(string)

	t := time.Now()

	cacheRes, err := database.CACHE.Value("s:" + userID)
	if err == nil {
		t = *cacheRes.Data().(*time.Time)
	}

	if helpers.CacheHeader(c, t) {
		c.AbortWithStatus(304)
		return
	}

	s.DB.Table(models.Subscription{}.TableName()).Where("user_id = ?", userID).
		Pluck("channel_id", &ids)

	if len(ids) > 0 {
		s.DB.Table(models.Channel{}.TableName()).Where("id in (?)", ids).Find(&subscriptions)
	}

	// for i, _ := range subscriptions {
	//	subscriptions[i].Uri = channel.FixUri()
	//	go subscriptions[i].SetSubscribed(userId)
	//	subscriptions[i].SetEpisodesIds()
	//	subscriptions[i].ToView = subscriptions[i].GetToView(s.DB, userId)
	//	subscriptions[i].Subscribed = true
	// }

	database.CACHE.Add("s:"+userID, ((7 * 24) * time.Hour), &t)
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

	c.JSON(200, gin.H{"subscription": channel})
}

// Delete TODO
func (s UserSubscriptionService) Delete(c *gin.Context) {
	var channel entities.Channel

	channelURI := c.Params.ByName("uri")
	_userID, _ := c.Get("user_id")
	userID := _userID.(string)

	database.CACHE.Delete("s:" + userID)

	if s.DB.Table(models.Channel{}.TableName()).Where("uri = ?", channelURI).First(&channel).Error != gorm.RecordNotFound {
		s.DB.Table(models.Subscription{}.TableName()).Where("channel_id = ? AND user_id = ?", channel.Id, userID).
			Delete(models.Subscription{})

		channel.Subscribed = false

		go helpers.NewEvent(userID, "unsubscribed", map[string]interface{}{"channel_id": channel.Id})
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
		params.Subscription.ChannelID = _channel.Uri
	}

	userIDInt, _ := helpers.GetUser(c)
	channel, ok = s.createByChannelID(userIDInt, params.Subscription.ChannelID)

	go channels.Ping(s.DB, channel.Id)

	if ok {
		c.JSON(200, gin.H{"subscription": channel})
	} else {
		c.AbortWithStatus(404)
	}
}

func (s UserSubscriptionService) createByChannelID(userID int, channelID string) (*entities.Channel, bool) {
	var channel entities.Channel

	if s.DB.Table(models.Channel{}.TableName()).Where("uri = ?", channelID).First(&channel).Error != gorm.RecordNotFound {
		subscription := models.Subscription{
			UserId:    int64(userID),
			ChannelId: channel.Id,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := s.DB.Table(models.Subscription{}.TableName()).Where("user_id = ?", userID).Where("channel_id = ?", channel.Id).First(&models.Subscription{}).Error; err == gorm.RecordNotFound {
			s.DB.Table(models.Subscription{}.TableName()).Save(&subscription)
		}

		channel.Subscribed = true

		go helpers.NewEvent(strconv.Itoa(userID), "subscribed", map[string]interface{}{"channel_id": channel.Id})

		return &channel, true
	}

	return nil, false
}
