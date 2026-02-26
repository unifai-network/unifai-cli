package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type ResolveOptions struct {
	ConfigPath   string
	FlagAPIKey   string
	FlagEndpoint string
	FlagTimeout  time.Duration
}

func Resolve(opts ResolveOptions) (EffectiveConfig, error) {
	configPath := opts.ConfigPath
	if strings.TrimSpace(configPath) == "" {
		var err error
		configPath, err = DefaultConfigPath()
		if err != nil {
			return EffectiveConfig{}, fmt.Errorf("determine default config path: %w", err)
		}
	}

	effective := EffectiveConfig{
		Endpoint:       DefaultEndpoint,
		Timeout:        defaultTimeout,
		ConfigPath:     configPath,
		APIKeySource:   SourceDefault,
		EndpointSource: SourceDefault,
		TimeoutSource:  SourceDefault,
	}

	if data, err := os.ReadFile(configPath); err == nil {
		effective.LoadedFromFile = true
		var fileCfg FileConfig
		if err := yaml.Unmarshal(data, &fileCfg); err != nil {
			return EffectiveConfig{}, fmt.Errorf("parse config file %s: %w", configPath, err)
		}
		if strings.TrimSpace(fileCfg.APIKey) != "" {
			effective.APIKey = strings.TrimSpace(fileCfg.APIKey)
			effective.APIKeySource = SourceConfig
		}
		if strings.TrimSpace(fileCfg.Endpoint) != "" {
			effective.Endpoint = strings.TrimRight(strings.TrimSpace(fileCfg.Endpoint), "/")
			effective.EndpointSource = SourceConfig
		}
		if fileCfg.TimeoutSeconds > 0 {
			effective.Timeout = time.Duration(fileCfg.TimeoutSeconds) * time.Second
			effective.TimeoutSource = SourceConfig
		}
	} else if !os.IsNotExist(err) {
		return EffectiveConfig{}, fmt.Errorf("read config file %s: %w", configPath, err)
	}

	if envKey := strings.TrimSpace(os.Getenv(APIKeyEnvVar)); envKey != "" {
		effective.APIKey = envKey
		effective.APIKeySource = SourceEnv
	}

	if strings.TrimSpace(opts.FlagAPIKey) != "" {
		effective.APIKey = strings.TrimSpace(opts.FlagAPIKey)
		effective.APIKeySource = SourceFlag
	}
	if strings.TrimSpace(opts.FlagEndpoint) != "" {
		effective.Endpoint = strings.TrimRight(strings.TrimSpace(opts.FlagEndpoint), "/")
		effective.EndpointSource = SourceFlag
	}
	if opts.FlagTimeout > 0 {
		effective.Timeout = opts.FlagTimeout
		effective.TimeoutSource = SourceFlag
	}

	if effective.Timeout <= 0 {
		return EffectiveConfig{}, fmt.Errorf("timeout must be greater than 0")
	}

	return effective, nil
}
