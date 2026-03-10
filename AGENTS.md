# AGENTS.md — Development Guide for AI Coding Agents

## Project Overview

`apx` is a Go CLI tool that runs programs inside OS-level sandboxes (macOS `sandbox-exec` / seatbelt).
Module: `github.com/UsingCoding/apx`, Go 1.25.

Architecture: CLI (`cmd/`) → use-case layer (`internal/core/`) → domain (`internal/app/`, `internal/sandbox/`).
A lightweight DI container (`internal/container/`) wires dependencies from CLI flags.
Sandbox backends self-register via `init()` side-effect imports, keeping core logic decoupled.

---

## Build, Lint, and Test Commands

The project uses [mise](https://mise.jdx.dev/) as the task runner. Tool versions are pinned in `mise.toml`.

| Task | Command | Description |
|---|---|---|
| Default (all) | `mise run` | Runs `modules` → `build` → `test` |
| Build | `mise run build` | Builds binary via goreleaser (snapshot, single target) |
| Tidy modules | `mise run modules` | `go mod tidy` |
| All tests | `mise run test` | Runs `test:go-test` and `test:lint` in parallel |
| Go tests | `mise run test:go-test` | `go test ./...` |
| Lint | `mise run test:lint` | `golangci-lint run` |

### Running a Single Test

```bash
# Run one test function by name
go test ./internal/sandbox/seatbelt -run TestSnapshot_Default -v

# Run all snapshot tests
go test ./internal/sandbox/seatbelt -run Snapshot -v

# Update golden snapshot files
UPDATE_SNAPSHOTS=1 go test ./internal/sandbox/seatbelt -run Snapshot -v

# Run all tests with verbose output
go test ./... -v
```

Test naming convention: `Test<Subject>_<Scenario>` (e.g., `TestSnapshot_NetworkAllowed`).

---

## Code Style Guidelines

### Formatting

- Use standard `gofmt` formatting. No exceptions.
- Use `goimports` for import organization (enforces three-group import style).

### Import Grouping

Imports are organized in exactly three groups, separated by blank lines:

```go
import (
    // 1. stdlib
    "context"
    "os"

    // 2. third-party
    "github.com/pkg/errors"
    "github.com/urfave/cli/v3"

    // 3. local (github.com/UsingCoding/apx/...)
    "github.com/UsingCoding/apx/internal/app"
    "github.com/UsingCoding/apx/internal/sandbox"
)
```

- Alias imports only when necessary (e.g., `zeroslog "github.com/samber/slog-zerolog"`).
- Use blank side-effect imports (`_ "pkg"`) for sandbox plugin registration and `go:embed`.
- No dot imports.

### Naming Conventions

- **Packages**: short, lowercase, single-word — `core`, `sandbox`, `seatbelt`, `shellenv`, `app`, `container`.
- **Types / Structs**: `PascalCase` — `Exec`, `Policy`, `Registry`, `Seatbelt`, `APXTOML`.
- **Acronyms**: uppercase when standalone — `APXTOML`, `CMD`; otherwise follow Go conventions.
- **Interfaces**: `PascalCase`; single-method interfaces named after the method they define.
- **Exported functions/methods**: `PascalCase` — `LoadRegistry`, `MergePolicies`.
- **Unexported functions/methods**: `camelCase` — `makeProfile`, `expandPaths`, `mergeApps`.
- **Variables and parameters**: `camelCase` — `userConfigDir`, `argv0`, `apxtoml`.
- **Struct fields**: always `PascalCase` — `ROPaths`, `RWPaths`, `DenyPaths`, `Logger`.
- **TOML/serialization tags**: `camelCase` — `"roPaths"`, `"fullDiskReadAccess"`.
- **File names**: `snake_case.go`; platform-specific files use `<base>_<goos>.go` pattern.
- **Test functions**: `Test<Subject>_<Scenario>`.
- **Golden/testdata files**: `snake_case` — `fs_full_disk_read.golden`.

### Types

- Prefer **value receivers** — pointer receivers are avoided throughout.
- Use **structs as use-case namespaces**: dependencies are struct fields; `Do(ctx)` is the entry point. No constructor functions.

  ```go
  return core.Exec{
      Reg:    c.ApxRegistry,
      CMD:    cmd.Args().Slice(),
      Logger: logger(cmd),
  }.Do(ctx)
  ```

- Define **interfaces only at extension points** (e.g., `Sandbox`); keep them minimal (one method preferred).
- Use **type aliases for semantic clarity** — `type Env map[string]string`.
- Prefer **concrete types** over `interface{}` / `any` — no `any` is used in this codebase.
- No generics.
- Use `fs.FS` to abstract filesystem sources (embedded vs real OS), enabling multi-source registry loading.
- Use anonymous structs for local one-off serialization (e.g., JSON output in `version.go`).
- Struct tags use `toml:"..."` for TOML fields.

### Error Handling

Use **`github.com/pkg/errors`** for all error creation and wrapping — not `fmt.Errorf` with `%w`.

```go
// Wrap with context (preferred)
return errors.Wrap(err, "load registry")
return errors.Wrapf(err, "absolute path for ROPath %s", roPath)

// Create new errors
return errors.New("no command specified")
return errors.Errorf("sandbox %q for %q not found", s.Type, argv0)

// Sentinel check
if err != nil && !errors.Is(err, os.ErrNotExist) { ... }
```

- Always return a **zero value** of the return type alongside an error — `return Policy{}, errors.Wrap(...)`.
- Use `panic` **only for programmer errors** (invariant violations), not for recoverable runtime errors.
- Suppress linter warnings sparingly with `//nolint:<linter>` and only with justification.

### Testing

- Test files live in the **same package** as the code under test (white-box style).
- Use **`github.com/stretchr/testify/assert`** for assertions (not `require`).
- Use `t.Helper()` in shared helper functions for correct line reporting.
- Use `t.Setenv()` for hermetic environment isolation in tests.
- Snapshot/golden file tests: compare output against files in `testdata/*.golden`; update with `UPDATE_SNAPSHOTS=1`.
- Non-fatal assertions (`assert`) are used with explicit early `return` when the rest of the test would be invalid.

### Linting

The project uses `golangci-lint` v2 with a strict configuration (`.golangci.yml`). Key enabled linters:

- `revive`, `staticcheck`, `gocritic` (tags: `experimental`, `opinionated`)
- `gocognit`, `gocyclo`, `nestif` — enforce low complexity
- `gosec` — security checks
- `dupl` — detect code duplication
- `misspell`, `asciicheck`, `whitespace`, `unconvert`
- `gochecknoinits` — `init()` functions require `//nolint:gochecknoinits` with justification

Always run `golangci-lint run` before committing.

---

## Architecture Notes

- **Plugin registration**: Sandbox backends register themselves via `init()` in `init_<goos>.go` files. The main binary imports them as blank side-effects in `main_<goos>.go`. This keeps the core decoupled from concrete implementations.
- **Embedded registry**: Built-in app definitions (`*.apx.toml`) are embedded via `go:embed` in `registry/embed.go`.
- **Multi-source registry**: App registries are loaded from multiple locations (built-in embed, user config dir, project-local file) and merged.
- **SBPL policy generation**: `internal/sandbox/seatbelt/profile.go` generates macOS sandbox profiles from `sandbox.Policy` structs by composing embedded base SBPL templates.
