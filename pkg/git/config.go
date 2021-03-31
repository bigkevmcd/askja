package git

import (
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// CommitOptionsFromConfig creates and returns a *git.CommitOptions
// with the Author field populated based on the user's ~/.gitconfig.
func CommitOptionsFromConfig(when time.Time) (*git.CommitOptions, error) {
	cfg, err := parseGitConfig()
	if err != nil {
		return nil, err
	}

	return &git.CommitOptions{
		Author: &object.Signature{
			Email: cfg.User.Email,
			Name:  cfg.User.Name,
			When:  when,
		},
	}, nil
}

func parseGitConfig() (*config.Config, error) {
	cfg := config.NewConfig()

	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("failed to get the current user: %w", err)
	}

	configPath := filepath.Join(usr.HomeDir, "/.gitconfig")
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %q: %w", configPath, err)
	}

	if err := cfg.Unmarshal(b); err != nil {
		return nil, fmt.Errorf("failed to unmarshal gitconfig in %q: %w", configPath, err)
	}

	return cfg, nil
}
