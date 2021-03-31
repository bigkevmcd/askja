package profiles

import (
	"os"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	p, err := ParseBytes(mustRead(t, "testdata/profile.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	want := &Profile{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Profile",
			APIVersion: "profiles.fluxcd.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx",
		},
		Spec: ProfileSpec{
			Description: "Profile for deploying nginx",
			Artifacts: []Artifact{
				{Name: "nginx-server", Path: "nginx/chart"},
			},
		},
	}
	if diff := cmp.Diff(want, p); diff != "" {
		t.Fatalf("failed to parse profile:\n%s", diff)
	}
}

func mustRead(t *testing.T, fname string) []byte {
	t.Helper()
	b, err := os.ReadFile(fname)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
