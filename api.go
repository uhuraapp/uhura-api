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

	r := _r.Group("/v2", middleware.Authentication())
	{
		r.GET("/subscriptions", subscriptions.Get)
		r.GET("/suggestions", suggestions.Get)
		r.GET("/channels/:uri", channels.Get)
	}
}
