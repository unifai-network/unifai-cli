package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"unifai-cli/internal/config"
	clierrors "unifai-cli/internal/errors"
	"unifai-cli/internal/output"
)

func newConfigCommand(global *GlobalOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage ucli configuration",
	}

	cmd.AddCommand(newConfigInitCommand(global))
	cmd.AddCommand(newConfigShowCommand(global))

	return cmd
}

func newConfigInitCommand(global *GlobalOptions) *cobra.Command {
	var (
		overridePath string
		force        bool
	)

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create a config template file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			path, err := resolveConfigPath(global.ConfigPath, overridePath)
			if err != nil {
				return err
			}

			if !force {
				if _, err := os.Stat(path); err == nil {
					return clierrors.NewUsageError("config already exists at %s (use --force to overwrite)", path)
				}
			}

			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return fmt.Errorf("create config directory: %w", err)
			}

			if err := os.WriteFile(path, []byte(config.ExampleConfigYAML()), 0o600); err != nil {
				return fmt.Errorf("write config file: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Wrote config template: %s\n", path)
			return nil
		},
	}

	cmd.Flags().StringVar(&overridePath, "path", "", "Config path override")
	cmd.Flags().BoolVar(&force, "force", false, "Overwrite existing config")
	return cmd
}

func newConfigShowCommand(global *GlobalOptions) *cobra.Command {
	var (
		overridePath string
		jsonOutput   bool
	)

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show effective config and its sources",
		RunE: func(cmd *cobra.Command, _ []string) error {
			path, err := resolveConfigPath(global.ConfigPath, overridePath)
			if err != nil {
				return err
			}

			cfg, err := config.Resolve(config.ResolveOptions{
				ConfigPath:   path,
				FlagAPIKey:   global.APIKey,
				FlagEndpoint: global.Endpoint,
				FlagTimeout:  global.Timeout,
			})
			if err != nil {
				return err
			}

			if jsonOutput {
				return output.PrintJSON(cmd.OutOrStdout(), cfg)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Config path: %s\n", cfg.ConfigPath)
			fmt.Fprintf(cmd.OutOrStdout(), "Config loaded: %t\n", cfg.LoadedFromFile)
			fmt.Fprintf(cmd.OutOrStdout(), "Endpoint: %s (source: %s)\n", cfg.Endpoint, cfg.EndpointSource)
			fmt.Fprintf(cmd.OutOrStdout(), "Timeout: %s (source: %s)\n", cfg.Timeout, cfg.TimeoutSource)

			apiKeyState := "missing"
			if cfg.APIKey != "" {
				apiKeyState = "present"
			}
			fmt.Fprintf(cmd.OutOrStdout(), "API key: %s (source: %s)\n", apiKeyState, cfg.APIKeySource)
			return nil
		},
	}

	cmd.Flags().StringVar(&overridePath, "path", "", "Config path override")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Print full config in JSON")
	return cmd
}

func resolveConfigPath(globalPath, overridePath string) (string, error) {
	if overridePath != "" {
		return overridePath, nil
	}
	if globalPath != "" {
		return globalPath, nil
	}
	return config.DefaultConfigPath()
}
