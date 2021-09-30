package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tastybug/github-pr-enforcer/internal/enforcer"
)

const hostPort = `0.0.0.0:9000`

type backend struct{}

// after setting up a new webhook receiver, there will be a ping event sent first
// https://docs.github.com/en/developers/webhooks-and-events/webhooks/webhook-events-and-payloads#ping
type pingEvent struct {
	Zen        string `json:"zen"`
	Repository struct {
		Name string `json:"name"`
		Id   int    `json:"id"`
	} `json:"repository"`
}

// https://docs.github.com/en/developers/webhooks-and-events/webhooks/webhook-events-and-payloads#pull_request
type pullRequestEvent struct {
	Action string `json:"action"`
	// pull request number
	Number     string `json:"number"`
	Repository struct {
		Name string `json:"name"`
		Id   int    `json:"id"`
	} `json:"repository"`
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

	if prEvent, err := readRequestBody(req); err != nil {
		http.Error(r, fmt.Sprintf("error understanding request: %s", err.Error()), 400)
	} else {
		if prEvent == nil { // we did not receive something that we can work with
			return
		}
		_, ok := enforcer.IsValid(prEvent.Repository.Name, prEvent.Number, sharedConfig)
		fmt.Fprintf(r, "Successful: %t", ok)
	}
}

func readRequestBody(req *http.Request) (*pullRequestEvent, error) {
	defer req.Body.Close()
	b := new(bytes.Buffer)
	if body, err := ioutil.ReadAll(req.Body); err != nil {
		return nil, err
	} else {
		if _, err := b.Write(body); err != nil {
			return nil, err
		} else {
			var prEvent pullRequestEvent
			if err := json.NewDecoder(b).Decode(&prEvent); err != nil {
				return nil, err
			}
			if !prEvent.valid() {
				var pingEvnt pingEvent
				if err := json.NewDecoder(b).Decode(&pingEvnt); err != nil {
					return nil, err
				} else if pingEvnt.valid() {
					fmt.Printf("Received ping event: %+v\n", pingEvnt)
					return nil, nil
				} else {
					return nil, fmt.Errorf("unexpected input received and discarded: %s", b.String())
				}
			} else {
				return &prEvent, nil
			}
		}
	}
}

func (e pingEvent) valid() bool {
	// TODO: do input validation here
	return e.Zen != ``
}

func (e pullRequestEvent) valid() bool {
	// TODO: do input validation here
	return e.Action != ``
}
