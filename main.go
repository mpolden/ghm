package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/google/go-github/github"
	flags "github.com/jessevdk/go-flags"
	"github.com/mpolden/ghm/git"
	gh "github.com/mpolden/ghm/github"
)

type CLI struct {
	GitPath     string `short:"g" long:"git" description:"Path to git executable" value-name:"PATH" default:"git"`
	Quiet       bool   `short:"q" long:"quiet" description:"Only print errors"`
	Dryrun      bool   `short:"n" long:"dryrun" description:"Print commands that would be run and exit"`
	Protocol    string `short:"p" long:"protocol" description:"Use the given protocol when mirroring" choice:"ssh" choice:"https" choice:"git" default:"ssh"`
	SkipFork    bool   `short:"s" long:"skip-fork" description:"Skip forked repositories"`
	Concurrency int    `short:"c" long:"concurrency" description:"Mirror COUNT repositories concurrently" value-name:"COUNT" default:"1"`
	Args        struct {
		Username string `description:"GitHub username" positional-arg-name:"github-user"`
		Path     string `description:"Path where repositories should be mirrored" positional-arg-name:"path"`
	} `positional-args:"yes" required:"yes"`
	mu sync.Mutex
}

func (c *CLI) run(cmd *exec.Cmd) error {
	if c.Dryrun {
		// Prevent overlapping output
		c.mu.Lock()
		defer c.mu.Unlock()
		fmt.Println(strings.Join(cmd.Args, " "))
		return nil
	}
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (c *CLI) sync(g *git.Git, r *github.Repository) error {
	repoURL, err := gh.CloneURL(c.Protocol, r)
	if err != nil {
		return err
	}
	localDir := git.LocalDir(c.Args.Path, *r.Name)
	syncCmd := g.Sync(repoURL, localDir)
	if err := c.run(syncCmd); err != nil {
		return err
	}
	return nil
}

func (c *CLI) syncAll(g *git.Git, repos []*github.Repository) {
	sem := make(chan bool, c.Concurrency)
	for _, r := range repos {
		if c.SkipFork && *r.Fork {
			continue
		}
		sem <- true
		go func(r *github.Repository) {
			defer func() { <-sem }()
			if err := c.sync(g, r); err != nil {
				log.Fatal(err)
			}
		}(r)
	}
	// Wait for remaining goroutines to finish
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
}

func main() {
	var cli CLI
	_, err := flags.ParseArgs(&cli, os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

	if cli.Concurrency < 1 {
		fmt.Fprintln(os.Stderr, "concurrency level must be >= 1")
		os.Exit(1)
	}

	gh := gh.New()
	repos, err := gh.ListAllRepositories(cli.Args.Username)
	if err != nil {
		log.Fatal(err)
	}

	g, err := git.New(cli.GitPath, !cli.Quiet)
	if err != nil {
		log.Fatal(err)
	}

	cli.syncAll(g, repos)
}
