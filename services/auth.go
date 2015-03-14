package services

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"time"

	"bitbucket.org/dukex/uhura-api/entities"
	"bitbucket.org/dukex/uhura-api/helpers"
	"bitbucket.org/dukex/uhura-api/models"
	login "github.com/dukex/login2"
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
	auth.SetReturnTo(c.Writer, c.Request, c.Request.Header.Get("Origin"))
	authorizer := auth.OAuthAuthorize(c.Params.ByName("provider"))
	authorizer(c.Writer, c.Request)
}

func (s AuthService) ByProviderCallback(c *gin.Context) {
	auth, _ := s.getAuth(c)
	userID, err := auth.OAuthCallback(c.Params.ByName("provider"), c.Request)
	if err == nil {
		userIdint := strconv.FormatInt(userID, 10)
		session := auth.Login(c.Request, userIdint)
		session.Save(c.Request, c.Writer)
		go helpers.NewEvent(userIdint, "login", map[string]interface{}{})
	}

	closeHTML := []byte("<html><head></head><body>Loading....<script>window.close()</script></body></html>")
	c.Data(200, "text/html", closeHTML)
}

func (s AuthService) GetUser(c *gin.Context) {
	var user entities.User
	auth, _ := s.getAuth(c)

	userId, _ := auth.CurrentUser(c.Request)
	err := s.DB.Table(models.User{}.TableName()).Where("id = ?", userId).First(&user).Error
	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	go func() {
		person := helpers.MixpanelPerson(userId)
		person.Update("$set", map[string]interface{}{
			"$last_login": time.Now(),
		})
	}()

	if user.ApiToken == "" {
		token := login.NewUserToken()
		hasher := md5.New()
		hasher.Write([]byte(token + user.Email))
		s.DB.Table(models.User{}.TableName()).Where("id = ?", userId).Update("api_token", hex.EncodeToString(hasher.Sum(nil)))
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

func (s AuthService) getAuth(c *gin.Context) (builder *login.Builder, err error) {
	var tempInterface interface{}
	tempInterface, err = c.Get("auth")
	builder = tempInterface.(*login.Builder)
	return
}
