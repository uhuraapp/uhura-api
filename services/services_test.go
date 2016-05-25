package services_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	d "github.com/uhuraapp/uhura-api/database"
	"github.com/uhuraapp/uhura-api/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

func PerformRequest(r http.Handler, method, path, since string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	req.Header.Set("If-Modified-Since", since)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func fakeAuth(userId string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", userId)
	}
}

func request(method string, fn gin.HandlerFunc, since string) *httptest.ResponseRecorder {
	r := gin.New()
	r.Use(fakeAuth("1"))
	r.Handle(method, "/test", fn)
	if since == "" {
		now := time.Now()
		since = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location()).Format(time.RFC1123)
	}
	return PerformRequest(r, method, "/test", since)
}

func botRequest(method string, fn gin.HandlerFunc) *httptest.ResponseRecorder {
	r := gin.New()
	r.Use(fakeAuth("bot"))
	r.Handle(method, "/test", fn)
	return PerformRequest(r, method, "/test", "")
}

func databaseTest() *gorm.DB {
	var database *gorm.DB

	databaseUrl, _ := pq.ParseURL(os.Getenv("TEST_DATABASE_URL"))
	database, _ = gorm.Open("postgres", databaseUrl)
	database.LogMode(os.Getenv("DEBUG") == "true")

	d.Migrations(database)

	return database
}

func resetDatabase() {
	db := databaseTest()
	db.Delete(models.Episode{})
	db.Delete(models.Listened{})
	db.Delete(models.Subscription{})
	db.Delete(models.Channel{})
}
