package services

import (
	"log"
	"os"

	"github.com/dukex/login2"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/uhuraapp/uhura-api/models"
)

type AuthService struct {
	DB    gorm.DB
	login *login2.Builder
}

func NewAuthService(db gorm.DB) AuthService {
	userHelper := models.UserHelpers{DB: db}

	login := login2.NewBuilder()
	login.UserSetupFn = userHelper.SetupFromOAuth

	getProviders(login)
	return AuthService{DB: db, login: login}
}

func (s AuthService) ByProvider(c *gin.Context) {
	authorizer := s.login.OAuthAuthorize(c.Params.ByName("provider"))
	authorizer(c.Writer, c.Request)
}

func (s AuthService) ByProviderCallback(c *gin.Context) {
	userID, err := s.login.OAuthCallback(c.Params.ByName("provider"), c.Request)
	log.Println(userID, err)
}

func getProviders(login *login2.Builder) {
	login.NewProvider(&login2.Provider{
		RedirectURL: os.Getenv("GOOGLE_CALLBACK_URL"),
		AuthURL:     "https://accounts.google.com/o/oauth2/auth",
		TokenURL:    "https://accounts.google.com/o/oauth2/token",
		Name:        "google",
		Key:         os.Getenv("GOOGLE_CLIENT_ID"),
		Secret:      os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scope:       "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile",
		UserInfoURL: "https://www.googleapis.com/oauth2/v1/userinfo?alt=json",
	})

	login.NewProvider(&login2.Provider{
		RedirectURL: os.Getenv("FACEBOOK_CALLBACK_URL"),
		AuthURL:     "https://www.facebook.com/dialog/oauth",
		TokenURL:    "https://graph.facebook.com/oauth/access_token",
		Name:        "facebook",
		Key:         os.Getenv("FACEBOOK_CLIENT_ID"),
		Secret:      os.Getenv("FACEBOOK_CLIENT_SECRET"),
		Scope:       "email",
		UserInfoURL: "https://graph.facebook.com/me",
	})
}
