# COSMIC Select

COSMIC Select is a small desktop utility for COSMIC Desktop on Wayland:

```text
Global shortcut -> Select screen area -> Local OCR -> Choose an action
```

The initial release is intentionally focused on four actions:

- Translate
- Copy
- Search
- Ask AI

COSMIC is the only supported and tested desktop environment. The application is built in Go, uses GTK4 for its UI, and integrates with COSMIC through D-Bus and XDG portals.

## Status

The repository currently contains the Go project skeleton. Product requirements and implementation boundaries are defined in [`docs/PRD - COSMIC Select.md`](docs/PRD%20-%20COSMIC%20Select.md), and repository-level guidance for coding agents is in [`AGENTS.md`](AGENTS.md).

The initial application foundation is documented in [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md). It provides a testable flow coordinator, explicit session states, local OCR text cleanup, minimal settings defaults, and interfaces for the future COSMIC portal adapters.

## Development

Run the basic checks with:

```sh
go test ./...
go build ./cmd/cosmic-select
```

GTK4, portal, and OCR dependencies will be introduced with the corresponding implementation milestones.
