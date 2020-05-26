package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/google/go-github/v28/github"

	"github.com/mpolden/ghm/git"
	gh "github.com/mpolden/ghm/github"
)

type syncer struct {
	quiet        bool
	protocol     string
	dryrun       bool
	concurrency  int
	skipFork     bool
	skipArchived bool
	localPath    string
	mu           sync.Mutex
}

func (s *syncer) run(cmd *exec.Cmd) error {
	if s.dryrun {
		// Prevent overlapping output
		s.mu.Lock()
		defer s.mu.Unlock()
		fmt.Println(strings.Join(cmd.Args, " "))
		return nil
	}
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (s *syncer) sync(g *git.Git, r *github.Repository) error {
	repoURL, err := gh.CloneURL(s.protocol, r)
	if err != nil {
		return err
	}
	localDir := git.LocalDir(s.localPath, *r.Name)
	syncCmd := g.Sync(repoURL, localDir)
	return s.run(syncCmd)
}

func (s *syncer) syncAll(g *git.Git, repos []*github.Repository) {
	sem := make(chan bool, s.concurrency)
	for _, r := range repos {
		if s.skipFork && *r.Fork {
			continue
		}
		if s.skipArchived && *r.Archived {
			continue
		}
		sem <- true
		go func(r *github.Repository) {
			defer func() { <-sem }()
			if err := s.sync(g, r); err != nil {
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
	log.SetPrefix("ghm: ")
	log.SetFlags(log.Lshortfile)

	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(out, "\nArguments:\n  <github-user>\tGitHub username")
		fmt.Fprintln(out, "  <path>\tPath where repositories should be mirrored")
	}
	quiet := flag.Bool("q", false, "Only print errors")
	dryrun := flag.Bool("n", false, "Print commands that would be run and exit")
	protocol := flag.String("p", "ssh", "Protocol to use for mirroring [ssh|https|git]")
	skipFork := flag.Bool("s", false, "Skip forked repositories")
	skipArchived := flag.Bool("a", false, "Skip archived repositories")
	concurrency := flag.Int("c", 1, "Number of repositories to mirror concurrently")
	flag.Parse()

	if *concurrency < 1 {
		log.Fatal("concurrency level must be positive")
	}

	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
	}
	username := args[0]
	path := args[1]

	gh := gh.New()
	repos, err := gh.ListAllRepositories(username)
	if err != nil {
		log.Fatal(err)
	}

	g, err := git.New(!*quiet)
	if err != nil {
		log.Fatal(err)
	}

	syncer := syncer{
		dryrun:       *dryrun,
		protocol:     *protocol,
		skipFork:     *skipFork,
		skipArchived: *skipArchived,
		concurrency:  *concurrency,
		localPath:    path,
	}
	syncer.syncAll(g, repos)
}
