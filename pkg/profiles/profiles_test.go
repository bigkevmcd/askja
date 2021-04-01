package profiles

import (
	"testing"

	helmv2beta1 "github.com/fluxcd/helm-controller/api/v2beta1"
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
		&helmv2beta1.HelmRelease{
			TypeMeta:   metav1.TypeMeta{Kind: "HelmRelease", APIVersion: "helm.toolkit.fluxcd.io/v2beta1"},
			ObjectMeta: metav1.ObjectMeta{Name: "subscription-helm-release-test-chart"},
			Spec: helmv2beta1.HelmReleaseSpec{
				Chart: helmv2beta1.HelmChartTemplate{
					Spec: helmv2beta1.HelmChartTemplateSpec{
						Chart: "artifacts",
						SourceRef: helmv2beta1.CrossNamespaceObjectReference{
							Kind: "GitRepository",
							Name: "subscription-testing-main",
						},
					},
				},
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
