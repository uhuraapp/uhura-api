package services

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"bitbucket.org/dukex/uhura-api/channels"
	"bitbucket.org/dukex/uhura-api/entities"
	"bitbucket.org/dukex/uhura-api/helpers"
	"bitbucket.org/dukex/uhura-api/models"
	"bitbucket.org/dukex/uhura-api/parser"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type ChannelsService struct {
	DB gorm.DB
}

func NewChannelsService(db gorm.DB) ChannelsService {
	return ChannelsService{DB: db}
}

func (s ChannelsService) getChannel(c *gin.Context) (channel entities.Channel, notFound bool, feedChannel *parser.Channel) {
	feedURL := c.Params.ByName("uri")
	uri := strings.Replace(feedURL, "/", "", 1)

	err := s.DB.Table(models.Channel{}.TableName()).Where("uri = ?", uri).First(&channel).Error

	if err == gorm.RecordNotFound {
		var url *url.URL
		url, err = helpers.ParseURL(feedURL)

		if err != nil {
			return channel, true, nil
		}

		log.Println(err, url)

		channelF, err := parser.URL(url)

		log.Println("getChannel err: ", err)
		if err != nil {
			return channel, true, nil
		}

		if uhuraID := FindUhuraID(s.DB, channelF); uhuraID != "" {
			err = s.DB.Table(models.Channel{}.TableName()).Where("uri = ?", uhuraID).First(&channel).Error
			if err == nil {
				log.Println("ERRR", err)

				return channel, false, nil
			}
		} else {
			// go channels.Create(s.DB, url.String())
		}

		channel = channels.TranslateFromFeedToEntity(channel, feedChannel)
	}

	return channel, err != nil, feedChannel
}

func (s ChannelsService) Get(c *gin.Context) {
	var episodes []*entities.Episode
	channel, notFound, feedChannel := s.getChannel(c)

	if notFound {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if feedChannel != nil {
		var ids []int64
		episodes, ids = channels.TranslateEpisodesFromFeedToEntity(feedChannel)
		channel.Episodes = ids
	} else {
		if helpers.CacheHeader(c, channel.UpdatedAt) {
			c.AbortWithStatus(304)
			return
		}

		userId, _ := helpers.GetUser(c)
		channel.Episodes, episodes = s.getEpisodes(channel.Id, channel.Uri, userId)

		if userId != 0 {
			channel.Subscribed = s.DB.Table(models.Subscription{}.TableName()).Where("user_id = ?", userId).
				Where("channel_id = ?", channel.Id).
				Find(&models.Subscription{}).Error != gorm.RecordNotFound
		}
	}

	c.JSON(200, gin.H{"channel": channel, "episodes": episodes})
}

func (s ChannelsService) Open(c *gin.Context) {
	var channel entities.Channel
	channelURI := c.Params.ByName("uri")
	s.DB.Table(models.Channel{}.TableName()).Update(&models.Channel{
		VisitedAt: time.Now().UTC(),
	}).Where("uri = ?", channelURI).First(&channel)

	c.JSON(200, gin.H{})
}

func (s ChannelsService) getEpisodes(channelID int64, channelUri string, userId int) (ids []int64, episodes []*entities.Episode) {
	s.DB.Table(models.Episode{}.TableName()).
		Where("channel_id = ?", channelID).
		Order("published_at DESC").
		Limit(20).
		Find(&episodes)

	for _, e := range episodes {
		ids = append(ids, e.Id)
	}

	models.SetListenAttributesToEpisode(s.DB, userId, episodes, channelUri)

	return
}
