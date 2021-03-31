package profiles

import (
	"testing"

	// helmv2beta1 "github.com/fluxcd/helm-controller/api/v2beta1"
	sourcev1beta1 "github.com/fluxcd/source-controller/api/v1beta1"
	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	testChartname  = "test-chart"
	testChartPath  = "artifacts"
	testProfileURL = "https://example.com/testing/testing.git"
)

func TestMakeArtifacts(t *testing.T) {
	p := makeTestProfile()
	want := []runtime.Object{
		&sourcev1beta1.GitRepository{
			TypeMeta:   metav1.TypeMeta{Kind: gitRepositoryKind, APIVersion: gitRepositoryAPIVersion},
			ObjectMeta: metav1.ObjectMeta{Name: "subscription-testing-main"},
			Spec: sourcev1beta1.GitRepositorySpec{
				URL:       "https://example.com/testing/testing.git",
				Reference: &sourcev1beta1.GitRepositoryRef{Branch: "main"},
			},
		},
	}

	o := MakeArtifacts(p, &ProfileOptions{
		ProfileURL: "https://example.com/testing/testing.git",
		Branch:     "main",
	})

	if diff := cmp.Diff(want, o); diff != "" {
		t.Fatalf("failed to make artifacts:\n%s", diff)
	}
}

func makeTestProfile() *Profile {
	return &Profile{
		Spec: ProfileSpec{
			Description: "foo",
			Artifacts: []Artifact{
				{
					Name: testChartname,
					Path: testChartPath,
				},
			},
		},
	}
}
