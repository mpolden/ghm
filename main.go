package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/martinp/ghm/git"
	"github.com/martinp/ghm/github"
)

type CLI struct {
	GitPath  string `short:"g" long:"git" description:"Path to git executable" value-name:"PATH" default:"git"`
	Quiet    bool   `short:"q" long:"quiet" description:"Only print errors"`
	Dryrun   bool   `short:"n" long:"dryrun" description:"Print commands that would be run and exit"`
	Protocol string `short:"p" long:"protocol" description:"Use the given protocol when mirroring" choice:"ssh" choice:"https" choice:"git" default:"ssh"`
	SkipFork bool   `short:"s" long:"skip-fork" description:"Skip forked repositories"`
	Args     struct {
		Username string `description:"GitHub username" positional-arg-name:"github-user"`
		Path     string `description:"Path where repositories should be mirrored" positional-arg-name:"path"`
	} `positional-args:"yes" required:"yes"`
}

func (c *CLI) Run(cmd *exec.Cmd) error {
	if c.Dryrun {
		fmt.Println(strings.Join(cmd.Args, " "))
		return nil
	}
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func main() {
	var cli CLI
	_, err := flags.ParseArgs(&cli, os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

	gh := github.New()
	repos, err := gh.ListAllRepositories(cli.Args.Username)
	if err != nil {
		log.Fatal(err)
	}

	g, err := git.New(cli.GitPath, !cli.Quiet)
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range repos {
		if cli.SkipFork && *r.Fork {
			continue
		}
		repoURL, err := github.CloneURL(cli.Protocol, r)
		if err != nil {
			log.Fatal(err)
		}
		localDir := git.LocalDir(cli.Args.Path, *r.Name)
		if err := cli.Run(g.Sync(repoURL, localDir)); err != nil {
			log.Fatal(err)
		}
	}
}
