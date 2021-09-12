package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// PodcastinationConfig holds all important config values needed in order to run the App.
type PodcastinationConfig struct {
	// StaticContentURL is the base url for accessing static content.
	StaticContentURL string `json:"static_content_url"`
	// PostgresDatasource is the datasource for the postgres database.
	PostgresDatasource string `json:"postgres_datasource"`
	// PullDir is the directory where tasks are placed.
	PullDir string `json:"pull_dir"`
	// PodcastDir is the directory where podcasts are stored.
	PodcastDir string `json:"podcast_dir"`
	// ImportInterval defines the duration in minutes after import tasks are retrieved.
	ImportInterval int `json:"import_interval"`
	// ServerAddr is the address the static file web_server will listen on (for example 127.0.0.1:8000).
	ServerAddr string `json:"server_addr"`
}

// ReadConfig reads a PodcastinationConfig from the given filepath.
func ReadConfig(filepath string) (PodcastinationConfig, error) {
	// Open config file.
	configFile, err := os.Open(filepath)
	if err != nil {
		return PodcastinationConfig{}, fmt.Errorf("could not open config file: %v", err)
	}
	// Read content.
	byteValue, err := ioutil.ReadAll(configFile)
	if err != nil {
		_ = configFile.Close()
		return PodcastinationConfig{}, fmt.Errorf("could not read content of config file: %v", err)
	}
	// Parse config.
	var config PodcastinationConfig
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		_ = configFile.Close()
		return PodcastinationConfig{}, fmt.Errorf("could not parse config file: %v", err)
	}
	// Close config.
	err = configFile.Close()
	if err != nil {
		return PodcastinationConfig{}, fmt.Errorf("could not close config file: %v", err)
	}
	return config, nil
}
