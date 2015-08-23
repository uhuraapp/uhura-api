package services

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"

	"bitbucket.org/dukex/uhura-api/entities"
	"bitbucket.org/dukex/uhura-api/helpers"
	"bitbucket.org/dukex/uhura-api/models"
	authenticator "github.com/dukex/go-auth"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type AuthService struct {
	DB gorm.DB
}

func NewAuthService(db gorm.DB) AuthService {
	return AuthService{DB: db}
}

func (s AuthService) ByProvider(c *gin.Context) {
	auth, _ := s.getAuth(c)
	authorizer := auth.Authorize(c.Params.ByName("provider"))
	authorizer(c.Writer, c.Request)
	return
}

func (s AuthService) ByEmailPassword(c *gin.Context) {
	auth, _ := s.getAuth(c)
	auth.SignIn(c.Writer, c.Request)
	return
}

func (s AuthService) ByProviderCallback(c *gin.Context) {
	auth, _ := s.getAuth(c)
	userId, err := auth.OAuthCallback(c.Params.ByName("provider"), c.Writer, c.Request)
	if err == nil {
		session := auth.Login(c.Request, userId)
		session.Save(c.Request, c.Writer)
		// 	go helpers.NewEvent(userIdint, "login", map[string]interface{}{})
	}

	closeHTML := []byte("<html><head></head><body>Loading....<script>window.close()</script></body></html>")
	c.Data(200, "text/html", closeHTML)
}

func (s AuthService) GetUser(c *gin.Context) {
	var user entities.User
	auth, _ := s.getAuth(c)

	userId, ok := auth.CurrentUser(c.Request)
	if !ok {
		c.AbortWithStatus(404)
		return
	}

	err := s.DB.Table(models.User{}.TableName()).Where("id = ?", userId).First(&user).Error
	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	if user.ApiToken == "" {
		token := authenticator.NewUserToken()
		hasher := md5.New()
		hasher.Write([]byte(token + user.Email))
		user.ApiToken = hex.EncodeToString(hasher.Sum(nil))
		s.DB.Table(models.User{}.TableName()).Where("id = ?", userId).Update("api_token", user.ApiToken)
	}

	s.DB.Table(models.User{}.TableName()).Where("id = ?", userId).Update("last_visited_at", time.Now().Format(time.RubyDate))

	c.JSON(200, user)
}

func (s AuthService) Logout(c *gin.Context) {
	auth, _ := s.getAuth(c)
	userId, _ := auth.CurrentUser(c.Request)

	go helpers.NewEvent(userId, "logout", map[string]interface{}{})
	session := auth.Logout(c.Request)
	session.Save(c.Request, c.Writer)

	c.Data(200, "", []byte(""))
}

func (s AuthService) getAuth(c *gin.Context) (auth *authenticator.Auth, err error) {
	var tempInterface interface{}
	var ok bool
	tempInterface, ok = c.Get("auth")
	if !ok {
		err = errors.New("not found auth")
	}
	auth = tempInterface.(*authenticator.Auth)
	return
}
