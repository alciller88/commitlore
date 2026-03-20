CONTEXT.md — shiplog

Documento de contexto para agentes de IA y colaboradores.
Describe el estado actual del proyecto, decisiones tomadas, convenciones
de trabajo y advertencias importantes.
Se actualiza al inicio de cada fase y cuando cambia algo relevante.


1. ¿Qué es shiplog?
CLI + app de escritorio en Go que analiza repos git (locales y GitHub) y
genera changelogs, narrativas e informes con tono configurable mediante un
sistema de estilos modular (archivos .shipstyle).
Consultar SPEC.md para la especificación completa.

2. Estado actual
CampoValorFase actualFase 1 — En progresoÚltima rama—Versiónv0.0.0Tests pasando—Lint limpio—

Actualizar esta tabla al iniciar y completar cada fase.


3. Decisiones técnicas tomadas
Estas decisiones están cerradas. No se debaten en cada sesión.
DecisiónElecciónMotivoLenguajeGo 1.22+Multiplataforma, binarios nativos, sin runtimeCLI frameworkCobra + ViperEstándar de facto en GoGit localgo-gitPuro Go, sin dependencia del binario gitGitHub APIgo-githubMantenida por Google, tipadaApp escritorioWails v2Go nativo + WebView del OS, sin ChromiumFrontendSvelte + TypeScriptCompila a vanilla JS, rendimiento nativo en WailsEstilos UITailwind + shadcn-svelteBase de componentes sin look genéricoLintergolangci-lintEstándar en Go, agrupa 50+ lintersTeststesting + testifyStdlib + aserciones legiblesPDFgofpdfPuro Go, sin dependencias externasVersionadoSemver (vMAJOR.MINOR.PATCH)Estándar universalRamasmain + dev + feat/* / fix/*Flujo claro, CI diferenciado

4. Convenciones de código
Go

Funciones de máximo 40 líneas. Si crece, se extrae.
Una responsabilidad por función/struct.
Errores siempre explícitos — nunca _ para ignorar un error.
Comentarios explican el porqué, nunca el qué.
Nombres en inglés, descriptivos, sin abreviaturas crípticas.
Paquetes en minúsculas, una palabra si es posible.

Tests

Un archivo _test.go por cada archivo de lógica en internal/.
Nombres de test: TestNombreFuncion_escenario (e.g. TestParseCommit_emptyMessage).
Cobertura mínima objetivo: 70% por paquete.
Usar testify/assert para aserciones.

Svelte / Frontend

Componentes en PascalCase.svelte.
Un componente = una responsabilidad.
Props tipadas con TypeScript.
Sin lógica de negocio en componentes — solo presentación y llamadas a bindings Wails.

Git

Commits en inglés, formato Conventional Commits:
feat:, fix:, chore:, docs:, test:, refactor:
Un commit = un cambio lógico. No mezclar refactors con features.
PRs pequeños y enfocados. No mezclar fases.


5. Flujo de trabajo por fase
Seguir este flujo sin excepciones:
1. Crear rama feat/<nombre> desde dev
2. Implementar el cambio mínimo de la fase
3. Escribir tests unitarios
4. Ejecutar: golangci-lint run ./...
5. Ejecutar: go test ./... -race
6. Si todo pasa → PR a dev
7. Revisar diff antes de mergear
8. Mergear a dev
9. Actualizar "Estado actual" en CONTEXT.md
10. Solo mergear dev → main cuando la fase está completa y estable
Nunca:

Implementar dos fases en un mismo PR.
Mergear con tests fallando.
Mergear con lint con errores.
Añadir funcionalidad no especificada en SPEC.md sin actualizar SPEC.md primero.


6. Estructura de directorios relevante
shiplog/
├── cmd/               # Punto de entrada de cada comando CLI
├── internal/          # Toda la lógica de negocio (testeable, sin dependencias de UI)
│   ├── git/
│   ├── github/
│   ├── changelog/
│   ├── narrative/
│   ├── renderer/
│   ├── llm/
│   └── styles/
├── app/               # App Wails
│   ├── frontend/      # Svelte
│   └── app.go         # Bindings Go ↔ frontend
├── styles/            # Estilos built-in (.shipstyle)
└── templates/         # Plantillas HTML/texto por formato
Regla de dependencias:

internal/ no importa nada de cmd/ ni de app/.
cmd/ y app/ importan de internal/.
Nunca crear dependencias circulares.


7. Variables de entorno
VariablePropósitoRequeridaSHIPLOG_LLM_PROVIDERProveedor LLM (anthropic, openai)NoSHIPLOG_LLM_API_KEYAPI key del proveedor LLMNoGITHUB_TOKENToken GitHub para repos privados/PRsNo

8. Instrucciones para agentes de IA
Si eres un agente trabajando en este proyecto, lee esto antes de escribir código:

Lee SPEC.md primero. No implementes nada que no esté especificado ahí.
Consulta la fase actual en la sección "Estado actual" de este documento.
No adelantes fases. Si ves algo que falta de una fase posterior, anótalo en un comentario // TODO(faseN): pero no lo implementes.
Tests primero o junto al código. No entregues código sin tests en internal/.
Ejecuta lint antes de terminar. El comando es golangci-lint run ./....
Funciones pequeñas. Si una función supera 40 líneas, divídela antes de continuar.
No cambies decisiones técnicas de la sección 3 sin consultar al humano.
Un cambio a la vez. Si necesitas refactorizar algo para implementar la fase, hazlo en un commit separado.
Actualiza este documento si el estado del proyecto cambia.
Ante la duda, pregunta. Es mejor pedir aclaración que implementar algo incorrecto.


9. Historial de fases completadas
FaseDescripciónFechaRama————

Añadir una fila aquí al completar cada fase.
