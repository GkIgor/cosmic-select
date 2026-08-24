# PRD — COSMIC Select

## 1. Overview

**COSMIC Select** is a lightweight desktop utility built exclusively for the COSMIC Desktop Environment.

It allows the user to press a global shortcut, select any rectangular area of the screen, extract visible text through OCR, and immediately perform one of four actions:

- Translate
- Copy
- Search
- Ask AI

The intended interaction is:

```text
Shortcut
   ↓
Select screen area
   ↓
OCR
   ↓
Choose action
```

COSMIC Select is not intended to become a general-purpose assistant or productivity platform.

Its scope is intentionally small.

---

## 2. Platform

COSMIC Select targets:

- COSMIC Desktop
- Wayland
- `xdg-desktop-portal-cosmic`

COSMIC is the **only supported and tested desktop environment**.

Cross-desktop compatibility is not a goal.

Generic Linux libraries may be used internally when they provide the required COSMIC/Wayland integration.

The application does not need to support:

- GNOME
- KDE Plasma
- X11
- Windows
- macOS

If COSMIC provides a better integration path than a generic abstraction, the COSMIC-specific approach should be preferred.

---

## 3. Implementation Stack

Primary implementation language:

```text
Go
```

UI toolkit:

```text
GTK4
```

Desktop integration:

```text
Wayland
D-Bus
xdg-desktop-portal-cosmic
XDG Screenshot Portal
XDG Global Shortcuts Portal
```

The project does not require Rust or `libcosmic`.

GTK4 is used as an implementation detail for rendering the application UI.

Using GTK4 does not imply GNOME support.

COSMIC remains the only officially supported desktop environment.

---

## 4. Problem

Desktop applications frequently display text that cannot be selected normally.

Examples include:

- games;
- images;
- videos;
- remote desktop sessions;
- rendered documents;
- screenshots;
- custom UI components;
- applications with broken or unavailable text selection.

Performing an action on this content currently requires several manual steps.

For example:

```text
Take screenshot
→ crop image
→ run OCR
→ copy text
→ open translator
→ paste
```

COSMIC Select reduces this to:

```text
Shortcut
→ select
→ Translate
```

The experience is inspired by features such as Google Circle to Search and Samsung AI Select, adapted specifically for COSMIC Desktop.

---

## 5. Goal

Allow the user to perform useful actions on any visible text on the screen without leaving the current workflow.

The core interaction must remain:

**shortcut → selection → action**

Anything that does not directly improve this flow should be treated as outside the initial scope.

---

## 6. Core Actions

The initial release contains exactly four primary actions:

1. Translate
2. Copy
3. Search
4. Ask AI

These actions operate on text extracted from the selected screen region.

---

## 7. Main User Flow

### 7.1 Activation

The user presses a configurable global shortcut.

Example:

```text
Super + Shift + Space
```

COSMIC Select activates while the user is inside any application.

The shortcut must work globally without requiring COSMIC Select to have focus.

---

### 7.2 Region Selection

The user clicks and drags to select a rectangular portion of the screen.

Example:

```text
┌───────────────────────────────────┐
│                                   │
│     ┌─────────────────────┐       │
│     │ Selected screen     │       │
│     │ region              │       │
│     └─────────────────────┘       │
│                                   │
└───────────────────────────────────┘
```

The application should use COSMIC's screenshot infrastructure through `xdg-desktop-portal-cosmic`.

COSMIC Select should not implement compositor-level screen capture itself unless the COSMIC portal lacks functionality required by the product.

---

## 8. OCR

After the region is selected, COSMIC Select runs OCR locally.

```text
Screen region
     ↓
Local OCR
     ↓
Extracted text
```

Requirements:

- local execution;
- no screenshot upload;
- English support;
- Portuguese support;
- automatic language detection where practical;
- good recognition of desktop UI and game text;
- fast enough for interactive usage.

Developing a custom OCR engine is outside the scope.

An existing OCR engine or executable should be integrated.

Potential implementations include:

- Tesseract
- RapidOCR
- another local OCR engine with suitable Linux support

The OCR layer should remain isolated from the rest of the application.

---

## 9. Action Menu

After OCR succeeds, COSMIC Select displays a small action overlay.

Example:

```text
┌────────────────────────────┐
│ Translate                  │
│ Copy                       │
│ Search                     │
│ Ask AI                     │
└────────────────────────────┘
```

The overlay should visually integrate well with COSMIC.

Requirements:

- support system light/dark appearance where possible;
- use compact desktop-native controls;
- avoid excessive decoration;
- disappear when cancelled;
- avoid stealing unnecessary screen space;
- remain usable on Wayland;
- behave correctly under COSMIC window management.

