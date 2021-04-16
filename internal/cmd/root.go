package cmd

import (
	"log"

	"github.com/bigkevmcd/askja/internal/cmd/helm"
	"github.com/bigkevmcd/askja/internal/cmd/install"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	cobra.OnInitialize(initConfig)
}

func makeRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "askja",
		Short:         "askja profiles installer",
		SilenceErrors: true,
	}
	cmd.AddCommand(install.MakeCmd())
	cmd.AddCommand(helm.MakeCmd())
	return cmd
}

func initConfig() {
	viper.AutomaticEnv()
}

// Execute is the main entry point into this component.
func Execute() {
	if err := makeRootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}

func logIfError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
