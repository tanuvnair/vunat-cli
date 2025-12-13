package cli

import "fmt"

func Run(args []string) error {
	if len(args) < 2 {
		return usageError()
	}

	switch args[1] {
	case "start":
		return runStart(args[2:])
	case "help":
		return printHelp()
	default:
		return fmt.Errorf("unknown command: %s", args[1])
	}
}

func usageError() error {
	return fmt.Errorf("usage: vunat start <project_name>")
}

func printHelp() error {
	fmt.Println("vunat-cli - your personal CLI for quick-starting development projects")
	fmt.Println()

	fmt.Println("usage:")
	fmt.Println("	vunat start <project_name>")

	return nil
}
