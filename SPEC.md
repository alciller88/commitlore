SPEC.md — CommitLore

Source of truth for the project. Any functionality change must be
reflected here before implementation. Nothing is implemented that is
not specified in this document.


1. Vision
CommitLore is a cross-platform tool (CLI + desktop app) written in Go
that analyzes git repositories — local and GitHub — and generates
changelogs, narratives, and reports about code history, with tone and
format configurable through a modular style system.
Tagline: Your repo has a story. CommitLore tells it.

2. Design Principles

No mandatory dependencies — works offline, no API key, no account required.
Optional LLM — if the user configures an API key, output improves; without it, works equally well via templates.
Cross-platform — native CLI binary and desktop app for Linux, macOS, and Windows.
Composable — multiple output formats for integration into existing pipelines.
Modular styles — tones are not hardcoded; they are loadable, exportable, and importable files.
Explicit over magic — all configuration is visible via flags or UI.
Sustainable code — small functions, single responsibility, no shortcuts.


3. Tech Stack
Backend / CLI
LayerTechnologyLanguageGo 1.22+CLI frameworkCobraConfig/flagsViperGitHub APIgo-githubLocal gitgo-gitHTMLhtml/template (stdlib)LLM (optional)Anthropic API / OpenAI API (configurable)Lintergolangci-lintTeststesting (stdlib) + testify
Desktop App
LayerTechnologyFrameworkWails v3 alphaFrontendSvelte + TypeScriptUI StylesTailwind CSS + shadcn-svelte (component base)Buildwails3 CLI (native binaries per platform)

The desktop app backend fully reuses the internal/ packages.
No duplicated logic.


4. Data Sources (v1)

Local git repository — via go-git, no authentication.
GitHub — via REST API. Optional token (without token: public repos only).


GitLab and other providers are out of scope for v1.


5. CLI Commands
All commands share the following global flags.
Global Flags
FlagValuesDefaultDescription--formatterminal, md, json, htmlterminalOutput format--stylename of loaded styleformalText tone--outputfile pathstdoutDestination file--llmanthropic, openai, nonenoneLLM to use for enriching output--llm-base-urlURL""Override API base URL (OpenAI-compatible endpoints)--llm-modelmodel name""Override LLM model (default per provider)

5.1 commitlore generate
Generates a changelog from commits and/or PRs.
commitlore generate [flags]
FlagDescription--sinceSince tag, commit SHA, or date (e.g. v1.2.0, 2024-01-01)--untilUntil tag or commit SHA (default: HEAD)--repoLocal path or GitHub URL (owner/repo)--include-prsInclude PR info (requires GitHub token)

5.2 commitlore story
Generates a complete narrative of the repository history.
commitlore story [flags]
FlagDescription--repoLocal path or GitHub URL (owner/repo)--fromStarting commit or tag (default: first commit)

5.3 commitlore history
Explore commits filtered by author, date, or range.
commitlore history [flags]
FlagDescription--repoLocal path or GitHub URL--authorFilter by author (name or email)--sinceSince date or tag--untilUntil date or tag--limitMax number of commits (default: 50)

5.4 commitlore contributors
Map of who has touched which parts of the code.
commitlore contributors [flags]
FlagDescription--repoLocal path or GitHub URL--sinceAnalysis period--topNumber of contributors to show (default: 10)

5.5 commitlore style
Management of the modular style system.
commitlore style <subcommand> [flags]
SubcommandDescriptionlistList available styles (built-in + installed)showShow a style definitioncreateCreate a new style (interactive wizard or flags)importImport a style from URL or local pathexportExport a style to .shipstyle filedeleteDelete an installed style (does not delete built-ins)

6. Modular Style System
Styles are .shipstyle files in YAML format that define tone, text
templates, and the full visual identity of the output.
Stored in ~/.config/commitlore/styles/ (on Windows: %APPDATA%\commitlore\styles\).

