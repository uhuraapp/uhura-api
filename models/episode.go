package models

import (
	"time"

	"github.com/dukex/uhura/core/helper"
)

type Episode struct {
	Id          int64
	ChannelId   int64
	Key         string `sql:"unique"`
	SourceUrl   string `sql:"not null;unique"`
	Title       string
	Description string
	PublishedAt time.Time `sql:"not null"`
	Duration    string
	Uri         string
	Type        string
	helper.Uriable
}

func (e Episode) TableName() string {
	return "items"
}
