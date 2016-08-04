package services

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/uhuraapp/uhura-api/helpers"
	"github.com/uhuraapp/uhura-api/models"
	"github.com/uhuraapp/uhura-api/parser"
)

type ParserService struct {
	DB *gorm.DB
}

func NewParser(db *gorm.DB) ParserService {
	return ParserService{db}
}

func (s ParserService) ByURL(c *gin.Context) {
	url, err := helpers.ParseURL(c.Request.URL.Query().Get("url"))
	if err != nil {
		c.JSON(500, map[string]string{"error": "URL invalid"})
	}

	channel, _, _err := parser.URL(url)

	channel.UhuraID = FindUhuraID(s.DB, channel)

	c.JSON(200, gin.H{
		"channel": channel,
		"errors":  _err.Error(),
	})
}

func FindUhuraID(db *gorm.DB, c *parser.Channel) string {
	var uris []string

	db.Table(models.Channel{}.TableName()).Where("url in (?)", c.Links).Pluck("uri", &uris)

	if len(uris) < 1 {
		return ""
	}

	return uris[0]
}