Complete structure of a .shipstyle file:

  name: "name"
  version: "1.0.0"
  description: "description"
  author: "author"
  tags: []                    # marketplace metadata
  preview_url: ""             # preview URL
  homepage: ""                # project URL

  llm_prompt: |               # LLM prompt (optional)
    ...

  templates:                  # text templates
    header, feature, fix, breaking, footer
    story_intro, story_milestone, story_peak, story_contributor, story_footer

  vocabulary:                 # substitutions without LLM (map[string]string)
    bug: "heresy"
    fix: "purge"

  theme:                      # HTML visual identity
    mode: "dark"              # dark | light
    colors:
      primary, secondary, background, surface, text, accent, border
    typography:
      font_family, font_size_base, font_size_header, font_size_code
    header_image: ""          # URL or base64
    logo: ""                  # URL, base64, or inline SVG string
    card_style: "bordered"    # minimal | bordered | glassmorphism
    animations: true
    custom_css: ""            # additional CSS injected at end
    window_controls:          # custom titlebar button colors
      default, close, minimize, maximize

  terminal:                   # terminal visual identity
    colors:
      header, feature, fix, breaking, footer   # ANSI color names
    decorators:
      separator, bullet, indent
    density: "normal"         # compact | normal | verbose

  ui_labels:                  # navigation and button label overrides (optional)
    dashboard: "Dashboard"
    generate: "Generate"
    generate_button: "Generate"
    story: "Story"
    story_button: "Tell the story"
    history: "History"
    contributors: "Contributors"
    styles: "Styles"
    settings: "Settings"

  icons:                      # per-style icon/emoji characters (optional)
    feature: "✦"
    fix: "✔"
    breaking: "⚠"
    chore: "⚙"
    docs: "📄"
    test: "🧪"
    story_peak: "🔥"
    bullet: "•"
    separator: "────────────────────────────────────────"

  html_template_changelog: |  # custom HTML template for changelogs (optional)
    <!DOCTYPE html>...        # complete self-contained HTML document
                              # receives HTMLTemplateContext with all style fields:
                              # {{.Theme.Colors.Primary}}, {{.Icons.Feature}}, etc.
                              # When empty, the default renderer is used.

  html_template_story: |      # custom HTML template for stories (optional)
    <!DOCTYPE html>...        # same context struct as changelog, but story fields
                              # (.Peaks, .Contributors, .Tags, .TotalCommits) populated
                              # When empty, the default renderer is used.

All vocabulary, theme, terminal, ui_labels, icons, and html_template_* fields
are optional with sensible zero-value defaults.

Built-in Styles (v1)

formal — technical and professional, neutral colors
patchnotes — video game style, purple/gold, animations
ironic — dry humor, muted colors, minimalist
epic — grand narrative, gold/dark, ornate

Behavior

Without LLM: templates and vocabulary fields from the .shipstyle are used.
With LLM: the llm_prompt field is used to instruct the model; templates as fallback.
Built-in styles are not modifiable or deletable.


7. Output Formats
FormatDescriptionterminalText with ANSI color direct to stdoutmdStandard Markdown, GitHub-compatiblejsonComplete data structure, suitable for pipelineshtmlSelf-contained HTML report with inline styles

PDF removed in favor of HTML. The generated HTML can be printed to PDF from any browser.

8. Optional LLM

The tool works completely without LLM using templates from the active style.
Environment variables: COMMITLORE_LLM_PROVIDER and COMMITLORE_LLM_API_KEY.
The --llm flag overrides the environment variable per command.
The --llm-base-url flag overrides the API endpoint (useful for OpenAI-compatible
providers: Ollama, Groq, LM Studio, Together AI, etc.).

Supported providers in v1: anthropic, openai.
OpenAI-compatible providers use --llm openai --llm-base-url <endpoint>.
Convenience aliases (map to openai adapter + default base URL):
  ollama  → http://localhost:11434/v1
  groq    → https://api.groq.com/openai/v1

