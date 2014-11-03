package main

import (
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	d "github.com/uhuraapp/uhura-api/database"
)

func PerformRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func fakeAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", "1")
	}
}

func request(method string, fn gin.HandlerFunc) *httptest.ResponseRecorder {
	r := gin.New()
	r.Use(fakeAuth())
	r.Handle(method, "/test", []gin.HandlerFunc{fn})
	return PerformRequest(r, method, "/test")
}

func databaseTest() gorm.DB {
	var database gorm.DB

	databaseUrl, _ := pq.ParseURL(os.Getenv("TEST_DATABASE_URL"))
	database, _ = gorm.Open("postgres", databaseUrl)
	database.LogMode(true)

	d.Migrations(database)

	return database
}
