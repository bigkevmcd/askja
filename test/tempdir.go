package test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// MakeTempDir creates a temporary directory that is automatically cleaned up
// after the test finishes.
func MakeTempDir(t *testing.T) string {
	t.Helper()
	dir, err := ioutil.TempDir("", "askja")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})
	return dir
}

// MakeTempRepository returns a git.Repository with an initial commit.
//
// The directory is deleted at the end of the test.
func MakeTempRepository(t *testing.T) *git.Repository {
	t.Helper()
	dir := MakeTempDir(t)
	r, err := git.PlainInit(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	w, err := r.Worktree()
	if err != nil {
		t.Fatal(err)
	}
	f, err := w.Filesystem.Create("README.md")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write([]byte("Just a test")); err != nil {
		t.Fatal(err)
	}
	_, err = w.Add("README.md")
	_, err = w.Commit("initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Testing",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// MakeTempGitRepo creates and returns a directory that has been initialised as a
// Git repository, it also returns rooted Filesystem that can be used to
// interact with the filesystem underneath.
//
// The directory is deleted at the end of the test.
func MakeTempGitRepo(t *testing.T) (string, billy.Filesystem) {
	t.Helper()
	dir := MakeTempDir(t)
	r, err := git.PlainInit(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	w, err := r.Worktree()
	if err != nil {
		t.Fatal(err)
	}
	f, err := w.Filesystem.Create("README.md")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write([]byte("Just a test")); err != nil {
		t.Fatal(err)
	}
	_, err = w.Add("README.md")
	_, err = w.Commit("initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Testing",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	return dir, osfs.New(dir)
}
