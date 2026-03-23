CHANGELOG — CommitLore
All notable changes to this project are documented here.
Format based on Keep a Changelog.
Versioning follows Semantic Versioning.

## [Unreleased]

### Fixed (Chart iframe timing and visual bugs)
- styles/builtin: all 14 `new Chart()` calls across 4 styles wrapped in `setTimeout(..., 100)` — fixes blank charts in Wails WebView iframe where canvas reports width:0 before layout completes
- styles/builtin/ironic.shipstyle: Clippy bubble changed from `position: fixed` to `position: sticky` and moved inside `.word-document` div — no longer covers content in iframe
- styles/builtin/ironic.shipstyle: removed duplicate "when people actually worked" chart title in story template (was showing as section label + chart title)
- styles/builtin/ironic.shipstyle: commit list now uses `{{.Icon}}` per-item instead of `{{$.Icons.Bullet}}` — breaking changes show `!` icon instead of generic `→`
- styles/builtin/formal.shipstyle: stat cells "First Commit" and "Most Active Month" inline `style="font-size:14px/16px"` replaced with `.rp-stat-value.is-date` CSS class (consistent 15px)
- Tests: 4 new test cases (setTimeout wrapping, ironic duplicate title, ironic icon distinction, formal inline font-size removal)

### Fixed (Charts, cache, and story peaks)
- styles/builtin: replaced all `color-mix()` in Chart.js `<script>` blocks with `hexToRgba()` JS helper across all 4 styles (11 occurrences) — Chart.js cannot parse CSS functions, was silently discarding colors causing blank charts
- app/story_app.go, cmd/story.go: `storyTopPeaks` increased from 3 to 12 — story charts now show up to 12 months of data for meaningful line/bar visualizations
- app/repo_cache.go: new in-memory commit cache keyed by repo ref — avoids re-fetching commits from GitHub API when user changes style on the same repo
- app/changelog_app.go, app/story_app.go: Generate/GenerateStory use commit cache for default opts (no author/since/until filters), bypass cache when filters are active
- app/git_app.go: ClearCommitCache() binding exposed for frontend to call when user switches repos
- Tests: TestBuiltinStyles_noColorMixInScripts (verifies no color-mix in script blocks), 6 commit cache tests (miss/hit/invalidate/shouldUseCache)

### Fixed (Visual and style-specific bugs)
- styles/builtin/formal.shipstyle: story banner meta no longer shows Generated date twice — restructured to show date range OR standalone date, not both
- styles/builtin/patchnotes.shipstyle: added 300ms setTimeout fallback in both changelog and story IntersectionObserver scripts — fixes blank page in Wails iframe when observer never fires
- styles/builtin/epic.shipstyle: changelog canvas IDs changed from generic "activityChart" to "epic-cl-type-chart" and "epic-cl-activity-chart" — prevents Chart.js ID collision
- styles/builtin/ironic.shipstyle: changelog title now shows "changelog" (ironic lowercase) instead of generic "Changelog", story title shows "a story i guess" instead of "Repository Story" — uses `eq` template function to detect and replace generic titles
- styles/builtin/ironic.shipstyle: story contributor table third column changed from redundant "Name · since Date" to just "since Date", header changed to "joined"
- Tests: 8 new test cases covering all 5 visual bugs

### Fixed (Renderer data bugs)
- internal/renderer: Version field now threaded from `until` param through Render/RenderWithTheme to buildChangelogContext — version badges ({{if .Version}}) now render when `until` is a semver tag
- internal/renderer: SemverFromString exported helper validates vMAJOR.MINOR.PATCH format, returns "" for non-semver strings
- app/changelog_app.go: Generate() passes `SemverFromString(until)` as version to renderChangelog
- internal/renderer: FirstAuthor fallback — if ch.FirstCommit.Author is empty, defaults to "an unknown contributor" instead of leaving blank
- internal/renderer/css.go: writeTypeBadgeCSS now includes .type-refactor, .type-docs, .type-test, .type-chore CSS rules (were missing, only feat/fix/breaking/other existed)
- internal/renderer/html.go: typeBadgeClass returns specific classes for refactor/docs/test/chore types instead of falling through to "type-other"
- styles/builtin/formal.shipstyle: added CSS rules for .cl-type-badge.type-refactor/docs/test/chore and .type-other in the custom HTML template
- internal/renderer/html_test.go: sampleStoryChronology contributors now include Count (42/20) matching real data
- Tests: 11 new test cases (semverFromString table test, version propagation, version badge rendering, contributor count context/output, firstAuthor fallback, type badges default+custom renderer)

