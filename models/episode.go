package models

import (
	"strconv"
	"time"

	"bitbucket.org/dukex/uhura-api/cache"
	"bitbucket.org/dukex/uhura-api/helpers"
	"github.com/jinzhu/gorm"
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
		var ok bool
		episodesCount, ok = cachedEpisodes.(int64)
		if !ok {
			episodesCount = int64(episodesCount)
		}
	} else {
		database.Table(e.TableName()).Where("channel_id = ?", channelId).Count(&episodesCount)
		defer cache.Set(key, episodesCount)
	}

	return episodesCount
}
