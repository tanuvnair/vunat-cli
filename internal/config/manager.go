package config

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Manager is responsible for reading, creating, and returning path to config.
//
// It provides a small abstraction around filesystem operations so higher-level
// code (commands, tests) can depend on the interface instead of the OS calls.
type Manager interface {
	// Path returns the absolute path to the config file (does not create it).
	Path() string
	// Ensure ensures directory and file exist and returns the path.
	Ensure() (string, error)
	// Read returns the config bytes.
	Read() ([]byte, error)
	// Write replaces the config file content.
	Write([]byte, fs.FileMode) error
}

// FSManager implements Manager using the OS filesystem.
type FSManager struct {
	path string
}

// NewFSManager constructs an FSManager. If the provided path is empty
// it will default to "$HOME/.vunat/config.json". If $HOME cannot be
// resolved, it falls back to "./.vunat/config.json".
func NewFSManager(path string) *FSManager {
	if path != "" {
		return &FSManager{path: path}
	}

	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		// Fallback to a relative path if home cannot be determined.
		path = filepath.Join(".", ".vunat", "config.json")
	} else {
		path = filepath.Join(home, ".vunat", "config.json")
	}
	return &FSManager{path: path}
}

// Path returns the configured path (may be absolute or relative).
func (m *FSManager) Path() string {
	return m.path
}

// Ensure makes sure the config directory exists and that a config file is present.
// If the file does not exist, it creates an initial empty config with a top-level
// "projects" map to match the expected structure.
func (m *FSManager) Ensure() (string, error) {
	configDir := filepath.Dir(m.path)
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	if _, err := os.Stat(m.path); os.IsNotExist(err) {
		initial := map[string]any{
			"projects": map[string]any{},
		}
		data, err := json.MarshalIndent(initial, "", "  ")
		if err != nil {
			// Unlikely, but surface the error.
			return "", fmt.Errorf("failed to marshal initial config: %w", err)
		}
		// Add trailing newline for nicer editing experience.
		data = append(data, '\n')
		if err := os.WriteFile(m.path, data, 0o644); err != nil {
			return "", fmt.Errorf("failed to create config file: %w", err)
		}
	}

	// If Stat returned an error other than NotExist, propagate it.
	if _, err := os.Stat(m.path); err != nil {
		return "", fmt.Errorf("failed to stat config file: %w", err)
	}

	return m.path, nil
}

// Read returns the contents of the config file.
func (m *FSManager) Read() ([]byte, error) {
	return os.ReadFile(m.path)
}

// Write replaces the config file content atomically (best-effort).
// Note: os.WriteFile is used for simplicity; callers may choose to implement
// more advanced atomic replace logic if desired.
func (m *FSManager) Write(b []byte, perm fs.FileMode) error {
	return os.WriteFile(m.path, b, os.FileMode(perm))
}
