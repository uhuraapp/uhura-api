package services

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uhuraapp/uhura-api/model"
)

type channel struct {
	Service
}

var Channel channel

func (s channel) Top(c *gin.Context) {
	channels := model.Channel(s.connection(c)).Top()

	c.JSON(200, gin.H{"channels": channels})
}

func (s channel) Get(c *gin.Context) {
	channel, episodes, _, found := model.Channel(s.connection(c)).Find(c.Params.ByName("uri"))

	log.Println(channel)
	log.Println(episodes)
	log.Println(found)
	if !found {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// 	userID, _ := helpers.GetUser(c)
	// 	episodes.SetPlays(userID)

	// 	if userId != 0 {
	// 		channel.Subscribed = s.DB.Table(models.Subscription{}.TableName()).Where("user_id = ?", userId).
	// 			Where("channel_id = ?", channel.Id).
	// 			Find(&models.Subscription{}).Error != gorm.ErrRecordNotFound
	// 	}

	// 	c.JSON(200, gin.H{"channel": channel, "episodes": episodes})
	// }

	// func (s ChannelsService) Open(c *gin.Context) {
	// 	var channel entities.Channel
	// 	channelURI := c.Params.ByName("uri")
	// 	s.DB.Table(models.Channel{}.TableName()).Update(&models.Channel{
	// 		VisitedAt: time.Now().UTC(),
	// 	}).Where("uri = ?", channelURI).First(&channel)

	// 	c.JSON(200, gin.H{})
}

// func (s ChannelsService) getPlayed(userId int, episodes entities.Episodes) entities.Episodes {
// 	var played []*models.Listened

// 	ids := episodes.IDs()

// 	if len(ids) > 0 {
// 		s.DB.Table(models.Listened{}.TableName()).
// 			Where("item_uid IN (?)", ids).
// 			Where("user_id = ?", userId).
// 			Find(&played)
// 	}

// 	mapPlayed := make(map[string]*models.Listened, 0)
// 	for _, play := range played {
// 		mapPlayed[play.ItemUID] = play
// 	}

// 	for _, episode := range episodes {
// 		if mapPlayed[episode.Id] != nil {
// 			episode.Listened = mapPlayed[episode.Id].Viewed
// 			episode.StoppedAt = &mapPlayed[episode.Id].StoppedAt
// 			if episode.Listened {
// 				z := int64(0)
// 				episode.StoppedAt = &z
// 			}
// 		}
// 	}

// 	return episodes
// }
