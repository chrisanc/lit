package config

import (
	"CLI_App/internal/domain"
	encoder "encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

/*
 * json.go - > manages the json files.
 */

type JsonAdapter struct {
	configPath string
}

func NewJSONAdapter() *JsonAdapter {
	return &JsonAdapter{}
}

func (json *JsonAdapter) GetConfig() *domain.Config {
	var dto *domain.Config
	err := encoder.Unmarshal(json.readJsonData(), &dto)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't read the JSON file. Are you sure it exists?")
		os.Exit(1)
	}
	return dto
}

func (json *JsonAdapter) SaveConfig(cfg *domain.Config) {
	data, err := encoder.Marshal(&cfg)
	if err != nil {
		os.Exit(1)
	}

	err = os.WriteFile(json.GetConfigPath(), data, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't write on the config.json file, are you sure it exists?")
		os.Exit(1)
	}
}

func (json *JsonAdapter) GetConfigPath() string {
	if json.configPath != "" {
		return json.configPath
	}

	path, err := os.Executable()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't get the path of the config file. Are you sure it exists?...")
		os.Exit(1)
	}
	json.configPath = filepath.Dir(path) + "/config.json"
	return json.configPath
}

func (json *JsonAdapter) readJsonData() []byte {
	file, err := os.ReadFile(json.GetConfigPath())
	if err != nil {
		os.Exit(1)
	}

	return file
}
