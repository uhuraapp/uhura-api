package services

import (
	"bitbucket.org/dukex/uhura-api/entities"
	"bitbucket.org/dukex/uhura-api/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type UserRecommendationService struct {
	DB gorm.DB
}

func NewUserRecommendationService(db gorm.DB) UserRecommendationService {
	return UserRecommendationService{DB: db}
}

// Index TODO
func (r UserRecommendationService) Index(c *gin.Context) {
	// _userID, _ := c.Get("user_id")
	// userID := _userID.(string)

	var ids []string
	channels := make([]entities.Subscription, 0)

	if len(ids) > 0 {
		r.DB.Table(models.Channel{}.TableName()).Where("id in (?)", ids).Find(&channels)
	}

	c.JSON(200, gin.H{"recommendations": channels})
}
