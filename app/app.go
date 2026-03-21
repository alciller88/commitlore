package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
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

	app := application.New(application.Options{
		Name: "CommitLore",
		Services: []application.Service{
			application.NewService(gitApp),
			application.NewService(changelogApp),
			application.NewService(storyApp),
			application.NewService(styleApp),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "CommitLore",
		Width:            1200,
		Height:           800,
		URL:              "/",
		BackgroundColour: application.NewRGB(13, 17, 23),
	})

	return app.Run()
}
