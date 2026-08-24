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

The repository contains the Go foundation and the first GTK4 window shell. Product requirements and implementation boundaries are defined in [`docs/PRD - COSMIC Select.md`](docs/PRD%20-%20COSMIC%20Select.md), and repository-level guidance for coding agents is in [`AGENTS.md`](AGENTS.md).

The initial application foundation is documented in [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md). It provides a testable flow coordinator, explicit session states, local OCR text cleanup, minimal settings defaults, and interfaces for the future COSMIC portal adapters.

## Development

Run the core checks with:

```sh
go test ./...
go build ./cmd/cosmic-select
```

The GTK4 window is built with the `gtk4` build tag and requires GTK4 development files (`gtk4`, `gobject-introspection`, and `pkg-config`) on the system:

```sh
go run -tags gtk4 ./cmd/cosmic-select
go test -tags gtk4 ./...
```

The current window is intentionally small: it is the UI shell for checking COSMIC portal support. Selection overlay, local OCR, and the action picker are added in their respective milestones.
