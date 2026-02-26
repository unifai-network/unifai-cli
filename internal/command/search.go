package command

import (
	"context"

	"github.com/spf13/cobra"

	clierrors "unifai-cli/internal/errors"
	"unifai-cli/internal/output"
	"unifai-cli/internal/unifai"
)

func newSearchCommand(global *GlobalOptions) *cobra.Command {
	var (
		query          string
		limit          int
		offset         int
		includeActions string
		jsonOutput     bool
	)

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search available services/actions",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if query == "" {
				return clierrors.NewUsageError("--query is required")
			}
			if limit <= 0 {
				return clierrors.NewUsageError("--limit must be greater than 0")
			}
			if offset < 0 {
				return clierrors.NewUsageError("--offset must be >= 0")
			}

			cfg, err := loadConfig(global, true)
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
			defer cancel()

			client := unifai.NewClient(cfg.Endpoint, cfg.APIKey, cfg.Timeout)
			resp, err := client.Search(ctx, unifai.SearchRequest{
				Query:          query,
				Limit:          limit,
				Offset:         offset,
				IncludeActions: splitCSV(includeActions),
			})
			if err != nil {
				return err
			}

			if jsonOutput {
				return output.PrintJSON(cmd.OutOrStdout(), resp)
			}
			return output.PrintSearch(cmd.OutOrStdout(), resp)
		},
	}

	cmd.Flags().StringVar(&query, "query", "", "Search query")
	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum results")
	cmd.Flags().IntVar(&offset, "offset", 0, "Result offset")
	cmd.Flags().StringVar(&includeActions, "include-actions", "", "Comma-separated action names to include")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Print raw JSON response")

	return cmd
}
