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
type upstreamGhPingEvent struct {
	Zen        string `json:"zen"`
	Repository struct {
		Name string `json:"name"`
		Id   int    `json:"id"`
	} `json:"repository"`
}

// https://docs.github.com/en/developers/webhooks-and-events/webhooks/webhook-events-and-payloads#pull_request
type upstreamGhPrEvent struct {
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
	mux.HandleFunc("/validate-pr", be.handleGithubEvent)
	if err := http.ListenAndServe(hostPort, mux); err != nil {
		log.Fatal(err)
	}
}

func (b *backend) handleGithubEvent(resp http.ResponseWriter, req *http.Request) {
	if err := extractAndProcess(req, resp); err != nil {
		http.Error(resp, fmt.Sprintf("error handling request: %s", err.Error()), 400)
	}
}

func extractAndProcess(req *http.Request, r http.ResponseWriter) error {
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	var pr upstreamGhPrEvent
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&pr); err != nil {
		return fmt.Errorf("decoding into pullRequestEvent: %s", err)
	}
	if !pr.valid() {
		// if it's not a PR event, maybe it's a ping event?
		var ping upstreamGhPingEvent
		if err := json.NewDecoder(bytes.NewReader(body)).Decode(&ping); err != nil {
			return fmt.Errorf("decoding into pingEvent, %s", err)
		} else if ping.valid() {
			fmt.Printf("Received ping event: %+v\n", ping)
			return ping.process(req, r)
		} else {
			return fmt.Errorf("unexpected input: %s", string(body))
		}
	} else {
		return pr.process(req, r)
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
			return nil, fmt.Errorf("given rule set broken: %s", err)
		} else {
			return enforcer.NewRules(paramRules.BannedLabels, paramRules.AnyOfTheseLabels), nil
		}
	}
	fmt.Printf("Going with default rule set.")
	return enforcer.DefaultRules(), nil
}

func (e upstreamGhPingEvent) valid() bool {
	// TODO: do input validation here
	return e.Zen != ``
}

func (p upstreamGhPingEvent) process(req *http.Request, resp http.ResponseWriter) error {
	if _, err := gatherRules(req); err != nil {
		return err
	}
	return nil
}

func (e upstreamGhPrEvent) valid() bool {
	// TODO: do input validation here
	return e.Action != ``
}

func (p upstreamGhPrEvent) process(req *http.Request, resp http.ResponseWriter) error {
	if rules, err := gatherRules(req); err != nil {
		return err
	} else {
		_, ok := enforcer.ValidatePullRequest(p.Repository.Name, p.Number, rules)
		fmt.Fprintf(resp, "Successful: %t", ok)
		return nil
	}
}
