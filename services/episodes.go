package services

import (
	"net/http"
	"strconv"
	"time"

	"bitbucket.org/dukex/uhura-api/entities"
	"bitbucket.org/dukex/uhura-api/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type EpisodeService struct {
	DB gorm.DB
}

func NewEpisodesService(db gorm.DB) EpisodeService {
	return EpisodeService{DB: db}
}

func (s EpisodeService) GetPaged(c *gin.Context) {
	var episodes []*entities.Episode

	params := c.Request.URL.Query()

	s.DB.Table(models.Episode{}.TableName()).
		Where("channel_id = ?", params.Get("channel_id")).
		Order("published_at DESC").
		Limit(params.Get("per_page")).
		Find(&episodes)

	c.JSON(200, map[string]interface{}{"episodes": episodes})
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
