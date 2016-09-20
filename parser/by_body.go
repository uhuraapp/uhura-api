package parser

import (
	"bytes"
	"log"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
	rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/uhuraapp/uhura-api/helpers"
	"golang.org/x/net/html/charset"
)

func ByBody(body []byte, URL string) (channel *Channel, ok bool) {
	log.Println("---------x----------")
	log.Println(string(body))

	if isHTML(body) {
		URLs := findLinks(body)

		for _, _URL := range URLs {
			parsedURL, err := helpers.ParseURL(_URL)
			if checkErr(err) {
				channel, ok = ByURL(*parsedURL)
				if ok {
					return channel, ok
				}
			}
		}

		return
	}

	// else {
	c := make(chan *Channel)

	go rss.New(0, true, _c, episodeHandler(URL, body, c)).
		FetchBytes(URL, body, charset.NewReaderLabel)

	go func() {
		<-time.After(5 * time.Second)
		log.Println("-> timeout")
		close(c)
	}()

	channel = <-c
	log.Println("============== channel =============", c)
	return channel, channel != nil
}

func episodeHandler(URL string, body []byte, c chan<- *Channel) rss.ItemHandlerFunc {
	return func(feed *rss.Feed, rssChannel *rss.Channel, episodes []*rss.Item) {
		channel := &Channel{
			Feed: rssChannel,
			// 			URL:  URL,
			// 			Body: body,
		}

		// 		if channel.HasNewURL() && channel.NewURL() != URL {
		// 			// RunFetcher(channel.NewURL(), p.channel, p.err, make(chan []byte))
		// 			// return
		// 		}

		c <- buildChannel(buildEpisodes(channel, episodes))
	}
}

func buildEpisodes(channel *Channel, _ []*rss.Item) *Channel {
	// 	// log.Debug("Adding episodes to channel")

	// 	for k := range e {
	// 		c.Episodes = append(c.Episodes, &Episode{Feed: e[k]})
	// 	}
	return channel
}

func buildChannel(channel *Channel) *Channel {
	channel.Build()
	// 	episodes := make([]*Episode, 0)
	// 	for _, e := range c.Episodes {
	// 		if e.Build() {
	// 			episodes = append(episodes, e)
	// 		}
	// 	}
	// 	c.Episodes = episodes
	return channel
}

func _c(feed *rss.Feed, channels []*rss.Channel) {}

func findLinks(body []byte) []string {
	links := make([]string, 0)

	r := bytes.NewReader(body)

	doc, err := goquery.NewDocumentFromReader(r)
	if err == nil {
		doc.Find("link[rel=alternate][type='application/rss+xml']").Each(func(_ int, s *goquery.Selection) {
			href, ok := s.Attr("href")
			if ok {
				links = append(links, href)
			}
		})
	}

	return links
}

var hasHTML = regexp.MustCompile(`<\/?html>`)

func isHTML(body []byte) bool {
	return hasHTML.Match(body)
}
