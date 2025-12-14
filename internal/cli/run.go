package cli

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/tanuvnair/vunat-cli/internal/cli/commands"
	"github.com/tanuvnair/vunat-cli/internal/config"
	"github.com/tanuvnair/vunat-cli/internal/launcher"
	"github.com/tanuvnair/vunat-cli/internal/runner"
)

// Run constructs the application's dependencies, registers commands into the
// Registry and dispatches the provided args to the selected command.
func Run(args []string) error {
	// Construct shared dependencies
	cfgMgr := config.NewFSManager("")
	osLauncher := launcher.NewOSLauncher(false)
	runr := runner.New()

	// Build registry and register commands
	reg := NewRegistry()

	// Register command implementations (start, list, config)
	reg.Register(commands.NewStartCommand(runr))
	reg.Register(commands.NewListCommand())
	reg.Register(commands.NewConfigCommand(cfgMgr, osLauncher))

	// Register dynamic help command. The provider builds help text from the
	// registry contents so help is always up-to-date.
	helpProvider := func() string {
		var b strings.Builder
		b.WriteString("vunat-cli - your personal CLI for quick-starting development projects\n\n")
		b.WriteString("usage:\n")
		b.WriteString("  vunat <command> [args]\n\n")
		b.WriteString("Available commands:\n")

		// Use tabwriter to align command names and their descriptions in columns.
		w := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)
		for _, c := range reg.Commands() {
			fmt.Fprintf(w, "  %s\t%s\n", c.Name(), c.Help())
		}
		_ = w.Flush()

		b.WriteString("\n")
		return b.String()
	}
	reg.Register(commands.NewHelpCommand(helpProvider))

	// Dispatch to the registry
	return reg.Run(args)
}
