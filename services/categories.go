package services

import (
	"bitbucket.org/dukex/uhura-api/entities"
	"bitbucket.org/dukex/uhura-api/helpers"
	"bitbucket.org/dukex/uhura-api/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type CategoriesService struct {
	DB         gorm.DB
	categories []*entities.Category
	channels   []entities.Channel
}

func NewCategoriesService(db gorm.DB) CategoriesService {
	c := CategoriesService{DB: db}
	return c
}

func (s *CategoriesService) Index(c *gin.Context) {
	s.cacheCategoriesAndChannels()
	c.JSON(200, gin.H{"categories": s.categories, "channels": s.channels})
}

func (s *CategoriesService) cacheCategoriesAndChannels() {
	var categories []*entities.Category
	var channels []entities.Channel

	if len(s.categories) == 0 {
		s.DB.Table(models.Category{}.TableName()).Find(&categories)

		categoriesIDs := helpers.MapInt(categories, fncategoryID)

		if len(categoriesIDs) > 0 {
			var channelsCategories []models.Categoriable

			s.DB.Select("DISTINCT(channel_id), category_id").Table(models.Categoriable{}.TableName()).
				Where("channel_id NOT IN (0)").
				Where("category_id IN (?)", categoriesIDs).
				Find(&channelsCategories)

			channelsIDs := make([]int64, 0)
			for _, cc := range channelsCategories {
				channelsIDs = append(channelsIDs, cc.ChannelId)
			}

			channelsIDs = helpers.RemoveDuplicates(channelsIDs)

			s.DB.Table(models.Channel{}.TableName()).
				Where("id IN (?)", channelsIDs).
				Find(&channels)

			for _, cc := range channelsCategories {
				category, fcategory := findCategory(categories, cc.CategoryId)
				channel, fchannel := findChannel(channels, cc.ChannelId)
				if fcategory && fchannel {
					category.ChannelsIDs = append(category.ChannelsIDs, channel.Uri)
				}
			}
		}
		s.categories = categories
		s.channels = channels
	}
}

func findCategory(categories []*entities.Category, id int64) (*entities.Category, bool) {
	for _, c := range categories {
		if c.Id == id {
			return c, true
		}
	}

	return nil, false
}

func findChannel(channels []entities.Channel, id int64) (entities.Channel, bool) {
	for _, c := range channels {
		if c.Id == id {
			return c, true
		}
	}

	return entities.Channel{}, false
}

// ------

func fncategoryID(c interface{}) int64 {
	return c.(*entities.Category).Id
}
