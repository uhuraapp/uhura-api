package services

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	authenticator "github.com/dukex/go-auth"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/uhuraapp/uhura-api/entities"
	"github.com/uhuraapp/uhura-api/models"
)

type AuthService struct {
	DB *gorm.DB
}

func NewAuthService(db *gorm.DB) AuthService {
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
	if userId, ok := auth.SignIn(c.Writer, c.Request); ok {
		go s.setAgree(userId)
	}
	return
}

func (s AuthService) ByProviderCallback(c *gin.Context) {
	auth, _ := s.getAuth(c)
	userId, err := auth.OAuthCallback(c.Params.ByName("provider"), c.Writer, c.Request)
	if err == nil {
		session := auth.Login(c.Request, userId)
		session.Save(c.Request, c.Writer)
		go s.setAgree(userId)
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
		s.createUserToken(user.Email)
	}

	user.OptIn = !user.OptInAt.IsZero()

	s.DB.Table(models.User{}.TableName()).Where("id = ?", userId).Update("last_visited_at", time.Now().Format(time.RubyDate))

	c.JSON(200, user)
}

func (s AuthService) UpdateUser(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var params struct {
		User struct {
			Name  string `json:"name"`
			OptIn bool   `json:"optin"`
		} `json:"user"`
	}

	decoder.Decode(&params)

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

	if !params.User.OptIn {
		user.OptInAt = time.Time{}
	} else if params.User.OptIn && user.OptInAt.IsZero() {
		user.OptInAt = time.Now()
	}

	s.DB.Table(models.User{}.TableName()).Where("id = ?", userId).Update(map[string]interface{}{
		"name":      params.User.Name,
		"opt_in_at": user.OptInAt,
	})
	s.DB.Table(models.User{}.TableName()).Where("id = ?", userId).First(&user)

	user.OptIn = !user.OptInAt.IsZero()
	c.JSON(200, user)
}

func (s AuthService) SignUp(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var params struct {
		User struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			Name     string `json:"name"`
		} `json:"user"`
	}

	decoder.Decode(&params)

	if params.User.Email == "" || params.User.Password == "" || params.User.Name == "" {
		c.JSON(422, entities.ErrorResponse{
			[]entities.Error{{
				Key: "fields_required", Message: "Email, Password and Name is required"},
			},
		})
		return
	}

	if len([]rune(params.User.Password)) < 6 {
		c.JSON(422, entities.ErrorResponse{
			[]entities.Error{{
				Key: "password_too_short", Message: "Password too short: Minimum amount of characters 6"},
			},
		},
		)
		return
	}

	password, err := authenticator.GenerateHash(params.User.Password)

	if err != nil {
		c.JSON(422, entities.ErrorResponse{
			[]entities.Error{{
				Key: "generate_hash_password_error", Message: "Internal server error: Could not encrypt your password ", Error: err.Error()},
			},
		})
		return
	}

	user := models.User{
		Email:         params.User.Email,
		Password:      password,
		Name:          params.User.Name,
		Provider:      "email",
		RememberToken: authenticator.NewUserToken(),
	}

	err = s.DB.Table(user.TableName()).Where("email = ?", user.Email).First(&models.User{}).Error
	if err == nil {
		c.JSON(422, entities.ErrorResponse{
			[]entities.Error{{
				Key: "already_registrated", Message: "Email already exists in the system, try log in or other email"},
			},
		})
		return
	}

	err = s.DB.Table(user.TableName()).Save(&user).Error
	if err != nil {
		c.JSON(422, entities.ErrorResponse{
			[]entities.Error{{
				Key: "internal_server_error", Message: "Internal server error: try again", Error: err.Error()},
			},
		})
		return
	}

	auth, _ := s.getAuth(c)
	userId := strconv.Itoa(int(user.Id))
	session := auth.Login(c.Request, userId)
	session.Save(c.Request, c.Writer)

	s.createUserToken(user.Email)
	go s.setAgree(strconv.Itoa(int(user.Id)))

	c.JSON(http.StatusCreated, struct {
		Id int64 `json:"id"`
	}{user.Id})
}

func (s AuthService) Logout(c *gin.Context) {
	auth, _ := s.getAuth(c)

	session := auth.Logout(c.Request)
	session.Save(c.Request, c.Writer)

	c.Data(200, "", []byte(""))
}

func (s AuthService) DeleteUser(c *gin.Context) {
	auth, _ := s.getAuth(c)
	userId, _ := auth.CurrentUser(c.Request)
	s.DB.Table(models.User{}.TableName()).Where("id = ?", userId).Delete(&models.User{})
	c.Data(200, "", []byte("{}"))
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

func (s AuthService) createUserToken(email string) {
	token := authenticator.NewUserToken()
	hasher := md5.New()
	hasher.Write([]byte(token + email))
	apiToken := hex.EncodeToString(hasher.Sum(nil))
	s.DB.Table(models.User{}.TableName()).Where("email = ?", email).Update("api_token", apiToken)
}

func (s AuthService) setAgree(id string) {
	var user models.User

	err := s.DB.Table(models.User{}.TableName()).
		Where("id = ?", id).
		First(&user).Error

	if err == nil && user.AgreeWithTheTermsAndPolicyAt.IsZero() {
		s.DB.Table(models.User{}.TableName()).Where("id = ?", id).Update("agree_with_the_terms_and_policy_at", time.Now())
	}
}
