package launcher

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Launcher opens a path in the user's editor/handler. Implementations should
// keep side-effects minimal so they can be replaced or mocked in tests.
type Launcher interface {
	// Open attempts to open the provided path with the platform default
	// application (editor/file-association). Path is typically a filesystem
	// path but may also be a URL on platforms that support it.
	Open(path string) error
}

// OSLauncher is a small, platform-aware Launcher implementation that uses
// the system's default commands to open files/URLs.
//
// Behavior:
// - Windows: executes `cmd /c start "" <path>` (uses cmd to invoke the shell builtin).
// - macOS: executes `open <path>`.
// - Linux/other: executes `xdg-open <path>`.
type OSLauncher struct {
	// Wait controls whether Open waits for the invoked command to exit.
	// - If true, Open calls cmd.Run() and returns only after the command exits.
	// - If false, Open calls cmd.Start() and returns immediately after starting.
	Wait bool
}

// NewOSLauncher creates an OSLauncher. Pass wait=true to block until the
// platform opener exits; pass wait=false to return immediately after starting.
func NewOSLauncher(wait bool) *OSLauncher {
	return &OSLauncher{Wait: wait}
}

// Open opens the provided path using the OS default opener. It returns a
// wrapped error if the platform is unsupported or if the underlying command
// fails to start/execute.
func (l *OSLauncher) Open(path string) error {
	if path == "" {
		return fmt.Errorf("launcher: empty path")
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		// 'start' is a cmd.exe builtin. The empty string argument after 'start'
		// is the window title; it's required if the path may begin with a quote.
		cmd = exec.Command("cmd", "/c", "start", "", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default:
		// Most desktop Linux environments provide xdg-open. If it's not present
		// the command will fail; callers can detect and provide an alternative.
		cmd = exec.Command("xdg-open", path)
	}

	if l.Wait {
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("launcher: failed to run opener: %w", err)
		}
		return nil
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("launcher: failed to start opener: %w", err)
	}
	return nil
}
