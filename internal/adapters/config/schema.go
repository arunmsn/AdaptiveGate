package config

// Config is the top-level typed struct for ixr.yaml.
// All fields have sane defaults; only api_key values are required.
type Config struct {
	Port      int                        `yaml:"port"`
	Providers map[string]ProviderConfig  `yaml:"providers"`
}

// ProviderConfig holds credentials and options for a single LLM provider.
type ProviderConfig struct {
	APIKey  string `yaml:"api_key"`
	BaseURL string `yaml:"base_url,omitempty"`
}
