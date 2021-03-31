package operations

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/bigkevmcd/askja/pkg/profiles"
	"github.com/bigkevmcd/askja/test"
)

func TestInstallProfile(t *testing.T) {
	client := newMockClient()
	dir, _ := test.MakeTempGitDir(t)
	DefaultClientFactory = func(s string) (Client, error) {
		if s == "https://github.com/weaveworks/nginx-profile.git" {
			return client, nil
		}
		return nil, nil
	}
	client.add("weaveworks/nginx-profile", "profile.yaml", "main", []byte(`
apiVersion: profiles.fluxcd.io/v1alpha1
kind: Profile
metadata:
  name: nginx
spec:
  description: Profile for deploying nginx
  version: v0.0.1
  artifacts:
    - name: nginx-server
      path: nginx/chart
`))

	if err := InstallProfile(context.TODO(), dir, &InstallOptions{ProfileOptions: &profiles.ProfileOptions{
		ProfileURL: "https://github.com/weaveworks/nginx-profile.git",
		Branch:     "main",
	},
		NewBranchName: "test-branch",
	}); err != nil {
		t.Fatal(err)
	}

	// TODO: test the remainder
}

func newMockClient() *mockClient {
	return &mockClient{contents: make(map[string][]byte)}
}

type mockClient struct {
	contents map[string][]byte
}

func (m *mockClient) FileContents(ctx context.Context, repo, path, ref string) ([]byte, error) {
	b, ok := m.contents[key(repo, path, ref)]
	if !ok {
		return nil, NewClientError(http.StatusNotFound, "file not found")
	}
	return b, nil
}

func (m *mockClient) add(repo, path, ref string, content []byte) {
	m.contents[key(repo, path, ref)] = content
}

func key(s ...string) string {
	return strings.Join(s, ":")
}
