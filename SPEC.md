SPEC.md — shiplog

Fuente de verdad del proyecto. Cualquier cambio en funcionalidad debe
reflejarse aquí antes de implementarse. No se implementa nada que no
esté especificado en este documento.


1. Visión
shiplog es una herramienta multiplataforma (CLI + app de escritorio) escrita
en Go que analiza repositorios git —locales y de GitHub— y genera changelogs,
narrativas e informes sobre la historia del código, con tono y formato
configurables mediante un sistema de estilos modular.
Tagline: Your repo has a story. shiplog tells it.

2. Principios de diseño

Sin dependencias obligatorias — funciona offline, sin API key, sin cuenta.
LLM opcional — si el usuario configura una API key, el output mejora; sin ella, funciona igualmente mediante plantillas.
Multiplataforma — binario nativo CLI y app de escritorio para Linux, macOS y Windows.
Composable — múltiples formatos de output para integrarse en pipelines existentes.
Estilos modulares — los tonos no están hardcodeados; son archivos cargables, exportables e importables.
Explícito sobre mágico — toda configuración es visible vía flags o UI.
Código sostenible — funciones pequeñas, responsabilidad única, sin atajos.


3. Stack técnico
Backend / CLI
CapaTecnologíaLenguajeGo 1.22+CLI frameworkCobraConfig/flagsViperGitHub APIgo-githubGit localgo-gitHTMLhtml/template (stdlib)LLM (opcional)Anthropic API / OpenAI API (configurable)Lintergolangci-lintTeststesting (stdlib) + testify
App de escritorio
CapaTecnologíaFrameworkWails v2FrontendSvelte + TypeScriptEstilosTailwind CSS + shadcn-svelte (base componentes)BuildWails CLI (binarios nativos por plataforma)

El backend de la app de escritorio reutiliza íntegramente los paquetes
de internal/. No hay lógica duplicada.


4. Fuentes de datos (v1)

Repositorio git local — via go-git, sin autenticación.
GitHub — via API REST. Token opcional (sin token: solo repos públicos).


GitLab y otros proveedores quedan fuera de scope en v1.


5. Comandos CLI
Todos los comandos comparten los flags globales siguientes.
Flags globales
FlagValoresDefaultDescripción--formatterminal, md, json, html, pdfterminalFormato del output--stylenombre de estilo cargadoformalTono del texto generado--outputruta de archivostdoutArchivo de destino--llmanthropic, openai, nonenoneLLM a usar para enriquecer el output

5.1 shiplog generate
Genera un changelog a partir de commits y/o PRs.
shiplog generate [flags]
FlagDescripción--sinceDesde tag, commit SHA o fecha (e.g. v1.2.0, 2024-01-01)--untilHasta tag o commit SHA (default: HEAD)--repoRuta local o URL de GitHub (owner/repo)--include-prsIncluye información de PRs (requiere token GitHub)

5.2 shiplog story
Genera una narrativa completa de la historia del repositorio.
shiplog story [flags]
FlagDescripción--repoRuta local o URL de GitHub (owner/repo)--fromCommit o tag de inicio (default: primer commit)

5.3 shiplog history
Explora commits filtrados por autor, fecha o rango.
shiplog history [flags]
FlagDescripción--repoRuta local o URL de GitHub--authorFiltra por autor (nombre o email)--sinceDesde fecha o tag--untilHasta fecha o tag--limitNúmero máximo de commits (default: 50)

5.4 shiplog contributors
Mapa de quién ha tocado qué partes del código.
shiplog contributors [flags]
FlagDescripción--repoRuta local o URL de GitHub--sincePeríodo de análisis--topNúmero de contributors a mostrar (default: 10)

5.5 shiplog style
Gestión del sistema de estilos modular.
shiplog style <subcommand> [flags]
SubcomandoDescripciónlistLista los estilos disponibles (built-in + instalados)showMuestra la definición de un estilocreateCrea un nuevo estilo (wizard interactivo o flags)importImporta un estilo desde URL o ruta localexportExporta un estilo a archivo .shipstyledeleteElimina un estilo instalado (no elimina built-ins)

6. Sistema de estilos modular
Los estilos son archivos .shipstyle en formato YAML que definen el tono,
las plantillas de texto y toda la identidad visual del output.
Se almacenan en ~/.config/shiplog/styles/.

