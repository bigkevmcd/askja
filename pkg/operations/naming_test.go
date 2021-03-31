package operations

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	profilesv1alpha1 "github.com/weaveworks/profiles/api/v1alpha1"
)

func TestFilenameFor(t *testing.T) {
	nameTests := []struct {
		typemeta   metav1.TypeMeta
		objectmeta metav1.ObjectMeta
		base       string
		want       string
	}{
		{
			typemeta: metav1.TypeMeta{
				Kind:       "Profile",
				APIVersion: "profiles.fluxcd.io/v1alpha1",
			},
			objectmeta: metav1.ObjectMeta{
				Name: "nginx",
			},
			want: "profile_nginx.yaml",
		},
		{
			typemeta: metav1.TypeMeta{
				Kind:       "HelmRelease",
				APIVersion: "helm.toolkit.fluxcd.io/v2beta1",
			},
			objectmeta: metav1.ObjectMeta{
				Name: "nginx-testing",
			},
			base: "test/project",
			want: "test/project/helmrelease_nginx-testing.yaml",
		},
	}

	for _, tt := range nameTests {
		t.Run(tt.want, func(t *testing.T) {
			o := &profilesv1alpha1.ProfileDefinition{
				TypeMeta:   tt.typemeta,
				ObjectMeta: tt.objectmeta,
			}

			s, err := filenameFrom(tt.base, o)
			if err != nil {
				t.Fatal(err)
			}

			if s != tt.want {
				t.Fatalf("filenameFrom() got %q, want %q", s, tt.want)
			}
		})
	}
}
