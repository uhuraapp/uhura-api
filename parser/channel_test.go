package parser

import (
	"crypto/md5"
	"encoding/hex"
	"testing"

	rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/stretchr/testify/assert"
)

func TestChannelBuild(t *testing.T) {
	feed := &rss.Channel{
		Title:       "My channel",
		Description: "My Description",
		Copyright:   "© 2010",
		PubDate:     "2016 out",
		Image: rss.Image{
			Url: "http://image.url",
		},
		Extensions: map[string]map[string][]rss.Extension{
			"" + ITUNES_EXT: {
				"subtitle": {{Value: "Subtitle"}},
				"summary":  {{Value: "Summary"}},
				"author":   {{Value: "Author"}},
			},
		},
		Links: []rss.Link{
			{Href: "http://link1.url"},
			{Href: "http://link2.url"},
		},
	}

	c := Channel{
		Feed:         feed,
		requestedURL: "http://request.url",
	}

	c.Build()

	assert.Equal(t, c.Title, feed.Title)
	assert.Equal(t, c.Description, feed.Description)
	assert.Equal(t, c.Copyright, feed.Copyright)
	assert.Equal(t, c.PubDate, feed.PubDate)
	assert.Equal(t, c.Subtitle, "Subtitle")
	assert.Equal(t, c.Summary, "Summary")
	assert.Equal(t, c.Author, "Author")
	assert.Equal(t, c.Image, feed.Image.Url)
	assert.Equal(t, c.Links, []string{c.requestedURL, feed.Links[0].Href, feed.Links[1].Href})

	c.Feed.Extensions[ITUNES_EXT]["image"] = []rss.Extension{
		{Attrs: map[string]string{"href": "http://itunes-image.url"}},
	}

	c.Build()

	assert.Equal(t, c.Image, "http://itunes-image.url")

	h := md5.New()
	h.Write([]byte(c.Title))
	id := hex.EncodeToString(h.Sum(nil))

	assert.Equal(t, c.ID, id)
}
