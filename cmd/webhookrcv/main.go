package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/tastybug/github-pr-enforcer/internal/enforcer"
)

const hostPort = `0.0.0.0:9000`

type backend struct{}

type requestBody struct {
	GhUser string `json:"ghUser"`
	GhRepo string `json:"ghRepo"`
	GhPrNo string `json:"ghPrNo"`
}

var sharedConfig = enforcer.NewRules(
	[]string{"wip", "do-not-merge"},
	[]string{"bug", "feature", "enabler", "rework"},
)

func main() {
	be := new(backend)

	fmt.Printf("Starting server on %s..\n", hostPort)

	mux := http.NewServeMux()
	mux.HandleFunc("/validate-pr", be.serveResult)
	if err := http.ListenAndServe(hostPort, mux); err != nil {
		log.Fatal(err)
	}
}

func (b *backend) serveResult(r http.ResponseWriter, req *http.Request) {

	if request, err := readRequestBody(req); err != nil {
		http.Error(r, fmt.Sprintf("error understanding request: %s", err.Error()), 400)
	} else {
		_, ok := enforcer.IsValid(request.GhUser, request.GhRepo, request.GhPrNo, sharedConfig)
		fmt.Fprintf(r, "Successful: %t", ok)
	}
}

func readRequestBody(req *http.Request) (*requestBody, error) {
	defer req.Body.Close()
	var result requestBody
	if err := json.NewDecoder(req.Body).Decode(&result); err != nil {
		return nil, err
	}
	// TODO here validate the request
	return &result, nil
}