Estructura completa de un archivo .shipstyle:

  name: "nombre"
  version: "1.0.0"
  description: "descripción"
  author: "autor"
  tags: []                    # metadatos marketplace
  preview_url: ""             # URL de preview
  homepage: ""                # URL del proyecto

  llm_prompt: |               # prompt para LLM (opcional)
    ...

  templates:                  # plantillas de texto
    header, feature, fix, breaking, footer
    story_intro, story_milestone, story_peak, story_contributor, story_footer

  vocabulary:                 # sustituciones sin LLM (map[string]string)
    bug: "heresy"
    fix: "purge"

  theme:                      # identidad visual HTML
    mode: "dark"              # dark | light
    colors:
      primary, secondary, background, surface, text, accent, border
    typography:
      font_family, font_size_base, font_size_header, font_size_code
    header_image: ""          # URL o base64
    logo: ""                  # URL o base64
    card_style: "bordered"    # minimal | bordered | glassmorphism
    animations: true
    custom_css: ""            # CSS adicional inyectado al final

  terminal:                   # identidad visual terminal
    colors:
      header, feature, fix, breaking, footer   # nombres de color ANSI
    decorators:
      separator, bullet, indent
    density: "normal"         # compact | normal | verbose

Todos los campos de vocabulary, theme y terminal son opcionales con
valores zero-value sensatos.

Estilos built-in (v1)

formal — técnico y profesional, colores neutros
patchnotes — estilo videojuego, púrpura/dorado, animaciones
ironic — humor seco, colores apagados, minimalista
epic — narrativo grandilocuente, dorado/oscuro, ornamentado

Comportamiento

Sin LLM: se usan los campos templates y vocabulary del .shipstyle.
Con LLM: se usa el campo llm_prompt para instruir al modelo; templates como fallback.
Los estilos built-in no son modificables ni eliminables.


7. Formatos de output
FormatoDescripciónterminalTexto con color ANSI directo a stdoutmdMarkdown estándar, compatible con GitHubjsonEstructura de datos completa, apta para pipelineshtmlInforme HTML autocontenido con estilos inline

PDF eliminado en favor de HTML. El HTML generado puede imprimirse a PDF desde cualquier navegador.

8. LLM opcional

La herramienta funciona completamente sin LLM usando plantillas del estilo activo.
Variables de entorno: SHIPLOG_LLM_PROVIDER y SHIPLOG_LLM_API_KEY.
El flag --llm sobreescribe la variable de entorno por comando.
Proveedores soportados en v1: anthropic, openai.


9. App de escritorio (Wails)
La app comparte toda la lógica de internal/. Solo añade una capa de UI encima.
Pantallas principales (v1)

Dashboard — resumen del repo activo: última actividad, contributors, tags.
Generate — formulario para configurar y generar un changelog.
Story — visualización narrativa de la historia del repo.
History — explorador de commits con filtros.
Contributors — mapa visual de contribuciones.
Styles — gestor visual: crear, importar, exportar, previsualizar estilos.

Identidad visual

Paleta oscura por defecto, con opción de tema claro.
Tipografía monoespaciada para outputs; sans-serif para navegación.
Estilo propio construido sobre shadcn-svelte como base de componentes.
Sin dependencia de librerías de UI genéricas (no Material, no Bootstrap).


10. Estructura del proyecto
shiplog/
├── cmd/                      # Comandos Cobra
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
│   ├── git/                  # Acceso a repos locales (go-git)
│   ├── github/               # Acceso a GitHub API
│   ├── changelog/            # Parsing y agrupación de commits
│   ├── narrative/            # Generación de texto
│   ├── renderer/             # Renderizado (terminal, md, json, html, pdf)
│   ├── llm/                  # Adaptadores LLM (Anthropic, OpenAI)
│   └── styles/               # Carga, validación y gestión de .shipstyle
│       └── builtin/          # Estilos built-in embebidos (.shipstyle)
├── app/                      # App Wails
│   ├── frontend/             # Svelte + Tailwind + shadcn-svelte
│   └── app.go                # Bindings Wails → internal/
├── styles/                   # Estilos del usuario (.shipstyle)
├── .github/
│   └── workflows/
│       ├── ci.yml            # Lint + tests en push/PR
│       └── release.yml       # Binarios multiplataforma en tag v*
├── CHANGELOG.md
├── SPEC.md
├── CONTEXT.md
├── main.go
└── README.md

Regla: ningún paquete en internal/ importa de cmd/ ni de app/.
El flujo de dependencias es siempre hacia adentro.


11. Reglas de código

Funciones de máximo 40 líneas. Si crece, se extrae.
Una responsabilidad por función/struct.
Sin comentarios obvios — los comentarios explican el porqué, no el qué.
Errores explícitos — nunca ignorar un error con _.
Tests unitarios obligatorios para todo el código en internal/.
Cobertura mínima objetivo: 70% por paquete.
Antes de cada PR: golangci-lint run ./... debe pasar sin errores.


