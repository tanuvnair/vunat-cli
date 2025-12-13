package projects

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type CommandGroup struct {
	Name         string   `json:"name"`
	AbsolutePath string   `json:"absolutePath"`
	Commands     []string `json:"commands"`
}

type Project []CommandGroup

type Config struct {
	Projects map[string]Project `json:"projects"`
}

var registry map[string]Project

func init() {
	registry = make(map[string]Project)
	loadConfig()
}

func loadConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".vunat", "config.json")

	// Create .vunat directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// If config file doesn't exist, create an empty one
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		emptyConfig := Config{Projects: make(map[string]Project)}
		data, _ := json.MarshalIndent(emptyConfig, "", "  ")
		if err := os.WriteFile(configPath, data, 0644); err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	registry = config.Projects
	return nil
}

func Get(name string) (Project, error) {
	// Reload config in case it was updated
	if err := loadConfig(); err != nil {
		return Project{}, err
	}

	project, ok := registry[name]
	if !ok {
		return Project{}, fmt.Errorf("unknown project: %s", name)
	}
	return project, nil
}

func GetAll() map[string]Project {
	loadConfig() // Reload to get latest
	return registry
}
