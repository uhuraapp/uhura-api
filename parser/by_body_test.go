package parser

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var feedMock = []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?><rss version=\"2.0\"><channel><title>Test Channel</title> <link>http://test.com</link> <description></description> <item> <title>Item 1</title> <description>Desc 1</description> <link>http://link.url</link> <pubDate>Tue, 16 Aug 2016 08:11:37 -0400</pubDate> </item> </channel> </rss>")

func FeedServer(response []byte) *httptest.Server {
	if response == nil {
		response = feedMock
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/xml")
		fmt.Fprintln(w, feedMock)
	}))

	return server
}

func TestByBodyTimeout(t *testing.T) {
	data := []byte("<x></x>")
	channel, ok := ByBody(data, "http://test.com")

	assert.False(t, ok)
	assert.Nil(t, channel)
}

func TestByBody(t *testing.T) {
	channel, ok := ByBody(feedMock, "http://test.com")

	assert.True(t, ok)
	assert.NotNil(t, channel)
	assert.Equal(t, channel.Title, "Test Channel")
}

func TestByBodyHTML(t *testing.T) {
	s := FeedServer(nil)
	defer s.Close()

	channel, ok := ByBody([]byte("<html><head><link rel='alternate' type='application/rss+xml' href="+s.URL+"></link></head></html>"), "http://test.com")

	assert.True(t, ok)
	assert.NotNil(t, channel)
	assert.Equal(t, channel.Title, "Test Channel")
}
