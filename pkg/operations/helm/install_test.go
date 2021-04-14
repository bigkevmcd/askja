package helm

import (
	"context"
	"testing"

	"github.com/bigkevmcd/askja/test"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	sourcev1beta1 "github.com/fluxcd/source-controller/api/v1beta1"
)

func TestInstallHelm(t *testing.T) {
	fs := osfs.New(test.MakeTempDir(t))
	files, err := InstallHelmChart(context.TODO(), fs, &InstallOptions{
		Chart: HelmChart{
			URL:     "https://charts.bitnami.com/bitnami",
			Name:    "bitnami/redis",
			Version: "6.2.1",
		},
		Profile:   "test-profile",
		Namespace: "test-namespace",
	})
	if err != nil {
		t.Fatal(err)
	}

	want := map[string]runtime.Object{
		"profiles/test-profile/bitnami/redis-chart-helm-release.yaml": &sourcev1beta1.HelmRepository{
			TypeMeta: metav1.TypeMeta{
				Kind:       sourcev1beta1.HelmRepositoryKind,
				APIVersion: sourcev1beta1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{Namespace: "test-namespace"},
			Spec:       sourcev1beta1.HelmRepositorySpec{URL: "https://charts.bitnami.com/bitnami"},
		},
	}

	if diff := cmp.Diff(want, files); diff != "" {
		t.Fatalf("failed to generate installation resources:\n%s", diff)
	}
}
