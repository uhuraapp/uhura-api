package services

import (
	//"io/ioutil"
	//"net/http"

	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/uhuraapp/uhura-api/channels"
	"github.com/uhuraapp/uhura-api/entities"
	"github.com/uhuraapp/uhura-api/helpers"
	"github.com/uhuraapp/uhura-api/models"
)

type EpisodeService struct {
	DB *gorm.DB
}

func NewEpisodesService(db *gorm.DB) EpisodeService {
	return EpisodeService{DB: db}
}

func (s EpisodeService) episodes(channelURI string) (entities.Episodes, int, bool) {
	if channelURI == "" {
		return entities.Episodes{}, 404, false
	}

	_, episodes, _, ok := channels.Find(s.DB, channelURI)

	if !ok {
		return entities.Episodes{}, 500, false
	}

	return episodes, 200, true
}

func (s EpisodeService) episode(channelURI, id string) (*entities.Episode, int, bool) {
	episodes, status, ok := s.episodes(channelURI)

	if !ok {
		return &entities.Episode{}, status, ok
	}

	episode := episodes.Find(id)

	// userId, _ := helpers.GetUser(c)
	// models.SetListenAttributesToEpisode(s.DB, userId, entities.Episodes{&episode}, channelURI[0])

	return episode, 200, true
}

func (s EpisodeService) Get(c *gin.Context) {
	params := c.Request.URL.Query()

	episode, status, ok := s.episode(params.Get("channel_id"), c.Params.ByName("id"))

	if !ok {
		c.AbortWithStatus(status)
		return
	}

	c.JSON(200, gin.H{"episode": episode})
}

func (s EpisodeService) Index(c *gin.Context) {
	params := c.Request.URL.Query()

	episodes, status, ok := s.episodes(params.Get("channel_id"))

	if !ok {
		c.AbortWithStatus(status)
		return
	}

	c.JSON(200, gin.H{"episodes": episodes})
}

func (s EpisodeService) Played(c *gin.Context) {
	channelID, _ := strconv.Atoi(c.Params.ByName("channel_id"))
	episodeUID := c.Params.ByName("id")
	userID, _ := helpers.GetUser(c)

	s.DB.Table(models.Listened{}.TableName()).Assign(&models.Listened{
		Viewed:    true,
		CreatedAt: time.Now(),
		StoppedAt: 0,
	}).Where(&models.Listened{
		UserId:    int64(userID),
		ChannelId: int64(channelID),
		ItemUID:   episodeUID,
	}).FirstOrCreate(&models.Listened{})

	c.JSON(201, gin.H{})
}

func (s EpisodeService) UnPlayed(c *gin.Context) {
	channelID, _ := strconv.Atoi(c.Params.ByName("channel_id"))
	episodeUID := c.Params.ByName("id")
	userID, _ := helpers.GetUser(c)

	s.DB.Table(models.Listened{}.TableName()).Where(&models.Listened{
		UserId:    int64(userID),
		ChannelId: int64(channelID),
		ItemUID:   episodeUID,
	}).Delete(&models.Listened{})

	c.AbortWithStatus(204)
}

func (s EpisodeService) Listen(c *gin.Context) {
	channelID, _ := strconv.Atoi(c.Params.ByName("channel_id"))
	episodeUID := c.Params.ByName("id")
	userID, _ := helpers.GetUser(c)

	var params struct {
		At int `json:"at"`
	}

	if c.BindJSON(&params) != nil {
		c.AbortWithStatus(500)
		return
	}

	at := params.At

	s.DB.Table(models.Listened{}.TableName()).Assign(&models.Listened{
		UpdatedAt: time.Now(),
		StoppedAt: int64(at),
	}).Where(&models.Listened{
		UserId:    int64(userID),
		ChannelId: int64(channelID),
		ItemUID:   episodeUID,
	}).FirstOrCreate(&models.Listened{})

	c.JSON(201, gin.H{})
	return
}

func (s EpisodeService) Download(c *gin.Context) {
	// var episode models.Episode
	// var body []byte

	// episodeId, _ := strconv.Atoi(c.Params.ByName("id"))

	// s.DB.Table(models.Episode{}.TableName()).Where("id = ?", episodeId).First(&episode)

	// if c.Request.Method == "HEAD" && episode.ContentLength > 0 {
	// } else {
	// 	response, err := http.Get(episode.SourceUrl)
	// 	if err != nil {
	// 		c.AbortWithStatus(500)
	// 		return
	// 	}
	// 	defer response.Body.Close()

	// 	body, err = ioutil.ReadAll(response.Body)

	// 	if err != nil {
	// 		c.AbortWithStatus(500)
	// 		return
	// 	}

	// 	episode.ContentLength = response.ContentLength
	// 	episode.ContentType = response.Header.Get("Content-Type")

	// 	go s.DB.Table(models.Episode{}.TableName()).Where("id = ?", episodeId).Update(map[string]interface{}{
	// 		"content_length": episode.ContentLength,
	// 		"content_type":   episode.ContentType,
	// 	})
	// }

	// c.Writer.Header().Set("Content-Length", strconv.Itoa(int(episode.ContentLength)))
	// c.Data(200, episode.ContentType, body)
}
