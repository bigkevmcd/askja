package helm

import (
	"github.com/spf13/cobra"

	"github.com/bigkevmcd/askja/pkg/operations/helm"
)

func MakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "helm",
		Short: "helm chart operations",
	}

	cmd.AddCommand(makeHelmInstallCmd())
	return cmd
}

func makeHelmInstallCmd() *cobra.Command {
	var opts helm.InstallOptions
	const (
		repositoryURLParam = "repository-url"
		chartNameParam     = "chart"
		chartVersionParam  = "version"
		profileParam       = "profile"
	)

	cmd := &cobra.Command{
		Use:   "install",
		Short: "add a helm chart to a profile",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	cmd.Flags().StringVar(
		&opts.Chart.URL,
		repositoryURLParam,
		"",
		"the chart repository URL e.g. https://charts.bitnami.com/bitnami",
	)
	cmd.MarkFlagRequired(repositoryURLParam)

	cmd.Flags().StringVar(
		&opts.Chart.Name,
		chartNameParam,
		"",
		"the chart to install e.g. bitnami/nginx",
	)
	cmd.MarkFlagRequired(chartNameParam)

	cmd.Flags().StringVar(
		&opts.Chart.Version,
		chartVersionParam,
		"",
		"the chart version to install e.g. v1.19.0",
	)
	cmd.MarkFlagRequired(chartVersionParam)

	cmd.Flags().StringVar(
		&opts.Profile,
		profileParam,
		"",
		"profile to install in e.g. demo, this will modify the profile in the profile directory relative to the current dir",
	)
	cmd.MarkFlagRequired(profileParam)
	return cmd
}
