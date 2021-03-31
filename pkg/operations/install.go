package operations

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/bigkevmcd/askja/pkg/git"
	"github.com/bigkevmcd/askja/pkg/profiles"
	"sigs.k8s.io/yaml"
)

const (
	defaultFileMode os.FileMode = 0644
)

// InstallOptions are passed to the install operation to provide information for
// bootstrappign resources.
type InstallOptions struct {
	*profiles.ProfileOptions
	NewBranchName string
}

// InstallProfile will generate the HelmRelease for a profile.
//
// TODO: could this take a git.Repository?
func InstallProfile(ctx context.Context, path string, options *InstallOptions) error {
	client, err := DefaultClientFactory(options.ProfileOptions.ProfileURL)
	if err != nil {
		return fmt.Errorf("failed to create a client for %q: %w",
			options.ProfileOptions.ProfileURL, err)
	}

	repo, err := extractRepo(options.ProfileOptions.ProfileURL)
	if err != nil {
		return err
	}
	bytes, err := client.FileContents(ctx, repo, "profile.yaml", "main")
	p, err := profiles.ParseBytes(bytes)
	if err != nil {
		return err
	}

	result := profiles.MakeArtifacts(p, options.ProfileOptions)

	g, err := git.New(path, options.NewBranchName)
	for _, v := range result {
		b, err := yaml.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal: %w", err)
		}
		output, err := filenameFrom("", v)
		if err != nil {
			return err
		}
		if err := g.WriteFile(output, b, defaultFileMode); err != nil {
			return fmt.Errorf("failed to write to file %q in %q: %w", output, path, err)
		}
	}
	return nil
}

// TODO: This is extremely naive.
func extractRepo(profileURL string) (string, error) {
	parsed, err := url.Parse(profileURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse profileURL %q: %w", profileURL, err)
	}
	return strings.TrimSuffix(parsed.Path, ".git"), nil
}
