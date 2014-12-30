package gpodder

import (
	"github.com/franela/goreq"
)

type Api struct {
	endpoint string
}

type Channel struct {
	Website             string `json:"website"`
	MygpoLink           string `json:"mygpo_link"`
	Description         string `json:"description"`
	Subscribers         int    `json:"subscribers"`
	Title               string `json:"title"`
	Url                 string `json:"url"`
	SubscribersLastWeek int    `json:"subscribers_last_week"`
	LogoUrl             string `json:"logo_url"`
	ScaledLogoUrl       string `json:"scaled_logo_url"`
}

var API Api

func init() {
	API = Api{
		endpoint: "https://gpodder.net/",
	}
}

func Top(limit string) (channels []Channel, err error) {
	body, err := API.request("GET", "toplist/"+limit)

	if err != nil {
		return
	}

	body.FromJsonTo(&channels)

	return
}

func (api Api) request(method string, path string) (*goreq.Body, error) {
	resp, err := goreq.Request{
		Method: method,
		Uri:    api.getUrlFor(path),
	}.Do()

	return resp.Body, err
}

func (api Api) getUrlFor(path string) string {
	return api.endpoint + path + ".json"
}