The UI does not need to use `libcosmic`.

GTK4 may be used to provide the interface as long as behavior remains fully compatible with COSMIC.

No permanent main window is required during normal usage.

---

## 10. Translate

The **Translate** action translates the extracted text into the configured target language.

Default example:

```text
Target language: Portuguese (Brazil)
```

Flow:

```text
OCR text
   ↓
Translation provider
   ↓
Translation overlay
```

Example:

```text
┌──────────────────────────────────┐
│ Original                         │
│                                  │
│ This weapon scales with          │
│ dexterity.                       │
│                                  │
│ Translation                      │
│                                  │
│ Esta arma escala com destreza.   │
└──────────────────────────────────┘
```

The original text should remain visible when displaying the translation.

The first version may use an external translation API.

---

## 11. Copy

The **Copy** action places the extracted text directly into the system clipboard.

Flow:

```text
OCR
 ↓
Clipboard
```

After copying, the overlay may display a brief confirmation:

```text
Text copied
```

No additional window should appear.

Clipboard integration must work correctly on Wayland under COSMIC.

---

## 12. Search

The **Search** action searches the extracted text using the user's configured search engine.

Flow:

```text
OCR text
   ↓
URL encoding
   ↓
Default browser
```

The result must open in the user's default browser.

The application should use the system's standard URI-opening mechanism rather than launching a specific browser executable.

Initial implementation may use:

```text
xdg-open
```

or an equivalent XDG-compatible mechanism.

---

## 13. Ask AI

The **Ask AI** action opens a compact prompt interface.

The OCR result becomes context for the request.

Example:

```text
┌────────────────────────────────────────┐
│ Selected text                          │
│                                        │
│ This weapon scales with dexterity.     │
│                                        │
│ Ask something                          │
│ ┌────────────────────────────────────┐ │
│ │ What does this mean in Dark Souls? │ │
│ └────────────────────────────────────┘ │
│                                  Send  │
└────────────────────────────────────────┘
```

The request contains:

```text
context = OCR extracted text
question = user input
```

The response is displayed inside the COSMIC Select UI.

The AI feature does not require:

- conversations;
- chat history;
- memory;
- tools;
- agents;
- RAG;
- file uploads.

Every request is independent.

---

## 14. Privacy

Screen captures remain local.

Default processing model:

```text
Screen
  ↓
Screenshot
  ↓
Local OCR
  ↓
Text
```

Only extracted text may be sent to external APIs.

The selected image must not be transmitted to translation or AI providers unless a future feature explicitly requires it and the user opts into that behavior.

No analytics or telemetry are required.

---

## 15. Configuration

COSMIC Select only needs a small settings interface.

### General

- Global shortcut
- Translation target language
- Search engine

### OCR

- OCR language preferences
- OCR engine configuration, if required

### Translation

- Provider
- API key

### AI

- Provider
- Model
- API key

Configuration is stored locally.

No account or cloud synchronization is required.

---

## 16. COSMIC Integration

COSMIC integration is a first-class requirement.

### 16.1 Screen Capture

Use:

```text
org.freedesktop.portal.Screenshot
```

through:

```text
xdg-desktop-portal-cosmic
```

Interactive area selection should be delegated to COSMIC whenever the portal provides the required functionality.

---

### 16.2 Global Shortcut

Global activation should use the XDG global shortcut infrastructure available under COSMIC.

Expected behavior:

```text
User is inside another application
          ↓
Global shortcut
          ↓
COSMIC Select activates
```

The application must not require focus to receive the shortcut.

---

### 16.3 D-Bus

D-Bus is the preferred mechanism for communication with XDG portals.

Go integration should use a stable D-Bus library rather than shelling out to command-line utilities for portal communication.

---

### 16.4 Wayland

COSMIC Select is Wayland-native.

X11 compatibility is not a requirement.

The application must not depend on:

```text
XGrabKey
Xlib screen capture
xdotool
```

or other X11-specific mechanisms.

---

### 16.5 Appearance

The UI should integrate visually with COSMIC as closely as practical.

This means:

- respecting system appearance;
- avoiding GNOME-specific assumptions;
- avoiding GNOME-only shell integrations;
- using normal Wayland windows and portals;
- ensuring overlays behave properly under COSMIC.

Visual similarity to native COSMIC applications is desirable, but use of the COSMIC Rust toolkit is not required.

Functional compatibility takes priority over exact native widget parity.

---

## 17. Proposed Architecture

