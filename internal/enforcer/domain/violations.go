package domain

import (
	"fmt"
	"strings"
)

type Violations []string

func (v Violations) String() string {
	return strings.Join(v, `, `)
}

func (pr *PullRequest) UID() string {
	return fmt.Sprintf(`%s:%d`, pr.RepoName, pr.Number)
}
