package entities

import "time"

type Channel struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	ImageUrl    string      `json:"image_url"`
	Url         string      `json:"url"`
	Id          int64       `json:"raw_id"`
	Uri         string      `json:"id"`
	ToView      int64       `json:"to_view"`
	Subscribed  interface{} `json:"subscribed"`
	Copyright   string      `json:"copyright"`
	Episodes    []int64     `json:"episodes"`
	UpdatedAt   time.Time   `json:"updated_at"`
}
