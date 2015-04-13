package parser

import (
	charset "code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
	rss "github.com/jteeuwen/go-pkg-rss"
)

type Fetcher struct {
	channel  chan<- *Channel
	episodes chan<- []*Episode
	err      chan<- error
	url      string
}

func NewFetcher(url string, c chan<- *Channel, e chan<- []*Episode, err chan<- error) *Fetcher {
	return &Fetcher{c, e, err, url}
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
		f.finish(channel, f.createEpisodes(episodes))
	}
}

func (f *Fetcher) rerun(newURL string) {
	f.url = newURL
	f.run()
}

func (f *Fetcher) finish(channel *Channel, episodes []*Episode) {
	channel.feedURL = f.url

	f.channel <- channel
	f.episodes <- episodes
	close(f.channel)
	close(f.episodes)
}

func (f *Fetcher) createEpisodes(e []*rss.Item) (episodes []*Episode) {
	for k, _ := range e {
		episodes = append(episodes, &Episode{Feed: e[k]})
	}
	return

}

func (f *Fetcher) _c(feed *rss.Feed, channels []*rss.Channel) {}