The --llm-base-url flag is also available as environment variable
COMMITLORE_LLM_BASE_URL.


9. Desktop App (Wails)
The app shares all logic from internal/. It only adds a UI layer on top.
Main Screens (v1)

Dashboard — summary of the active repo: latest activity, contributors, tags.
Generate — form to configure and generate a changelog.
Story — narrative visualization of the repo history.
History — commit explorer with filters.
Contributors — visual contribution map.
Styles — read-only browser: preview styles, import, export, set active. No in-app editing.
Settings — LLM provider configuration, API key management via OS keychain, app style selection.

Repo Picker (all screens)

Empty state with three entry points: native folder picker, drag & drop, GitHub owner/repo input.
Recent repos list persisted in ~/.config/commitlore/config.yml (max 10 entries).
Visual feedback to reopen recent repos quickly.

Output Display (Generate, Story)

HTML preview rendered inline + Copy button + Save as file button.
Uses internal/renderer HTML output.

Style Selection (Settings screen)

Active app style selected from dropdown (built-in + user styles).
Persisted in ~/.config/commitlore/config.yml as active_style field.
Theme colors and typography injected as CSS variables across the entire UI.
Default style: formal.

LLM Configuration (Settings screen)

API key stored in OS keychain via go-keyring (Windows Credential Manager, macOS Keychain, Linux Secret Service).
Key never written to disk in plaintext.
UI shows key status: configured / not configured.

Visual Identity

Frameless window with custom titlebar — macOS-style window controls in sidebar header.
Dark palette by default, with light theme option via active style.
Design system tokens for spacing, typography, radius, transitions (design.css).
Compact sidebar (220px) with VIEWS/SYSTEM nav sections, 32px row height.
Monospaced typography for outputs; sans-serif for navigation.
Custom style built on shadcn-svelte as component base.
No dependency on generic UI libraries (no Material, no Bootstrap).


10. Project Structure
commitlore/
├── cmd/                      # Cobra commands
│   ├── generate.go
│   ├── story.go
│   ├── history.go
│   ├── contributors.go
│   └── style/
│       ├── list.go
│       ├── create.go
│       ├── import.go
│       ├── export.go
│       └── delete.go
├── internal/
│   ├── git/                  # Local repo access (go-git)
│   ├── github/               # GitHub API access
│   ├── changelog/            # Commit parsing and grouping
│   ├── narrative/            # Text generation
│   ├── renderer/             # Rendering (terminal, md, json, html)
│   ├── llm/                  # LLM adapters (Anthropic, OpenAI)
│   └── styles/               # Loading, validation, management of .shipstyle
│       └── builtin/          # Embedded built-in styles (.shipstyle)
├── assets/
│   └── logo.svg              # Official CommitLore logo
├── app/                      # Wails app
│   ├── frontend/             # Svelte + Tailwind + shadcn-svelte
│   └── app.go                # Wails bindings → internal/
├── styles/                   # User styles (.shipstyle)
├── .github/
│   └── workflows/
│       ├── ci.yml            # Lint + tests on push/PR
│       └── release.yml       # Cross-platform binaries on tag v*
├── CHANGELOG.md
├── SPEC.md
├── CONTEXT.md
├── main.go
└── README.md

Rule: no package in internal/ imports from cmd/ or app/.
The dependency flow is always inward.


11. Code Rules

Functions of maximum 40 lines. If it grows, extract.
One responsibility per function/struct.
No obvious comments — comments explain the why, not the what.
Explicit errors — never ignore an error with _.
Mandatory unit tests for all code in internal/.
Minimum coverage target: 70% per package.
Before each PR: golangci-lint run ./... must pass without errors.


