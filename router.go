package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/uhuraapp/uhura-api/services"
)

func main() {
	router := gin.Default()

	Mount(router.Group("/"))

	log.Println("Listening", os.Getenv("PORT"), "...")
	log.Println(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}

func Mount(router *gin.RouterGroup) {
	router.Use(services.Connection())

	router.Group("/v2")
	{
		router.OPTIONS("*action", func(c *gin.Context) { c.Data(200, "", []byte{}) })

		router.GET("/top/channels", services.Channel.Top)

		//     r.GET("/channels/*uri", channels.Get)

		//     r.GET("/parser", parser.ByURL)

		//     r.GET("/auth/:provider", auth.ByProvider)
		//     r.GET("/auth/:provider/callback", auth.ByProviderCallback)

		//     r.GET("/user", auth.GetUser)
		//     r.PUT("/user", auth.UpdateUser)
		//     r.DELETE("/user", auth.DeleteUser)
		//     r.POST("/users/sign_in", auth.ByEmailPassword)
		//     r.GET("/users/logout", auth.Logout)
		//     r.POST("/users", auth.SignUp)

		//     r.GET("/users/subscriptions", needAuth, userSubscriptions.Index)
		//     r.POST("/users/subscriptions", needAuth, userSubscriptions.Create)
		//     r.GET("/users/subscriptions/:uri", needAuth, userSubscriptions.Show)
		//     r.DELETE("/users/subscriptions/:uri", needAuth, userSubscriptions.Delete)

		//     r.GET("/episodes", episodes.Index)
		//     r.GET("/episodes/:id", episodes.Get)

		//     r.POST("/channels/:channel_id/episodes/:id/played", needAuth, episodes.Played)
		//     r.DELETE("/channels/:channel_id/episodes/:id/played", needAuth, episodes.UnPlayed)
		//     r.PUT("/channels/:channel_id/episodes/:id/listen", needAuth, episodes.Listen)
		//     // r.GET("/channels/:channel_id/episodes/:id/download", needAuth, episodes.Download(c))
		//     //r.HEAD("/channels/:channel_id/episodes/:id/download", needAuth, episodes.Download)

		//     r.GET("/categories", categories.Index)
		//     r.GET("/categories/:uri", categories.Get)

		//     r.GET("/export/:id", export.Get)
	}
}
