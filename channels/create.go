package channels

import (
	"sort"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/uhuraapp/uhura-api/entities"
	"github.com/uhuraapp/uhura-api/helpers"
	"github.com/uhuraapp/uhura-api/models"
	"github.com/uhuraapp/uhura-api/parser"
)

func Create(database *gorm.DB, url string) (*entities.Channel, bool) {
	channelURL, _ := helpers.ParseURL(url)
	channelF, err := parser.URL(channelURL)

	if err != nil {
		log.Debug("error: %s", err)
		return nil, false
	}

	var ok bool
	var channel entities.Channel

	log.Debug("no error found")

	log.Debug("channel UhuraID: %s", channelF.UhuraID)
	if channelF.UhuraID != "" {
		ok = database.Table(models.Channel{}.TableName()).Where("uri = ?", channelF.UhuraID).
			First(&channel).Error != gorm.ErrRecordNotFound
	} else {
		model := translateChannel(channelF)
		log.Debug("channel: %s", model)

		ok = database.Table(models.Channel{}.TableName()).Save(&model).Error == nil
		log.Debug("is ok: %s", ok)

		if ok {
			database.Table(models.Channel{}.TableName()).First(&channel, model.Id)
			go CreateLinks(channelF.Links, channel.Id, database)
			// TODO: add to sync queue
		}
	}

	return &channel, ok
}

func TranslateFromFeedToEntity(entity entities.Channel, channel *parser.Channel) entities.Channel {
	entity.Title = channel.Title
	entity.Description = channel.Description
	entity.Copyright = channel.Copyright
	entity.ImageUrl = channel.Image
	entity.Uri = helpers.MakeUri(channel.Title)
	entity.UpdatedAt = time.Now()
	entity.Enabled = true
	return entity
}

func TranslateEpisodesFromFeedToEntity(channel *parser.Channel) (entities.Episodes, []string) {
	episodes := make(entities.Episodes, 0)
	ids := make([]string, 0)

	for _, episode := range channel.Episodes {
		s := int64(0)
		id := helpers.MakeUri(episode.Title)
		episodes = append(episodes, &entities.Episode{
			Id:          id,
			Title:       episode.Title,
			Description: episode.Description,
			PublishedAt: episode.PublishedAt,
			SourceUrl:   episode.Source,
			StoppedAt:   &s,
		})
	}

	sort.Sort(episodes)

	for i := range episodes {
		ids = append(ids, episodes[i].Id)
	}

	return episodes, ids
}

func TranslateFromFeed(model models.Channel, channel *parser.Channel) models.Channel {
	model.Title = channel.Title
	model.Description = channel.Description
	model.Copyright = channel.Copyright
	model.ImageUrl = channel.Image
	model.Uri = helpers.MakeUri(channel.Title)
	model.Language = channel.Language
	model.UpdatedAt = time.Now()
	model.LastBuildDate = channel.LastBuildDate
	model.Url = channel.URL
	return model
}

func translateChannel(channel *parser.Channel) models.Channel {
	model := models.Channel{}
	model.CreatedAt = time.Now()
	return TranslateFromFeed(model, channel)
}

func CreateLinks(links []string, channelId int64, database *gorm.DB) {
	for i := 0; i < len(links); i++ {
		u := models.ChannelURL{}
		database.Table(models.ChannelURL{}.TableName()).
			FirstOrCreate(&u, models.ChannelURL{
				ChannelId: channelId,
				Url:       links[i],
			})
	}
}
