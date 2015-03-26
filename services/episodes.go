package services

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"bitbucket.org/dukex/uhura-api/entities"
	"bitbucket.org/dukex/uhura-api/helpers"
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
	var (
		episodes  []*entities.Episode
		ids       []int64
		listeneds []models.Listened
		userId    string
	)

	_userId, err := c.Get("user_id")
	if err == nil {
		userId = _userId.(string)
	}

	params := c.Request.URL.Query()

	s.
		DB.Table(models.Episode{}.TableName()).
		Where("channel_id = ?", params.Get("channel_id")).
		Where("published_at < ?", params.Get("since")).
		Order("published_at DESC").
		Limit(params.Get("per_page")).
		Find(&episodes).
		Pluck("id", &ids)

	if len(ids) > 0 {
		s.DB.Table(models.Listened{}.TableName()).
			Where("item_id IN (?)", ids).
			Where("viewed = true").
			Where("user_id = ?", userId).
			Find(&listeneds)
	} else {
		listeneds = make([]models.Listened, 0)
	}

	for _, episode := range episodes {
		episode.ChannelUri = params.Get("channel_id")
		listened, ok := helpers.Returns(listeneds, "ItemId", episode.Id).(models.Listened)
		if ok {
			episode.Listened = listened.Viewed
			episode.StoppedAt = listened.StoppedAt
		}
	}

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
		StoppedAt: 0,
	}).Where(&models.Listened{
		UserId:    int64(userId),
		ItemId:    int64(episodeId),
		ChannelId: episode.ChannelId,
	}).FirstOrCreate(&models.Listened{})

	go helpers.NewEvent(_userId.(string), "listened", map[string]interface{}{"episode_id": episode.Id, "channel_id": episode.ChannelId})
}

func (s EpisodeService) Listen(c *gin.Context) {
	var episode models.Episode
	episodeId, _ := strconv.Atoi(c.Params.ByName("id"))
	userId, _ := helpers.GetUser(c)

	c.Request.ParseForm()
	at, err := strconv.Atoi(c.Request.Form.Get("at"))
	if err != nil {
		c.Abort()
		return
	}

	err = s.DB.Table(models.Episode{}.TableName()).Where("id = ?", episodeId).First(&episode).Error
	if err != nil {
		c.Abort()
		return
	}

	s.DB.Table(models.Listened{}.TableName()).Assign(&models.Listened{
		UpdatedAt: time.Now(),
		StoppedAt: int64(at),
	}).Where(&models.Listened{
		UserId:    int64(userId),
		ItemId:    int64(episodeId),
		ChannelId: episode.ChannelId,
	}).FirstOrCreate(&models.Listened{})

	c.Data(204, "", []byte(""))
	return
}

func (s EpisodeService) Unlistened(c *gin.Context) {
	var episode models.Episode
	episodeId, _ := strconv.Atoi(c.Params.ByName("id"))
	_userId, _ := c.Get("user_id")
	userId, _ := strconv.Atoi(_userId.(string))

	s.DB.Table(models.Episode{}.TableName()).Where("id = ?", episodeId).First(&episode)
	s.DB.Table(models.Listened{}.TableName()).Where(&models.Listened{
		UserId:    int64(userId),
		ItemId:    int64(episodeId),
		ChannelId: episode.ChannelId,
	}).Delete(&models.Listened{})

	go helpers.NewEvent(_userId.(string), "unlistened", map[string]interface{}{"episode_id": episode.Id, "channel_id": episode.ChannelId})
}

func (s EpisodeService) Download(c *gin.Context) {
	var episode models.Episode
	var body []byte

	episodeId, _ := strconv.Atoi(c.Params.ByName("id"))

	s.DB.Table(models.Episode{}.TableName()).Where("id = ?", episodeId).First(&episode)

	if c.Request.Method == "HEAD" && episode.ContentLength > 0 {
	} else {
		response, err := http.Get(episode.SourceUrl)
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
		defer response.Body.Close()

		body, err = ioutil.ReadAll(response.Body)

		if err != nil {
			c.AbortWithStatus(500)
			return
		}

		episode.ContentLength = response.ContentLength
		episode.ContentType = response.Header.Get("Content-Type")

		go s.DB.Table(models.Episode{}.TableName()).Where("id = ?", episodeId).Update(map[string]interface{}{
			"content_length": episode.ContentLength,
			"content_type":   episode.ContentType,
		})
	}

	c.Writer.Header().Set("Content-Length", strconv.Itoa(int(episode.ContentLength)))
	c.Data(200, episode.ContentType, body)
}
