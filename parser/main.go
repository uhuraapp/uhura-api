package parser

import (
	"net/url"

	_ "code.google.com/p/go-charset/data"
)

func URL(url *url.URL) (channels *Channel, _error error) {
	log.Debug("creating fetcher")

	c := make(chan *Channel)
	err := make(chan error)

	fetcher := NewFetcher([]string{url.String()}, c, err)

	go fetcher.run()

	log.Debug("channels: %s", channels)
	log.Debug("_error: %s", _error)

	return <-c, <-err
}
