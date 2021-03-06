package config

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type Config struct {
	Podcasts             []string `toml:"podcasts"`
	PathForTVShows       string   `toml:"path_for_tv_shows"`
	HoursToRefresh       int      `toml:"hours_to_refresh"`
	ConcurrencyDownloads int      `toml:"concurrency_downloads"`
}

const PathConfig string = "config.toml"

// LoadConfig loads TOML configuration from a file path
func LoadConfig() (*Config, error) {
	config := Config{}

	_, err := toml.DecodeFile(PathConfig, &config)

	if err != nil {
		return nil, errors.Wrap(err, "failed to load config file")
	}

	return &config, nil
}
