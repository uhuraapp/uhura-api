package models

import (
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/uhuraapp/uhura-api/cache"
	"github.com/uhuraapp/uhura-api/helpers"
)

type Episode struct {
	Id          int64
	ChannelId   int64
	Key         string `sql:"unique"`
	SourceUrl   string `sql:"not null;unique"`
	Title       string
	Description string
	PublishedAt time.Time `sql:"not null"`
	Duration    string
	Uri         string
	Type        string
	helpers.Uriable
}

func (e Episode) TableName() string {
	return "items"
}

func (e Episode) CountByChannel(database gorm.DB, channelId int64) int64 {
	var (
		key           = "c:e:" + strconv.Itoa(int(channelId))
		episodesCount int64
	)

	cachedEpisodes, err := cache.Get(key, episodesCount)
	if err == nil {
		episodesCount = cachedEpisodes.(int64)
	} else {
		database.Table(e.TableName()).Where("channel_id = ?", channelId).Count(&episodesCount)
		defer cache.Set(key, episodesCount)
	}

	return episodesCount
}