```text
               COSMIC Desktop
                     │
              Global Shortcut
                     │
                     ▼
           COSMIC Select Process
                  Go + GTK4
                     │
                     ▼
       xdg-desktop-portal-cosmic
                     │
              Area Selection
                     │
                     ▼
                  Image
                     │
                     ▼
                Local OCR
                     │
                     ▼
              Extracted Text
                     │
       ┌─────────────┼──────────────┐
       │             │              │
       ▼             ▼              ▼
      Copy       Translate        Search
       │             │              │
       │             │              ▼
       │             │        Default Browser
       │             │
       │             ▼
       │       Translation API
       │
       └─────────────────────┐
                             │
                             ▼
                           Ask AI
                             │
                             ▼
                          LLM API
```

---

## 18. Internal Components

```text
cosmic-select
├── app
├── capture
├── shortcut
├── ocr
├── overlay
├── clipboard
├── translate
├── search
├── ai
└── config
```

### `app`

Application lifecycle and coordination.

### `capture`

Handles D-Bus integration with the COSMIC screenshot portal.

### `shortcut`

Handles global shortcut registration and activation.

### `ocr`

Converts the selected image into text.

### `overlay`

Provides the action menu, translation result, AI prompt and AI response UI.

### `clipboard`

Writes extracted text to the Wayland clipboard.

### `translate`

Communicates with the configured translation provider.

### `search`

Opens the default browser with a search query.

### `ai`

Communicates with the configured LLM provider.

### `config`

Stores local application settings.

---

## 19. Provider Interfaces

External services should use small internal interfaces to avoid coupling the application to a single provider.

### Translation

Conceptual interface:

```text
Translate(text, sourceLanguage?, targetLanguage)
```

### AI

Conceptual interface:

```text
Ask(context, question)
```

No generalized provider framework is required.

The abstraction exists only to prevent provider-specific HTTP and authentication logic from leaking into unrelated parts of the application.

---

## 20. Go Responsibilities

Go should handle:

- application lifecycle;
- D-Bus communication;
- portal requests;
- OCR process coordination;
- HTTP APIs;
- configuration;
- JSON;
- clipboard orchestration;
- browser launching;
- action routing;
- logging;
- error handling.

GTK4 should remain primarily responsible for presentation and user interaction.

The architecture should avoid pushing application logic into the UI layer.

---

## 21. Performance

COSMIC Select should feel immediate.

### Activation

The selection interface should appear effectively immediately after the shortcut is pressed.

### OCR

For normal UI-sized selections:

```text
target: < 1 second
```

### Idle State

The application should consume negligible CPU while idle.

Memory usage should remain reasonable for a small desktop utility.

Network requests only occur when the user explicitly selects an action requiring an external service.

---

## 22. Error Handling

Errors should remain simple and visible.

Examples:

```text
No text detected
```

```text
Translation failed
```

```text
AI provider unavailable
```

```text
API key not configured
```

```text
Unable to access screenshot portal
```

The application should never present large diagnostic screens during normal interaction.

Detailed errors may be written to logs.

---

## 23. Out of Scope

The following features are explicitly outside the initial project:

- image search;
- object recognition;
- image understanding;
- screen recording;
- screenshot history;
- OCR history;
- AI conversation history;
- accounts;
- authentication;
- cloud synchronization;
- team features;
- browser extension;
- mobile application;
- plugins;
- workflow automation;
- agent system;
- background screen monitoring;
- RAG;
- vector databases;
- remote control;
- GNOME support;
- KDE support;
- X11 support;
- Rust requirement;
- `libcosmic` requirement.

These features should not influence the initial architecture.

---

## 24. MVP

The MVP is complete when the user can:

1. install COSMIC Select on COSMIC;
2. configure a global shortcut;
3. press the shortcut from any application;
4. select a rectangular screen region;
5. extract text through local OCR;
6. copy extracted text;
7. translate extracted text;
8. search extracted text in the default browser;
9. ask an AI about the extracted text;
10. configure the required API credentials.

That is the product.

---

## 25. Success Criteria

COSMIC Select succeeds if this interaction feels natural:

```text
Super + Shift + Space

        ↓

Select something on screen

        ↓

┌────────────────────────┐
│ Translate              │
│ Copy                   │
│ Search                 │
│ Ask AI                 │
└────────────────────────┘
```

The user should not need to:

- manually create a screenshot;
- save an image;
- open another OCR application;
- manually copy OCR results;
- open a translator;
- manually open a search engine;
- manually transfer context into an AI application.

The tool exists to remove those steps.

Nothing more is required.