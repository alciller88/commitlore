# SPEC.md — CommitLore

Source of truth for what the product does and what comes next.
Contains the full functional specification and the prioritized backlog.

---

## 1. Vision

CommitLore is a cross-platform tool (CLI + desktop app) written in Go
that analyzes git repositories — local and GitHub — and generates
changelogs, narratives, and reports about code history, with tone and
format configurable through a modular style system.

**Tagline:** Your repo has a story. CommitLore tells it.

---

## 2. Design Principles

- **No mandatory dependencies** — works offline, no API key, no account required.
- **Optional LLM** — if the user configures an API key, output improves; without it, works equally well via templates.
- **Cross-platform** — native CLI binary and desktop app for Linux, macOS, and Windows.
- **Composable** — multiple output formats for integration into existing pipelines.
- **Modular styles** — tones are not hardcoded; they are loadable, exportable, and importable files.
- **Explicit over magic** — all configuration is visible via flags or UI.
- **Sustainable code** — small functions, single responsibility, no shortcuts.

---

## 3. Data Sources

- **Local git repository** — via go-git, no authentication.
- **GitHub** — via REST API. Optional token (without token: public repos only).

GitLab and other providers are out of scope for v1.

---

## 4. CLI Commands

All commands share the following global flags.

### Global Flags

| Flag | Values | Default | Description |
|------|--------|---------|-------------|
| `--format` | `terminal`, `md`, `json`, `html` | `terminal` | Output format |
| `--style` | name of loaded style | `formal` | Text tone |
| `--output` | file path | stdout | Destination file |
| `--llm` | `anthropic`, `openai`, `none` | `none` | LLM to use for enriching output |
| `--llm-base-url` | URL | `""` | Override API base URL (OpenAI-compatible endpoints) |
| `--llm-model` | model name | `""` | Override LLM model (default per provider) |

### 4.1 commitlore generate

Generates a changelog from commits and/or PRs.

```
commitlore generate [flags]
```

| Flag | Description |
|------|-------------|
| `--since` | Since tag, commit SHA, or date (e.g. `v1.2.0`, `2024-01-01`) |
| `--until` | Until tag or commit SHA (default: HEAD) |
| `--repo` | Local path or GitHub URL (`owner/repo`) |
| `--include-prs` | Include PR info (requires GitHub token) |

### 4.2 commitlore story

Generates a complete narrative of the repository history.

```
commitlore story [flags]
```

| Flag | Description |
|------|-------------|
| `--repo` | Local path or GitHub URL (`owner/repo`) |
| `--from` | Starting commit or tag (default: first commit) |

### 4.3 commitlore history

Explore commits filtered by author, date, or range.

```
commitlore history [flags]
```

| Flag | Description |
|------|-------------|
| `--repo` | Local path or GitHub URL |
| `--author` | Filter by author (name or email) |
| `--since` | Since date or tag |
| `--until` | Until date or tag |
| `--limit` | Max number of commits (default: 50) |

### 4.4 commitlore contributors

Map of who has touched which parts of the code.

```
commitlore contributors [flags]
```

| Flag | Description |
|------|-------------|
| `--repo` | Local path or GitHub URL |
| `--since` | Analysis period |
| `--top` | Number of contributors to show (default: 10) |

### 4.5 commitlore style

Management of the modular style system.

```
commitlore style <subcommand> [flags]
```

| Subcommand | Description |
|------------|-------------|
| `list` | List available styles (built-in + marketplace-installed) |
| `show` | Show a style definition |
| `delete` | Delete a marketplace-installed style (does not affect built-ins) |

---

## 5. Modular Style System

Styles are `.shipstyle` files in YAML format that define tone, text
templates, and the full visual identity of the output.

