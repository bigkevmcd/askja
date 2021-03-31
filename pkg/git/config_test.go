package git

import (
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/go-cmp/cmp"
)

func TestCommitOptionsFromConfig(t *testing.T) {
	cfg, err := parseGitConfig()
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()

	c, err := CommitOptionsFromConfig(now)
	if err != nil {
		t.Fatal(err)
	}

	want := &git.CommitOptions{
		Author: &object.Signature{
			Email: cfg.User.Email,
			Name:  cfg.User.Name,
			When:  now,
		},
	}
	if diff := cmp.Diff(want, c); diff != "" {
		t.Fatalf("failed to create the commit options:\n%s", diff)
	}
}

func Test_parseUserConfig(t *testing.T) {
	// TODO: Skip if no ~/.gitconfig
	c, err := parseGitConfig()
	if err != nil {
		t.Fatal(err)
	}

	if c.User.Email == "" || c.User.Name == "" {
		t.Fatal("failed to parse a user from the config")
	}
}
