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
internal/ui         GTK4 window shell (build tag: gtk4)
internal/portal     COSMIC/XDG portal adapters
internal/cosmicshortcut COSMIC-native shortcut fallback
```

## Boundaries

- `internal/domain` has no desktop or network dependencies.
- `internal/app` coordinates the flow and never renders UI or calls D-Bus directly.
- `internal/ports` describes external capabilities without choosing an implementation.
- `internal/ocr` receives local OCR output and returns normalized text.
- Portal adapters must pass selected image bytes directly to OCR and must not persist screenshots.
- External translation and AI calls will receive extracted text only after an explicit user action.
- GTK4 code is isolated in `internal/ui`; the domain and coordinator remain toolkit-independent.
- When the Global Shortcuts portal is unavailable, the user may explicitly install a COSMIC-native `Super + Shift + S` custom shortcut that launches `cosmic-select --activate`.

## Next implementation step

The first portal foundation is now in `internal/portal`: it validates COSMIC/Wayland, checks Screenshot and GlobalShortcuts capabilities, delegates interactive area capture, and registers one global shortcut. Screenshot portal v2 is supported through its interactive mode when the v3 `AvailableTargets` property is absent. GTK4 UI should connect to these adapters through the existing coordinator rather than moving application logic into widgets.

Portal failures are returned as explicit errors. The application does not silently fall back to X11, compositor-specific capture, or another desktop environment.
