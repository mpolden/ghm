package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
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

func (s *syncer) sync(g *git, r repository) error {
	var url string
	switch s.protocol {
	case "ssh":
		url = r.SSHURL
	case "git":
		url = r.GitURL
	default:
		url = r.CloneURL
	}
	localDir := repositoryPath(s.localPath, r.Name)
	syncCmd := g.sync(url, localDir)
	return s.run(syncCmd)
}

func (s *syncer) syncAll(g *git, repos []repository) {
	sem := make(chan bool, s.concurrency)
	for _, r := range repos {
		if s.skipFork && r.Fork {
			continue
		}
		if s.skipArchived && r.Archived {
			continue
		}
		sem <- true
		go func(r repository) {
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
	protocol := flag.String("p", "ssh", "Protocol to use for mirroring [git|https|ssh]")
	skipFork := flag.Bool("s", false, "Skip forked repositories")
	skipArchived := flag.Bool("a", false, "Skip archived repositories")
	concurrency := flag.Int("c", 1, "Number of repositories to mirror concurrently")
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		return
	}
	if *concurrency < 1 {
		log.Fatal("invalid concurrency level")
	}

	g, err := gitCommand(!*quiet)
	if err != nil {
		log.Fatal(err)
	}
	githubUser := args[0]
	path := args[1]
	syncer := syncer{
		dryrun:       *dryrun,
		protocol:     *protocol,
		skipFork:     *skipFork,
		skipArchived: *skipArchived,
		concurrency:  *concurrency,
		localPath:    path,
	}
	repos, err := listAllRepositories(githubUser)
	if err != nil {
		log.Fatal(err)
	}
	syncer.syncAll(g, repos)
}
