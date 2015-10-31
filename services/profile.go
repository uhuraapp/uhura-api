package services

import (
	"strconv"

	"bitbucket.org/dukex/uhura-api/entities"
	"bitbucket.org/dukex/uhura-api/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type ProfileService struct {
	DB gorm.DB
}

func NewProfileService(db gorm.DB) ProfileService {
	return ProfileService{DB: db}
}

func (s ProfileService) Me(c *gin.Context) {
	_userID, _ := c.Get("user_id")
	userID := _userID.(string)

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

	c.JSON(200, gin.H{"profile": profile})
}

func (s ProfileService) Get(c *gin.Context) {
	id := c.Params.ByName("key")

	var profile entities.Profile
	err := s.DB.Table(models.Profile{}.TableName()).Where("key = ?", id).First(&profile).Error

	if err == nil {
		c.JSON(200, gin.H{"profile": profile})
		return
	}

	c.AbortWithStatus(404)
}
