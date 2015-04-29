package services

import (
	"bitbucket.org/dukex/uhura-api/helpers"
	"bitbucket.org/dukex/uhura-api/models"
	"bitbucket.org/dukex/uhura-api/parser"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type ParserService struct {
	DB gorm.DB
}

func NewParser(db gorm.DB) ParserService {
	return ParserService{db}
}

func (s ParserService) ByURL(c *gin.Context) {
	url, err := helpers.ParseURL(c.Request.URL.Query().Get("url"))
	if err != nil {
		c.JSON(500, map[string]string{"error": "URL invalid"})
	}

	channel, err := parser.URL(url)

	if err != nil {
		c.JSON(500, map[string]interface{}{"error": err})
	}

	channel.UhuraId = s.findUhuraID(channel)
	c.JSON(200, map[string]interface{}{
		"channel": channel,
	})
}

func (s ParserService) findUhuraID(c *parser.Channel) string {
	var channels []models.Channel

	s.DB.Table(models.Channel{}.TableName()).Where("url in (?)", c.Links).Find(&channels)

	if len(channels) < 1 {
		return ""
	}

	return channels[0].Uri
}
