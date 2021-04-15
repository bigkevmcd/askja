package askja

import (
	"fmt"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"

	"github.com/bigkevmcd/askja/test"
)

func CommitFiles(wt *git.Worktree, files map[string]runtime.Object, msg string, opts *git.CommitOptions) error {
	for k, v := range files {
		f, err := wt.Filesystem.Create(k)
		if err != nil {
			return fmt.Errorf("failed to create file %q: %w", k, err)
		}
		b, err := yaml.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal %v: %w", v, err)
		}
		if _, err := f.Write(b); err != nil {
			return fmt.Errorf("failed to write to %q: %w", k, err)
		}
		f.Close()
		if _, err = wt.Add(k); err != nil {
			return fmt.Errorf("failed to add %q to a commit: %w", k, err)
		}
	}
	_, err := wt.Commit(msg, opts)
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}
	return nil
}

func TestCommitFiles(t *testing.T) {
	files := map[string]runtime.Object{
		"testing.yaml": &corev1.Secret{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Secret",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testing",
				Namespace: "testing",
			},
		},
	}
	r, dir := test.MakeTempRepository(t)
	wt, err := r.Worktree()
	wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName("test-files"),
		Create: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	opts := &git.CommitOptions{}
	if err := CommitFiles(wt, files, "test commit", opts); err != nil {
		t.Fatal(err)
	}

	head, err := r.Head()
	if err != nil {
		t.Fatal(err)
	}
	headCommit, err := r.CommitObject(head.Hash())
	if err != nil {
		t.Fatal(err)
	}
	want := map[string][]byte{
		"testing.yaml": []byte("apiVersion: v1\nkind: Secret\nmetadata:\n  creationTimestamp: null\n  name: testing\n  namespace: testing\n"),
	}
	committed := test.GetFilesInCommit(t, headCommit, dir)
	if diff := cmp.Diff(want, committed); diff != "" {
		t.Fatalf("failed to commit files:\n%s", diff)
	}
}

func TestUpdatingFiles(t *testing.T) {
}
