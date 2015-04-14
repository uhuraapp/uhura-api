package parser

import (
	"net/url"

	_ "code.google.com/p/go-charset/data"
)

func URL(url *url.URL) (*Channel, error) {
	var _error error

	c := make(chan *Channel)
	err := make(chan error)

	log.Debug("creating fetcher")
	fetcher := NewFetcher(url.String(), c, err)
	go fetcher.run()

	// 	//if isFeed(body) {
	// 	//		xml := ParserXML(body)
	// 	// channel := FindOrCreateChannel(xml)
	// 	// UpdateChannel(channel, xml)
	// 	// CacheImage(channel)
	// 	// episodes := FindOrCreateEpisodes(channel, xml)
	// 	// GetDelayBetweenEpisodes(episodes)
	// 	// SetNewRun(channel)
	// 	//}
	// 	return nil
	// }

	go func() {
		_error = <-err
		log.Error("%s", _error)
		close(c)
		close(err)
	}()

	channel := <-c
	channel.requestedURL = url.String()
	log.Debug("finishing %s - %s", channel.Title, _error)

	return channel, _error
}

// func ProcessChannel(channel *rss.Channel) {
// 	CacheImage(channel)
// 	SaveOnDatabase(channel)
// }
//
// func CacheImage(channel *rss.Channel) {
// 	imageURL := GetChannelImageURL(channel)
// 	originalImage, err := RequestURL(imageURL)
//
// 	log.Println(imageURL)
// 	if err != nil {
// 		//	notifyWrongImage
// 	}
// 	image, err := imageproxy.Transform(originalImage, resizeOptions())
// 	if err != nil {
// 		log.Println(err)
// 	}
//
// 	str := base64.StdEncoding.EncodeToString(image)
// 	fmt.Println(str)
// }
//

// func resizeOptions() imageproxy.Options {
// 	return imageproxy.ParseOptions("250x")
// }
//
// func SaveOnDatabase(channel *rss.Channel) {
// }

////

// func GetChannelImageURL(channel *rss.Channel) string {
// 	var imageURL string
// 	// 	if itunesImage := channel.Extensions[ITUNES_EXT]["image"]; itunesImage != nil {
// 	// 		imageURL = itunesImage[0].Attrs["href"]
// 	// 	} else {
// 	imageURL = channel.Image.Url
// 	// 	}
// 	return imageURL
// }

////////

// func GetDelayBetweenEpisodes() {
//
// }
//
// func SetNewRun() {
// }

// If isFeed
//    RequestURL -> ParserXML
//    FindChannel or CreateChannel
//    UpdateChannel
//    CacheImage
//    FindEpisodes or CreateEpisodes
//      SetContentType
//      SetContentLengh
//    GetDelayBetweenEpisodes
//    SetNewRun
// else
//    FindChannelOnWebsite
