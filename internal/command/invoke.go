package command

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	clierrors "unifai/internal/errors"
	"unifai/internal/output"
	"unifai/internal/retry"
	"unifai/internal/unifai"
)

func newInvokeCommand(global *GlobalOptions) *cobra.Command {
	var (
		action        string
		payloadInput  string
		payloadFormat string
		maxRetries    int
		jsonOutput    bool
	)

	cmd := &cobra.Command{
		Use:   "invoke",
		Short: "Invoke a specific service/action",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if strings.TrimSpace(action) == "" {
				return clierrors.NewUsageError("--action is required")
			}
			if maxRetries < 0 {
				return clierrors.NewUsageError("--max-retries must be >= 0")
			}

			cfg, err := loadConfig(global, true)
			if err != nil {
				return err
			}

			payload, err := parsePayload(payloadInput, payloadFormat)
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
			defer cancel()

			client := unifai.NewClient(cfg.Endpoint, cfg.APIKey, cfg.Timeout)
			request := unifai.InvokeRequest{
				Action:  action,
				Payload: payload,
			}

			resp, err := retry.Do(ctx, maxRetries, time.Second, unifai.IsRetryableError, func() (any, error) {
				return client.Invoke(ctx, request)
			})
			if err != nil {
				return err
			}

			if jsonOutput {
				return output.PrintJSON(cmd.OutOrStdout(), resp)
			}

			return output.PrintValue(cmd.OutOrStdout(), output.NormalizeInvokeResponse(resp))
		},
	}

	cmd.Flags().StringVar(&action, "action", "", "Action ID/name to invoke")
	cmd.Flags().StringVar(&payloadInput, "payload", "", "Payload JSON string or @file path")
	cmd.Flags().StringVar(&payloadFormat, "payload-format", "auto", "Payload format: auto|object|string")
	cmd.Flags().IntVar(&maxRetries, "max-retries", 1, "Retry count for retryable failures")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Print raw JSON response")

	return cmd
}

func parsePayload(input, format string) (any, error) {
	format = strings.ToLower(strings.TrimSpace(format))
	if format == "" {
		format = "auto"
	}
	if format != "auto" && format != "object" && format != "string" {
		return nil, clierrors.NewUsageError("--payload-format must be one of: auto, object, string")
	}

	raw, err := loadPayloadInput(input)
	if err != nil {
		return nil, err
	}
	if raw == "" {
		return nil, nil
	}

	switch format {
	case "string":
		return raw, nil
	case "object":
		var v any
		if err := json.Unmarshal([]byte(raw), &v); err != nil {
			return nil, clierrors.NewUsageError("payload is not valid JSON object/string: %v", err)
		}
		return v, nil
	default: // auto
		var v any
		if err := json.Unmarshal([]byte(raw), &v); err == nil {
			return v, nil
		}
		return raw, nil
	}
}

func loadPayloadInput(input string) (string, error) {
	if strings.TrimSpace(input) == "" {
		return "", nil
	}
	trimmed := strings.TrimSpace(input)
	if strings.HasPrefix(trimmed, "@") {
		path := strings.TrimSpace(strings.TrimPrefix(trimmed, "@"))
		if path == "" {
			return "", clierrors.NewUsageError("payload file path after @ cannot be empty")
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("read payload file %q: %w", path, err)
		}
		return string(data), nil
	}
	return input, nil
}