### Fixed (Renderer critical fixes)
- internal/renderer: writeSiteHeaderLogo now detects inline SVG logos (starts with `<svg`) and writes them raw into a `<div class="logo-svg">` instead of escaping into an `<img>` tag — fixes patchnotes/epic/ironic logos rendering as escaped text
- styles/builtin: all `{{.Theme.Logo}}` references in patchnotes.shipstyle (3 occurrences) and formal.shipstyle (1 occurrence) replaced with `{{safeHTML .Theme.Logo}}` to prevent auto-escaping in html/template
- internal/renderer: `Generated` field in buildChangelogContext/buildStoryContext changed from hardcoded "Generated by CommitLore" to `time.Now().Format("2 Jan 2006")` — templates now show actual generation date
- internal/renderer: `mul`, `divf`, and `divi` template functions added to FuncMap — fixes patchnotes story template parse errors for percentage bar calculations
- styles/builtin/patchnotes.shipstyle: `$pct` initialized as `0.0` (float64) instead of `0` (int), `ge .Count (divf ...)` changed to `ge .Count (divi ...)` for type-safe comparison
- Tests: 10 new test cases covering all 4 bugs (SVG logo rendering, URL logo, safeHTML in custom templates, Generated date format, FontSizeH default, mul/divf/divi functions, divf-by-zero, patchnotes story render)

### Added (Story HTML templates)
- internal/styles: HTMLTemplate field split into HTMLTemplateChangelog + HTMLTemplateStory — separate templates for changelog and story output
- html_template_changelog: renamed from html_template, keeps existing changelog templates
- html_template_story: new field with unique story templates per style:
  - formal: Executive Activity Report — line chart for commit arc, vertical dot timeline, contributor ranking bar chart, executive summary stats box
  - patchnotes: Dev Diary — blog-post style monthly entries, season stats card, activity bar chart, dev team grid, fade-in animations
  - epic: The Saga — 40px Cinzel title, prologue box, chapter-per-month with Roman numerals, heraldic fellowship circles, full ornamental treatment
  - ironic: "A Story I Guess" — fake Word titlebar (astory.doc), ridiculously understated chapter titles, resigned commentary, Clippy with story-specific message
- Template function added: `add` (integer addition for chapter numbering)
- Tests: TestRenderStory_usesHTMLTemplateStory, TestRenderStory_fallsBackToDefault, existing tests updated

