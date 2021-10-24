package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

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

type urlParamRuleset struct {
	BannedLabels     []string
	AnyOfTheseLabels []string
}

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
		http.Error(r, fmt.Sprintf("error parsing request: %s", err.Error()), 400)
	} else {
		if prEvent == nil { // we did not receive something that we can work with
			return
		}
		rules, err := gatherRules(req)
		if err != nil {
			http.Error(r, fmt.Sprintf("error parsing request: %s", err.Error()), 400)
		}
		_, ok := enforcer.ValidatePullRequest(prEvent.Repository.Name, prEvent.Number, rules)
		fmt.Fprintf(r, "Successful: %t", ok)
	}
}

func readRequestBody(req *http.Request) (*pullRequestEvent, error) {
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	var prEvent pullRequestEvent
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&prEvent); err != nil {
		return nil, fmt.Errorf("decoding into pullRequestEvent: %s", err)
	}
	if !prEvent.valid() {
		// if it's not a PR event, maybe it's a ping event?
		var pingEvnt pingEvent
		if err := json.NewDecoder(bytes.NewReader(body)).Decode(&pingEvnt); err != nil {
			return nil, fmt.Errorf("decoding into pingEvent, %s", err)
		} else if pingEvnt.valid() {
			fmt.Printf("Received ping event: %+v\n", pingEvnt)
			return nil, nil
		} else {
			return nil, fmt.Errorf("unexpected input: %s", string(body))
		}
	} else {
		return &prEvent, nil
	}
}

// Return the applicable RuleConfig. This can either come from a request param ('rules') or, as fallback,
// the default RuleConfig canonically provided by `enforcer.DefaultRules()`.
func gatherRules(req *http.Request) (*enforcer.RuleConfig, error) {
	if rules := req.URL.Query()[`rules`]; len(rules) > 0 {
		givenViaUrl := rules[0]
		var paramRules urlParamRuleset
		fmt.Printf("Decoding rule set: %s", givenViaUrl)
		if err := json.NewDecoder(strings.NewReader(givenViaUrl)).Decode(&paramRules); err != nil {
			return nil, fmt.Errorf("parsing URL param for rule set, %s", err)
		} else {
			return enforcer.NewRules(paramRules.BannedLabels, paramRules.AnyOfTheseLabels), nil
		}
	}
	fmt.Printf("Going with default rule set.")
	return enforcer.DefaultRules(), nil
}

func (e pingEvent) valid() bool {
	// TODO: do input validation here
	return e.Zen != ``
}

func (e pullRequestEvent) valid() bool {
	// TODO: do input validation here
	return e.Action != ``
}
