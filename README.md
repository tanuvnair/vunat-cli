# vunat-cli

A small, opinionated CLI to quickly start development projects composed of multiple process groups (frontend, backend, workers, etc.). vunat-cli reads a JSON-based project registry and can:

- Start groups of commands and supervise their lifecycle
- Stream stdout/stderr prefixed with the group name
- Open or create the per-user config (`~/.vunat/config.json`)
- Provide a simple command registry and help output

This repository contains a compact, testable Go implementation with clear separation of concerns: config manager, launcher, runner, and CLI registry.

---

## Contents

- `cmd/vunat` — CLI entrypoint
- `internal/cli` — registry + command implementations
- `internal/config` — filesystem-backed config manager
- `internal/projects` — config loader / project registry
- `internal/runner` — process supervision and output streaming
- `internal/launcher` — OS-aware opener for files/URLs

---

## Quick links

- Example config: `config.example.json` (in the repo root)

Example snippet:
```json
{
  "projects": {
    "gradepoint": [
      {
        "name": "frontend",
        "absolutePath": "",
        "commands": ["npm run dev"]
      },
      {
        "name": "backend",
        "absolutePath": "",
        "commands": [
          "go run ./cmd/api/main.go",
          "go run ./cmd/scheduler/main.go",
          "npx prisma@6 studio -y"
        ]
      }
    ]
  }
}
```

---

## Build

- Unix / macOS:
```sh
./build.sh
```

- Windows:
```cmd
build.bat
```

You can also run in development mode with the Go toolchain:
```sh
go run ./cmd/vunat help
```

---

## Install (recommended workflow)

To be able to call the tool as `vunat` from any shell, put the compiled binary in a small per-user bin directory and add it to your PATH. This README recommends `~/.vunat` as the install directory so it
 is colocated with the config.

1. Build the binaries using the included scripts (or build for your OS with `go build`).
2. Move or copy the appropriate binary into `~/.vunat` and name it `vunat`.

Example (Linux/macOS):
```sh
# create install dir
mkdir -p "$HOME/.vunat"

# move the built binary into the install dir and make executable
mv ./bin/vunat-linux-amd64 "$HOME/.vunat/vunat"
chmod +x "$HOME/.vunat/vunat"
```

Example (Windows PowerShell):
```powershell
# create install dir (in %USERPROFILE%)
New-Item -ItemType Directory -Path "$env:USERPROFILE\.vunat" -Force

# move the built binary (example)
Move-Item .\bin\vunat-windows-amd64.exe "$env:USERPROFILE\.vunat\vunat.exe"
```

3. Add `~/.vunat` (or `%USERPROFILE%\.vunat` on Windows) to your PATH so that the shell can find `vunat`.

Temporarily (current shell session) — Unix/macOS:
```sh
export PATH="$HOME/.vunat:$PATH"
```

Persistently (add to `~/.bashrc`, `~/.zshrc`, or other shell rc):
```sh
echo 'export PATH
="$HOME/.vunat:$PATH"' >> ~/.bashrc
# then reload or restart your shell:
source ~/.bashrc
```

Windows (PowerShell, persistent):
```powershell
# Add to user PATH using setx (may require reopening terminals)
setx PATH "$env:USERPROFILE\.vunat;$env:PATH"
```

After the directory is in your PATH you can run the CLI directly:

```sh
vunat help
vunat list
vunat start <project_name>
vunat config
```

---

## Usage

- Show help (dynamic, generated from registered commands):
```sh
vunat help
```

- List registered projects:
```sh
vunat list
```

- Start a project:
```sh
vunat start <project_name>
```

- Open or create the config file:
```sh
vunat config
```

---

## Configuration

- The per-user configuration file is `~/.vunat/config.json`.
- Use the provided `config.example.json` as a template.
- The config format:
  - Top-level `projects` object
  - Each key under `projects` is a project name that maps to an array of command groups.
  - A command group contains:
    - `name` — human-readable group name
    - `absolutePath` — directory where the commands will run (empty allowed)
    - `commands` — array of shell command strings

---

## Editing the config

- If the `EDITOR` environment variable is set, `vunat config` will invoke that editor and wait for it to exit.
- Otherwise, `vunat config` uses the platform default opener:
  - Windows: `cmd /c start "" <path>`
  - macOS: `open <path>`
  - Linux: `xdg-open <path>`

---

## Architecture notes

- CLI registry (`internal/cli`)
  - Small `Command` interface with `Name()`, `Run(args)`, and `Help()`.
  - `Registry` holds commands and dispatches based on `os.Args`.

- Config manager (`internal/config`)
  - Exposes a `Manager` interface and `FSManager` implementation that ensures config directory/file exist and reads/writes the JSON file.

- Projects loader (`internal/projects`)
  - Reads `~/.vunat/config.json`, unmarshals into typed structs, and exposes `Get` and `GetAll` helpers.

- Runner (process supervision) (`internal/runner`)
  - Starts groups sequentially and commands in a group concurrently.
  - Streams each process' stdout/stderr prefixed with the group name.
  - Cancels remaining processes on first failure and attempts to kill already-started children.

- Launcher (`internal/launcher`)
  - Provides `OSLauncher` to open files/URLs using platform-specific commands, with an option to wait for the opener to exit.
