package app

import (
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

// App is the main application struct for Wails bindings.
type App struct {
	ctx context.Context
}

// NewApp creates a new App instance.
func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Run starts the Wails application with all bindings registered.
func Run() error {
	app := NewApp()
	gitApp := NewGitApp()
	changelogApp := NewChangelogApp()
	storyApp := NewStoryApp()
	styleApp := NewStyleApp()

	return wails.Run(createOptions(app, gitApp, changelogApp, storyApp, styleApp))
}

func createOptions(
	app *App,
	gitApp *GitApp,
	changelogApp *ChangelogApp,
	storyApp *StoryApp,
	styleApp *StyleApp,
) *options.App {
	return &options.App{
		Title:  "CommitLore",
		Width:  1200,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
			gitApp,
			changelogApp,
			storyApp,
			styleApp,
		},
	}
}
