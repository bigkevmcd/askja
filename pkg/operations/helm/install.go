package helm

import (
	"context"
	"path/filepath"
	"strings"

	sourcev1beta1 "github.com/fluxcd/source-controller/api/v1beta1"
	"github.com/go-git/go-billy/v5"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type HelmChart struct {
	// URL is the Chart repository
	URL string
	// Name of the chart at the remote repository
	Name string
	// Version defines the version of the chart at the remote repository
	Version string
}

type InstallOptions struct {
	Chart     HelmChart
	Profile   string
	Namespace string
}

func InstallHelmChart(ctx context.Context, fs billy.Filesystem, opts *InstallOptions) (map[string]runtime.Object, error) {
	files := map[string]runtime.Object{
		filepath.Join("profiles", opts.Profile, makeFilename(opts, "helm-release")): makeHelmRepository(opts),
	}

	return files, nil
}

func makeFilename(opts *InstallOptions, suffix string) string {
	return strings.Join([]string{opts.Chart.Name, "chart", suffix}, "-") + ".yaml"
}

func flatten(s string) string {
	return strings.ReplaceAll(strings.ToLower(s), "/", "-")
}

func makeHelmRepository(o *InstallOptions) *sourcev1beta1.HelmRepository {
	return &sourcev1beta1.HelmRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "",
			Namespace: o.Namespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       sourcev1beta1.HelmRepositoryKind,
			APIVersion: sourcev1beta1.GroupVersion.String(),
		},
		Spec: sourcev1beta1.HelmRepositorySpec{
			URL: o.Chart.URL,
		},
	}
}
