CHANGELOG — shiplog
Todos los cambios notables de este proyecto se documentan aquí.
Formato basado en Keep a Changelog.
Versionado siguiendo Semantic Versioning.

## [Unreleased]

### Fixed
- internal/git: commit Message now contains only subject line (first line), not full body
- internal/narrative: vocabulary replacements now match whole words only (no partial matches)
- internal/renderer: animations gate — theme.Animations=false disables all CSS animations
- internal/renderer: terminal bullet, indent, and per-section colors driven by style config
- internal/renderer: compact density strips author/date details from terminal output

### Removed
- Formato PDF eliminado en favor de HTML (mejor calidad visual, compatible con impresión desde navegador)

### Added
- internal/styles: schema .shipstyle extendido con vocabulary, theme, terminal y metadatos marketplace (Fase 6.5)
- internal/renderer: HTML dinámico con colores, tipografía, card_style y custom_css del tema (Fase 6.5)
- internal/renderer: terminal.go con colores ANSI y decoradores del estilo (Fase 6.5)
- internal/narrative: ApplyVocabulary() para sustituciones de palabras sin LLM (Fase 6.5)
- internal/styles: validación de card_style, density y mode (Fase 6.5)
- internal/renderer: formato HTML autocontenido con dark theme, badges por tipo y barras de actividad (Fase 6)
- cmd/story: comando story con flags --repo, --from, --style, --format, --output (Fase 5)
- internal/git: chronology, tags, and activity peaks for story command (Fase 5)
- internal/narrative: GenerateStory() con plantillas story por estilo (Fase 5)
- internal/renderer: RenderStory() con soporte terminal ANSI y JSON (Fase 5)
- internal/styles: campos story_intro, story_milestone, story_peak, story_contributor, story_footer en .shipstyle (Fase 5)
- internal/git: acceso a repos locales con go-git (Fase 2)
- internal/changelog: parsing y agrupación de commits por tipo (Fase 3)
- cmd/history: comando history con filtros --author, --since, --until, --limit (Fase 2)
- cmd/contributors: comando contributors con flags --repo, --since, --top (Fase 3)
- cmd/generate: comando generate con flags --repo, --since, --until, --style, --format, --output (Fase 4)
- internal/narrative: generación de texto por estilo con plantillas embebidas (Fase 4)
- internal/styles: sistema de estilos modular con formato .shipstyle en YAML (Fase 4)
- internal/renderer: renderizado por formato (terminal, md, json) (Fase 4)

### Changed
- Enriched built-in styles with full visual identity and creative templates
- Estructura de plantillas migrada de .tmpl planos a formato .shipstyle (YAML)
- Separación de responsabilidades entre internal/narrative y internal/renderer
- Improved built-in style templates for clearer tone differentiation

### Fixed
- Eliminada duplicación de plantillas entre raíz y internal/narrative/templates/

## [0.0.0] — 2026-03-20

### Added
- Estructura inicial del proyecto
- SPEC.md — especificación completa del proyecto
- CONTEXT.md — contexto para agentes y colaboradores
- CHANGELOG.md — este archivo
