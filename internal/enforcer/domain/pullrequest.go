package domain

type PullRequest struct {
	RepoName string
	Number   int
	Labels   []Label
}
