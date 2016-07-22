package github

import (
	"fmt"

	"github.com/google/go-github/github"
)

type Client struct {
	*github.Client
}

func New() *Client {
	return &Client{github.NewClient(nil)}
}

func (c *Client) ListAllRepositories(username string) ([]*github.Repository, error) {
	opt := &github.RepositoryListOptions{Type: "owner", ListOptions: github.ListOptions{PerPage: 100}}
	var repos []*github.Repository
	for done := false; !done; {
		rs, response, err := c.Repositories.List(username, opt)
		if err != nil {
			return nil, err
		}
		repos = append(repos, rs...)
		opt.ListOptions.Page = response.NextPage
		done = response.NextPage == 0
	}
	return repos, nil
}

func CloneURL(protocol string, r *github.Repository) (string, error) {
	switch protocol {
	case "https":
		return *r.CloneURL, nil
	case "git":
		return *r.GitURL, nil
	case "ssh":
		return *r.SSHURL, nil
	}
	return "", fmt.Errorf("unknown protocol: %s", protocol)
}
