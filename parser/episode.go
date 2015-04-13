package parser

import rss "github.com/jteeuwen/go-pkg-rss"

type Episode struct {
	Title      string
	Enclosures []*rss.Enclosure
	PubDate    string
	ID         string
	Subtitle   string
	Summary    string
	Duration   string

	iTunes
	Feed *rss.Item `json:"-"`
}

func (e *Episode) Build() {
	e.Title = e.Feed.Title
	e.Enclosures = e.Feed.Enclosures
	e.PubDate = e.Feed.PubDate
	e.ID = *e.Feed.Guid
	e.Subtitle = e.value(e, "subtitle")
	e.Summary = e.value(e, "summary")
	e.Duration = e.value(e, "duration")
}

func (e *Episode) GetExtensions(ext string) map[string][]rss.Extension {
	if e.Feed != nil && len(e.Feed.Extensions) > 0 {
		return e.Feed.Extensions[ext]
	}
	return nil
}
