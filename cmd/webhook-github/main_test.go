package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func TestPing(t *testing.T) {
	// given
	handler := WebhookHandler()
	server := httptest.NewServer(handler)
	defer server.Close()
	e := httpexpect.New(t, server.URL)

	ping := map[string]interface{}{
		"zen": "Wisdom goes here",
		"repository": map[string]interface{}{
			"name": "github-pr-enforcer",
			"id":   12345,
		},
	}

	// when, then
	e.GET("/validate-pr").WithJSON(ping).
		Expect().
		Status(http.StatusOK).NoContent()
}