Stored in `~/.config/commitlore/styles/` (on Windows: `%APPDATA%\commitlore\styles\`).

### Complete `.shipstyle` Schema

```yaml
name: "name"
version: "1.0.0"
description: "description"
author: "author"
language: ""                # language code: "en" (default), "es". Empty = "en".
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
  marketplace: "Marketplace"
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
```

All `vocabulary`, `theme`, `terminal`, `ui_labels`, `icons`, `window_controls`,
`html_template_changelog`, and `html_template_story` fields are optional with
sensible zero-value defaults.

### Built-in Styles

| Style | Character |
|-------|-----------|
| **formal** | Technical and professional, neutral colors, Inter font |
| **patchnotes** | Video game style, purple/gold, Rajdhani font, animations |
| **ironic** | Dry humor, coral/teal palette, Comic Sans, Clippy |
| **epic** | Grand medieval narrative, gold/dark, Cinzel font, ornate |

### Behavior

- **Without LLM:** templates and vocabulary fields from the `.shipstyle` are used.
- **With LLM:** the `llm_prompt` field is used to instruct the model; templates serve as fallback.
- Built-in styles are not modifiable or deletable.

---

## 6. Output Formats

| Format | Description |
|--------|-------------|
| `terminal` | Text with ANSI color direct to stdout |
| `md` | Standard Markdown, GitHub-compatible |
| `json` | Complete data structure, suitable for pipelines |
| `html` | Self-contained HTML report with inline styles |

PDF was removed in favor of HTML. The generated HTML can be printed to PDF from any browser.

When a style defines `html_template_changelog` or `html_template_story`, the HTML
output uses that custom template instead of the default renderer. Custom templates
receive an `HTMLTemplateContext` containing all style fields (theme, icons, ui_labels)
plus the report data. Template helper functions available: `upper`, `lower`, `initials`,
`add`, `mul`, `divf`, `divi`, `safeHTML`.

All four built-in styles include unique HTML templates with Chart.js visualizations,
style-appropriate typography and layout, and colors driven by `.shipstyle` theme values.

---

## 7. LLM Integration

The tool works completely without LLM using templates from the active style.

### Configuration

| Method | Detail |
|--------|--------|
| Environment variable | `COMMITLORE_LLM_PROVIDER` — `anthropic`, `openai`, `ollama`, `groq` |
| Environment variable | `COMMITLORE_LLM_API_KEY` — provider API key |
| Environment variable | `COMMITLORE_LLM_BASE_URL` — override API base URL |
| CLI flag | `--llm` overrides the environment variable per command |
| CLI flag | `--llm-base-url` overrides the API endpoint |
| CLI flag | `--llm-model` overrides the default model |
| Desktop app | Settings screen with OS keychain storage for API key |

### Supported Providers

| Provider | Default Model | Notes |
|----------|---------------|-------|
| `anthropic` | claude-haiku-4-5-20251001 | Native adapter |
| `openai` | gpt-4o-mini | Native adapter |

### Convenience Aliases

Aliases map to the `openai` adapter with a preset base URL:

| Alias | Base URL |
|-------|----------|
| `ollama` | `http://localhost:11434/v1` |
| `groq` | `https://api.groq.com/openai/v1` |

Any OpenAI-compatible provider works via `--llm openai --llm-base-url <endpoint>`.

### Anti-Hallucination

All built-in style `llm_prompt` fields include an explicit instruction requiring the LLM
to base output exclusively on the provided commit data and never fabricate content.

---

## 8. Desktop App

The app shares all logic from `internal/`. It only adds a UI layer on top.

### Tech Stack

| Layer | Technology |
|-------|------------|
| Framework | Wails v3 alpha |
| Frontend | Svelte + TypeScript |
| UI styles | Tailwind CSS + shadcn-svelte (component base) |
| Build | wails3 CLI (native binaries per platform) |

### Screens

| Screen | Description |
|--------|-------------|
| **Dashboard** | Summary of the active repo: latest activity, contributors, tags. Three entry points: native folder picker, drag & drop, GitHub owner/repo input. |
| **Generate** | Form to configure and generate a changelog. Two-column layout: filter sidebar + HTML preview iframe. |
| **Story** | Narrative visualization of the repo history. Same two-column layout as Generate. |
| **History** | Commit explorer with filters. Dense table rows with click-to-copy hash. |
| **Contributors** | Visual contribution map with avatar initials and activity bars. |
| **Styles** | Read-only browser: preview styles, install from marketplace, delete, set active. No in-app editing. |
| **Marketplace** | Browse and install community styles from the official commitlore-styles repository. Requires internet connection. Built-in styles are always available offline regardless of connectivity. |
| **Settings** | LLM provider configuration, API key management via OS keychain, app style selection, language selector (English / Spanish). |

### Repo Picker

- Empty state with three entry points: native folder picker, drag & drop, GitHub owner/repo input.
- GitHub connection via modal dialog with owner/repo input and optional token (session-only, not persisted to disk).
- Recent repos list persisted in `~/.config/commitlore/config.yml` (max 10 entries, MRU order).
- Global repo store: all screens read the active repo from a shared Svelte store; only Dashboard writes to it.

### Output Display (Generate, Story)

- HTML preview rendered inline in iframe + Copy button + Save as file button.
- Narrative content rendered via goldmark (markdown to HTML with XSS protection).
- In-memory commit cache avoids re-fetching from GitHub API when switching styles on the same repo.

### Style Selection (Settings screen)

- Active app style selected from dropdown (built-in + user styles).
- Persisted in `~/.config/commitlore/config.yml` as `active_style` field.
- Theme colors, typography, logo, and ui_labels injected as CSS variables across the entire UI.
- Default style: `formal`.

### Language Selector (Settings screen)

- Two options: English / Español, rendered as radio buttons.
- Persisted in `~/.config/commitlore/config.yml` as `language` field.
- On change: validates that the active style has a variant for the selected language.
  - If not available: shows inline error, does not change the selection.
  - If available: saves language, reloads theme and ui_labels from the language variant.
- Language variants are separate `.shipstyle` files with `<name>.<lang>.shipstyle` naming.
- Default language: `"en"`.
- All 4 built-in styles have Spanish variants. Marketplace styles may optionally provide them.

### Styles Screen

- Two-column layout: style card list + read-only detail panel.
- Detail panel shows: logo, name/version/author, description, theme color circles, font preview, mode badge, UI labels (if custom), icons (if custom), collapsible LLM prompt.
- Actions available:
  - **Set as active** — applies the style to the entire app immediately.
  - **Delete** — marketplace-installed styles only, with confirmation dialog. Built-in styles cannot be deleted.
  - **Browse marketplace** — opens the in-app marketplace screen.
- Styles are installed exclusively from the marketplace. Import and export are not supported.

### Marketplace Architecture

The marketplace is backed by a private GitHub repository (`commitlore-styles`)
controlled by the CommitLore team. The app fetches a catalog file directly
from the repository — no backend or additional infrastructure required.

**Repository structure:**
```
commitlore-styles/
├── index.json          # catalog with metadata for all available styles
├── <style-name>/
│   ├── style.shipstyle
│   └── preview.png
└── ...
```

**index.json format:**
```json
[
  {
    "name": "style-name",
    "description": "One-line description",
    "author": "author",
    "version": "1.0.0",
    "tags": ["dark", "colorful"],
    "preview": "https://raw.githubusercontent.com/alciller88/commitlore-styles/main/<name>/preview.png",
    "download": "https://raw.githubusercontent.com/alciller88/commitlore-styles/main/<name>/style.shipstyle",
    "downloads": {
      "es": "https://raw.githubusercontent.com/alciller88/commitlore-styles/main/<name>/<name>.es.shipstyle"
    }
  }
]
```

The `downloads` field is a map of language code to download URL for language variants.
One catalog entry per style — language variants are not listed as separate entries.

**App behavior:**
- On opening the Marketplace screen, the app fetches `index.json` and displays
  the catalog with name, description, author, tags, and preview image.
- Installing a style fetches the base `.shipstyle` from `download` URL and all
  language variants from the `downloads` map, validates each, and saves them to
  `~/.config/commitlore/styles/`. Variants are saved as `<name>.<lang>.shipstyle`.
- If the fetch fails (no connection, rate limit), the screen shows an error with
  a retry option. Already-installed styles remain available offline.
- The four built-in styles are always embedded in the binary and are never
  fetched from the marketplace.

**Security:**
- Only `.shipstyle` files from the official `commitlore-styles` repository are
  installable via the marketplace UI.
- The app validates the full `.shipstyle` schema after download before saving,
  rejecting unknown fields.
- The `llm_prompt` field in marketplace styles triggers the existing warning
  to the user before the style is activated with an LLM.

### LLM Configuration (Settings screen)

- API key stored in OS keychain via go-keyring (Windows Credential Manager, macOS Keychain, Linux Secret Service).
- Key never written to disk in plaintext.
- UI shows key status: configured / not configured.
- LLM settings auto-read by Generate and Story screens — no redundant selectors.

### Visual Identity

- Frameless window with custom titlebar — macOS-style window controls (close/minimize/maximize circles) in sidebar header.
- Per-style window control colors via `window_controls` in `.shipstyle`.
- Dark palette by default, with light theme option via active style.
- Design system tokens for spacing, typography, radius, transitions.
- Compact sidebar (220px) with VIEWS/SYSTEM nav sections, 32px row height.
- Per-style navigation labels via `ui_labels`.
- Per-style logo in sidebar header (inline SVG or image URL, fallback to CommitLore wordmark).
- Repo indicator in content topbar (always visible, all screens).
- Monospaced typography for outputs; sans-serif for navigation.
- Custom style built on shadcn-svelte as component base. No dependency on generic UI libraries.

---

## 9. Security Model

### Read-Only Guarantee

CommitLore is a read-only tool. It never performs write operations on any repository, local or remote.

**Permitted operations:**
- Reading commits, tags, branches, diffs (go-git, read-only).
- GET calls to GitHub API (public and private repos with token).
- Writing output files only to paths explicitly specified by the user via `--output`.
- Writing to `~/.config/commitlore/` (app's own configuration and styles).

**Prohibited operations (never implement):**
- `git push`, `git commit`, `git add`, `git rm` on user repos.
- Writing files inside analyzed repo directories.
- Modifying git configuration (`.git/config`, hooks, etc.).
- Creating or modifying branches, tags, or refs in user repos.
- GitHub API calls with POST/PUT/PATCH/DELETE methods on user repos.
- Executing arbitrary shell commands.

### Prompt Injection Protection

- Commit messages, file names, author names, and any data from a repository are untrusted content.
- Never execute or evaluate commit content as code or instructions.
- Never pass commit content directly to an LLM without prior sanitization.
- Mandatory sanitization before passing repo data to an LLM:
  - Truncate commit messages to 500 characters maximum.
  - Escape control characters.
  - Add explicit delimiters in the LLM prompt to separate instructions from data (`---DATA START---` / `---DATA END---`).
- The `llm_prompt` from an imported `.shipstyle` is potentially untrusted content — warn the user before using it with an LLM.
- Basic pattern detection for known injection patterns ("ignore previous", "exfiltrate", "reveal system prompt", etc.) rejects prompts with clear error before sending to LLM.

### Credentials

- `COMMITLORE_LLM_API_KEY` and `GITHUB_TOKEN` are never logged, never appear in output, never included in reports.
- Tokens are read only from environment variables or OS keychain, never from analyzed repo files.
- If an imported `.shipstyle` contains fields that look like credentials, ignore them and warn the user.
- GitHub token from the connection modal is session-only — not persisted to disk.

### Input Validation

- Repository paths: validate they exist and are git directories before operating.
- GitHub URLs: validate format before calling the API.
- `--output` flags: validate that the destination path is outside any `.git/` directory.
- Imported `.shipstyle` files: validate complete schema before loading, reject unknown fields.
- Style names: validated on creation to prevent filesystem issues.

---

## 10. Development Phases

| Phase | Scope | Status |
|-------|-------|--------|
| Phase 1 | Project setup, base structure, CI pipeline, branches | Completed |
| Phase 2 | `internal/git` — local repo access + `history` command | Completed |
| Phase 3 | `internal/changelog` — commit parsing + `contributors` command | Completed |
| Phase 4 | `generate` command (no LLM, templates) | Completed |
| Phase 5 | `story` command (no LLM, templates) | Completed |
| Phase 6 | `internal/renderer` — md, json, html formats | Completed |
| Phase 7 | `internal/styles` — modular system + `style` command | Completed |
| Phase 8 | `internal/github` — GitHub API integration | Completed |
| Phase 9 | `internal/llm` — optional LLM integration (Anthropic + OpenAI) | Completed |
| Phase 10 | Wails app — base structure + bindings | Completed |
| Phase 11 | Wails app — screens and complete UI | Completed |
| Phase 12 | Release pipeline + cross-platform binaries | In Progress |
| Phase 13 | Polish, docs, README, examples | In Progress |
| Phase 14 | Marketplace — style catalog, in-app browser, install/delete | Completed |
| Phase 15 | Internationalisation (i18n) — English/Spanish language support | Completed |

> **Phase 12 — In Progress, UNDOCUMENTED:** Release workflow exists (`.github/workflows/release.yml`) with CLI binary builds for 3 platforms on `v*` tag push. **Missing:** Wails desktop app builds per platform; automatic CHANGELOG.md update on release. Not yet documented in CHANGELOG.md.

> **Phase 13 — In Progress, UNDOCUMENTED:** README.md exists with features, quick start, styles table, LLM and desktop app sections, contributing guidelines. **Missing:** screenshots of each desktop app screen; license not yet defined. Not yet documented in CHANGELOG.md.

---

## 11. Backlog

### P0 — Must fix before next release

No P0 items at this time.

### P1 — Next planned phase

**Phase 12: Release pipeline + cross-platform binaries** _(In Progress)_

- ~~Build and release CLI binaries for linux/amd64, darwin/arm64, windows/amd64 on tag push.~~ ✅ Done
- Build and attach Wails desktop app binaries per platform to the GitHub Release.
  - _Acceptance:_ Release workflow builds Wails app for linux, darwin, windows and attaches all binaries alongside CLI builds.
- Automatic CHANGELOG.md update on release.
  - _Acceptance:_ Release workflow or tag-push hook updates CHANGELOG.md with the new version heading.

**Phase 13: Polish, docs, README, examples** _(In Progress)_

- ~~Complete README with features, quick start, styles table, LLM and desktop app sections.~~ ✅ Done
- Add screenshots of each desktop app screen to README.
  - _Acceptance:_ At least one screenshot per screen (Dashboard, Generate, Story, History, Contributors, Styles, Settings).
- Define and add project license.
  - _Acceptance:_ LICENSE file in repo root; license badge in README links to it.

### P2 — Planned but not scheduled

**Per-style navigation icons**

- Styles control sidebar navigation icons via inline SVG strings.
  - _Acceptance:_ New `ui_icons` block in `.shipstyle` with fields for each nav item (`dashboard`, `generate`, `story`, `history`, `contributors`, `styles`, `settings`, `local_repo`, `github_repo`). Icons are inline SVG strings. Falls back to current default icons when not defined.

**`contributors --with-files` for remote repos**

- For remote repos, the `--with-files` flag makes an additional API call per commit to obtain file diffs.
  - _Acceptance:_ Disabled by default due to rate limit cost. Without the flag, the TOP FILES column remains empty for remote repos. With the flag, file-level contribution data appears.

**Story: richer content**

- More milestones, activity metrics per period in story output.
  - _Acceptance:_ Story output includes additional temporal breakdowns and milestone detection beyond current implementation.

**Generate: style-influenced structure**

- Active style influences output structure, not just colors.
  - _Acceptance:_ Different styles produce structurally different changelogs (e.g., grouped differently, different section ordering) beyond template text and visual theming.

### Future — Roadmap / out of scope v1

- **GitLab support** — GitLab API integration as an additional data source.
- **VS Code / Cursor plugin** — editor extension for generating changelogs from within the IDE.
- **Slack / Discord integration** — publish changelogs automatically to messaging platforms.
- **`commitlore watch`** — daemon mode that generates changelog automatically when a tag is created.