12. CI/CD
Branch Strategy
BranchPurposemainProduction. Only receives merges from dev via PR.devIntegration. Base branch for features.feat/*Feature branches. Opened from dev.fix/*Bugfix branches. Opened from dev.
CI Pipeline — .github/workflows/ci.yml
Trigger: push to dev, PR to main.

golangci-lint run ./...
go test ./... -race -coverprofile=coverage.out
go build ./...

Release Pipeline — .github/workflows/release.yml
Trigger: tag v* pushed to main.

Full CI (lint + tests)
Build CLI binaries:

GOOS=linux GOARCH=amd64
GOOS=darwin GOARCH=arm64
GOOS=windows GOARCH=amd64


Build Wails app per platform
Create GitHub Release with all artifacts attached
Automatic CHANGELOG.md update


13. Semantic Versioning
Tag format: vMAJOR.MINOR.PATCH (e.g. v1.2.0)

MAJOR — incompatible changes (removed flags, broken behavior).
MINOR — new commands, flags, or screens, backward compatible.
PATCH — bugfixes and internal improvements without interface change.


14. Development Phases
PhaseScopePhase 1Project setup, base structure, CI pipeline, branchesPhase 2internal/git — local repo access + history commandPhase 3internal/changelog — commit parsing + contributors commandPhase 4generate command (no LLM, templates)Phase 5story command (no LLM, templates)Phase 6internal/renderer — md, json, html formatsPhase 7internal/styles — modular system + style commandPhase 8internal/github — GitHub API integrationPhase 9internal/llm — optional LLM integration (Anthropic + OpenAI)Phase 10Wails app — base structure + bindingsPhase 11Wails app — screens and complete UIPhase 12Release pipeline + cross-platform binariesPhase 13Polish, docs, README, examples

Each phase ends with passing tests and clean lint before merging to dev.
Only merge to main when a complete phase is stable in dev.


15. Security

### Fundamental Principle
CommitLore is a READ-ONLY tool. It never performs write operations on any repository, local or remote, under any circumstances.

### Prohibited Operations (never implement)
- git push, git commit, git add, git rm on user repos
- Writing files inside analyzed repo directories
- Modifying git configuration (.git/config, hooks, etc.)
- Creating or modifying branches, tags, or refs in user repos
- GitHub API calls with POST/PUT/PATCH/DELETE methods on user repos
- Executing arbitrary shell commands

### Permitted Operations
- Reading commits, tags, branches, diffs (go-git, read-only)
- GET calls to GitHub API (public and private repos with token)
- Writing output files ONLY to paths explicitly specified by the user via --output
- Writing to ~/.config/commitlore/ (app's own configuration and styles)

### Prompt Injection Protection
- Commit messages, file names, author names, and any data from a repository are UNTRUSTED content
- Never execute or evaluate commit content as code or instructions
- Never pass commit content directly to an LLM without prior sanitization
- Mandatory sanitization before passing repo data to an LLM:
  - Truncate commit messages to 500 characters maximum
  - Escape control characters
  - Add explicit delimiters in the LLM prompt to separate instructions from data: use "---DATA START---" and "---DATA END---"
- The llm_prompt from an imported .shipstyle is potentially untrusted content — warn the user before using it with an LLM

### Tokens and Credentials
- COMMITLORE_LLM_API_KEY and GITHUB_TOKEN are never logged, never appear in output, never included in reports
- Tokens are read only from environment variables, never from analyzed repo files
- If an imported .shipstyle contains fields that look like credentials, ignore them and warn the user

### Input Validation
- Repository paths: validate they exist and are git directories before operating
- GitHub URLs: validate format before calling the API
- --output flags: validate that the destination path is outside any .git/ directory
- Imported .shipstyle files: validate complete schema before loading, reject unknown fields


16. Roadmap (out of scope v1)

Interactive style create wizard — interactive stdin-based style creation (currently flags-only).
Style marketplace — public repository of community .shipstyle files.
GitLab support.
VS Code / Cursor plugin.
Slack / Discord integration to publish changelogs automatically.
commitlore watch — daemon mode that generates changelog automatically when a tag is created.
