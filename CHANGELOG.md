CHANGELOG — shiplog
Todos los cambios notables de este proyecto se documentan aquí.
Formato basado en Keep a Changelog.
Versionado siguiendo Semantic Versioning.

## [Unreleased]

### Added
- internal/git: acceso a repos locales con go-git (Fase 2)
- internal/changelog: parsing y agrupación de commits por tipo (Fase 3)
- cmd/history: comando history con filtros --author, --since, --until, --limit (Fase 2)
- cmd/contributors: comando contributors con flags --repo, --since, --top (Fase 3)
- cmd/generate: comando generate con flags --repo, --since, --until, --style, --format, --output (Fase 4)
- internal/narrative: generación de texto por estilo con plantillas embebidas (Fase 4)
- internal/styles: sistema de estilos modular con formato .shipstyle en YAML (Fase 4)
- internal/renderer: renderizado por formato (terminal, md, json) (Fase 4)

### Changed
- Estructura de plantillas migrada de .tmpl planos a formato .shipstyle (YAML)
- Separación de responsabilidades entre internal/narrative y internal/renderer

### Fixed
- Eliminada duplicación de plantillas entre raíz y internal/narrative/templates/

## [0.0.0] — 2026-03-20

### Added
- Estructura inicial del proyecto
- SPEC.md — especificación completa del proyecto
- CONTEXT.md — contexto para agentes y colaboradores
- CHANGELOG.md — este archivo
