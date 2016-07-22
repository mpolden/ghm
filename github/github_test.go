package github

import (
	"testing"

	"github.com/google/go-github/github"
)

func TestCloneURL(t *testing.T) {
	httpsURL := "https://github.com/octocat/Hello-World.git"
	gitURL := "git://github.com/octocat/Hello-World.git"
	sshURL := "git@github.com:octocat/Hello-World.git"
	r := github.Repository{
		CloneURL: &httpsURL,
		GitURL:   &gitURL,
		SSHURL:   &sshURL,
	}
	var tests = []struct {
		in  string
		out string
	}{
		{"https", httpsURL},
		{"git", gitURL},
		{"ssh", sshURL},
	}
	for _, tt := range tests {
		got, err := CloneURL(tt.in, &r)
		if err != nil {
			t.Fatal(err)
		}
		if got != tt.out {
			t.Errorf("got %s for %s, want %s", got, tt.in, tt.out)
		}
	}
	if _, err := CloneURL("", &r); err == nil {
		t.Error("want error")
	}
}
