package parser

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/uhuraapp/uhura-api/helpers"
)

func ByURL(URL url.URL) (*Channel, bool) {
	body, ok := request(URL.String())

	if !ok {
		return &Channel{}, false
	}

	return ByBody(body, URL.String())
}

func checkURL(URL string) bool {
	_, err := helpers.ParseURL(URL)

	return checkErr(err)
}

func checkErr(err error) bool {
	return err == nil
}

func request(URL string) ([]byte, bool) {
	if !checkURL(URL) {
		return nil, false
	}

	response, err := http.Get(URL)
	defer response.Body.Close()

	if !checkErr(err) && response.StatusCode != 200 {
		return nil, false
	}

	body, err := ioutil.ReadAll(response.Body)

	return body, checkErr(err)
}
