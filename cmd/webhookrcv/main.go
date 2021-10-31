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
	Number     int `json:"number"`
	Repository struct {
		FullName string `json:"full_name"`
		Id       int    `json:"id"`
	} `json:"repository"`
	PullRequest struct {
		Labels []struct {
			Name string `json:"name"`
		} `json:"labels"`
	} `json:"pull_request"`
}

type urlParamRuleset struct {
	BannedLabels     []string
	AnyOfTheseLabels []string
}

func main() {
	fmt.Printf("Starting server on %s..\n", hostPort)

	if err := http.ListenAndServe(hostPort, WebhookHandler()); err != nil {
		log.Fatal(err)
	}
}

func WebhookHandler() http.Handler {
	be := new(backend)

	mux := http.NewServeMux()
	mux.HandleFunc("/validate-pr", be.handleGithubEvent)

	return mux
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
		return fmt.Errorf("decoding into upstreamGhPrEvent: %s", err)
	}
	if !pr.valid() {
		// if it's not a PR event, maybe it's a ping event?
		var ping upstreamGhPingEvent
		if err := json.NewDecoder(bytes.NewReader(body)).Decode(&ping); err != nil {
			return fmt.Errorf("decoding into pingEvent, %s", err)
		} else if ping.valid() {
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
	fmt.Println("Going with default rule set.")
	return enforcer.DefaultRules(), nil
}

func (e upstreamGhPingEvent) valid() bool {
	// TODO: do input validation here
	return e.Zen != ``
}

func (p upstreamGhPingEvent) process(req *http.Request, resp http.ResponseWriter) error {
	fmt.Printf("Received ping event: %+v\n", p)
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
		innerPr := p.toInnerPr()
		fmt.Printf("Checking %+v\n", innerPr)
		_, ok := enforcer.IsValidPr(innerPr, rules)
		fmt.Printf("Result: %s/%d is ok=%t\n", innerPr.RepoName, innerPr.Number, ok)

		fmt.Fprintf(resp, "Successful: %t", ok)
		return nil
	}
}

func (p upstreamGhPrEvent) toInnerPr() *enforcer.InternalPullRequest {

	innerPr := new(enforcer.InternalPullRequest)
	for _, label := range p.PullRequest.Labels {
		innerPr.Labels = append(innerPr.Labels, enforcer.InternalLabel{
			Name: label.Name,
		})
	}
	innerPr.Number = p.Number
	innerPr.RepoName = p.Repository.FullName
	return innerPr
}
