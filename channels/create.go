package channels

import (
	"time"

	"bitbucket.org/dukex/uhura-api/entities"
	"bitbucket.org/dukex/uhura-api/helpers"
	"bitbucket.org/dukex/uhura-api/models"
	"bitbucket.org/dukex/uhura-api/parser"
	"github.com/jinzhu/gorm"
)

func Create(database gorm.DB, url string) (*entities.Channel, bool) {
	channelURL, _ := helpers.ParseURL(url)
	channels, errors := parser.URL(channelURL)

	if len(errors) > 0 {
		log.Debug("errors: %s", errors)
		return nil, false
	}

	if len(channels) < 1 {
		log.Debug("error: no channel found")
		return nil, false
	}

	var ok bool
	var channel entities.Channel

	log.Debug("no error found")
	channelF := channels[0] // TODO: fix it

	log.Debug("channel UhuraID: %s", channelF.UhuraId)
	if channelF.UhuraId != "" {
		ok = database.Table(models.Channel{}.TableName()).Where("uri = ?", channelF.UhuraId).
			First(&channel).Error != gorm.RecordNotFound
	} else {
		model := translateChannel(channelF)
		log.Debug("channel: %s", model)

		ok = database.Table(models.Channel{}.TableName()).Save(&model).Error == nil
		log.Debug("is ok: %s", ok)

		if ok {
			database.Table(models.Channel{}.TableName()).First(&channel, model.Id)
			go createLinks(channelF.Links, channel, database)
		}
	}

	return &channel, ok
}

func translateChannel(channel *parser.Channel) *models.Channel {
		model := models.Channel{}
		model.Title = channel.Title
		model.Description = channel.Description
		model.Copyright = channel.Copyright
		model.ImageUrl = channel.Image
		model.Uri = helpers.MakeUri(channel.Title)
		model.Language = channel.Language
		model.UpdatedAt = time.Now()
		model.CreatedAt = time.Now()
		model.LastBuildDate = channel.LastBuildDate
		model.Url = channel.Links[0]
		return &model
}

func createLinks(links []string, channel entities.Channel, database gorm.DB) {
	for i := 0; i < len(links); i++ {
		u := models.ChannelURL{}
		database.Table(models.ChannelURL{}.TableName()).
			FirstOrCreate(&u, models.ChannelURL{
				ChannelId: channel.Id,
				Url:       links[i],
			})
	}
}
