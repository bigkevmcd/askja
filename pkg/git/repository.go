package git

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

// Repository is a struct that provides a simplified interface to a git
// repository.
type Repository struct {
	*git.Repository
	wt *git.Worktree
}

// New creates and returns a new Repository rooted at the provided base path.
func New(root string) (*Repository, error) {
	r, err := git.PlainOpen(root)
	if err != nil {
		return nil, fmt.Errorf("failed to open path %q as a git repository: %w", root, err)
	}
	w, err := r.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to create a Worktree: %w", err)
	}
	return &Repository{Repository: r, wt: w}, nil
}

// CreateAndSwitchBranch switches from the current branch to the
// one with the name provided.
func (r *Repository) CreateAndSwitchBranch(name string) error {
	branchRef := plumbing.NewBranchReferenceName(name)
	if err := r.Repository.CreateBranch(&config.Branch{
		Name:  name,
		Merge: branchRef,
	}); err != nil {
		return fmt.Errorf("failed to create branch %q: %w", name, err)
	}
	h, err := r.Head()
	if err != nil {
		return fmt.Errorf("failed to get the HEAD: %w", err)
	}
	ref := plumbing.NewHashReference(branchRef, h.Hash())
	err = r.Storer.SetReference(ref)
	if err != nil {
		return fmt.Errorf("failed to SetReference to %s: %w", ref, err)
	}
	if err := r.wt.Checkout(&git.CheckoutOptions{Branch: branchRef}); err != nil {
		return fmt.Errorf("failed to switch to branch %q: %w", name, err)
	}
	return nil
}

// WriteFile writes data to the named file, creating it if necessary.
// If the file does not exist, WriteFile creates it with permissions perm
// (before umask); otherwise WriteFile truncates it before writing, without
// changing permissions.
//
// The file is also "git added" to the current worktree.
func (r *Repository) WriteFile(name string, data []byte, perm os.FileMode) error {
	f, err := r.wt.Filesystem.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %w", name, err)
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data to file %q: %w", name, err)
	}
	_, err = r.wt.Add(name)
	if err != nil {
		return fmt.Errorf("failed to add file %q: %w", name, err)
	}
	return nil
}

// Commit creates a new commit in the git repository.
//
// It returns the sha of the commit.
func (r *Repository) Commit(msg string, opts *git.CommitOptions) (string, error) {
	c, err := r.wt.Commit(msg, opts)
	if err != nil {
		return "", fmt.Errorf("failed to commit: %w", err)
	}
	return c.String(), nil
}
