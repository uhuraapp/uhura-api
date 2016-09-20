package services_test

import (
	"log"
	"net/http"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/jinzhu/gorm"
	"github.com/uhuraapp/uhura-api/database"
	"github.com/uhuraapp/uhura-api/models"
	"github.com/uhuraapp/uhura-api/services"
	. "gopkg.in/check.v1"
)

type SubscriptionsSuite struct {
	DB *gorm.DB
}

var _ = Suite(&SubscriptionsSuite{})

func (x *SubscriptionsSuite) SetUpSuite(c *C) {
	x.DB = databaseTest()
}

func (_ *SubscriptionsSuite) SetUpTest(c *C) {
	database.CACHE.Flush()
}

func (_ *SubscriptionsSuite) TearDownTest(c *C) {
	resetDatabase()
}

func (x SubscriptionsSuite) TestGetIsOk(c *C) {
	s := services.NewUserSubscriptionService(x.DB)
	r := request("GET", s.Index, "")

	c.Assert(r.Code, Equals, http.StatusOK)
}

func (x SubscriptionsSuite) TestGetReturnSubscriptions(c *C) {
	channel := models.Channel{Title: "Meu Podcast", Language: "PT"}

	x.DB.Create(&channel)
	log.Println("CHNNAE", channel)
	x.DB.Create(&models.Subscription{ChannelId: channel.Id, UserId: 1})

	s := services.NewUserSubscriptionService(x.DB)
	r := request("GET", s.Index, "")

	data, _ := simplejson.NewJson(r.Body.Bytes())
	log.Println("SSASAS", data)
	currentTitle, _ := data.GetPath("subscriptions").GetIndex(0).Get("title").String()

	c.Assert(channel.Title, Equals, currentTitle)
	c.Assert(r.Code, Equals, http.StatusOK)
	cached, _ := database.CACHE.Value("s:1")
	c.Assert(cached.Data().(*time.Time), NotNil)
}

func (x SubscriptionsSuite) TestGetReturnSubscriptionsCached(c *C) {
	channel := models.Channel{Title: "Meu Podcast"}

	x.DB.Create(&channel)

	now := time.Now()
	t := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()-2, 0, 0, 0, now.Location())
	database.CACHE.Add("s:1", 1*time.Minute, &t)

	s := services.NewUserSubscriptionService(x.DB)
	r := request("GET", s.Index, time.Now().Add(-3*time.Hour).Format(time.RFC1123))

	c.Assert(string(r.Body.Bytes()), Equals, "")
	c.Assert(r.Code, Equals, http.StatusNotModified)
}

func (x SubscriptionsSuite) TestGetIsOkBot(c *C) {
	s := services.NewUserSubscriptionService(x.DB)
	r := botRequest("GET", s.Index)

	data, _ := simplejson.NewJson(r.Body.Bytes())
	subscriptions, _ := data.GetPath("subscriptions").Array()

	c.Assert(len(subscriptions), Equals, 0)
}
