package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/uhuraapp/uhura-api/database"
	"github.com/uhuraapp/uhura-api/middleware"
	"github.com/uhuraapp/uhura-api/services"
)

func main() {
	r := gin.Default()
	r.Use(middleware.CORS())

	Mount(r.Group("/"))

	log.Println("Listening...")
	http.ListenAndServe(":"+os.Getenv("PORT"), r)
}

func Mount(_r *gin.RouterGroup) {
	DB := database.New()

	subscriptions := services.NewSubscriptionService(DB)
	suggestions := services.NewSuggestionsService(DB)
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
		r.GET("/user/logout", auth.Logout)

		r.GET("/subscriptions", needAuth, subscriptions.Get)
		r.GET("/suggestions", needAuth, suggestions.Get)
		r.GET("/episodes/:id/listened", needAuth, episodes.Listened)
		r.GET("/episodes/:id/download", needAuth, episodes.Download)
	}
}
