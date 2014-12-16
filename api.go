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

	apiRouter := r.Group("/api")
	Mount(apiRouter)

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

	r := _r.Group("/v2")
	{
		r.GET("/channels/:uri", channels.Get)

		r.GET("/auth/:provider", auth.ByProvider)
		r.GET("/auth/:provider/callback", auth.ByProviderCallback)

		r.GET("/subscriptions", middleware.Authentication(), subscriptions.Get)
		r.GET("/suggestions", middleware.Authentication(), suggestions.Get)
		r.GET("/episodes/:id/listened", middleware.Authentication(), episodes.Listened)
		r.GET("/episodes/:id/download", middleware.Authentication(), episodes.Download)
	}
}
