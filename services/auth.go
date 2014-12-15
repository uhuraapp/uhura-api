package services

import (
	"github.com/dukex/login2"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type AuthService struct {
	DB    gorm.DB
	login *login2.Builder
}

func NewAuthService(db gorm.DB) AuthService {
	login := login2.NewBuilder()

	return AuthService{DB: db, login: login}
}

func (s AuthService) ByProvider(c *gin.Context) {

}
