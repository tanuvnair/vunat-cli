package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/tanuvnair/vunat-cli/internal/config"
	"github.com/tanuvnair/vunat-cli/internal/launcher"
)

// ConfigCommand implements the 'config' subcommand.
//
// Behavior:
// - Validates there are no extra args (usage: `vunat config`).
// - Ensures the config file exists via config.Manager.Ensure().
// - If $EDITOR is set, launches the editor and waits for it to exit (so the user can edit).
// - Otherwise uses the injected launcher.Launcher to open the config file with the platform default.
type ConfigCommand struct {
	cfg      config.Manager
	launcher launcher.Launcher
	// WaitForEditor controls whether we wait for the platform opener to exit.
	// When using $EDITOR we always wait. For the platform launcher we rely on the
	// Launcher's own behavior (it may start and return immediately).
	// This field is provided for tests or alternate behaviors.
	WaitForPlatform bool
}

// NewConfigCommand constructs a ConfigCommand with the provided dependencies.
func NewConfigCommand(cfg config.Manager, l launcher.Launcher) *ConfigCommand {
	return &ConfigCommand{
		cfg:             cfg,
		launcher:        l,
		WaitForPlatform: false,
	}
}

func (c *ConfigCommand) Name() string {
	return "config"
}

func (c *ConfigCommand) Help() string {
	return "Open the config file in your editor"
}

func (c *ConfigCommand) Run(args []string) error {
	// No args expected
	if len(args) != 0 {
		return fmt.Errorf("usage: vunat config")
	}

	if c.cfg == nil {
		return fmt.Errorf("config manager not provided")
	}

	// Ensure config exists and get path
	path, err := c.cfg.Ensure()
	if err != nil {
		return fmt.Errorf("failed to ensure config file: %w", err)
	}

	// If the user has EDITOR set, prefer that and wait for it to exit.
	editor := strings.TrimSpace(os.Getenv("EDITOR"))
	if editor != "" {
		parts := strings.Fields(editor)
		// Append the config path as the last argument to the editor command.
		cmdArgs := append(parts[1:], path)
		cmd := exec.Command(parts[0], cmdArgs...)
		// Attach stdio so the editor can interact with the terminal.
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("$EDITOR command failed: %w", err)
		}
		return nil
	}

	// Otherwise use the platform launcher
	if c.launcher == nil {
		return fmt.Errorf("no launcher available to open the config file")
	}

	// If the launcher should be waited for, and the launcher implementation
	// does not itself block, consider wrapping here. We call Open and return
	// its error (the OSLauncher has an option to wait if constructed so).
	if err := c.launcher.Open(path); err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}

	// If we want to block until the platform opener exits but the launcher
	// doesn't support waiting, there's no portable way here without changing
	// the launcher interface. The OSLauncher provided elsewhere exposes a
	// Wait flag if needed.
	_ = c.WaitForPlatform // kept for clarity and future extension

	return nil
}
