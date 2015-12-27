package services

import (
	"log"
	"strings"

	"bitbucket.org/dukex/uhura-api/entities"
	"bitbucket.org/dukex/uhura-api/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type UserService struct {
	DB gorm.DB
}

func NewUserService(db gorm.DB) UserService {
	return UserService{DB: db}
}

func (s UserService) All(c *gin.Context) {
	var users []entities.User

	email := c.Request.URL.Query().Get("email")

	query := s.DB.Table(models.User{}.TableName())

	if email != "" {
		query = query.Where("email = ?", email)
	}

	err := query.Find(&users).Error

	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	c.JSON(200, users)
}

func (s UserService) Get(c *gin.Context) {
	var user entities.User
	log.Println(c.Params)
	userId := strings.Replace(c.Params.ByName("id"), ".json", "", 1)

	err := s.DB.Table(models.User{}.TableName()).Where("id = ?", userId).First(&user).Error
	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, gin.H{"user": user})
}
