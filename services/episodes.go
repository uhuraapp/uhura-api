package services

import (
	//"io/ioutil"
	//"net/http"
	// "strconv"
	//"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	// "github.com/uhuraapp/uhura-api/entities"
	// "github.com/uhuraapp/uhura-api/helpers"
	//"github.com/uhuraapp/uhura-api/models"
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

func (s EpisodeService) GetPaged(c *gin.Context) {
	// var episodes entities.Episodes
	// var channelURI []string

	// userId, _ := helpers.GetUser(c)
	// params := c.Request.URL.Query()
	// perPage, _ := strconv.Atoi(params.Get("per_page"))

	// s.
	// 	DB.Table(models.Episode{}.TableName()).
	// 	Where("channel_id = ?", params.Get("channel_id")).
	// 	Where("published_at < ?", params.Get("since")).
	// 	Order("published_at DESC").
	// 	Limit(perPage).
	// 	Find(&episodes)

	// s.
	// 	DB.Table(models.Channel{}.TableName()).
	// 	Where("id = ?", params.Get("channel_id")).
	// 	Pluck("uri", &channelURI)

	// models.SetListenAttributesToEpisode(s.DB, userId, episodes, channelURI[0])

	// c.JSON(200, map[string]interface{}{"episodes": episodes})
}

func (s EpisodeService) Played(c *gin.Context) {
	// var episode models.Episode
	// episodeId, _ := strconv.Atoi(c.Params.ByName("id"))
	// _userId, _ := c.Get("user_id")
	// userId, _ := strconv.Atoi(_userId.(string))

	// s.DB.Table(models.Episode{}.TableName()).Where("id = ?", episodeId).First(&episode)
	// s.DB.Table(models.Listened{}.TableName()).Assign(&models.Listened{
	// 	Viewed:    true,
	// 	CreatedAt: time.Now(),
	// 	StoppedAt: 0,
	// }).Where(&models.Listened{
	// 	UserId:    int64(userId),
	// 	ItemId:    int64(episodeId),
	// 	ChannelId: episode.ChannelId,
	// }).FirstOrCreate(&models.Listened{})

	c.JSON(201, gin.H{})
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
	// 	UserId:    int64(userId),
	// 	ItemId:    int64(episodeId),
	// 	ChannelId: episode.ChannelId,
	// }).FirstOrCreate(&models.Listened{})

	// c.JSON(201, gin.H{})
	return
}

func (s EpisodeService) Unlistened(c *gin.Context) {
	// var episode models.Episode
	// episodeId, _ := strconv.Atoi(c.Params.ByName("id"))
	// _userId, _ := c.Get("user_id")
	// userId, _ := strconv.Atoi(_userId.(string))

	// s.DB.Table(models.Episode{}.TableName()).Where("id = ?", episodeId).First(&episode)
	// s.DB.Table(models.Listened{}.TableName()).Where(&models.Listened{
	// 	UserId:    int64(userId),
	// 	ItemId:    int64(episodeId),
	// 	ChannelId: episode.ChannelId,
	// }).Delete(&models.Listened{})

	// c.AbortWithStatus(204)
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
