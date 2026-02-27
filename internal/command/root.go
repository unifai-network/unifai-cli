package command

import (
	"time"

	"github.com/spf13/cobra"

	clierrors "unifai/internal/errors"
)

type GlobalOptions struct {
	ConfigPath string
	APIKey     string
	Endpoint   string
	Timeout    time.Duration
}

func NewRootCommand() *cobra.Command {
	opts := &GlobalOptions{}

	root := &cobra.Command{
		Use:           "unifai",
		Short:         "CLI for Unifai search_services and invoke_service",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	root.PersistentFlags().StringVar(&opts.ConfigPath, "config", "", "Path to config file (default: ~/.config/unifai-cli/config.yaml)")
	root.PersistentFlags().StringVar(&opts.APIKey, "api-key", "", "Unifai API key")
	root.PersistentFlags().StringVar(&opts.Endpoint, "endpoint", "", "API endpoint base URL")
	root.PersistentFlags().DurationVar(&opts.Timeout, "timeout", 0, "Request timeout (default: 50s)")
	root.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return clierrors.NewUsageError(err.Error())
	})

	root.AddCommand(newSearchCommand(opts))
	root.AddCommand(newInvokeCommand(opts))
	root.AddCommand(newConfigCommand(opts))
	root.AddCommand(newVersionCommand())

	return root
}
