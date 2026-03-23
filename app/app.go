package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	gitApp := NewGitApp()
	changelogApp := NewChangelogApp()
	storyApp := NewStoryApp()
	styleApp := NewStyleApp()
	marketplaceApp := NewMarketplaceApp()
	configApp := NewConfigApp()

	app := application.New(application.Options{
		Name: "CommitLore",
		Services: []application.Service{
			application.NewService(gitApp),
			application.NewService(changelogApp),
			application.NewService(storyApp),
			application.NewService(styleApp),
			application.NewService(marketplaceApp),
			application.NewService(configApp),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
	})

	win := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "CommitLore",
		Width:            1200,
		Height:           800,
		URL:              "/",
		Frameless:        true,
		BackgroundColour: application.NewRGB(13, 17, 23),
		EnableFileDrop:   true,
	})

	win.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
		files := event.Context().DroppedFiles()
		if len(files) > 0 {
			app.Event.Emit("file-dropped", files[0])
		}
	})

	return app.Run()
}
