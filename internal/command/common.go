package command

import (
	"strings"

	"unifai/internal/config"
	clierrors "unifai/internal/errors"
)

func loadConfig(opts *GlobalOptions, requireAPIKey bool) (config.EffectiveConfig, error) {
	cfg, err := config.Resolve(config.ResolveOptions{
		ConfigPath:   opts.ConfigPath,
		FlagAPIKey:   opts.APIKey,
		FlagEndpoint: opts.Endpoint,
		FlagTimeout:  opts.Timeout,
	})
	if err != nil {
		return config.EffectiveConfig{}, err
	}

	if requireAPIKey && strings.TrimSpace(cfg.APIKey) == "" {
		return config.EffectiveConfig{}, clierrors.NewUsageError("missing API key: set --api-key, %s, or config file", config.APIKeyEnvVar)
	}

	return cfg, nil
}

func splitCSV(input string) []string {
	if strings.TrimSpace(input) == "" {
		return nil
	}
	parts := strings.Split(input, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}
