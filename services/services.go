package services

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/uhuraapp/uhura-api/database"
)

type Service struct {
}

func (s *Service) connection(c *gin.Context) *gorm.DB {
	connection, _ := c.Get("connection")
	return connection.(*gorm.DB)
}

func Connection() gin.HandlerFunc {
	connection := database.NewPostgresql()

	return func(c *gin.Context) {
		c.Set("connection", connection)
		c.Next()
	}
}
