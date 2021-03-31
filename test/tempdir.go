package test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
)

// MakeTempGitDir creates and returns a directory that has been initialised as a
// Git repository, it also returns rooted Filesystem that can be used to
// interact with the filesystem underneath.
//
// The directory is deleted at the end of the test.
func MakeTempGitDir(t *testing.T) (string, billy.Filesystem) {
	t.Helper()
	dir, err := ioutil.TempDir("", "askja")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})
	// This is needed to ensure that the .git directory is created inside the
	// repo dir.
	_, err = git.PlainInit(filepath.Join(dir, ".git"), true)
	if err != nil {
		t.Fatal(err)
	}
	return dir, osfs.New(dir)
}
