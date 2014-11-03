package main

import (
	"net/http"
	"testing"

	"github.com/uhuraapp/uhura-api/services"
)

func TestSubscriptionsGetIsOk(t *testing.T) {
	s := services.NewSubscriptionService(databaseTest())
	r := request("GET", s.Get)

	if r.Code != http.StatusOK {
		t.Errorf("Status code should be %v, was %d", http.StatusOK, r.Code)
	}
}
