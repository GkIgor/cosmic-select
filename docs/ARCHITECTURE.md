# Initial Architecture

The application is being built around the product flow:

```text
Global shortcut -> screen selection -> local OCR -> one of four actions
```

The current foundation deliberately has no GTK4 or D-Bus dependency. Those integrations belong in adapters and will be added when their corresponding milestones are implemented. This keeps the core flow testable without requiring a running COSMIC session.

## Packages

```text
cmd/cosmic-select   process entry point
internal/app        flow coordinator
internal/config     local core-flow settings
internal/domain     session state and supported actions
internal/ocr        OCR text cleanup
internal/ports      contracts for portal, shortcut, and OCR adapters
```

## Boundaries

- `internal/domain` has no desktop or network dependencies.
- `internal/app` coordinates the flow and never renders UI or calls D-Bus directly.
- `internal/ports` describes external capabilities without choosing an implementation.
- `internal/ocr` receives local OCR output and returns normalized text.
- Portal adapters must pass selected image bytes directly to OCR and must not persist screenshots.
- External translation and AI calls will receive extracted text only after an explicit user action.

## Next implementation step

Add capability checks and portal adapters for COSMIC's screenshot and global shortcut portals. GTK4 UI should then connect to the existing coordinator rather than moving application logic into widgets.
