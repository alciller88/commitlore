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
CampoValorFase actualFase 2 — CompletadaÚltima ramadevVersiónv0.0.0Tests pasandoSí (12 tests, 97% cobertura en internal/git)Lint limpioSí

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


9. Flujo de trabajo Git — instrucciones para agentes

El agente es responsable del ciclo completo de Git al finalizar cualquier tarea.
No debe pedir al usuario que haga merges, PRs ni pushes manualmente.

### Flujo obligatorio por fase

1. Crear rama desde dev:
   git checkout dev
   git pull origin dev
   git checkout -b feat/<nombre-fase>

2. Implementar la fase en commits pequeños y atómicos (Conventional Commits).

3. Antes de abrir PR, verificar que pasa todo:
   golangci-lint run ./...
   go test ./... -count=1

4. Push de la rama:
   git push -u origin feat/<nombre-fase>

5. Abrir PR de feat/<nombre-fase> → dev usando gh CLI:
   gh pr create \
     --base dev \
     --head feat/<nombre-fase> \
     --title "feat: <descripción de la fase>" \
     --body "<resumen de cambios, qué se ha implementado, qué tests cubren>"

6. Esperar a que el CI pase:
   gh pr checks --watch

7. Si el CI pasa, mergear el PR:
   gh pr merge --squash --delete-branch

8. Volver a dev y sincronizar:
   git checkout dev
   git pull origin dev

9. Actualizar CONTEXT.md:
   - "Fase actual" → siguiente fase o "completada"
   - Añadir fila en "Historial de fases completadas"

10. Commitear y pushear el CONTEXT.md actualizado:
    git add CONTEXT.md
    git commit -m "chore: update CONTEXT.md — phase <N> completed"
    git push origin dev

### Reglas estrictas

- NUNCA abrir PR directamente a main. Siempre feat/* → dev.
- NUNCA mergear si el CI no ha pasado.
- NUNCA mergear dev → main manualmente. Eso lo hace el humano cuando decide que la fase está estable.
- go test debe ejecutarse sin -race en Windows (CGO no disponible). La CI en Ubuntu lo ejecutará con -race.
- gh CLI está instalado y autenticado. Usarlo siempre para PRs y checks.

### Merge dev → main

Solo el humano decide cuándo mergear dev → main. El agente nunca lo hace.
Cuando el humano quiera hacerlo, ejecutará:
   gh pr create --base main --head dev --title "release: <version>" --body "<resumen>"


10. Historial de fases completadas
FaseDescripciónFechaRamaFase 1Setup del proyecto, estructura base, CI pipeline, ramas2026-03-20devFase 2internal/git — acceso a repos locales + comando history2026-03-20feat/phase-2-history

Añadir una fila aquí al completar cada fase.
