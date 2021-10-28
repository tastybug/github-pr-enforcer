package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func TestPing(t *testing.T) {
	// create http.Handler
	handler := WebhookHandler()

	// run server using httptest
	server := httptest.NewServer(handler)
	defer server.Close()

	// create httpexpect instance
	e := httpexpect.New(t, server.URL)

	ping := map[string]interface{}{
		"zen": "Wisdom goes here",
		"repository": map[string]interface{}{
			"name": "github-pr-enforcer",
			"id":   12345,
		},
	}

	e.GET("/validate-pr").WithJSON(ping).
		Expect().
		Status(http.StatusOK).NoContent()
}
