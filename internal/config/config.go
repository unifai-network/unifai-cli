package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	DefaultEndpoint    = "https://app.uniclaw.ai/api/v1/unifai"
	APIKeyEnvVar       = "UNIFAI_AGENT_API_KEY"
	EndpointEnvVar     = "UNIFAI_ENDPOINT"
)

const defaultTimeout = 50 * time.Second

type FileConfig struct {
	APIKey         string `yaml:"apiKey"`
	Endpoint       string `yaml:"endpoint"`
	TimeoutSeconds int    `yaml:"timeoutSeconds"`
}

type ValueSource string

const (
	SourceDefault ValueSource = "default"
	SourceConfig  ValueSource = "config"
	SourceEnv     ValueSource = "env"
	SourceFlag    ValueSource = "flag"
)

type EffectiveConfig struct {
	APIKey string `json:"apiKey"`

	Endpoint string        `json:"endpoint"`
	Timeout  time.Duration `json:"timeout"`

	ConfigPath     string `json:"configPath"`
	LoadedFromFile bool   `json:"loadedFromFile"`

	APIKeySource   ValueSource `json:"apiKeySource"`
	EndpointSource ValueSource `json:"endpointSource"`
	TimeoutSource  ValueSource `json:"timeoutSource"`
}

func DefaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "unifai-cli", "config.yaml"), nil
}

func ExampleConfigYAML() string {
	return fmt.Sprintf(`# unifai-cli configuration
# Fill apiKey with your real key, or set %s in your shell.
# Set endpoint via %s environment variable or keep default.
apiKey: ""
endpoint: %q
timeoutSeconds: 50
`, APIKeyEnvVar, EndpointEnvVar, DefaultEndpoint)
}
