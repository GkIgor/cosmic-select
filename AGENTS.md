# AGENTS.md

## First Rule

Before planning, editing, refactoring, or generating implementation tasks, read:

```text
docs/PRD - COSMIC Select.md
```

Treat the PRD as the product source of truth. This file is only a repository-level context guide for coding agents. If this file and the PRD disagree, follow the PRD and update this file only when explicitly asked.

## Project Summary

COSMIC Select is a small desktop utility for COSMIC Desktop on Wayland.

The user presses a global shortcut, selects a rectangular screen area, the app captures that area, runs local OCR, and then lets the user choose one of four actions:

- Translate
- Copy
- Search
- Ask AI

The intended flow is:

```text
Shortcut -> Select screen area -> OCR -> Choose action
```

Keep the product focused. Do not turn it into a general assistant, launcher, note-taking app, clipboard manager, screenshot manager, or cross-desktop Linux utility.

## Platform And Stack

- Supported desktop: COSMIC only
- Display server: Wayland
- Primary language: Go
- UI toolkit: GTK4
- Desktop integration: D-Bus and XDG portals
- Screenshot capture: XDG Screenshot Portal, preferably through `xdg-desktop-portal-cosmic`
- Global shortcut: XDG Global Shortcuts Portal, preferably through `xdg-desktop-portal-cosmic`
- OCR: local OCR only

COSMIC is the only supported and tested desktop environment. Cross-desktop compatibility is not a goal.

GTK4 is an implementation detail for rendering UI. Using GTK4 does not imply GNOME support.

## Privacy Constraints

- OCR must run locally.
- Do not upload screenshots by default.
- Do not upload raw screen captures to third-party services.
- Only send extracted text to external services when the user explicitly chooses an action that requires it, such as Translate, Search, or Ask AI.
- Keep network behavior obvious, minimal, and action-driven.
- Do not add background sync, analytics, telemetry, remote history, or cloud storage unless the PRD is explicitly changed.

## Architecture Boundaries

Keep the code organized around clear boundaries:

- Portal integration: screenshot capture, global shortcuts, permission/session handling.
- Selection UI: screen overlay, rectangle selection, cancel/confirm behavior.
- OCR pipeline: image preprocessing, OCR execution, text cleanup.
- Actions: Translate, Copy, Search, Ask AI.
- App state and configuration: shortcut, preferences, provider choices if needed.
- UI shell: menus, action picker, errors, settings.

Avoid mixing portal calls, OCR logic, action execution, and GTK event handling in the same module.

Prefer small interfaces around external systems so portal, OCR, clipboard, browser/search, and AI/translation behavior can be tested or replaced without rewriting the app.

## Out Of Scope

Do not implement these unless the PRD is explicitly changed:

- GNOME, KDE Plasma, X11, Windows, or macOS support
- Rust or `libcosmic` rewrite
- General-purpose AI assistant features
- Chat history
- Screenshot library or screenshot manager
- Full clipboard manager
- Document translation workflows
- Browser extension
- Mobile app
- Cloud OCR
- Automatic background OCR
- Always-on screen monitoring
- Telemetry or usage analytics
- Plugin marketplace
- Multi-user accounts

## Implementation Principles

- Build the narrow user flow first.
- Prefer COSMIC/Wayland/XDG portal correctness over portability.
- Keep permissions explicit and user-driven.
- Fail clearly when required portals are unavailable.
- Keep UI lightweight and visually compatible with COSMIC.
- Make cancellation fast and reliable at every step.
- Keep local-only behavior as the default.
- Add configuration only when it supports the core flow.
- Do not introduce framework-heavy abstractions before the product shape requires them.

## Coding Guardrails

- Use Go for application code unless there is a documented reason not to.
- Keep packages small and named by responsibility, not by vague layers.
- Keep GTK-specific code near the UI boundary.
- Keep D-Bus/XDG portal code isolated from business logic.
- Wrap external commands, OCR engines, and network providers behind interfaces.
- Prefer structured errors with user-appropriate messages at the UI boundary.
- Do not silently fall back to unsupported desktop environments.
- Do not add dependencies casually; each dependency should serve the PRD.
- Write tests around OCR text cleanup, action routing, configuration, and portal/OCR abstractions where practical.
- Do not store screenshots unless the PRD explicitly adds that behavior.
- Do not log extracted text or screenshot paths by default.

## Suggested Initial Milestones

1. Repository skeleton: Go module, package layout, basic build/test commands, minimal GTK app window.
2. Portal capability checks: detect COSMIC/Wayland assumptions, screenshot portal availability, global shortcuts portal availability.
3. Global shortcut prototype: register shortcut and trigger an app callback.
4. Selection overlay prototype: draw and confirm/cancel a rectangular selection.
5. Screenshot capture: capture selected area through the XDG Screenshot Portal.
6. Local OCR pipeline: run OCR on captured image and return cleaned text.
7. Action picker: show Translate, Copy, Search, and Ask AI options after OCR.
8. Implement Copy and Search first, because they are simpler and validate the flow.
9. Implement Translate and Ask AI with explicit user-triggered network boundaries.
10. Settings and polish: shortcut configuration, provider configuration, error states, COSMIC visual integration.

## Agent Behavior

When working in this repository:

- Read the PRD first.
- Keep changes aligned with the current milestone.
- Prefer small, reviewable changes.
- State when a proposed change expands scope.
- Do not infer new product requirements from implementation convenience.
- If a decision is unclear, choose the smallest option that preserves the PRD.
