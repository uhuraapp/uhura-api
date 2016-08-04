package channels

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uhuraapp/uhura-api/helpers"
	"github.com/uhuraapp/uhura-api/parser"
)

var (
	channel = parser.Channel{
		Title:         "Example Test",
		Description:   "The long Description",
		Copyright:     "(c) 2011",
		Image:         "http://image.com/my",
		Language:      "en",
		LastBuildDate: "2016",
		URL:           "http://channel.com/my",
		Body:          []byte{},
		Episodes: []*parser.Episode{
			&parser.Episode{
				Title:       "First Episode",
				Description: "Description First Episode",
				PublishedAt: time.Date(2009, time.January, 1, 1, 0, 0, 0, time.UTC),
				Source:      "http://source.com/my1",
			},
			&parser.Episode{
				Title:       "Last Episode",
				Description: "Description Last Episode",
				PublishedAt: time.Now(),
				Source:      "http://source.com/my2",
			},
		},
	}
)

func TestChannelEntityFromFeed(t *testing.T) {
	entity := ChannelEntityFromFeed(&channel)

	assert.Equal(t, entity.Title, channel.Title)
	assert.Equal(t, entity.Description, channel.Description)
	assert.Equal(t, entity.Copyright, channel.Copyright)
	assert.Equal(t, entity.ImageUrl, channel.Image)
	assert.Equal(t, entity.Uri, helpers.MakeUri(channel.Title))
	assert.Equal(t, entity.Enabled, true)
}

func TestChannelModelFromFeed(t *testing.T) {
	model := ChannelModelFromFeed(&channel)
	assert.Equal(t, model.Title, channel.Title)
	assert.Equal(t, model.Description, channel.Description)
	assert.Equal(t, model.Copyright, channel.Copyright)
	assert.Equal(t, model.ImageUrl, channel.Image)
	assert.Equal(t, model.Uri, helpers.MakeUri(channel.Title))
	assert.Equal(t, model.Language, channel.Language)
	assert.Equal(t, model.LastBuildDate, channel.LastBuildDate)
	assert.Equal(t, model.Url, channel.URL)
	assert.Equal(t, model.Body, channel.Body)
}

func TestEpisodesEntityFromFeed(t *testing.T) {
	episodes, ids := EpisodesEntityFromFeed(&channel)

	assert.Equal(t, ids, []string{"lastepisode", "firstepisode"}, "returns ordered ids")

	for i, episode := range episodes {
		index := map[int]int{0: 1, 1: 0}[i]
		channelEpisodes := channel.Episodes[index]
		z := int64(0)

		assert.Equal(t, episode.Id, helpers.MakeUri(channelEpisodes.Title))
		assert.Equal(t, episode.Title, channelEpisodes.Title)
		assert.Equal(t, episode.Description, channelEpisodes.Description)
		assert.Equal(t, episode.PublishedAt, channelEpisodes.PublishedAt)
		assert.Equal(t, episode.SourceUrl, channelEpisodes.Source)
		assert.Equal(t, episode.StoppedAt, &z)
	}
}
