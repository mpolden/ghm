package git

import "testing"

func TestLocalDir(t *testing.T) {
	parentDir := "/foo/bar"
	repoName := "baz"

	got := LocalDir(parentDir, repoName)
	if want := "/foo/bar/baz.git"; got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
