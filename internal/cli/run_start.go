package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/tanuvnair/vunat-cli/internal/projects"
)

func runStart(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: vunat start <project_name>")
	}

	projectName := args[0]

	project, err := projects.Get(projectName)
	if err != nil {
		return err
	}

	fmt.Printf("Starting project: %s\n", project.Name)

	for _, cmdArgs := range project.Commands {
		fmt.Printf("→ %s\n", cmdArgs)

		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
