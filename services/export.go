package services

import (
	"encoding/xml"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/uhuraapp/uhura-api/models"
)

type ExportSevice struct {
	DB *gorm.DB
}

func NewExportService(db *gorm.DB) ExportSevice {
	return ExportSevice{DB: db}
}

type OPMLHead struct {
	Title        string    `xml:"title"`
	DateCreated  time.Time `xml:"dateCreated"`
	DateModified time.Time `xml:"dateModified"`
	OwnerName    string    `xml:"ownerName"`
	OwnerEmail   string    `xml:"ownerEmail"`
}

type OPMLBody struct {
	Outline []OPMLOutline `xml:"outline"`
}

type OPMLOutline struct {
	Title  string `xml:"title,attr"`
	Text   string `xml:"text,attr"`
	URL    string `xml:"url,attr"`
	XMLURL string `xml:"xmlUrl,attr"`
}

type OPML struct {
	XMLName xml.Name `xml:"opml"`
	Head    OPMLHead `xml:"head"`
	Body    OPMLBody `xml:"body"`
}

func (s ExportSevice) Get(c *gin.Context) {
	userID := c.Param("id")
	auth := c.Query("auth_key")

	var user models.User
	err := s.DB.Table(models.User{}.TableName()).
		Unscoped().
		Where("export_key = ?", auth).
		Where("id = ?", userID).First(&user).Error

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	var channels []models.Channel
	s.DB.Table(models.Channel{}.TableName()).Select("DISTINCT channels.*").Joins("JOIN user_channels s ON channels.id = s.channel_id").Where("s.user_id = ?", userID).Find(&channels)

	outlines := make([]OPMLOutline, 0)
	for _, s := range channels {
		outlines = append(outlines, OPMLOutline{
			s.Title,
			s.Title,
			s.Url,
			s.Url,
		})
	}

	c.XML(200, OPML{
		Head: OPMLHead{
			Title:        user.Email + " Subscriptions",
			DateCreated:  time.Now(),
			DateModified: time.Now(),
			OwnerName:    user.Name,
			OwnerEmail:   user.Email,
		},
		Body: OPMLBody{outlines},
	})
}
