package install

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/bigkevmcd/askja/pkg/operations"
	"github.com/bigkevmcd/askja/pkg/profiles"
)

const (
	profileURLParam = "profile-url"
	branchNameParam = "profile-branch"
)

func MakeCmd() *cobra.Command {
	opts := &profiles.ProfileOptions{}

	cmd := &cobra.Command{
		Use:   "install",
		Short: "install a WeaveWorks profile",
		Run: func(cmd *cobra.Command, args []string) {
			if err := generateProfileResources(opts); err != nil {
				log.Fatalf("failed to generate profile resources: %s", err)
			}
		},
	}

	cmd.Flags().StringVar(
		&opts.ProfileURL,
		profileURLParam,
		"",
		"URL for fetching the profile from e.g. https://github.com/weaveworks/nginx-profile.git",
	)
	cmd.MarkFlagRequired(profileURLParam)

	cmd.Flags().StringVar(
		&opts.Branch,
		branchNameParam,
		"main",
		"branch name within the profile repo to use for the HelmRelease",
	)

	return cmd
}

func generateProfileResources(opts *profiles.ProfileOptions) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get the working directory: %w", err)
	}
	return operations.InstallProfile(context.TODO(), cwd, &operations.InstallOptions{
		ProfileOptions: opts,
		NewBranchName:  "test-branch",
	})
}
