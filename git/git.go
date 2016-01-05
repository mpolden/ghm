package git

import (
	"os"
	"os/exec"
	"path/filepath"
)

type Git struct {
	path      string
	inheritIO bool
}

func New(path string, inheritIO bool) (*Git, error) {
	p, err := exec.LookPath(path)
	if err != nil {
		return nil, err
	}
	return &Git{path: p, inheritIO: inheritIO}, nil
}

func LocalDir(parentDir, repoName string) string {
	return filepath.Join(parentDir, repoName+".git")
}

func (g *Git) command(args ...string) *exec.Cmd {
	cmd := exec.Command(g.path, args...)
	if g.inheritIO {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd
}

func (g *Git) Mirror(repoURL, localDir string) *exec.Cmd {
	return g.command("clone", "--mirror", repoURL, localDir)
}

func (g *Git) Update(localDir string) *exec.Cmd {
	return g.command("-C", localDir, "fetch", "--prune")
}

func (g *Git) Sync(repoURL, localDir string) *exec.Cmd {
	if _, err := os.Stat(localDir); os.IsNotExist(err) {
		return g.Mirror(repoURL, localDir)
	}
	return g.Update(localDir)
}
