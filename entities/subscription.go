package entities

import (
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/uhuraapp/uhura-api/cache"
	"github.com/uhuraapp/uhura-api/models"
)

type Subscription struct {
	Id       int64  `json:"raw_id"`
	Uri      string `json:"id"`
	Title    string `json:"title"`
	ImageUrl string `json:"image_url"`
	ToView   int64  `json:"to_view"`
}

func (s *Subscription) GetToView(database gorm.DB, userId string) int64 {
	var (
		listenedCount int64
	)

	database.Table(models.Listened{}.TableName()).
		Where("channel_id = ? AND user_id = ?", s.Id, userId).
		Count(&listenedCount)

	return getEpisodesCount(database, s.Id) - listenedCount
}

func getEpisodesCount(database gorm.DB, channelId int64) int64 {
	var (
		key           = "c:e:" + strconv.Itoa(int(channelId))
		episodesCount int64
	)

	cachedEpisodes, err := cache.Get(key, episodesCount)
	if err == nil {
		episodesCount = cachedEpisodes.(int64)
	} else {
		database.Table(models.Episode{}.TableName()).Where("channel_id = ?", channelId).Count(&episodesCount)
		defer cache.Set(key, episodesCount)
	}
	return episodesCount
}
