package git

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/utils/merkletrie"
	"github.com/google/go-cmp/cmp"
)

const (
	testFilename = "path/to/testing.yaml"
	testBranch   = "main"
)

func TestWriteFile(t *testing.T) {
	tmpDir, bfs := makeTempDir(t)
	g, err := New(tmpDir, testBranch)
	if err != nil {
		t.Fatal(err)
	}

	if err := g.WriteFile(testFilename, []byte(`testing: value\n`), 0644); err != nil {
		t.Fatal(err)
	}

	src, err := bfs.Open(testFilename)
	if err != nil {
		t.Fatal(err)
	}
	defer src.Close()
	b, err := ioutil.ReadAll(src)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff([]byte(`testing: value\n`), b); diff != "" {
		t.Fatalf("file written didn't match:\n%s", diff)
	}

	fi, err := bfs.Stat(testFilename)
	if err != nil {
		t.Fatal(err)
	}
	if m := fi.Mode(); m != 0644 {
		t.Fatalf("file mode got %#o, want 0644", m)
	}
}

func TestCommit(t *testing.T) {
	tmpDir, _ := makeTempDir(t)
	g, err := New(tmpDir, testBranch)
	if err != nil {
		t.Fatal(err)
	}
	otherFilename := "test/path.yaml"
	if err := g.WriteFile(testFilename, []byte(`testing: value\n`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := g.WriteFile(otherFilename, []byte(`testing: value\n`), 0600); err != nil {
		t.Fatal(err)
	}

	opts := &git.CommitOptions{
		Author: &object.Signature{
			Email: "test@example.com",
			Name:  "Testing",
			When:  time.Now(),
		},
	}
	sha, err := g.Commit("test commit", opts)
	if err != nil {
		t.Fatal(err)
	}
	if sha == "" {
		t.Fatalf("commit did not return a SHA")
	}
	commit, err := g.CommitObject(plumbing.NewHash(sha))
	if err != nil {
		t.Fatal(err)
	}
	want := []string{testFilename, otherFilename}
	assertFilesCommitted(t, commit, want)

	if commit.Message != "test commit" {
		t.Errorf("failed to commit with the correct message, got %q, want %q", commit.Message, "test commit")
	}

	if commit.Author.Name != opts.Author.Name {
		t.Errorf("got commit Author.Name %#v, want %#v", commit.Author.Name, opts.Author.Name)
	}
	if commit.Author.Email != opts.Author.Email {
		t.Errorf("got commit Author.Email %#v, want %#v", commit.Author.Email, opts.Author.Email)
	}
}

func TestCommitNewBranch(t *testing.T) {
	tmpDir, _ := makeTempDir(t)
	opts := &git.CommitOptions{
		Author: &object.Signature{
			Email: "test@example.com",
			Name:  "Testing",
			When:  time.Now(),
		},
	}

	g, err := New(tmpDir, testBranch)
	if err != nil {
		t.Fatal(err)
	}
	if err := g.WriteFile(testFilename, []byte(`testing: value\n`), 0644); err != nil {
		t.Fatal(err)
	}
	_, err = g.Commit("test commit", opts)
	if err != nil {
		t.Fatal(err)
	}

	otherFilename := "test/path.yaml"
	g, err = New(tmpDir, "other")
	if err != nil {
		t.Fatal(err)
	}
	if err := g.WriteFile(otherFilename, []byte(`second: value\n`), 0600); err != nil {
		t.Fatal(err)
	}

	sha, err := g.Commit("second commit", opts)
	if err != nil {
		t.Fatal(err)
	}
	if sha == "" {
		t.Fatalf("commit did not return a SHA")
	}
	commit, err := g.CommitObject(plumbing.NewHash(sha))
	if err != nil {
		t.Fatal(err)
	}
	want := []string{otherFilename}
	assertFilesCommitted(t, commit, want)

	if commit.Message != "second commit" {
		t.Errorf("failed to commit with the correct message, got %q, want %q", commit.Message, "second commit")
	}
}

func assertFilesCommitted(t *testing.T, commit *object.Commit, want []string) {
	t.Helper()
	found := []string{}

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
			found = append(found, f.Name)
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
				found = append(found, ch.To.Name)
			}
		}
	}

	if diff := cmp.Diff(want, found); diff != "" {
		t.Fatalf("failed to get the committed files:\n%s", diff)
	}
}

func makeTempDir(t *testing.T) (string, billy.Filesystem) {
	t.Helper()
	dir, err := ioutil.TempDir("", "askja")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})
	r, err := git.PlainInit(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	mainRef := plumbing.NewBranchReferenceName(testBranch)
	err = r.Storer.SetReference(plumbing.NewSymbolicReference(plumbing.HEAD, mainRef))
	if err != nil {
		t.Fatal(err)
	}
	err = r.CreateBranch(&config.Branch{
		Name: testBranch,
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = r.Branch(testBranch)
	if err != nil {
		t.Fatal(err)
	}
	return dir, osfs.New(dir)
}
