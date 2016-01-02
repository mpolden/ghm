package git

import (
	"os"
	"os/exec"
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
	cmd := g.command("fetch", "--prune")
	cmd.Dir = localDir
	return cmd
}
