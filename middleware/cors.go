package middleware

import (
	"net/url"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		referer, _ := url.Parse(origin)

		allow := referer.Scheme + "://" + referer.Host

		c.Writer.Header().Set("Access-Control-Allow-Origin", allow)
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS, GET, POST, PUT")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
	}
}
