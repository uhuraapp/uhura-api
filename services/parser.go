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

	channels, _err := parser.URL(url)

	for _, channel := range channels {
		channel.UhuraID = s.findUhuraID(channel)
	}

	c.JSON(200, gin.H{
		"channels": channels,
		"errors":   _err.Error(),
	})
}

func (s ParserService) findUhuraID(c *parser.Channel) string {
	var uris []string

	s.DB.Table(models.Channel{}.TableName()).Where("url in (?)", c.Links).Pluck("uri", &uris)

	if len(uris) < 1 {
		return ""
	}

	return uris[0]
}
