package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stvp/rollbar"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// Recovery from error and send to rollbar
func Recovery() gin.HandlerFunc {
	rollbar.Token = os.Getenv("ROLLBAR_ACCESS_TOKEN")
	rollbar.Environment = os.Getenv("ROLLBAR_ENV")

	return recovery()
}

func recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("ENABLE_ROLLBAR") == "true" {
			defer func() {
				if r := recover(); r != nil {
					err, _ := r.(error)
					rollbar.ErrorWithStackSkip(rollbar.ERR, err, 5)
					c.AbortWithStatus(500)
					rollbar.Wait()
				}
			}()
		}
		c.Next()
	}
}
