package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tastybug/github-pr-enforcer/internal/enforcer/domain"
	"github.com/tastybug/github-pr-enforcer/internal/enforcer/service"
	"io/ioutil"
	"log"
	"net/http"
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

func gatherRules(req *http.Request) (*domain.RuleConfig, error) {
	var ruleConfiguration *domain.RuleConfig
	if rules := req.URL.Query()[`rules`]; len(rules) > 0 {
		// if there is a custom ruleset, collect that
		if config, err := service.GetRulesFromJsonString(rules[0]); err != nil {
			return nil, err
		} else {
			ruleConfiguration = config
		}
	}
	// let there be a decision on what ruleset to use
	return service.GetRules(ruleConfiguration), nil
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
		log.Printf("Validating %+v\n", innerPr)
		if violations, isValid := service.ValidatePr(innerPr, rules); !isValid {
			log.Printf("Result: %s is BAD (report: '%s')\n", innerPr.UID(), violations.String())
			fmt.Fprintf(resp, "%s invalid", innerPr.UID())
		} else {
			log.Printf("Result: %s is GOOD (report: '%s')\n", innerPr.UID(), violations.String())
			fmt.Fprintf(resp, "%s valid", innerPr.UID())
		}
		return nil
	}
}

func (p upstreamGhPrEvent) toInnerPr() *domain.PullRequest {

	innerPr := new(domain.PullRequest)
	for _, label := range p.PullRequest.Labels {
		innerPr.Labels = append(innerPr.Labels, domain.Label{
			Name: label.Name,
		})
	}
	innerPr.Number = p.Number
	innerPr.RepoName = p.Repository.FullName
	return innerPr
}
