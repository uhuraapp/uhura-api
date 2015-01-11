package middleware

import (
	"net/url"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		referer, _ := url.Parse(c.Request.Referer())

		c.Writer.Header().Set("Access-Control-Allow-Origin", referer.Scheme+"://"+referer.Host)
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS, GET, POST")
	}
}
