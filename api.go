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
	r.Use(middleware.Recovery())

	Mount(r.Group("/"))

	log.Println("Listening", os.Getenv("PORT"), "...")
	log.Println(http.ListenAndServe(":"+os.Getenv("PORT"), middleware.Gzip(r)))
}

func Mount(_r *gin.RouterGroup) {
	DB := database.NewPostgresql()

	userSubscriptions := services.NewUserSubscriptionService(DB)
	userRecommendations := services.NewUserRecommendationService(DB)
	channels := services.NewChannelsService(DB)
	episodes := services.NewEpisodesService(DB)
	auth := services.NewAuthService(DB)
	categories := services.NewCategoriesService(DB)
	parser := services.NewParser(DB)
	users := services.NewUserService(DB)
	export := services.NewExportService(DB)

	_r.Use(middleware.Authentication(DB))

	needAuth := middleware.Protected()
	apiAuth := middleware.ApiProtected()

	r := _r.Group("/v2")
	{
		r.OPTIONS("*action", func(c *gin.Context) { c.Data(200, "", []byte{}) })

		r.GET("/top/channels", channels.Top)

		r.GET("/channels", channels.Index)
		r.GET("/channels/*uri", channels.Get)
		// r.GET("/sync/:uri", channels.Sync)

		r.GET("/parser", parser.ByURL)

		r.GET("/auth/:provider", auth.ByProvider)
		r.GET("/auth/:provider/callback", auth.ByProviderCallback)

		r.GET("/user", auth.GetUser)
		r.PUT("/user", auth.UpdateUser)
		r.DELETE("/user", auth.DeleteUser)
		r.POST("/users/sign_in", auth.ByEmailPassword)
		r.GET("/users/logout", auth.Logout)
		r.POST("/users", auth.SignUp)

		r.GET("/users/subscriptions", needAuth, userSubscriptions.Index)
		r.POST("/users/subscriptions", needAuth, userSubscriptions.Create)
		r.GET("/users/subscriptions/:uri", needAuth, userSubscriptions.Show)
		r.DELETE("/users/subscriptions/:uri", needAuth, userSubscriptions.Delete)
		r.GET("/users/recommendations", needAuth, userRecommendations.Index)

		// r.GET("/episodes", episodes.GetPaged)
		//r.GET("/episodes/:id", episodes.Get)

		r.POST("/channels/:channel_id/episodes/:id/played", needAuth, episodes.Played)
		r.DELETE("/channels/:channel_id/episodes/:id/played", needAuth, episodes.UnPlayed)

		//r.GET("/episodes/:id/download", needAuth, episodes.Download)
		//r.HEAD("/episodes/:id/download", needAuth, episodes.Download)

		//r.PUT("/episodes/:id/listen", needAuth, episodes.Listen)

		r.GET("/categories", categories.Index)
		r.GET("/categories/:uri", categories.Get)

		r.GET("/export/:id", export.Get)
	}

	v := _r.Group("/v3")
	{
		v.GET("/users.json", apiAuth, users.All)
		v.GET("/users/:id", apiAuth, users.Get)
	}
}
