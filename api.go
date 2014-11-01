package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/uhuraapp/uhura-api/middleware"
	"github.com/uhuraapp/uhura-api/subscriptions"
	"github.com/uhuraapp/uhura-api/suggestions"
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
	r := _r.Group("/v2", middleware.Authentication())
	{
		r.GET("/subscriptions", subscriptions.Get)
		r.GET("/suggestions", suggestions.Get)
	}
}
