package profiles

import (
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	helmv2beta1 "github.com/fluxcd/helm-controller/api/v2beta1"
	sourcev1beta1 "github.com/fluxcd/source-controller/api/v1beta1"
)

const (
	gitRepositoryKind       = "GitRepository"
	gitRepositoryAPIVersion = "source.toolkit.fluxcd.io/v1beta1"
	helmReleaseKind         = "HelmRelease"
	helmReleaseAPIVersion   = "helm.toolkit.fluxcd.io/v2beta1"
)

// ProfileOptions is a set of configuration options to use when creating the
// Profile artifacts.
type ProfileOptions struct {
	ProfileURL string
	Branch     string
}

// MakeArtifacts creates and returns the artifacts necessary to deploy a Profile.
func MakeArtifacts(p *Profile, opts *ProfileOptions) []runtime.Object {
	objects := []runtime.Object{}
	objects = append(objects, createGitRepository(p, opts))
	objects = append(objects, createHelmRelease(p, opts))
	return objects
}

func createGitRepository(p *Profile, opts *ProfileOptions) *sourcev1beta1.GitRepository {
	return &sourcev1beta1.GitRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name: makeGitRepoName(opts.ProfileURL, opts.Branch),
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       gitRepositoryKind,
			APIVersion: gitRepositoryAPIVersion,
		},
		Spec: sourcev1beta1.GitRepositorySpec{
			URL: opts.ProfileURL,
			Reference: &sourcev1beta1.GitRepositoryRef{
				Branch: opts.Branch,
			},
		},
	}
	// err := controllerutil.SetOwnerReference(&p.subscription, &gitRepo, p.client.Scheme())
	// if err != nil {
	// 	return fmt.Errorf("failed to set resource ownership: %w", err)
	// }
}

func createHelmRelease(p *Profile, opts *ProfileOptions) *helmv2beta1.HelmRelease {
	return &helmv2beta1.HelmRelease{
		ObjectMeta: metav1.ObjectMeta{
			Name: makeHelmReleaseName(p),
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       helmReleaseKind,
			APIVersion: helmReleaseAPIVersion,
		},
		Spec: helmv2beta1.HelmReleaseSpec{
			Chart: helmv2beta1.HelmChartTemplate{
				Spec: helmv2beta1.HelmChartTemplateSpec{
					// TODO obvs don't rely on index 0
					Chart: p.Spec.Artifacts[0].Path,
					SourceRef: helmv2beta1.CrossNamespaceObjectReference{
						Kind: gitRepositoryKind,
						Name: makeGitRepoName(opts.ProfileURL, opts.Branch),
					},
				},
			},
		},
	}
	// 	err := controllerutil.SetControllerReference(&p.subscription, &helmRelease, p.client.Scheme())
	// 	if err != nil {
	// 		return fmt.Errorf("failed to set resource ownership: %w", err)
	// 	}

	// 	p.log.Info("creating HelmRelease", "resource", helmReleasename)
	// 	return p.client.Create(ctx, &helmRelease)
}

func makeHelmReleaseName(p *Profile) string {
	return join("subscription", "helm-release", p.Spec.Artifacts[0].Name)
}

// TODO: error if this is more than 63 chars long.
func makeGitRepoName(profileURL, branch string) string {
	repoParts := strings.Split(profileURL, "/")
	repoName := strings.TrimSuffix(repoParts[len(repoParts)-1], ".git")
	return join("subscription", repoName, branch)
}

func join(s ...string) string {
	return strings.Join(s, "-")
}
