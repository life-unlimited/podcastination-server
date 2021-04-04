package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type PodcastinationConfig struct {
	PostgresDatasource string `json:"postgres_datasource"`
	PullDirectory      string `json:"pull_directory"`
	PodcastDirectory   string `json:"podcast_directory"`
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

	return config, configFile.Close()
}
