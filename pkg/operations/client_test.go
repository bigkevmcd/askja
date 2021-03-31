package operations

import (
	"context"
	"net/http"
	"testing"

	"gopkg.in/h2non/gock.v1"
)

var _ Client = (*RawGitHubClient)(nil)

func TestRawGitHubClient(t *testing.T) {
	body := "testing"
	defer gock.Off() // Flush pending mocks after test execution
	gock.New("https://raw.githubusercontent.com").
		Get("test/repo/main/profile.yaml").
		Reply(200).
		BodyString(body)

	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)
	c := NewRawGitHubClient(client)

	b, err := c.FileContents(context.TODO(), "test/repo", "profile.yaml", "main")
	if err != nil {
		t.Fatal(err)
	}

	if string(b) != body {
		t.Fatalf("got %s, want %s", b, body)
	}
}

func TestRawGitHubClientFactory(t *testing.T) {
	t.Skip()
}
