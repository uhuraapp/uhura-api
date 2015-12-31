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
	r.Use(middleware.Recovery())

	Mount(r.Group("/"))

	log.Println("Listening", os.Getenv("PORT"), "...")
	http.ListenAndServe(":"+os.Getenv("PORT"), middleware.Gzip(r))
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
	profile := services.NewProfileService(DB)
	users := services.NewUserService(DB)

	_r.Use(middleware.Authentication(DB))

	needAuth := middleware.Protected()
	apiAuth := middleware.ApiProtected()

	r := _r.Group("/v2")
	{
		r.OPTIONS("*action", func(c *gin.Context) { c.Data(200, "", []byte{}) })

		r.GET("/top/channels", channels.Top)

		r.GET("/channels/*uri", channels.Get)
		// r.GET("/channels/:uri/open", channels.Open)

		r.GET("/parser", parser.ByURL)

		r.GET("/auth/:provider", auth.ByProvider)
		r.GET("/auth/:provider/callback", auth.ByProviderCallback)

		r.GET("/user", auth.GetUser)
		r.PUT("/user", auth.UpdateUser)
		r.DELETE("/user", auth.DeleteUser)
		r.POST("/users/sign_in", auth.ByEmailPassword)
		r.GET("/users/logout", auth.Logout)
		r.POST("/users", auth.SignUp)

		r.POST("/profiles", needAuth, profile.Create)
		r.GET("/me", needAuth, profile.Me)
		r.GET("/profiles/:key", profile.Get)

		r.GET("/users/subscriptions", needAuth, userSubscriptions.Index)
		r.POST("/users/subscriptions", needAuth, userSubscriptions.Create)
		r.GET("/users/subscriptions/:uri", needAuth, userSubscriptions.Show)
		r.DELETE("/users/subscriptions/:uri", needAuth, userSubscriptions.Delete)
		r.GET("/users/recommendations", needAuth, userRecommendations.Index)

		r.GET("/episodes", episodes.GetPaged)
		r.GET("/episodes/:id", episodes.Get)

		r.POST("/episodes/:id/played", needAuth, episodes.Played)
		r.DELETE("/episodes/:id/played", needAuth, episodes.Unlistened)

		r.GET("/episodes/:id/download", needAuth, episodes.Download)
		r.HEAD("/episodes/:id/download", needAuth, episodes.Download)

		r.PUT("/episodes/:id/listen", needAuth, episodes.Listen)

		r.GET("/categories", categories.Index)
		r.GET("/categories/:uri", categories.Get)
	}

	v := _r.Group("/v3")
	{
		v.GET("/users.json", apiAuth, users.All)
		v.GET("/users/:id", apiAuth, users.Get)
	}
}
