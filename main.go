package main

import (
	"log"
	"os"
	"path/filepath"

	flags "github.com/jessevdk/go-flags"
	"github.com/martinp/git-mirror/git"
	"github.com/martinp/git-mirror/github"
)

type CLI struct {
	GitPath string `short:"g" long:"git" description:"Path to git executable" default:"git"`
	Quiet   bool   `short:"q" long:"quiet" description:"Only print errors"`
	Args    struct {
		Username string `description:"GitHub username" positional-arg-name:"USER"`
		Path     string `description:"Path where repositories should be mirrored" positional-arg-name:"PATH"`
	} `positional-args:"yes" required:"yes"`
}

func (c *CLI) Log(format string, v ...interface{}) {
	if !c.Quiet {
		log.Printf(format, v...)
	}
}

func (c *CLI) localDir(name string) string {
	return filepath.Join(c.Args.Path, name+".git")
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
		localDir := cli.localDir(*r.Name)
		if _, err := os.Stat(localDir); os.IsNotExist(err) {
			cli.Log("%s does not exist, mirroring", localDir)
			if err := g.Mirror(*r.SSHURL, localDir); err != nil {
				log.Fatal(err)
			}
		} else {
			cli.Log("%s already exists, updating", localDir)
			if err := g.Update(localDir); err != nil {
				log.Fatal(err)
			}
		}
	}
}
