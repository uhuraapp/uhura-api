package parser

import (
	charset "code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
	rss "github.com/jteeuwen/go-pkg-rss"
)

type Fetcher struct {
	channel chan<- *Channel
	err     chan<- error
	url     string
}

func NewFetcher(url string, c chan<- *Channel, err chan<- error) *Fetcher {
	return &Fetcher{c, err, url}
}

func (f *Fetcher) run() error {
	body, err := f.request()
	if err != nil {
		f.err <- err
	}

	rss.New(0, true, f._c, f.episodeHandler).
		FetchBytes(f.url, body, charset.NewReader)

	return nil
}

func (f *Fetcher) request() ([]byte, error) {
	return RequestURL(f.url)
}

func (f *Fetcher) episodeHandler(feed *rss.Feed, rssChannel *rss.Channel, episodes []*rss.Item) {
	channel := &Channel{Feed: rssChannel}
	if channel.HasNewURL() && channel.NewURL() != f.url {
		log.Debug("has new URL: %s != %s", channel.NewURL(), feed.Url)
		f.rerun(channel.NewURL())
	} else {
		f.appendEpisodes(channel, episodes)
		f.finish(channel)
	}
}

func (f *Fetcher) rerun(newURL string) {
	f.url = newURL
	f.run()
}

func (f *Fetcher) finish(channel *Channel) {
	channel.feedURL = f.url

	f.buildRecords(channel)

	f.channel <- channel
	close(f.channel)
}

func (f *Fetcher) appendEpisodes(c *Channel, e []*rss.Item) {
	for k, _ := range e {
		c.Episodes = append(c.Episodes, &Episode{Feed: e[k]})
	}
	return
}

func (f *Fetcher) buildRecords(c *Channel) {
	c.Build()
	for _, e := range c.Episodes {
		e.Build()
	}
}

func (f *Fetcher) _c(feed *rss.Feed, channels []*rss.Channel) {}
