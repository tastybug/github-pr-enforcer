package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type PullRequest struct {
	labels []string
}

func main() {
	fmt.Println("Hi")
}

func fetchAndFilterByAge(noOlderThan time.Time) (*PullRequest, error) {
	// https://docs.github.com/en/rest/reference/pulls#get-a-pull-request
	searchTerms := []string{"repo:golang/go", "is:open", "json", "decoder"}

	ghUser := "tastybug"
	ghRepo := "gorki"
	fmt.Sprintf("https://github.com/repos/%s/%s/", url.PathEscape("tastybug"), url.PathEscape())
	url := "https://github.com/repos/?q=" + url.QueryEscape(strings.Join(searchTerms, ` `))
	resp, err := http.Get(url)
	if err != nil {
		resp.Body.Close()
		return nil, err
	}
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()

	var inRange []*Issue
	for _, issue := range result.Items {
		if issue.CreatedAt.After(noOlderThan) {
			inRange = append(inRange, issue)
		}
	}
	result.Items = inRange

	return &result, nil
}
