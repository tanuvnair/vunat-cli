package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/tanuvnair/vunat-cli/internal/runner"
)

// StartCommand starts a named project using the provided Runner.
//
// Usage: vunat start <project_name>
type StartCommand struct {
	Runner *runner.Runner
}

// NewStartCommand constructs a StartCommand. If r is nil a default runner is created.
func NewStartCommand(r *runner.Runner) *StartCommand {
	if r == nil {
		r = runner.New()
	}
	return &StartCommand{Runner: r}
}

func (c *StartCommand) Name() string { return "start" }
func (c *StartCommand) Help() string { return "Start a project" }

func (c *StartCommand) Run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: vunat start <project_name>")
	}
	projectName := args[0]

	if c.Runner == nil {
		c.Runner = runner.New()
	}

	// Create a context which is cancelled on SIGINT/SIGTERM (Ctrl+C).
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	fmt.Printf("Starting project: %s\n\n", projectName)

	// Start the project by name. StartByName will load the project config and
	// launch the processes. It blocks until processes exit or the context is cancelled.
	if err := c.Runner.StartByName(ctx, projectName); err != nil {
		// Best-effort shutdown if StartByName returned an error.
		_ = c.Runner.Shutdown()
		return err
	}

	return nil
}
