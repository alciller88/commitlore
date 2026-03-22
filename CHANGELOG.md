CHANGELOG — CommitLore
All notable changes to this project are documented here.
Format based on Keep a Changelog.
Versioning follows Semantic Versioning.

## [Unreleased]

### Changed (P3 formal style v2)
- formal.shipstyle: bumped to v2.0.0 with improved colors (#0969DA primary, #0550AE accent, pure white background), tighter typography (13px base, 20px headers, Inter system stack), bordered card style, refined LLM prompt with executive summary instruction, enriched templates with hash+author metadata, bullet decorator changed to •

### Fixed (P3 UI bugs round 2)
- App.svelte .content: changed from overflow-y:auto to overflow:hidden flex column — was the root cause of all table scroll issues (height:100% inside scrollable parent never constrains)
- All screens: replaced height:100% with flex:1 to work with the new flex column .content container
- History/Contributors: inline styles for scroll containment (flex:1; min-height:0; overflow-y:auto on table wrapper)
- Dashboard: active-view gets overflow-y:auto for long recent lists
- Settings: screen gets overflow-y:auto for content that exceeds viewport
- Generate/Story: two-col uses flex:1 + min-height:0 instead of height:100%
- Window controls: wired to Wails v3 runtime (Window.Minimise/ToggleMaximise/Close) — were previously no-ops

### Fixed (P3 UI bugs)
- History: scrollbar no longer overlaps Fetch button — overflow: hidden on container, min-height: 0 on table
- App.svelte: repo indicator moved from sidebar bottom to 32px content topbar (always visible, all screens)
- App.svelte: scroll SVG removed from sidebar header — only style logo (if defined) or wordmark alone
- Window controls: per-style colors via WindowControls struct in .shipstyle (default, close, minimize, maximize)
- internal/styles: WindowControls struct added to Theme
- app/style_app.go: StyleTheme includes winDefault/winClose/winMinimize/winMaximize with fallback defaults
- theme.ts: injects --cl-win-default/close/minimize/maximize CSS variables
- All 4 built-in styles have window_controls fields with style-appropriate colors
- Styles editor: Window Controls color pickers added to Images & Icons tab

### Added (P3.1 + P3.5 Visual overhaul)
- app/app.go: frameless window (Wails v3 Frameless: true) for custom titlebar
- app/frontend/src/lib/design.css: design system with --space-*, --text-*, --radius-*, --transition-* tokens, global reset, scrollbar styling, button micro-interactions
- app/frontend/src/App.svelte: custom titlebar with macOS-style window controls (close/min/max circles), sidebar as drag region, VIEWS/SYSTEM nav sections, compact 220px sidebar with 32px nav items
- app/frontend/src/screens/Dashboard.svelte: three entry cards (folder/drop/GitHub), stat cards (commits, contributors, last commit), dense recent repos list
- app/frontend/src/screens/Generate.svelte: 240px sidebar with FILTERS label, accent action button, dot LLM indicator, right-aligned toolbar
- app/frontend/src/screens/Story.svelte: same pattern as Generate
- app/frontend/src/screens/History.svelte: single 40px filter row, 36px dense table rows with alternating background, monospace hash with click-to-copy
- app/frontend/src/screens/Contributors.svelte: avatar initials circles (24px), 4px activity bars, 44px rows
- app/frontend/src/screens/Settings.svelte: section headers with bottom borders, outline/primary button split, 20px color swatches

### Added (P3.2 Styles screen overhaul)
- app/style_app.go: GetStyleDetail binding — returns all .shipstyle fields as typed StyleDetail struct
- app/style_app.go: SaveStyleDetail binding — saves all fields, rejects built-in styles
- app/style_app.go: IsStyleBuiltIn binding — returns true for built-in style names
- app/style_app.go: StyleDetail, TemplatesDetail, ThemeDetail, ColorsDetail, TypoDetail, TerminalDetail, TermColorsDetail, DecorDetail structs
- app/frontend/src/screens/Styles.svelte: complete overhaul — two-column layout with style cards (left) and 5-tab editor (right)
- Styles screen: Colors tab — 7 color pickers with hex input, mode toggle (dark/light)
- Styles screen: Typography tab — font family with live preview, font sizes
- Styles screen: Images tab — logo/header image URL fields, card style radio, animations toggle
- Styles screen: Templates tab — textarea per template field with variable hints, vocabulary key-value editor
- Styles screen: Advanced tab — terminal colors, decorators, density, custom CSS, LLM prompt
- Styles screen: built-in styles read-only with info banner, user styles fully editable with Save/Delete/Export
- Styles screen: new style creation flow with name validation and duplicate detection
- app/style_app_test.go: TestGetStyleDetail_formal, TestSaveStyleDetail_user, TestSaveStyleDetail_builtinRejected, TestIsStyleBuiltIn_formal, TestIsStyleBuiltIn_userStyle

### Fixed (P2 style cleanup)
- app/frontend: removed redundant style selector from Generate and Story — active style now read from global store
- app/frontend: added read-only style pill in Generate/Story sidebar showing active style name with accent color
- internal/renderer: added RenderWithTheme/RenderStoryWithTheme with HTMLTheme override — HTML output reflects active style colors
- app/changelog_app.go, app/story_app.go: pass active style's theme as HTMLTheme override to renderer

### Added (P2 UI theming)
- app/config_app.go: GetActiveStyle/SetActiveStyle — active style persisted in config.yml, defaults to "formal"
- app/style_app.go: GetStyleTheme binding returns typed StyleTheme struct with fallback defaults for missing theme fields
- app/style_app.go: StyleTheme struct with primary, secondary, background, surface, text, accent, border, fontFamily, fontSize, mode, logo fields
- app/frontend/src/lib/theme.ts: applyTheme() fetches StyleTheme and injects CSS variables on :root
- app/frontend/src/lib/store.ts: activeStyle writable store
- app/frontend: all screens use CSS variables (--cl-background, --cl-surface, --cl-text, --cl-accent, --cl-secondary, --cl-border, --cl-font-family) — no hardcoded structural colors
- app/frontend/src/App.svelte: branded sidebar header with style logo + "Commit"/"Lore" split wordmark colored by primary/accent
- app/frontend/src/screens/Settings.svelte: Appearance section with style dropdown and color swatch preview
- app/style_app_test.go: tests for GetStyleTheme (formal, patchnotes, missing fields, unknown style)
- app/config_app_test.go: tests for GetActiveStyle/SetActiveStyle

### Fixed (Markdown rendering)
- internal/renderer: narrative content now rendered as HTML via goldmark — was displaying raw markdown (## headings, * bullets) in iframe
- internal/renderer: goldmark safe mode omits raw HTML tags (XSS protection)
- go.mod: github.com/yuin/goldmark promoted to direct dependency

### Added (P1 UX global)
- app/frontend/src/lib/store.ts: Svelte writable stores for activeRepo and repoSummary — single source of truth across all screens
- app/frontend: sidebar repo indicator with SVG icons (folder/GitHub), repo name, truncated path
- app/frontend: Dashboard shows cached summary on return (no loading spinner), "Change repository" button in top-right
- app/frontend: Generate and Story two-column layout (280px form sidebar + flex iframe preview)
- app/frontend: History and Contributors compact horizontal filter rows
- app/frontend: repo picker removed from Generate, Story, History, Contributors — all screens read from global store
- app/frontend: "Select a repository in Dashboard" banner when no repo active
- app/frontend: App.svelte preloads most recent repo on startup from GetRecentRepos()
- app/frontend: all nav icons replaced with inline SVGs (no emoji)

### Fixed (Story/Changelog output parity)
- internal/renderer: HTML output now includes narrative content — was previously discarded, producing bare data-only HTML while terminal/markdown got rich styled narrative
- internal/renderer: narrative text rendered as HTML paragraphs in .narrative div above structured data section
- internal/renderer: CSS for .narrative and .data-section layout

### Fixed (LLM Settings)
- app/changelog_app.go, app/story_app.go: Generate and GenerateStory now auto-read LLM config from Settings (config.yml + keychain) — removed redundant llmProvider/llmModel parameters
- app/frontend: removed LLM selectors from Generate and Story screens — Settings is the single source of truth for LLM configuration
- app/frontend: added read-only LLM status indicator on Generate and Story screens

### Fixed (Phase 11 bugs)
- internal/git: added JSON tags to Commit struct — frontend was receiving PascalCase fields but expected camelCase
- app/frontend: replaced HTML entity icons (&#9729; &#128194;) with inline SVGs — WebView didn't render emoji entities
- app/app.go: enabled native file drop via Wails v3 EnableFileDrop + WindowFilesDropped event — HTML5 drag events don't receive OS files in WebView
- app/changelog_app.go: resolveAPIKey() checks OS keychain via go-keyring when env var is empty — UI-configured LLM keys now reach the provider

### Added (Phase 11)
- app/config_app.go: ConfigApp binding — GetRecentRepos, AddRecentRepo, SetLLMConfig, GetLLMConfig, ClearLLMKey
- app/config_app.go: recent repos persisted in ~/.config/commitlore/config.yml (max 10 entries, MRU order)
- app/config_app.go: LLM API key storage via OS keychain (go-keyring: Windows Credential Manager, macOS Keychain, Linux Secret Service)
- app/git_app.go: OpenFolderPicker() — native folder dialog via Wails v3 Dialog API
- app/frontend/src/screens/Dashboard.svelte: repo picker (folder, drag & drop, GitHub input) + recent repos + repo summary
- app/frontend/src/screens/Generate.svelte: changelog form with style/LLM selectors, HTML preview iframe, copy/save
- app/frontend/src/screens/Story.svelte: narrative form with iframe preview, copy/save
- app/frontend/src/screens/History.svelte: commit table with hash copy, filters, limit slider
- app/frontend/src/screens/Contributors.svelte: contribution table with relative activity bars
- app/frontend/src/screens/Styles.svelte: style manager with built-in/user badges, import/export/delete/create, detail panel
- app/frontend/src/screens/Settings.svelte: LLM provider config, API key modal (OS keychain), about section
- app/config_app_test.go: unit tests for ConfigApp (recent repos limit, deduplication, LLM config, persistence)
- SPEC.md §9: repo picker, output display, LLM configuration, Settings screen specifications
- go.mod: github.com/zalando/go-keyring dependency

### Changed
- app/: migrated from Wails v2 to Wails v3 alpha for Go 1.26 compatibility
- app/app.go: v3 API — application.New() with Services, Window.NewWithOptions(), AssetOptions.Handler
- SPEC.md §3: Framework → Wails v3 alpha, Build → wails3 CLI

### Added
- app/: Wails desktop app base structure with Svelte + TypeScript frontend (Phase 10)
- app/app.go: main Wails entry point with Run(), binds all domain structs (Phase 10)
- app/git_app.go: GitApp binding — History() and Contributors() for local and GitHub repos (Phase 10)
- app/changelog_app.go: ChangelogApp binding — Generate() with optional LLM enrichment (Phase 10)
- app/story_app.go: StoryApp binding — GenerateStory() with chronology and LLM support (Phase 10)
- app/style_app.go: StyleApp binding — ListStyles, ShowStyle, ImportStyle, ExportStyle, DeleteStyle, CreateStyle (Phase 10)
- app/frontend/: Svelte frontend with sidebar navigation and 6 placeholder screens (Phase 10)
- internal/llm: optional LLM integration with Anthropic and OpenAI adapters (Phase 9)
- internal/llm: provider.go — Provider interface with Enrich() method (Phase 9)
- internal/llm: anthropic.go — Anthropic Messages API adapter, default model claude-haiku-4-5-20251001 (Phase 9)
- internal/llm: openai.go — OpenAI-compatible Chat Completions adapter, default model gpt-4o-mini (Phase 9)
- internal/llm: aliases.go — convenience aliases: ollama → localhost:11434, groq → api.groq.com (Phase 9)
- internal/llm: factory.go — New() factory with provider validation and 30s timeout (Phase 9)
- internal/llm: sanitize.go — SanitizeRepoData() truncates to 500 chars, escapes control chars (Phase 9)
- internal/llm: prompt.go — BuildPrompt() wraps data with ---DATA START/END--- delimiters (Phase 9)
- cmd/generate: --llm and --llm-base-url flags for LLM enrichment (Phase 9)
- cmd/story: --llm and --llm-base-url flags for LLM enrichment (Phase 9)
- SPEC.md §8: updated with --llm-base-url, aliases (ollama, groq), env var COMMITLORE_LLM_BASE_URL (Phase 9)
- internal/github: GitHub API integration via go-github — client, commits, and PRs (Phase 8)
- internal/github: client.go — GitHub client with optional GITHUB_TOKEN authentication (Phase 8)
- internal/github: repo.go — fetch commits from remote repos with --since, --until, --author, --limit filters (Phase 8)
- internal/github: prs.go — fetch merged PRs with date range filters, read-only GET only (Phase 8)
- internal/github: errors.go — user-friendly error messages for 404, 401, 403, 429 (Phase 8)
- cmd/reposource.go: shared repo source detection — auto-routes to local git or GitHub API (Phase 8)
- cmd/generate: --repo owner/repo and GitHub URL support, --include-prs with PR data in changelog (Phase 8)
- cmd/history: --repo owner/repo and GitHub URL support (Phase 8)
- cmd/contributors: --repo owner/repo and GitHub URL support (Phase 8)
- cmd/story: --repo owner/repo and GitHub URL support (Phase 8)
- internal/changelog: AppendCommit() and InferTypeFromMessage() for PR integration (Phase 8)
- cmd/style: style command with subcommands list, show, create, import, export, delete (Phase 7)
- internal/styles: user style management — load, save, list, delete, import from file/URL, export (Phase 7)
- internal/styles: user styles directory at ~/.config/commitlore/styles/ with cross-platform support (Phase 7)

### Changed
- assets/logo.svg: replaced with square scroll icon (400x400 viewBox, commit hashes only)
- internal/renderer: header shows icon-only at 100x100px, removed "CommitLore" text

### Fixed
- internal/git: commit Message now contains only subject line (first line), not full body
- internal/narrative: vocabulary replacements now match whole words only (no partial matches)
- internal/renderer: animations gate — theme.Animations=false disables all CSS animations
- internal/renderer: terminal bullet, indent, and per-section colors driven by style config
- internal/renderer: compact density strips author/date details from terminal output

### Removed
- PDF format removed in favor of HTML (better visual quality, compatible with browser printing)

### Added
- assets/logo.svg: official CommitLore logo (scroll + branch design, adaptive dark/light)
- internal/renderer: embedded logo in HTML header (48px) and footer (24px) as default
- internal/styles: extended .shipstyle schema with vocabulary, theme, terminal, and marketplace metadata (Phase 6.5)
- internal/renderer: dynamic HTML with colors, typography, card_style, and custom_css from theme (Phase 6.5)
- internal/renderer: terminal.go with ANSI colors and decorators from style (Phase 6.5)
- internal/narrative: ApplyVocabulary() for word substitutions without LLM (Phase 6.5)
- internal/styles: validation for card_style, density, and mode (Phase 6.5)
- internal/renderer: self-contained HTML format with dark theme, type badges, and activity bars (Phase 6)
- cmd/story: story command with flags --repo, --from, --style, --format, --output (Phase 5)
- internal/git: chronology, tags, and activity peaks for story command (Phase 5)
- internal/narrative: GenerateStory() with story templates per style (Phase 5)
- internal/renderer: RenderStory() with terminal ANSI and JSON support (Phase 5)
- internal/styles: story_intro, story_milestone, story_peak, story_contributor, story_footer fields in .shipstyle (Phase 5)
- internal/git: local repo access with go-git (Phase 2)
- internal/changelog: commit parsing and grouping by type (Phase 3)
- cmd/history: history command with filters --author, --since, --until, --limit (Phase 2)
- cmd/contributors: contributors command with flags --repo, --since, --top (Phase 3)
- cmd/generate: generate command with flags --repo, --since, --until, --style, --format, --output (Phase 4)
- internal/narrative: text generation by style with embedded templates (Phase 4)
- internal/styles: modular style system with .shipstyle YAML format (Phase 4)
- internal/renderer: format rendering (terminal, md, json) (Phase 4)

### Changed
- SPEC.md, CONTEXT.md, CHANGELOG.md translated to English
- Enriched built-in styles with full visual identity and creative templates
- Template structure migrated from plain .tmpl to .shipstyle (YAML) format
- Separation of responsibilities between internal/narrative and internal/renderer
- Improved built-in style templates for clearer tone differentiation

### Fixed
- Removed duplicate templates between root and internal/narrative/templates/

## [0.0.0] — 2026-03-20

### Added
- Initial project structure
- SPEC.md — complete project specification
- CONTEXT.md — context for agents and collaborators
- CHANGELOG.md — this file
