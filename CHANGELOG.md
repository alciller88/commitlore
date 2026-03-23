# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **Marketplace screen**: Browse and install community styles from the official `commitlore-styles` repository. Displays style cards with preview image, name, author, description, tags, and version. Install button with loading/installed states. Full-screen error state with retry button for offline scenarios.
- **MarketplaceApp binding**: `FetchCatalog`, `InstallStyle`, `IsInstalled` — fetches catalog from the official repository, downloads and validates `.shipstyle` files with strict schema checking (rejects unknown fields), saves to user styles directory.
- **`marketplace` ui_label field**: Styles can now customize the Marketplace navigation label via `ui_labels.marketplace` in `.shipstyle` schema.
- **`commitlore-styles` repository**: Official marketplace repository created at `alciller88/commitlore-styles` with `index.json` catalog.
- **Cyberpunk style v1.0.0**: Night City aesthetic with glitch effects, neon colors (yellow/cyan/pink), Orbitron + Share Tech Mono typography, HUD-style cards with cut corners, scanlines overlay, and full changelog + story HTML templates with Chart.js visualizations. Includes cyberpunk vocabulary, themed UI labels (NIGHT CITY, COMPILE, NEURAL LINK, RIPPERDOC, etc.), and anti-hallucination LLM prompt.
- **Marketplace end-to-end flow verified**: Catalog fetch, style install, theme activation, and UI label propagation all functional.

### Fixed

- **Cyberpunk style**: Completed all missing fields — `preview_url`, `homepage`, `custom_css` (Google Fonts import), full LLM prompt with `{{.Data}}` delimiters, vocabulary entries for `error`/`warning`, contributors section ("Known Operatives") in changelog template, 5th stat card (Last Scan) in story template, contributor ranking chart canvas id

### Removed

- **`commitlore style create`, `import`, `export` subcommands**: Styles are now managed exclusively via the marketplace. Only `list`, `show`, and `delete` remain.

### Previously added

- **HTML template system**: Each built-in style now has a unique, self-contained HTML report template with Chart.js visualizations, theme-driven colors, and client-side rendering. Styles can define separate templates for changelog and story output.
  - Formal: Stripe/GitHub docs aesthetic with stats cards, commit tables, and bar/donut charts
  - Patchnotes: Steam/game studio look with gradient header, highlights cards, and scroll-triggered fade-in animations
  - Epic: Medieval chronicle with parchment texture, drop caps, Roman numeral chapters, and ornamental separators
  - Ironic: Fake Word 95 interface with Comic Sans, toolbar chrome, and a Clippy speech bubble
- **Icons system**: Per-style icon sets for commit types, story elements, and decorators, defined in the `.shipstyle` schema
- **UI label customization**: Styles can override all navigation labels and button text (e.g., "Deploy Patch" instead of "Generate")
- **Contributor ranking charts**: Horizontal bar charts for contributor rankings added to all four built-in styles in story output
- **Epic story stats**: Story stats expanded from 3 to 5 cards, adding first commit date and most active month
- **Ironic changelog stats row**: Displays total commits, features added, and fixes count
- **Ironic story enhancements**: Most active month in subtitle, milestones/tags table, contributor ranking chart
- **GitHub Connect modal**: Centered modal for owner/repo input with optional session-only token, replacing the inline card
- **In-memory commit cache**: Avoids redundant GitHub API calls when switching styles on the same repository; cache is cleared on repo change
- **Template helper functions**: `add`, `mul`, `divf`, `divi`, `upper`, `lower`, `initials`, `safeHTML` available in custom HTML templates

### Changed

- **Styles screen simplified**: Replaced the full multi-tab editor with a read-only detail panel showing logo, colors, fonts, icons, and LLM prompt. Styles are now managed exclusively via import/export. In-app creation and editing removed.
- **Story data depth**: Story charts now show up to 12 months of activity data instead of 3, producing more meaningful visualizations

### Fixed

- **Chart rendering in desktop app**: All Chart.js calls wrapped in a 100ms delay to prevent blank charts when the WebView iframe has not yet completed layout
- **Chart colors**: Replaced CSS `color-mix()` in Chart.js script blocks with a JavaScript helper; Chart.js cannot parse CSS functions, which silently caused blank charts
- **Donut chart visibility**: Replaced palette with 7 visually distinct colors so legend entries no longer disappear when theme colors match the background
- **Donut chart centering**: Charts are now horizontally centered in their container
- **SVG logo rendering**: Inline SVG logos are now written as raw HTML instead of being escaped into an `<img>` tag
- **Version badge rendering**: The `until` parameter is now threaded through the renderer so semver version badges display correctly
- **Type badge CSS**: Added missing CSS rules for refactor, docs, test, and chore commit type badges
- **First author fallback**: Defaults to "an unknown contributor" when the first commit has no author
- **Generated date**: Reports now show the actual generation date instead of a hardcoded string
- **Formal style**: Stat row layout overflow fixed with flex layout; date values use consistent sizing via CSS class instead of inline styles; story banner no longer shows the generated date twice
- **Epic style**: Canvas IDs made unique to prevent Chart.js collisions; consecutive separator rendering fixed; footer text centered between separators
- **Ironic style**: Clippy bubble repositioned to avoid covering content in iframes; per-commit icons now reflect the actual type (e.g., `!` for breaking changes); redundant narrative title heading hidden; narrative text left-aligned for deadpan aesthetic; list style overrides prevent commit CSS from bleeding into narrative
- **Patchnotes style**: Added fallback timer for intersection observer animations in the desktop app iframe; LLM prompt rewritten to infer categories from commit content instead of using hardcoded gaming terms
- **Patchnotes template math**: Fixed type-safe percentage calculations in story template

