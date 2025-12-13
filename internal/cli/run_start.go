package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

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

	fmt.Printf("Starting project: %s\n\n", projectName)

	// Iterate through each command group sequentially
	for _, group := range project {
		fmt.Printf("[%s] Starting in: %s\n", group.Name, group.AbsolutePath)

		// Run all commands in this group concurrently
		var wg sync.WaitGroup
		var mu sync.Mutex
		var firstError error

		for _, cmdStr := range group.Commands {
			wg.Add(1)
			go func(cmdString string, groupName string, workDir string) {
				defer wg.Done()

				// Parse command string into executable and args
				parts := strings.Fields(cmdString)
				if len(parts) == 0 {
					return
				}

				execCmd := exec.Command(parts[0], parts[1:]...)
				if workDir != "" {
					execCmd.Dir = workDir
				}

				// Create pipes for stdout and stderr to prefix with group name
				stdoutPipe, err := execCmd.StdoutPipe()
				if err != nil {
					mu.Lock()
					if firstError == nil {
						firstError = err
					}
					mu.Unlock()
					return
				}

				stderrPipe, err := execCmd.StderrPipe()
				if err != nil {
					mu.Lock()
					if firstError == nil {
						firstError = err
					}
					mu.Unlock()
					return
				}

				// Start the command
				if err := execCmd.Start(); err != nil {
					mu.Lock()
					if firstError == nil {
						firstError = err
					}
					mu.Unlock()
					return
				}

				// Prefix output with group name
				prefix := fmt.Sprintf("[%s] ", groupName)

				// Read and print stdout with prefix
				go func() {
					scanner := bufio.NewScanner(stdoutPipe)
					for scanner.Scan() {
						fmt.Printf("%s%s\n", prefix, scanner.Text())
					}
				}()

				// Read and print stderr with prefix
				go func() {
					scanner := bufio.NewScanner(stderrPipe)
					for scanner.Scan() {
						fmt.Fprintf(os.Stderr, "%s%s\n", prefix, scanner.Text())
					}
				}()

				// Wait for command to finish
				if err := execCmd.Wait(); err != nil {
					mu.Lock()
					if firstError == nil {
						firstError = err
					}
					mu.Unlock()
				}
			}(cmdStr, group.Name, group.AbsolutePath)
		}

		// Wait for all commands in this group to finish
		wg.Wait()

		if firstError != nil {
			return fmt.Errorf("error in group %s: %w", group.Name, firstError)
		}

		fmt.Printf("[%s] Completed\n\n", group.Name)
	}

	fmt.Printf("All groups completed successfully!\n")
	return nil
}
