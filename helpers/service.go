package helpers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) (int, error) {
	_userId, _ := c.Get("user_id")
	userId, err := strconv.Atoi(_userId.(string))
	return userId, err
}