### Removed

- **Style editor tabs**: Colors, Typography, Icons, Images, Templates, HTML, and Advanced editor tabs removed from the Styles screen in favor of import/export workflow

## [0.0.0] - 2026-03-20

### Added

- Initial project structure with SPEC.md, CONTEXT.md, and CHANGELOG.md
- **Local git access**: Repository analysis via go-git with commit history, contributors, chronology, tags, and activity peaks
- **GitHub API integration**: Remote repository support via go-github with commit and PR fetching, date/author filters, and user-friendly error messages for rate limits and auth failures
- **CLI commands**: `history`, `contributors`, `generate`, `story`, and `style` with full flag support for filtering, formatting, and output
- **Changelog generation**: Commit parsing and grouping by conventional commit type, with optional LLM enrichment for narrative content
- **Story generation**: Repository narrative built from chronology, milestones, activity peaks, and contributor data
- **Modular style system**: `.shipstyle` YAML format with vocabulary, theme (colors, typography, card style, animations), terminal formatting, icons, UI labels, and LLM prompt configuration
- **Four built-in styles**:
  - Formal v2: Clean documentation aesthetic (Inter font, blue primary palette)
  - Patchnotes v2: Gaming studio aesthetic (Rajdhani font, purple/gold palette)
  - Epic v2: Medieval fantasy aesthetic (Cinzel font, gold/crimson on dark palette)
  - Ironic v2: Deadpan minimalist aesthetic (Comic Sans, coral/teal palette)
- **User style management**: Create, import (file/URL), export, and delete custom styles stored in `~/.config/commitlore/styles/`
- **Output formats**: Terminal (ANSI), Markdown, JSON, and self-contained HTML with dark theme, type badges, and activity visualizations
- **Narrative rendering**: Markdown-to-HTML conversion via goldmark with XSS protection; vocabulary substitutions without requiring an LLM
- **LLM integration**: Optional enrichment via Anthropic, OpenAI, Ollama, and Groq with sanitized input, prompt delimiters, and anti-hallucination instructions in all built-in style prompts
- **Wails v3 desktop app**: Frameless window with custom titlebar, macOS-style window controls, and six screens (Dashboard, Generate, Story, History, Contributors, Settings)
- **Desktop app UI**: Design system with CSS variable tokens; all screens use theme-driven colors from the active style; sidebar with style logo and branded wordmark
- **Dashboard**: Repository picker (folder dialog, drag-and-drop, GitHub input), recent repos list, and repository summary cards
- **Generate and Story screens**: Two-column layout with filter sidebar and live HTML preview iframe; read-only style and LLM status indicators
- **History screen**: Dense commit table with monospace hash, click-to-copy, filters, and limit slider
- **Contributors screen**: Contribution table with avatar initials and relative activity bars
- **Styles screen**: Two-column layout with style card list and read-only detail panel; import, export, delete, and set-as-active actions
- **Settings screen**: LLM provider configuration with API key storage via OS keychain (Windows Credential Manager, macOS Keychain, Linux Secret Service); appearance section with style selector
- **App theming**: Active style themes the entire app via CSS variables, persisted in config.yml; per-style window control colors
- **Global state**: Svelte stores for active repository, summary, style, and UI labels; repo selected once on Dashboard, shared across all screens
- **Official logo**: Scroll + git branch SVG design with adaptive dark/light support, embedded in HTML report headers

### Changed

- Migrated desktop app framework from Wails v2 to Wails v3 alpha for Go 1.26 compatibility
- Template structure migrated from plain `.tmpl` files to `.shipstyle` YAML format
- Separated responsibilities between `internal/narrative` (text generation) and `internal/renderer` (format output)

### Fixed

- Commit message field now contains only the subject line, not the full body
- Vocabulary replacements match whole words only, preventing partial substitutions
- Theme animations gate properly disables all CSS animations when set to false

### Removed

- PDF output format removed in favor of HTML (better visual quality, compatible with browser printing)
