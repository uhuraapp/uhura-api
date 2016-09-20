package model

import (
	self "github.com/uhuraapp/uhura-api/model/channel"

	"github.com/jinzhu/gorm"
	"github.com/uhuraapp/uhura-api/serializer"
)

type channel struct {
	connection *gorm.DB
}

func Channel(connection *gorm.DB) *channel {
	return &channel{
		connection,
	}
}

func (c *channel) TableName() string {
	return "channels"
}

func (c *channel) Top() []*serializer.Channel {
	var channels []*serializer.Channel

	c.connection.Table("user_channels").
		Select("channels.title, channels.image_url, channels.description, channels.uri, channels.id, COUNT(*) AS subscribers_count").
		Joins("INNER JOIN channels ON user_channels.channel_id = channels.id").
		Group("1,2,3,4,5 ORDER BY subscribers_count DESC").
		Limit(5).
		Find(&channels)

	return channels
}

func (c *channel) Find(id string) {
	return self.Find(c.connection, id)
}
