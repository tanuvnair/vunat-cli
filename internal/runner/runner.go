package runner

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/tanuvnair/vunat-cli/internal/projects"
)

// Runner supervises processes started for a project.
type Runner struct {
	mu    sync.Mutex
	procs []*exec.Cmd
}

// New creates a new Runner.
func New() *Runner {
	return &Runner{
		procs: make([]*exec.Cmd, 0, 8),
	}
}

// StartByName looks up the project by name and starts it via Runner.Start.
// This is a small convenience method so callers (like the start command) can
// directly request starting a project by its registered name.
func (r *Runner) StartByName(ctx context.Context, name string) error {
	proj, err := projects.Get(name)
	if err != nil {
		return err
	}
	return r.Start(ctx, proj)
}

// Start launches all command groups in the provided project.
// - Groups are started sequentially; commands within a group are started concurrently.
// - Each command is started with exec.CommandContext so it is cancelled when ctx is done.
// - Returns nil if all processes exit cleanly, or the first non-nil error encountered.
func (r *Runner) Start(ctx context.Context, proj projects.Project) error {
	// derive cancellable context so we can cancel on first error
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// channel for first process error
	errCh := make(chan error, 1)

	var wg sync.WaitGroup

	// Helper to stream output
	startStream := func(rc io.ReadCloser, out io.Writer, prefix string) {
		go func() {
			defer rc.Close()
			scanner := bufio.NewScanner(rc)
			for scanner.Scan() {
				fmt.Fprintln(out, prefix+scanner.Text())
			}
			// ignore scanner error here; process Wait will surface failure
		}()
	}

	// Start groups sequentially
	for _, group := range proj {
		fmt.Printf("[%s] Starting in: %s\n", group.Name, group.AbsolutePath)

		for _, cmdStr := range group.Commands {
			parts := splitFields(cmdStr)
			if len(parts) == 0 {
				continue
			}

			cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
			if group.AbsolutePath != "" {
				cmd.Dir = group.AbsolutePath
			}

			stdoutPipe, err := cmd.StdoutPipe()
			if err != nil {
				// immediate failure to obtain pipe
				_ = r.shutdownAll()
				return fmt.Errorf("failed to obtain stdout pipe for %q: %w", cmdStr, err)
			}
			stderrPipe, err := cmd.StderrPipe()
			if err != nil {
				_ = r.shutdownAll()
				return fmt.Errorf("failed to obtain stderr pipe for %q: %w", cmdStr, err)
			}

			if err := cmd.Start(); err != nil {
				_ = r.shutdownAll()
				return fmt.Errorf("failed to start command %q: %w", cmdStr, err)
			}

			// record process for later shutdown
			r.addProc(cmd)

			// prefix output with group name
			prefix := fmt.Sprintf("[%s] ", group.Name)
			startStream(stdoutPipe, os.Stdout, prefix)
			startStream(stderrPipe, os.Stderr, prefix)

			// wait for process in background
			wg.Add(1)
			go func(cmd *exec.Cmd, desc string) {
				defer wg.Done()
				if err := cmd.Wait(); err != nil {
					// Try to send the first error observed; do not block if channel already has an error.
					select {
					case errCh <- fmt.Errorf("process %q exited with error: %w", desc, err):
						// cancel context so other processes get signaled
						cancel()
						// attempt to shutdown remaining processes
						_ = r.shutdownAll()
					default:
						// already have an error, nothing to do
					}
				}
			}(cmd, strings.Join(parts, " "))
		}

		fmt.Printf("[%s] Started\n\n", group.Name)
	}

	// Wait for all processes to finish or for an error/cancellation.
	doneCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	select {
	case <-ctx.Done():
		// Context was cancelled (either externally or due to a process error).
		// Return context error if no specific process error was reported.
		select {
		case e := <-errCh:
			return e
		default:
			return ctx.Err()
		}
	case e := <-errCh:
		// First process-level error
		return e
	case <-doneCh:
		// All processes exited normally
		return nil
	}
}

// Shutdown attempts to kill all started processes.
func (r *Runner) Shutdown() error {
	return r.shutdownAll()
}

// addProc records a started process (thread-safe).
func (r *Runner) addProc(cmd *exec.Cmd) {
	r.mu.Lock()
	r.procs = append(r.procs, cmd)
	r.mu.Unlock()
}

// shutdownAll kills all recorded processes and clears the list.
// It returns the first error encountered while killing processes, if any.
func (r *Runner) shutdownAll() error {
	r.mu.Lock()
	procs := make([]*exec.Cmd, len(r.procs))
	copy(procs, r.procs)
	r.procs = r.procs[:0]
	r.mu.Unlock()

	var firstErr error
	for _, cmd := range procs {
		if cmd == nil || cmd.Process == nil {
			continue
		}
		if err := cmd.Process.Kill(); err != nil {
			if firstErr == nil {
				firstErr = fmt.Errorf("failed to kill process %d: %w", cmd.Process.Pid, err)
			}
		}
	}
	return firstErr
}

func splitFields(s string) []string {
	// Naive splitting; this can be enhanced to handle quotes if needed.
	return strings.Fields(s)
}
