package channels

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/uhuraapp/uhura-api/entities"
	"github.com/uhuraapp/uhura-api/helpers"
	"github.com/uhuraapp/uhura-api/models"
	"github.com/uhuraapp/uhura-api/parser"
)

func Find(database *gorm.DB, uidORurl string) (channel entities.Channel,
	episodes entities.Episodes,
	feed *parser.Channel,
	ok bool) {
	// By default uidORurl the feedURL, hat is the first guess is this arguments
	// is a URL
	var feedURL = uidORurl

	// Try to find uidORurl in database
	err := database.Table(models.Channel{}.TableName()).
		Where("uri = ?", prepareURI(uidORurl)).
		First(&channel).Error

	if err != gorm.ErrRecordNotFound {
		// If Find the feed URL is the channel URL
		feedURL = channel.Url
	}

	url, urlErr := helpers.ParseURL(feedURL)

	if urlErr != nil {
		return
	}

	// Find Channel with this URL
	if time.Now().Sub(channel.UpdatedAt) < (time.Hour*5) && len(channel.Body) > 5 {
		// Use the body is valid
		feed, err = parser.Body(channel.Body, channel.Url)
	} else {
		// Find xml with URL
		var body []byte
		feed, body, err = parser.URL(url)
		go Save(database, feed, body)
	}

	if err != nil {
		return channel, episodes, feed, false
	}

	channel = ChannelEntityFromFeed(feed)
	episodes, ids := EpisodesEntityFromFeed(feed)
	channel.Episodes = ids

	return channel, episodes, feed, true
}

func prepareURI(uri string) string {
	return strings.Replace(uri, "/", "", 1)
}
