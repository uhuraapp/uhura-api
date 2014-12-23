package services

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/uhuraapp/uhura-api/models"
)

type EpisodeService struct {
	DB gorm.DB
}

func NewEpisodesService(db gorm.DB) EpisodeService {
	return EpisodeService{DB: db}
}

func (s EpisodeService) Listened(c *gin.Context) {
	var episode models.Episode
	episodeId, _ := strconv.Atoi(c.Params.ByName("id"))
	_userId, _ := c.Get("user_id")
	userId, _ := strconv.Atoi(_userId.(string))

	s.DB.Table(models.Episode{}.TableName()).Where("id = ?", episodeId).First(&episode)
	s.DB.Table(models.Listened{}.TableName()).Assign(&models.Listened{
		Viewed:    true,
		CreatedAt: time.Now(),
	}).Where(&models.Listened{
		UserId:    int64(userId),
		ItemId:    int64(episodeId),
		ChannelId: episode.ChannelId,
	}).FirstOrCreate(&models.Listened{})
}

func (s EpisodeService) Download(c *gin.Context) {
	var episode models.Episode
	episodeId, _ := strconv.Atoi(c.Params.ByName("id"))

	s.DB.Table(models.Episode{}.TableName()).Where("id = ?", episodeId).First(&episode)
	c.Redirect(http.StatusMovedPermanently, episode.SourceUrl)
}
