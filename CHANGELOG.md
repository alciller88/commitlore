CHANGELOG — CommitLore
All notable changes to this project are documented here.
Format based on Keep a Changelog.
Versioning follows Semantic Versioning.

## [Unreleased]

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
