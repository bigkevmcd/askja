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

type gitRepositoryRefFunc func(*sourcev1beta1.GitRepositoryRef)

func branch(n string) gitRepositoryRefFunc {
	return func(o *sourcev1beta1.GitRepositoryRef) {
		o.Branch = n
	}
}

func testMakeGitRepository(name, repoURL string, opts ...gitRepositoryRefFunc) *sourcev1beta1.GitRepository {
	ref := &sourcev1beta1.GitRepositoryRef{}
	for _, o := range opts {
		o(ref)
	}

	return &sourcev1beta1.GitRepository{
		TypeMeta:   metav1.TypeMeta{Kind: gitRepositoryKind, APIVersion: gitRepositoryAPIVersion},
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec: sourcev1beta1.GitRepositorySpec{
			URL:       repoURL,
			Reference: ref,
		},
	}
}

type helmReleaseSpecFunc func(*helmv2beta1.HelmReleaseSpec)

func gitRepositorySourceRef(chart, repositoryName string) helmReleaseSpecFunc {
	return func(o *helmv2beta1.HelmReleaseSpec) {
		o.Chart = helmv2beta1.HelmChartTemplate{
			Spec: helmv2beta1.HelmChartTemplateSpec{
				Chart: chart,
				SourceRef: helmv2beta1.CrossNamespaceObjectReference{
					Kind: "GitRepository",
					Name: repositoryName,
				},
			},
		}
	}
}

func testMakeHelmRelease(name string, opts ...helmReleaseSpecFunc) *helmv2beta1.HelmRelease {
	spec := helmv2beta1.HelmReleaseSpec{}
	for _, o := range opts {
		o(&spec)
	}

	return &helmv2beta1.HelmRelease{
		TypeMeta:   metav1.TypeMeta{Kind: "HelmRelease", APIVersion: "helm.toolkit.fluxcd.io/v2beta1"},
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec:       spec,
	}
}

func TestMakeArtifacts(t *testing.T) {
	artifactTests := []struct {
		name      string
		profile   *Profile
		artifacts []runtime.Object
	}{
		{
			name:    "git repository helm release",
			profile: makeTestProfile(Artifact{Name: testChartname, Path: testChartPath}),
			artifacts: []runtime.Object{
				testMakeGitRepository("subscription-testing-main", testProfileURL, branch("main")),
				testMakeHelmRelease("subscription-helm-release-test-chart", gitRepositorySourceRef(testChartPath, "subscription-testing-main")),
			},
		},
	}

	for _, tt := range artifactTests {
		t.Run(tt.name, func(t *testing.T) {
			o := MakeArtifacts(tt.profile, &ProfileOptions{
				ProfileURL: testProfileURL,
				Branch:     "main",
			})

			if diff := cmp.Diff(tt.artifacts, o); diff != "" {
				t.Fatalf("failed to make artifacts:\n%s", diff)
			}
		})
	}
}

func makeTestProfile(a ...Artifact) *Profile {
	return &Profile{
		Spec: ProfileSpec{
			Description: "foo",
			Artifacts:   a,
		},
	}
}