### Changed (Styles screen simplified)
- Styles.svelte: complete rework — removed full editor (Colors/Typography/Icons/Images/Templates/HTML/Advanced tabs), create form, save functionality
- New read-only detail panel: logo (48px), name/version badge, author, description, 7 color circles (20px) with hex tooltip, font preview, mode badge, UI labels (if custom), icons (if custom), collapsible LLM prompt
- Actions: "Set as active" (changes app theme immediately), Export (all styles), Delete (user only with confirmation)
- Left column bottom: "Import style" (file picker), "Get more styles" (opens https://commitlore.dev/styles in browser)
- Styles are managed via import/export — no in-app creation or editing
- Removed bindings: SaveStyleDetail no longer called from Styles screen

### Added (HTML template system)
- internal/styles: HTMLTemplate string field added to Style struct — per-style custom HTML templates using Go html/template syntax
- internal/renderer: HTMLTemplateContext with all style fields (Theme, Icons, UILabels), unified for both changelog and story contexts
- internal/renderer: renderCustomChangelogHTML / renderCustomStoryHTML — uses HTMLTemplate when set, falls back to default renderer when empty
- Template helper functions: upper, lower, initials (extracts 1-2 initials from a name)
- formal.shipstyle: Stripe/GitHub docs template — Inter font, 60px header, stats cards, commits table, Chart.js bar+donut charts, contributors with initials circles
- patchnotes.shipstyle: Steam/game studio template — Rajdhani font, 120px gradient header, highlights cards, Chart.js charts, fade-in animations on scroll
- epic.shipstyle: Medieval chronicle template — Cinzel font, parchment texture, drop cap, decree scroll entries, separator ornaments, "Chronicles of Activity" chart, Fellowship contributors
- ironic.shipstyle: Fake Word 95 template — Comic Sans, titlebar+toolbar chrome, Clippy speech bubble, bulleted list (no table), horizontal bar chart, Word-style table
- All 4 templates: Chart.js from CDN, all colors from .shipstyle theme (no hardcoded hex), client-side chart rendering
- Styles.svelte: new "HTML" tab with monospace textarea (min 20 rows) for raw HTML template editing, template variable syntax warning
- app/style_app.go: HTMLTemplate field in StyleDetail, styleToDetail/detailToStyle mappings
- Tests: TestRenderChangelog_usesHTMLTemplate, TestRenderChangelog_fallsBackToDefault, TestHTMLTemplateContext_allFieldsPopulated
- SPEC.md: html_template field documented in .shipstyle schema

### Changed (GitHub Connect modal)
- Dashboard: GitHub entry card replaced with clean button that opens a centered modal (400px) with owner/repo input, optional token (password field, session-only), format validation, keyboard support (Enter/Escape)
- app/git_app.go: SetGitHubToken() binding — sets GITHUB_TOKEN env var for session without persisting to disk

### Added (Icons system)
- internal/styles: Icons struct (feature, fix, breaking, chore, docs, test, story_peak, bullet, separator) added to Style
- app/style_app.go: IconsDetail in StyleTheme and StyleDetail, buildIcons() with fallback defaults
- All 4 built-in styles: icons block with per-style characters (formal: minimal, patchnotes: gaming, epic: medieval, ironic: deadpan arrows)
- Styles.svelte: new "Icons" tab with Commit icons, Story icons, and Decorators sections, emoji preview per field

### Added (UI button labels + Styles editor)
- internal/styles: GenerateButton, StoryButton fields added to UILabels struct
- app/style_app.go: UILabelsDetail extended with generateButton/storyButton, defaults "Generate"/"Tell the story"
- Generate.svelte: button text uses $uiLabels.generateButton reactively
- Story.svelte: button text uses $uiLabels.storyButton reactively
- Styles.svelte: Templates tab gains "UI Labels" section with 9 editable fields (dashboard, generate, generate button, story, story button, history, contributors, styles, settings)
- All 4 built-in styles updated with button labels: formal (Generate/Tell the story), patchnotes (Deploy Patch/Write the Dev Diary), epic (Inscribe the Chronicle/Begin the Saga), ironic (Fine, generate/Sure, a story)

### Fixed (patchnotes llm_prompt)
- patchnotes.shipstyle: llm_prompt rewritten to infer categories from actual commit content instead of using predefined gaming categories

### Fixed (Dashboard GitHub card)
- Dashboard: GitHub repo picker card restructured — input and Connect button now in a joined horizontal row (shared border, no gap), input takes full width, button never wraps, Enter key triggers connect, loading shows "..."

### Changed (P3 ironic style v2)
- ironic.shipstyle v2.0.0: deadpan minimalist aesthetic (monochrome grays on #1A1A1A), JetBrains Mono font, Phosphor Minus SVG logo, minimal card style, ui_labels (whatever, fine, a story i guess, stuff that happened, people, aesthetics, knobs), anti-hallucination LLM prompt, simplified vocabulary

### Changed (P3 epic style v2)
- epic.shipstyle v2.0.0: medieval fantasy aesthetic (gold #C9A84C / crimson #8B1A1A on dark #0F0A05), Cinzel serif font, Phosphor Sword inline SVG logo, bordered card style, ui_labels (The Keep, The Chronicle, The Saga, The Scrolls, The Fellowship, The Wardrobe, The Forge), refined LLM prompt with anti-hallucination instruction, simplified vocabulary

### Added (P3 patchnotes style v2 + UILabels)
- internal/styles: UILabels struct added to .shipstyle schema — per-style navigation label overrides (dashboard, generate, story, history, contributors, styles, settings)
- app/style_app.go: UILabelsDetail in StyleTheme and StyleDetail, buildUILabels() with English defaults fallback
- app/frontend: uiLabels writable store, applyTheme() updates labels reactively, App.svelte nav uses $uiLabels
- App.svelte: inline SVG logo support via {@html} when logo field starts with `<svg`
- patchnotes.shipstyle v2.0.0: professional gaming aesthetic (purple #7C6FCD / gold #F0A500), Rajdhani font, Phosphor GameController SVG logo, ui_labels overrides (Hub, Patch Notes, Dev Diary, Commit Log, Dev Team, Themes, Options), refined LLM prompt

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
