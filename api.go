package main

import (
	"log"
	"net/http"
	"os"

	"bitbucket.org/dukex/uhura-api/database"
	"bitbucket.org/dukex/uhura-api/middleware"
	"bitbucket.org/dukex/uhura-api/services"
	"github.com/gin-gonic/gin"
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

	_r.Use(middleware.Authentication(DB))

	needAuth := middleware.Protected()

	r := _r.Group("/v2")
	{
		r.OPTIONS("*action", func(c *gin.Context) { c.Data(200, "", []byte{}) })
		r.GET("/channels/:uri", channels.Get)

		r.GET("/auth/:provider", auth.ByProvider)
		r.GET("/auth/:provider/callback", auth.ByProviderCallback)

		r.GET("/user", needAuth, auth.GetUser)
		r.GET("/users/logout", auth.Logout)
		r.GET("/users/subscriptions", needAuth, userSubscriptions.Get)
		r.GET("/users/suggestions", needAuth, suggestions.Get)

		r.GET("/episodes", needAuth, episodes.GetPaged)
		r.GET("/episodes/:id/listened", needAuth, episodes.Listened)
		r.GET("/episodes/:id/download", needAuth, episodes.Download)

		r.GET("/subscriptions/top", subscriptions.Top)
	}
}
