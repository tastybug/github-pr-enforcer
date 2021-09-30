package enforcer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Label struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Default     bool   `json:"default"`
}

type PullRequest struct {
	Labels []Label `json:"labels"`
}

func fetchPrViaFullName(repoFullName, prNumber string) (*PullRequest, error) {

	// path elements should already be safe, but better be safe here and escape it
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/pulls/%s",
		repoFullName,
		url.PathEscape(prNumber))

	fmt.Println(url)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		resp.Body.Close()
		return nil, err
	}
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	var result PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()

	return &result, nil
}
