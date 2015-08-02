package middleware

import (
	"os"

	"bitbucket.org/dukex/uhura-api/models"

	authenticator "github.com/dukex/go-auth"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func Authentication(db gorm.DB) gin.HandlerFunc {
	auth := authenticator.NewAuth()
	auth.Helper = models.NewUserHelper(db)

	getProviders(auth)

	return func(c *gin.Context) {
		c.Set("auth", auth)

		if userId, ok := auth.CurrentUser(c.Request); ok {
			c.Set("user_id", userId)
		}

		c.Next()
	}
}

func Protected() gin.HandlerFunc {
	return func(c *gin.Context) {
		_auth, _ := c.Get("auth")
		auth := _auth.(*authenticator.Auth)
		if userId, ok := auth.CurrentUser(c.Request); ok {
			c.Set("user_id", userId)
			c.Next()
		} else {
			c.AbortWithStatus(403)
		}
	}
}

func getProviders(auth *authenticator.Auth) {
	auth.NewProvider(authenticator.EmailPasswordProvider)

	auth.NewProvider(authenticator.Provider{
		RedirectURL: os.Getenv("GOOGLE_CALLBACK_URL"),
		AuthURL:     "https://accounts.google.com/o/oauth2/auth",
		TokenURL:    "https://accounts.google.com/o/oauth2/token",
		Name:        "google",
		Key:         os.Getenv("GOOGLE_CLIENT_ID"),
		Secret:      os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:      []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		UserInfoURL: "https://www.googleapis.com/oauth2/v1/userinfo?alt=json",
	})

	auth.NewProvider(authenticator.Provider{
		RedirectURL: os.Getenv("FACEBOOK_CALLBACK_URL"),
		AuthURL:     "https://www.facebook.com/dialog/oauth",
		TokenURL:    "https://graph.facebook.com/oauth/access_token",
		Name:        "facebook",
		Key:         os.Getenv("FACEBOOK_CLIENT_ID"),
		Secret:      os.Getenv("FACEBOOK_CLIENT_SECRET"),
		Scopes:      []string{"email"},
		UserInfoURL: "https://graph.facebook.com/me",
	})
}
