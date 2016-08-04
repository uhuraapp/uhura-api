package services

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/uhuraapp/uhura-api/channels"
	"github.com/uhuraapp/uhura-api/entities"
	"github.com/uhuraapp/uhura-api/helpers"
	"github.com/uhuraapp/uhura-api/models"
)

type ChannelsService struct {
	DB *gorm.DB
}

func NewChannelsService(db *gorm.DB) ChannelsService {
	return ChannelsService{DB: db}
}

func (s ChannelsService) Top(c *gin.Context) {
	var channels []entities.Channel
	s.DB.Table(models.Subscription{}.TableName()).
		Select("channels.title, channels.image_url, channels.description, channels.uri, channels.id, COUNT(*) AS subscribers_count").
		Joins("INNER JOIN channels ON user_channels.channel_id = channels.id").
		Group("1,2,3,4,5 ORDER BY subscribers_count DESC").
		Limit(5).
		Find(&channels)

	c.JSON(200, gin.H{"channels": channels})
}

func (s ChannelsService) Get(c *gin.Context) {
	channel, episodes, _, found := channels.Find(s.DB, c.Params.ByName("uri"))

	if !found {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	userId, _ := helpers.GetUser(c)
	episodes = s.getPlayed(userId, episodes)

	// 	if userId != 0 {
	// 		channel.Subscribed = s.DB.Table(models.Subscription{}.TableName()).Where("user_id = ?", userId).
	// 			Where("channel_id = ?", channel.Id).
	// 			Find(&models.Subscription{}).Error != gorm.ErrRecordNotFound
	// 	}
	// }

	c.JSON(200, gin.H{"channel": channel, "episodes": episodes})
}

func (s ChannelsService) Open(c *gin.Context) {
	var channel entities.Channel
	channelURI := c.Params.ByName("uri")
	s.DB.Table(models.Channel{}.TableName()).Update(&models.Channel{
		VisitedAt: time.Now().UTC(),
	}).Where("uri = ?", channelURI).First(&channel)

	c.JSON(200, gin.H{})
}

func (s ChannelsService) getPlayed(userId int, episodes entities.Episodes) entities.Episodes {
	var played []*models.Listened

	ids := episodes.IDs()

	if len(ids) > 0 {
		s.DB.Table(models.Listened{}.TableName()).
			Where("item_uid IN (?)", ids).
			Where("user_id = ?", userId).
			Find(&played)
	}

	mapPlayed := make(map[string]*models.Listened, 0)
	for _, play := range played {
		mapPlayed[play.Id] = play
	}

	for _, episode := range episodes {
		if mapPlayed[episode.Id] != nil {
			episode.Listened = mapPlayed[episode.Id].Viewed
			episode.StoppedAt = &mapPlayed[episode.Id].StoppedAt
			if episode.Listened {
				z := int64(0)
				episode.StoppedAt = &z
			}
		}
	}

	return episodes
}

func (s ChannelsService) Index(c *gin.Context) {
	LIMIT := 25
	channels := make([]entities.Channel, 0)
	query := c.Request.URL.Query()
	q := query.Get("q")
	if q != "" {
		page, err := strconv.Atoi(query.Get("page"))
		if err != nil {
			page = 0
		} else {
			page = page - 1
		}
		offset := page * LIMIT

		ids := make([]int64, 0)

		db := s.DB.DB()

		rids, err := db.Query("SELECT id FROM (SELECT id, tsv FROM channels, plainto_tsquery('" + q + "') AS q WHERE (tsv @@ q)) AS t1 ORDER BY ts_rank_cd(t1.tsv, plainto_tsquery('" + q + "')) DESC")
		log.Println("Error query", err)
		if err == nil {
			for rids.Next() {
				var id int64
				err = rids.Scan(&id)
				if err == nil {
					ids = append(ids, id)
				}
			}
		}

		if len(ids) > 0 {
			s.DB.Table(models.Channel{}.TableName()).
				Where("id in (?)", ids).
				Limit(LIMIT).
				Offset(offset).
				Find(&channels)
		}
	}

	for i := range channels {
		channels[i].Episodes = make([]string, 0)
	}

	c.JSON(200, gin.H{"channels": channels, "episodes": []int64{}})
}
