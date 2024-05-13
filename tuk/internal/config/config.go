package config

import (
	"context"
	"encoding/json"
	"log"
	"os"
)

// Config is the configuration for tuk
type Config struct {
	Paths []Path `json:"paths"`
	Args  Args   `json:"args"`
}

// Path represents a path to watch
type Path struct {
	// Raw is the raw path
	Raw string `json:"raw"`
	// Recursice sets whether to watch the path recursively
	Recursive bool `json:"recursive"` // fsnotify doesn't support it yet
}

// Args represents the arguments for tuk
type Args struct {
	// Path is the path to directory(s) to be watched. Set using the option: -p, --path
	Path string `json:"path"`
	// Recursive sets whether to watch the path recursively. Set using the option: -r, --recursive
	Recursive bool `json:"recursive"`
}

var (
	// DefaultConfig is the default configuration for tuk
	DefaultConfigPath = "config.json"
)

// LoadConfig loads the configuration from the default path
func LoadConfig(ctx context.Context, args *Args) (*Config, error) {
	// Load the configuration from the default path
	var config = &Config{}
	config.Args = *args

	log.Printf("config path: %#v\n", config)
	bytes, err := os.ReadFile(DefaultConfigPath)
	if err != nil {
		if os.IsNotExist(err) && args.Path != "" {
			return config, nil
		}
		return nil, err
	}

	// Unmarshal the configuration
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
