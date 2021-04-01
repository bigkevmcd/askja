package test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/utils/merkletrie"
)

// GetFilesInCommit returns a map of the files committed in the provided commit.
func GetFilesInCommit(t *testing.T, commit *object.Commit, base string) map[string][]byte {
	t.Helper()
	found := map[string][]byte{}
	currentDirState, err := commit.Tree()
	if err != nil {
		t.Fatalf("failed to get tree for commit: %s", err)
	}

	prevCommitObject, err := commit.Parents().Next()
	if err != nil && err != io.EOF {
		t.Fatalf("failed to get the previousCommit: %s", err)
	}

	if prevCommitObject == nil {
		files := currentDirState.Files()
		defer files.Close()
		files.ForEach(func(f *object.File) error {
			b, err := os.ReadFile(f.Name)
			if err != nil {
				t.Fatalf("failed to read file %s: %w", f.Name, err)
			}
			found[f.Name] = b
			return nil
		})
	}

	if prevCommitObject != nil {
		prevDirState, err := prevCommitObject.Tree()
		if err != nil {
			t.Fatalf("could not get tree from previous commit: %s", err)
		}
		changes, err := prevDirState.Diff(currentDirState)
		if err != nil {
			t.Fatalf("failed to get previous dir state diff: %s", err)
		}
		for _, ch := range changes {
			action, err := ch.Action()
			if err != nil {
				t.Fatalf("could not get the action for change %s: %s", ch, err)
			}
			if action == merkletrie.Insert {
				b, err := os.ReadFile(filepath.Join(base, ch.To.Name))
				if err != nil {
					t.Fatalf("failed to read file %s: %w", ch.To.Name, err)
				}
				found[ch.To.Name] = b
			}
		}
	}
	return found
}
