package services

import (
	"bitbucket.org/dukex/uhura-api/entities"
	"bitbucket.org/dukex/uhura-api/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type CategoriesService struct {
	DB gorm.DB
}

func NewCategoriesService(db gorm.DB) CategoriesService {
	return CategoriesService{DB: db}
}

func (s CategoriesService) Index(c *gin.Context) {
	var tmpCategories []entities.Category
	categories := make([]entities.Category, 0)

	s.DB.Table(models.Category{}.TableName()).Find(&tmpCategories)

	channelsIds := []int64{0}

	for i, _ := range tmpCategories {
		var id []int64
		s.DB.Table(models.Categoriable{}.TableName()).
			Where("channel_id NOT IN (?)", channelsIds).
			Where("category_id = ?", tmpCategories[i].Id).
			Order("channel_id DESC").
			Limit(1).
			Pluck("channel_id", &id)

		if len(id) > 0 {
			s.DB.Table(models.Channel{}.TableName()).Where("id = ?", id[0]).First(&tmpCategories[i].Channel)
			categories = append(categories, tmpCategories[i])
		}

		channelsIds = append(channelsIds, id...)
	}

	c.JSON(200, gin.H{"categories": categories})
}
