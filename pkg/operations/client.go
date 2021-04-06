package operations

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// DefaultClientFactory is the default client factory implementation.
var DefaultClientFactory ClientFactory = RawGitHubClientFactory

// Client implementations access git hosts to fetch files.
type Client interface {
	FileContents(ctx context.Context, repo, path, ref string) ([]byte, error)
}

// ClientFactory implementations should return a Client interface ready to be
// used to access the provided repoURL.
type ClientFactory func(repoURL string) (Client, error)

// ClientError is an error from Client implementations.
type ClientError struct {
	StatusCode int
	Message    string
}

func (c ClientError) Error() string {
	return fmt.Sprintf("client error: code %d, message: %s", c.StatusCode, c.Message)
}

// NewClientError creates and returns a new Git client error.
func NewClientError(code int, message string) ClientError {
	return ClientError{StatusCode: code, Message: message}
}

// RawGitHubClient is a very naive client that fetches files via
// raw.githubusercontent.com, without authentication.
type RawGitHubClient struct {
	*http.Client
}

// NewRawGitHubClient returns an implementation of the client that can fetch using
// the raw GitHub access.
//
// This will not work for authenticated requests, or anything other than GitHub.
func NewRawGitHubClient(c *http.Client) *RawGitHubClient {
	return &RawGitHubClient{Client: c}
}

// FileContents implements the Client interface.
func (c RawGitHubClient) FileContents(ctx context.Context, repo, path, ref string) ([]byte, error) {
	fileURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", repo, ref, path)
	resp, err := c.Client.Get(fileURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch file from %s: %w", fileURL, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch file from %s: status code %d", fileURL, resp.StatusCode)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file contents from %s: %w", fileURL, err)
	}
	return b, nil
}

// RawGitHubClientFactory is a very simple client that only supports fetching
// via github.com (and only unauthenticated requests).
func RawGitHubClientFactory(repoURL string) (Client, error) {
	parsed, err := url.Parse(repoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse repo URL %q: %w", repoURL, err)
	}

	if parsed.Host != "github.com" {
		return nil, errors.New("unsupported git provider, only github.com is currently supported")
	}
	return NewRawGitHubClient(http.DefaultClient), nil
}
