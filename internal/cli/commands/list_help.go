package commands

import (
	"fmt"

	"github.com/tanuvnair/vunat-cli/internal/projects"
)

// ListCommand lists all registered projects from the config.
type ListCommand struct{}

// NewListCommand constructs a ListCommand.
func NewListCommand() *ListCommand {
	return &ListCommand{}
}

func (c *ListCommand) Name() string { return "list" }
func (c *ListCommand) Help() string { return "List all registered projects" }

func (c *ListCommand) Run(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("usage: vunat list")
	}

	allProjects := projects.GetAll()

	if len(allProjects) == 0 {
		fmt.Println("No projects registered. Add projects to ~/.vunat/config.json")
		return nil
	}

	fmt.Println("Registered projects:")
	for name, project := range allProjects {
		fmt.Printf("  %s\n", name)
		for _, group := range project {
			fmt.Printf("    [%s] in %s\n", group.Name, group.AbsolutePath)
			for _, cmd := range group.Commands {
				fmt.Printf("      → %s\n", cmd)
			}
		}
	}
	return nil
}

// HelpCommand prints dynamic help text supplied by the CLI registry.
// The commands package purposely avoids importing the registry to prevent import cycles;
// instead, the registry should construct a HelpCommand and pass a provider function.
type HelpCommand struct {
	// Provider returns the help text to print. If nil, HelpCommand falls back to a simple message.
	Provider func() string
}

// NewHelpCommand constructs a HelpCommand. Pass a provider that returns the full help text.
func NewHelpCommand(provider func() string) *HelpCommand {
	return &HelpCommand{Provider: provider}
}

func (h *HelpCommand) Name() string { return "help" }
func (h *HelpCommand) Help() string { return "Show this help message" }

func (h *HelpCommand) Run(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("usage: vunat help")
	}
	if h.Provider == nil {
		// Basic fallback help
		fmt.Println("vunat-cli - your personal CLI for quick-starting development projects")
		fmt.Println()
		fmt.Println("usage:")
		fmt.Println("  vunat start <project_name>  Start a project")
		fmt.Println("  vunat list                  List all registered projects")
		fmt.Println("  vunat config                Open the config file in your default editor")
		fmt.Println("  vunat help                  Show this help message")
		return nil
	}
	fmt.Print(h.Provider())
	return nil
}
