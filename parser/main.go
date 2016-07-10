package parser

import (
  "net/url"
)

func URL(url *url.URL) (channels *Channel, _error error) {
  log.Debug("creating fetcher")

  c := make(chan *Channel)
  err := make(chan error)

  go RunFetcher(url.String(), c, err)

  return <-c, <-err
}
