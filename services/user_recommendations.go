package services

import (
	"log"
	"os"

	"bitbucket.org/dukex/uhura-api/entities"
	"bitbucket.org/dukex/uhura-api/models"
	"github.com/gin-gonic/gin"
	"github.com/hjr265/too"
	"github.com/jinzhu/gorm"
)

type UserRecommendationService struct {
	DB gorm.DB
	e  *too.Engine
}

func NewUserRecommendationService(db gorm.DB) UserRecommendationService {
	e, err := too.New(os.Getenv("REDIS_URL"), "channels")

	if err != nil {
		log.Fatal(err)
	}
	return UserRecommendationService{DB: db, e: e}
}

// Index TODO
func (r UserRecommendationService) Index(c *gin.Context) {
	_userID, _ := c.Get("user_id")
	userID := _userID.(string)

	var ids []string
	channels := make([]entities.Subscription, 0)

	items, _ := r.e.Suggestions.For(too.User(userID), 5)

	for _, id := range items {
		ids = append(ids, string(id))
	}

	if len(ids) > 0 {
		r.DB.Table(models.Channel{}.TableName()).Where("id in (?)", ids).Find(&channels)
	}

	c.JSON(200, gin.H{"recommendations": channels})
}
