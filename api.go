package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"bitbucket.org/dukex/uhura-api/database"
	"bitbucket.org/dukex/uhura-api/middleware"
	"bitbucket.org/dukex/uhura-api/models"
	"bitbucket.org/dukex/uhura-api/services"
	"github.com/gin-gonic/gin"
	"github.com/hjr265/too"
	"github.com/jinzhu/gorm"
)

func main() {
	r := gin.Default()
	r.Use(middleware.CORS())

	Mount(r.Group("/"))

	log.Println("Listening", os.Getenv("PORT"), "...")
	http.ListenAndServe(":"+os.Getenv("PORT"), r)
}

func Mount(_r *gin.RouterGroup) {
	DB := database.New()

	userSubscriptions := services.NewUserSubscriptionService(DB)
	subscriptions := services.NewSubscriptionService(DB)
	channels := services.NewChannelsService(DB)
	episodes := services.NewEpisodesService(DB)
	auth := services.NewAuthService(DB)
	categories := services.NewCategoriesService(DB)

	_r.Use(middleware.Authentication(DB))

	needAuth := middleware.Protected()

	r := _r.Group("/v2")
	{
		r.OPTIONS("*action", func(c *gin.Context) { c.Data(200, "", []byte{}) })

		r.GET("/channels/:uri", channels.Get)

		r.GET("/auth/:provider", auth.ByProvider)
		r.GET("/auth/:provider/callback", auth.ByProviderCallback)

		r.GET("/user", auth.GetUser)
		r.GET("/users/logout", auth.Logout)
		r.GET("/users/subscriptions", needAuth, userSubscriptions.Index)
		r.POST("/users/subscriptions", needAuth, userSubscriptions.Create)
		r.GET("/users/subscriptions/:uri", needAuth, userSubscriptions.Show)
		r.DELETE("/users/subscriptions/:uri", needAuth, userSubscriptions.Delete)

		r.GET("/episodes", needAuth, episodes.GetPaged)
		r.GET("/episodes/:id/listened", needAuth, episodes.Listened)
		r.GET("/episodes/:id/download", needAuth, episodes.Download)

		r.GET("/subscriptions/top", subscriptions.Top)
		r.GET("/categories", categories.Index)
	}

	// GetRecommendations(DB)
}

func GetRecommendations(DB gorm.DB) {
	var subscriptions []models.Subscription
	DB.Table(models.Subscription{}.TableName()).Find(&subscriptions)

	redisAddr, _ := net.ResolveTCPAddr("tcp", os.Getenv("REDIS_URL"))
	te, _ := too.New(redisAddr, "channels")

	for _, subscription := range subscriptions {
		var channel models.Channel
		if gorm.RecordNotFound == DB.Table(models.Channel{}.TableName()).Where("id = ?", subscription.ChannelId).Find(&channel).Error {
			log.Panic(subscription.ChannelId)
		} else {
			te.Likes.Add(too.User(subscription.UserId), too.Item(channel.Id))
		}
	}

	log.Println("Finished recommendations!")
}
