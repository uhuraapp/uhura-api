package services

import (
	//"io/ioutil"
	//"net/http"

	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/uhuraapp/uhura-api/channels"
	"github.com/uhuraapp/uhura-api/entities"
	"github.com/uhuraapp/uhura-api/models"
	// "github.com/uhuraapp/uhura-api/helpers"
)

type EpisodeService struct {
	DB *gorm.DB
}

func NewEpisodesService(db *gorm.DB) EpisodeService {
	return EpisodeService{DB: db}
}

func (s EpisodeService) Played(c *gin.Context) {
	episodeUID := c.Params.ByName("id")
	channelID, _ := strconv.Atoi(c.Params.ByName("channel_id"))
	_userId, _ := c.Get("user_id")
	userId, _ := strconv.Atoi(_userId.(string))

	s.DB.Table(models.Listened{}.TableName()).Assign(&models.Listened{
		Viewed:    true,
		CreatedAt: time.Now(),
		StoppedAt: 0,
	}).Where(&models.Listened{
		UserId:    int64(userId),
		ChannelId: int64(channelID),
		ItemUID:   episodeUID,
	}).FirstOrCreate(&models.Listened{})

	c.JSON(201, gin.H{})
}

func (s EpisodeService) UnPlayed(c *gin.Context) {
	episodeUID := c.Params.ByName("id")
	channelID, _ := strconv.Atoi(c.Params.ByName("channel_id"))
	_userId, _ := c.Get("user_id")
	userId, _ := strconv.Atoi(_userId.(string))

	s.DB.Table(models.Listened{}.TableName()).Where(&models.Listened{
		UserId:    int64(userId),
		ChannelId: int64(channelID),
		ItemUID:   episodeUID,
	}).Delete(&models.Listened{})

	c.AbortWithStatus(204)
}

func (s EpisodeService) Download(c *gin.Context) {
	var episode *entities.Episode
	var found bool

	channelURI := c.Params.ByName("channel_id")
	episodeID := c.Params.ByName("id")

	_, episodes, _, ok := channels.Find(s.DB, channelURI)

	if !ok {
		c.AbortWithStatus(500)
		return
	}

	for _, _episode := range episodes {
		if _episode.Id == episodeID {
			episode = _episode
			found = true
		}
	}

	if !found {
		c.AbortWithStatus(404)
		return
	}

	c.Redirect(http.StatusFound, episode.SourceUrl)
	return
}
