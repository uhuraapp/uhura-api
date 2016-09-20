package channels

import (
	"log"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/uhuraapp/uhura-api/helpers"
	"github.com/uhuraapp/uhura-api/parser"
)

func Find(connection *gorm.DB, idOrUrl string) bool { //(channel serializer.Channel,
	// 	episodes entities.Episodes,
	// 	feed *parser.Channel,
	// 	ok bool) {

	URL, urlErr := helpers.ParseURL(feedURL(connection, idOrUrl))

	if urlErr != nil {
		return false
	}

	// if channelIsOld(channel) {
	feed, ok := parser.ByURL(URL)

	log.Println(feed)
	log.Println(ok)
	// 	} else {
	// 		var body []byte
	// 		feed, body, err = parser.URL(url)
	// 		go Save(database, feed, body)
	// 	}

	// 	if err != nil {
	// 		errr(err, "Find finished")
	// 		return channel, episodes, feed, false
	// 	}

	// 	info("+ translate feed to channel")
	// 	channelFeed := ChannelEntityFromFeed(feed)
	// 	episodes, ids := EpisodesEntityFromFeed(feed)
	// 	channelFeed.Episodes = ids

	// 	channel.Uri = channelFeed.Uri
	// 	database.Table(models.Channel{}.TableName()).
	// 		Where("id = ?", channel.Id).
	// 		UpdateColumns(models.Channel{
	// 			Title: channelFeed.Title,
	// 			Uri:   channelFeed.Uri,
	// 		})

	// 	channelFeed.Id = channel.Id

	// 	return channelFeed, episodes, feed, true
}

func prepareURI(uri string) string {
	return strings.Replace(uri, "/", "", 1)
}

// func channelIsOld(channel string) bool {
// 	return time.Now().Sub(channel.UpdatedAt) < (time.Hour*5) && len(channel.Body) > 5
// }

func feedURL(connection *gorm.DB, idOrUrl string) string {
	var feedURL = idOrUrl
	var channelURL string

	err := connection.Table("channels").
		Where("uri = ?", prepareURI(idOrUrl)).
		Pluck("url", &channelURL).Error

	if err != gorm.ErrRecordNotFound {
		feedURL = channelURL
	}

	return feedURL
}