12. CI/CD
Estrategia de ramas
RamaPropósitomainProducción. Solo recibe merges desde dev via PR.devIntegración. Rama base para features.feat/*Feature branches. Se abren desde dev.fix/*Bugfix branches. Se abren desde dev.
Pipeline CI — .github/workflows/ci.yml
Trigger: push a dev, PR hacia main.

golangci-lint run ./...
go test ./... -race -coverprofile=coverage.out
go build ./...

Pipeline Release — .github/workflows/release.yml
Trigger: tag v* pusheado a main.

CI completo (lint + tests)
Build binarios CLI:

GOOS=linux GOARCH=amd64
GOOS=darwin GOARCH=arm64
GOOS=windows GOARCH=amd64


Build app Wails por plataforma
Creación de GitHub Release con todos los artefactos adjuntos
Actualización automática de CHANGELOG.md


13. Versionado semántico
Formato de tag: vMAJOR.MINOR.PATCH (e.g. v1.2.0)

MAJOR — cambios incompatibles (flags eliminados, comportamiento roto).
MINOR — nuevos comandos, flags o pantallas, compatibles hacia atrás.
PATCH — bugfixes y mejoras internas sin cambio de interfaz.


14. Fases de desarrollo
FaseScopeFase 1Setup del proyecto, estructura base, CI pipeline, ramasFase 2internal/git — acceso a repos locales + comando historyFase 3internal/changelog — parsing de commits + comando contributorsFase 4Comando generate (sin LLM, plantillas)Fase 5Comando story (sin LLM, plantillas)Fase 6internal/renderer — formatos md, json, html, pdfFase 7internal/styles — sistema modular + comando styleFase 8internal/github — integración GitHub APIFase 9internal/llm — integración LLM opcional (Anthropic + OpenAI)Fase 10App Wails — estructura base + bindingsFase 11App Wails — pantallas y UI completaFase 12Pipeline release + binarios multiplataformaFase 13Pulido, docs, README, ejemplos

Cada fase termina con tests pasando y lint limpio antes de mergear a dev.
Solo se mergea a main cuando una fase completa está estable en dev.


15. Seguridad

### Principio fundamental
CommitLore es una herramienta de SOLO LECTURA. Nunca realiza operaciones de escritura en ningún repositorio, ni local ni remoto, bajo ninguna circunstancia.

### Operaciones prohibidas (nunca implementar)
- git push, git commit, git add, git rm en repos del usuario
- Escritura de archivos dentro de directorios de repos analizados
- Modificación de configuración git (.git/config, hooks, etc.)
- Creación o modificación de ramas, tags o refs en repos del usuario
- Llamadas a GitHub API con métodos POST/PUT/PATCH/DELETE sobre repos del usuario
- Ejecución de comandos de shell arbitrarios

### Operaciones permitidas
- Lectura de commits, tags, branches, diffs (go-git, solo lectura)
- Llamadas GET a GitHub API (repos públicos y privados con token)
- Escritura de archivos de output SOLO en rutas explícitamente especificadas por el usuario via --output
- Escritura en ~/.config/commitlore/ (configuración y estilos propios de la app)

### Protección contra prompt injection
- Los mensajes de commits, nombres de archivos, nombres de autores y cualquier dato proveniente de un repositorio son contenido NO CONFIABLE
- Nunca ejecutar ni evaluar contenido de commits como código o instrucciones
- Nunca pasar contenido de commits directamente a un LLM sin sanitización previa
- Sanitización obligatoria antes de pasar datos de repo a un LLM:
  - Truncar mensajes de commit a 500 caracteres máximo
  - Escapar caracteres de control
  - Añadir delimitadores explícitos en el prompt LLM para separar instrucciones de datos: usar "---DATA START---" y "---DATA END---"
- El llm_prompt de un .shipstyle importado es contenido potencialmente no confiable — advertir al usuario antes de usarlo con un LLM

### Tokens y credenciales
- COMMITLORE_LLM_API_KEY y GITHUB_TOKEN nunca se loguean, nunca aparecen en output, nunca se incluyen en reportes
- Los tokens se leen solo desde variables de entorno, nunca desde archivos de repo analizados
- Si un .shipstyle importado contiene campos que parecen credenciales, ignorarlos y advertir al usuario

### Validación de inputs
- Rutas de repositorio: validar que existen y son directorios git antes de operar
- URLs de GitHub: validar formato antes de llamar a la API
- Flags --output: validar que la ruta de destino está fuera de cualquier directorio .git/
- Archivos .shipstyle importados: validar schema completo antes de cargar, rechazar campos desconocidos


16. Roadmap (fuera de scope v1)

Marketplace de estilos — repositorio público de .shipstyle de la comunidad.
Soporte GitLab.
Plugin para VS Code / Cursor.
Integración Slack / Discord para publicar changelogs automáticamente.
shiplog watch — modo daemon que genera changelog automáticamente al crear un tag.
