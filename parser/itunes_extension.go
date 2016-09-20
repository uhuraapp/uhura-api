package parser

import (
	rss "github.com/jteeuwen/go-pkg-rss"
)

const ITUNES_EXT = "http://www.itunes.com/dtds/podcast-1.0.dtd"

type iTunesExtension struct{}

type iTunesExtensiable interface {
	GetExtensions(string) map[string][]rss.Extension
}

func (c iTunesExtension) ivalue(f iTunesExtensiable, k string) string {
	if i := c._get(f, k); len(i) > 0 {
		return i[0].Value
	}
	return ""
}

func (c iTunesExtension) iattr(f iTunesExtensiable, k, attr string) string {
	attrs := c._attrs(f, k, attr)
	if len(attrs) > 0 {
		return attrs[0]
	}
	return ""
}

func (c iTunesExtension) _attrs(f iTunesExtensiable, k, attr string) []string {
	items := make([]string, 0)
	if i := c._get(f, k); len(i) > 0 {
		items = append(items, i[0].Attrs[attr])
	}
	return items
}

func (c iTunesExtension) _get(f iTunesExtensiable, k string) []rss.Extension {
	if f != nil &&
		f.GetExtensions(ITUNES_EXT) != nil &&
		f.GetExtensions(ITUNES_EXT)[k] != nil &&
		len(f.GetExtensions(ITUNES_EXT)[k]) > 0 {
		return f.GetExtensions(ITUNES_EXT)[k]
	}
	return nil
}
