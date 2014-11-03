package services_test

import (
	"net/http"
	"strconv"

	"github.com/bitly/go-simplejson"
	"github.com/jinzhu/gorm"
	"github.com/uhuraapp/uhura-api/cache"
	"github.com/uhuraapp/uhura-api/models"
	"github.com/uhuraapp/uhura-api/services"
	. "gopkg.in/check.v1"
)

type SubscriptionsSuite struct {
	DB gorm.DB
}

var _ = Suite(&SubscriptionsSuite{})

func (x *SubscriptionsSuite) SetUpSuite(c *C) {
	x.DB = databaseTest()
}

func (_ *SubscriptionsSuite) SetUpTest(c *C) {
	cache.DATA = nil
}

func (_ *SubscriptionsSuite) TearDownTest(c *C) {
	resetDatabase()
}

func (x SubscriptionsSuite) TestGetIsOk(c *C) {
	s := services.NewSubscriptionService(x.DB)
	r := request("GET", s.Get)

	c.Assert(r.Code, Equals, http.StatusOK)
}

func (x SubscriptionsSuite) TestGetReturnSubscriptions(c *C) {
	channel := models.Channel{Title: "Meu Podcast"}

	x.DB.Create(&channel)
	x.DB.Create(&models.Subscription{ChannelId: channel.Id, UserId: 1})

	s := services.NewSubscriptionService(x.DB)
	r := request("GET", s.Get)

	data, _ := simplejson.NewJson(r.Body.Bytes())
	currentTitle, _ := data.GetPath("subscriptions").GetIndex(0).Get("title").String()

	c.Assert(channel.Title, Equals, currentTitle)
	cached, _ := cache.Get("s:ids:1", []int64{})
	c.Assert(cached.([]int64)[0], Equals, channel.Id)
}

func (x SubscriptionsSuite) TestGetReturnSubscriptionsCached(c *C) {
	channel := models.Channel{Title: "Meu Podcast"}

	x.DB.Create(&channel)

	cache.Set("s:ids:1", []int64{channel.Id})

	s := services.NewSubscriptionService(x.DB)
	r := request("GET", s.Get)

	data, _ := simplejson.NewJson(r.Body.Bytes())
	currentTitle, _ := data.GetPath("subscriptions").GetIndex(0).Get("title").String()

	c.Assert(channel.Title, Equals, currentTitle)
}

func (x SubscriptionsSuite) TestGetReturnSubscriptionsCachedError(c *C) {
	channel := models.Channel{Title: "Meu Podcast"}
	x.DB.Create(&channel)
	x.DB.Create(&models.Subscription{ChannelId: channel.Id, UserId: 1})
	cache.Set("s:ids:1", "ABC")
	s := services.NewSubscriptionService(x.DB)
	r := request("GET", s.Get)

	data, _ := simplejson.NewJson(r.Body.Bytes())
	currentTitle, _ := data.GetPath("subscriptions").GetIndex(0).Get("title").String()

	c.Assert(channel.Title, Equals, currentTitle)
	cached, _ := cache.Get("s:ids:1", []int64{})
	c.Assert(cached.([]int64)[0], Equals, channel.Id)
}

func (x SubscriptionsSuite) TestGetToView(c *C) {
	channel := models.Channel{Title: "Meu Podcast"}
	x.DB.Create(&channel)
	x.DB.Create(&models.Subscription{ChannelId: channel.Id, UserId: 1})
	x.DB.Create(&models.Episode{Id: 1, SourceUrl: "a", Key: "a", ChannelId: channel.Id})
	x.DB.Create(&models.Episode{Id: 2, SourceUrl: "b", Key: "b", ChannelId: channel.Id})
	x.DB.Create(&models.Listened{ChannelId: channel.Id, UserId: 1})

	s := services.NewSubscriptionService(x.DB)
	r := request("GET", s.Get)

	data, _ := simplejson.NewJson(r.Body.Bytes())
	toView, _ := data.GetPath("subscriptions").GetIndex(0).Get("to_view").Int64()

	c.Assert(toView, Equals, int64(1))
	cached, _ := cache.Get("c:e:"+strconv.Itoa(int(channel.Id)), 0)
	c.Assert(cached, Equals, int64(2))
}

func (x SubscriptionsSuite) TestGetToViewCached(c *C) {
	channel := models.Channel{Title: "Meu Podcast"}
	x.DB.Create(&channel)
	x.DB.Create(&models.Subscription{ChannelId: channel.Id, UserId: 1})
	x.DB.Create(&models.Listened{ChannelId: channel.Id, UserId: 1})

	cache.Set("c:e:"+strconv.Itoa(int(channel.Id)), int64(2))

	s := services.NewSubscriptionService(x.DB)
	r := request("GET", s.Get)

	data, _ := simplejson.NewJson(r.Body.Bytes())
	toView, _ := data.GetPath("subscriptions").GetIndex(0).Get("to_view").Int64()

	c.Assert(toView, Equals, int64(1))
}
