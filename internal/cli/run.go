package cli

import (
	"fmt"

	"github.com/tanuvnair/vunat-cli/internal/projects"
)

func Run(args []string) error {
	if len(args) < 2 {
		return usageError()
	}

	switch args[1] {
	case "start":
		return runStart(args[2:])
	case "list":
		return listProjects()
	case "help":
		return printHelp()
	default:
		return fmt.Errorf("unknown command: %s", args[1])
	}
}

func usageError() error {
	return fmt.Errorf("usage: vunat <command> [args]\n\nCommands:\n  start <project_name>  Start a project\n  list                  List all projects\n  help                  Show help")
}

func printHelp() error {
	fmt.Println("vunat-cli - your personal CLI for quick-starting development projects")
	fmt.Println()
	fmt.Println("usage:")
	fmt.Println("  vunat start <project_name>  Start a project")
	fmt.Println("  vunat list                  List all registered projects")
	fmt.Println("  vunat help                  Show this help message")
	return nil
}

func listProjects() error {
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
