package git

import (
	"fmt"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

// Repository is a struct that provides a simplified interface to a git
// repository.
type Repository struct {
	*git.Repository
}

// New creates and returns a new Repository rooted at the provided base path.
func New(p, branch string) (*Repository, error) {
	r, err := git.PlainOpen(p)
	if err != nil {
		return nil, fmt.Errorf("failed to open path %q as a git repository: %w", p, err)
	}
	branchRef := plumbing.NewBranchReferenceName(branch)
	// // If we get an error trying to get a Head reference this is likely an empty
	// // repository, with no commits, so it creates a new reference.
	var h *plumbing.Reference
	headRef, err := r.Head()
	if err != nil {
		if err != plumbing.ErrReferenceNotFound {
			return nil, fmt.Errorf("failed to get the HEAD for the repository: %w", err)
		}
		h = plumbing.NewSymbolicReference(plumbing.HEAD, branchRef)
	} else {
		h = plumbing.NewHashReference(branchRef, headRef.Hash())
	}

	if err := r.Storer.SetReference(h); err != nil {
		return nil, fmt.Errorf("failed to store the branch %q: %w", branch, err)
	}

	// TODO: extract to a function
	_, err = r.Branch(branch)
	if err != nil {
		if err != git.ErrBranchNotFound {
			return nil, fmt.Errorf("failed to query for branch %q: %w", branch, err)
		}
		err = r.CreateBranch(&config.Branch{
			Name:  branch,
			Merge: branchRef,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create branch %q: %w", branch, err)
		}
	}

	w, err := r.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to create a Worktree: %w", err)
	}

	log.Printf("KEVIN!!!!!! %s and %s but %s", headRef, branchRef, h)

	if err := w.Checkout(&git.CheckoutOptions{Branch: branchRef}); err != nil {
		return nil, fmt.Errorf("failed to checkout the branch %q (%s): %w", branch, branchRef, err)
	}
	return &Repository{Repository: r}, nil
}

// WriteFile writes data to the named file, creating it if necessary.
// If the file does not exist, WriteFile creates it with permissions perm
// (before umask); otherwise WriteFile truncates it before writing, without
// changing permissions.
//
// The file is also "git added" to the current worktree.
func (r *Repository) WriteFile(name string, data []byte, perm os.FileMode) error {
	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("failed to create worktree: %w", err)
	}
	f, err := w.Filesystem.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %w", name, err)
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data to file %q: %w", name, err)
	}
	_, err = w.Add(name)
	if err != nil {
		return fmt.Errorf("failed to add file %q: %w", name, err)
	}
	return nil
}

// Commit creates a new commit in the git repository.
//
// It returns the sha of the commit.
func (r *Repository) Commit(msg string, opts *git.CommitOptions) (string, error) {
	w, err := r.Worktree()
	if err != nil {
		return "", fmt.Errorf("failed to get a Worktree: %w", err)
	}
	c, err := w.Commit(msg, opts)
	if err != nil {
		return "", fmt.Errorf("failed to commit: %w", err)
	}
	return c.String(), nil
}
