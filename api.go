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
	DB := database.NewPostgresql()

	userSubscriptions := services.NewUserSubscriptionService(DB)
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
		r.GET("/channels/:uri/open", channels.Open)

		r.GET("/auth/:provider", auth.ByProvider)
		r.GET("/auth/:provider/callback", auth.ByProviderCallback)

		r.GET("/user", auth.GetUser)
		r.GET("/users/logout", auth.Logout)
		r.GET("/users/subscriptions", needAuth, userSubscriptions.Index)
		r.POST("/users/subscriptions", needAuth, userSubscriptions.Create)
		r.GET("/users/subscriptions/:uri", needAuth, userSubscriptions.Show)
		r.DELETE("/users/subscriptions/:uri", needAuth, userSubscriptions.Delete)

		r.GET("/episodes", needAuth, episodes.GetPaged)
		r.POST("/episodes/:id/listened", needAuth, episodes.Listened)
		r.PUT("/episodes/:id/listen", needAuth, episodes.Listen)
		r.DELETE("/episodes/:id/listened", needAuth, episodes.Unlistened)
		r.GET("/episodes/:id/download", needAuth, episodes.Download)
		r.HEAD("/episodes/:id/download", needAuth, episodes.Download)

		r.GET("/categories", categories.Index)
	}
}
