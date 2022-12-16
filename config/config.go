package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
)

var configPath = fmt.Sprintf("%s/.config/mr-reviewer/config.json", os.Getenv("HOME"))

type Config struct {
	BasePath     string        `json:"gitlab_base_path"`
	Token        string        `json:"gitlab_token"`
	Repositories *[]Repository `json:"repositories"`
}

type Repository struct {
	Name  string `json:"name"`
	Desc  string `json:"description"`
	Route string `json:"route"`
}

func (r Repository) Title() string       { return r.Name }
func (r Repository) Description() string { return r.Desc }
func (r Repository) FilterValue() string { return r.Name }

func Read() (*Config, error) {
	fileContents, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(fileContents, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *Config) RepositoriesAsItems() []list.Item {
	var items []list.Item
	for _, repo := range *c.Repositories {
		items = append(items, repo)
	}
	return items
}
