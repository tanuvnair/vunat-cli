package cli

import (
	"fmt"
	"sort"
	"strings"
)

// Command represents a single CLI subcommand.
// Implementations should be small and focused (SRP) and take dependencies
// via fields so they can be unit tested.
type Command interface {
	// Name returns the literal subcommand name, e.g. "start", "list", "config"
	Name() string
	// Run executes the command with the provided args (arguments after the command).
	// It returns an error for any failure; callers can wrap/print the error.
	Run(args []string) error
	// Help returns a short one-line usage/help string for help output.
	Help() string
}

// ErrUsage is returned when the CLI is invoked without a command.
var ErrUsage = fmt.Errorf("usage: vunat <command> [args]")

// Registry holds registered commands and dispatches invocation to them.
// Use NewRegistry to create a registry and Register to add commands.
type Registry struct {
	commands map[string]Command
}

// NewRegistry constructs a registry and optionally registers the provided commands.
func NewRegistry(cmds ...Command) *Registry {
	r := &Registry{
		commands: make(map[string]Command),
	}
	for _, c := range cmds {
		if c == nil {
			continue
		}
		r.commands[c.Name()] = c
	}
	return r
}

// Register adds or replaces a command in the registry.
func (r *Registry) Register(c Command) {
	if c == nil {
		return
	}
	r.commands[c.Name()] = c
}

// Commands returns a slice of registered commands (sorted by name).
func (r *Registry) Commands() []Command {
	names := make([]string, 0, len(r.commands))
	for n := range r.commands {
		names = append(names, n)
	}
	sort.Strings(names)

	out := make([]Command, 0, len(names))
	for _, n := range names {
		out = append(out, r.commands[n])
	}
	return out
}

// Run dispatches the args to the appropriate command.
// args is expected to be os.Args (or an equivalent slice).
// Behavior:
// - If no subcommand is provided, returns ErrUsage or runs the "help" command if present.
// - If a known subcommand is provided, calls its Run with the remaining args.
// - If unknown command, returns an error indicating the unknown command.
func (r *Registry) Run(args []string) error {
	if len(args) < 2 {
		// If a help command is registered, show help by default.
		if helpCmd, ok := r.commands["help"]; ok {
			return helpCmd.Run(nil)
		}
		return ErrUsage
	}

	name := args[1]
	if cmd, ok := r.commands[name]; ok {
		return cmd.Run(args[2:])
	}

	// Unknown command: return informative error and suggestion of available commands.
	available := r.availableCommandNames()
	return fmt.Errorf("unknown command: %s\n\nAvailable commands:\n%s", name, available)
}

func (r *Registry) availableCommandNames() string {
	names := make([]string, 0, len(r.commands))
	for n := range r.commands {
		names = append(names, n)
	}
	sort.Strings(names)
	var b strings.Builder
	for _, n := range names {
		if cmd := r.commands[n]; cmd != nil {
			fmt.Fprintf(&b, "  %s\t%s\n", n, cmd.Help())
		} else {
			fmt.Fprintf(&b, "  %s\n", n)
		}
	}
	return b.String()
}
