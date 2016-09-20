package parser

import (
	"crypto/md5"
	"encoding/hex"

	rss "github.com/jteeuwen/go-pkg-rss"
)

// Channel parsed
type Channel struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Copyright   string   `json:"copyright"`
	PubDate     string   `json:"pub_date"`
	Subtitle    string   `json:"subtitle"`
	ID          string   `json:"id"`
	Summary     string   `json:"summary"`
	Author      string   `json:"author"`
	Image       string   `json:"image_url"`
	Links       []string `json:"links"`
	// Categories    []*Category `json:"categories"`
	// Episodes []*Episode `json:"episodes"`
	// UhuraID       string      `json:"uhura_id"`
	// LastBuildDate string      `json:"last_build_date"`
	// URI           string      `json:"uri"`

	requestedURL string
	// URL string

	// Body []byte       `json:"-"`
	Feed *rss.Channel `json:"-"`
	iTunesExtension
}

// HasNewURL check if channel has a new URL
// TODO: test
// func (c *Channel) HasNewURL() bool {
// 	return c.NewURL() != ""
// }

// NewURL get feed new URL
// TODO: test
// func (c *Channel) NewURL() string {
// 	return c.value(c, "new-feed-url")
// }

// Build channel from parsed rss
func (c *Channel) Build() {
	c.Title = c.Feed.Title
	c.Description = c.Feed.Description
	c.Copyright = c.Feed.Copyright
	c.PubDate = c.Feed.PubDate
	c.ID = c.GenerateID()
	c.Subtitle = c.ivalue(c, "subtitle")
	c.Summary = c.ivalue(c, "summary")
	c.Author = c.ivalue(c, "author")
	c.Image = c.FixImage()
	c.Links = c.GetLinks()
	// 	c.LastBuildDate = c.Feed.LastBuildDate
	// 	c.URI = helpers.MakeUri(c.Feed.Title)

	// 	categories := c.attrs(c, "category", "text")
	// 	for _, category := range categories {
	// 		c.Categories = append(c.Categories, &Category{
	// 			Name: category,
	// 		})
	// 	}

	// 	// log.Debug("%s", c.Feed.Links)
	// 	// log.Debug("channel build finished: %s", c.Title)
}

// FixImage get correct image
func (c *Channel) FixImage() string {
	if image := c.iattr(c, "image", "href"); image != "" {
		return image
	}

	return c.Feed.Image.Url
}

// // GetExtensions from rss
func (c *Channel) GetExtensions(ext string) map[string][]rss.Extension {
	if c.Feed != nil && len(c.Feed.Extensions) > 0 {
		return c.Feed.Extensions[ext]
	}
	return nil
}

// GetLinks from urls requested and feed data
func (c *Channel) GetLinks() []string {
	links := make([]string, 0)
	links = append(links, c.requestedURL)

	for _, link := range c.Feed.Links {
		links = append(links, link.Href)
	}

	return links
}

// GenerateID md5 from title
func (c *Channel) GenerateID() string {
	h := md5.New()
	h.Write([]byte(c.Title))
	return hex.EncodeToString(h.Sum(nil))
}
