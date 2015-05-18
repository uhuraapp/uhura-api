package core

import (
	"strings"
	"time"
)

const (
	itunesExt = "http://www.itunes.com/dtds/podcast-1.0.dtd"
)

var noItemFn = func(feed *rss.Feed, ch *rss.Channel, items []*rss.Item) {}

func FetchChannel(url string) {
	go fetcher(url, 5)
}

func channelFetchHandler(feed *rss.Feed, channels []*rss.Channel) {
	for _, channelData := range channels {
		var channel Channel
		database.Save(&channel)

		if itunesCategory := channelData.Extensions[itunesExt]["category"]; itunesCategory != nil {
			for _, category := range itunesCategory {
				var categoryDB Category
				database.Where(&Category{Name: category.Attrs["text"]}).FirstOrCreate(&categoryDB)
				database.Where(&ChannelCategories{ChannelId: int64(channel.Id), CategoryId: categoryDB.Id}).FirstOrCreate(&ChannelCategories{})
			}
		}
	}
}
