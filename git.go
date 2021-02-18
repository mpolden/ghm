package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type git struct {
	path      string
	inheritIO bool
}

type repository struct {
	Name     string `json:"name"`
	SSHURL   string `json:"ssh_url"`
	GitURL   string `json:"git_url"`
	CloneURL string `json:"clone_url"`
	Fork     bool   `json:"fork"`
	Archived bool   `json:"archived"`
}

func listRepositories(user string) ([]repository, error) {
	res, err := http.Get("https://api.github.com/users/" + user + "/repos")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	dec := json.NewDecoder(res.Body)
	var repos []repository
	if err := dec.Decode(&repos); err != nil {
		return nil, err
	}
	return repos, nil
}

func gitCommand(inheritIO bool) (*git, error) {
	p, err := exec.LookPath("git")
	if err != nil {
		return nil, err
	}
	return &git{path: p, inheritIO: inheritIO}, nil
}

func repositoryPath(parentDir, repoName string) string {
	return filepath.Join(parentDir, repoName+".git")
}

func (g *git) command(args ...string) *exec.Cmd {
	cmd := exec.Command(g.path, args...)
	if g.inheritIO {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd
}

func (g *git) sync(url, path string) *exec.Cmd {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return g.command("clone", "--mirror", url, path)
	} else if err != nil {
		log.Fatal(err)
	}
	return g.command("-C", path, "fetch", "--prune")
}
