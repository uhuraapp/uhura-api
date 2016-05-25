package services

import (
	"strconv"

	"github.com/uhuraapp/uhura-api/database"
	"github.com/uhuraapp/uhura-api/entities"
	"github.com/uhuraapp/uhura-api/helpers"
	"github.com/uhuraapp/uhura-api/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type ProfileService struct {
	DB *gorm.DB
}

func NewProfileService(db *gorm.DB) ProfileService {
	return ProfileService{DB: db}
}

func (s ProfileService) Me(c *gin.Context) {
	_userID, _ := c.Get("user_id")
	userID := _userID.(string)

	if helpers.UseHTTPCache("me:"+userID, database.CACHE, c) {
		return
	}

	var profile entities.Profile

	err := s.DB.Table(models.Profile{}.TableName()).Where("user_id = ?", userID).First(&profile).Error
	if err == nil {
		c.JSON(200, gin.H{"profile": profile})
		return
	}
	c.AbortWithStatus(404)
}

func (s ProfileService) Create(c *gin.Context) {
	_userID, _ := c.Get("user_id")
	userID := _userID.(string)

	var profile entities.Profile

	err := s.DB.Table(models.Profile{}.TableName()).Where("user_id = ?", userID).First(&profile).Error

	if err != nil {
		userIDInt, _ := strconv.Atoi(userID)
		profileModel := models.Profile{
			Key:    strconv.FormatInt(int64(userIDInt), 36),
			UserID: int64(userIDInt),
		}

		err = s.DB.Table(models.Profile{}.TableName()).Save(&profileModel).Error
		if err == nil {
			err = s.DB.Table(models.Profile{}.TableName()).Where("user_id = ?", userID).First(&profile).Error
		}
	}

	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	subscriptions, ids := helpers.UserSubscriptions(profile.UserID, s.DB, models.Subscription{}.TableName(), models.Channel{}.TableName(), profile.Key)
	profile.Channels = ids

	c.JSON(200, gin.H{"profile": profile, "channels": subscriptions})
}

func (s ProfileService) Get(c *gin.Context) {
	id := c.Params.ByName("key")

	if helpers.UseHTTPCache("profile-get:"+id, database.CACHE, c) {
		return
	}

	var profile entities.Profile
	err := s.DB.Table(models.Profile{}.TableName()).Where("key = ?", id).First(&profile).Error

	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	subscriptions, ids := helpers.UserSubscriptions(profile.UserID, s.DB, models.Subscription{}.TableName(), models.Channel{}.TableName(), id)
	profile.Channels = ids

	c.JSON(200, gin.H{"profile": profile, "channels": subscriptions})
}
