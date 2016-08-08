package services

import (
	//"io/ioutil"
	//"net/http"

	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/uhuraapp/uhura-api/channels"
	"github.com/uhuraapp/uhura-api/models"
	// "github.com/uhuraapp/uhura-api/entities"
	// "github.com/uhuraapp/uhura-api/helpers"
)

type EpisodeService struct {
	DB *gorm.DB
}

func NewEpisodesService(db *gorm.DB) EpisodeService {
	return EpisodeService{DB: db}
}

func (s EpisodeService) Get(c *gin.Context) {
	// var episode entities.Episode
	// var channelURI []string
	// episodeId, _ := strconv.Atoi(c.Params.ByName("id"))
	// userId, _ := helpers.GetUser(c)

	// err := s.DB.Table(models.Episode{}.TableName()).Where("id = ?", episodeId).First(&episode).Error
	// if err != nil {
	// 	c.AbortWithStatus(404)
	// 	return
	// }

	// s.
	// 	DB.Table(models.Channel{}.TableName()).
	// 	Where("id = ?", episode.ChannelId).
	// 	Pluck("uri", &channelURI)

	// models.SetListenAttributesToEpisode(s.DB, userId, entities.Episodes{&episode}, channelURI[0])

	// c.JSON(200, map[string]interface{}{"episode": episode})
}

func (s EpisodeService) Index(c *gin.Context) {
	params := c.Request.URL.Query()
	channelURI := params.Get("channel_id")

	if channelURI == "" {
		c.AbortWithStatus(404)
		return
	}

	// userId, _ := helpers.GetUser(c)

	_, episodes, _, ok := channels.Find(s.DB, channelURI)

	if !ok {
		c.AbortWithStatus(500)
	}

	// models.SetListenAttributesToEpisode(s.DB, userId, episodes, channelURI)

	c.JSON(200, map[string]interface{}{"episodes": episodes})
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

func (s EpisodeService) Listen(c *gin.Context) {
	// var episode models.Episode
	// episodeId, _ := strconv.Atoi(c.Params.ByName("id"))
	// userId, _ := helpers.GetUser(c)

	// var params struct {
	// 	At int `json:"at"`
	// }

	// if c.BindJSON(&params) != nil {
	// 	c.AbortWithStatus(500)
	// 	return
	// }

	// at := params.At

	// err := s.DB.Table(models.Episode{}.TableName()).Where("id = ?", episodeId).First(&episode).Error
	// if err != nil {
	// 	c.AbortWithStatus(404)
	// 	return
	// }

	// s.DB.Table(models.Listened{}.TableName()).Assign(&models.Listened{
	// 	UpdatedAt: time.Now(),
	// 	StoppedAt: int64(at),
	// }).Where(&models.Listened{
	// UserId:    int64(userId),
	// ChannelId: int64(channelID),
	// ItemUID:   episodeUID,
	// }).FirstOrCreate(&models.Listened{})

	// c.JSON(201, gin.H{})
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
