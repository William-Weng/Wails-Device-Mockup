package main

import (
	"embed"
	"wails-device-mockup/backend"

	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {
}

func main() {

	videoPreview := &backend.VideoPreviewService{}
	imagePreview := &backend.ImagePreviewService{}
	videoCompose := &backend.VideoComposeService{}
	fileManager := &backend.FileManagerService{}

	assetHandler := backend.NewAssetHandler(assets, videoPreview, imagePreview)

	app := application.New(application.Options{
		Name:        "Device Mockup",
		Description: "一個將影片與圖片外框結合的小工具",
		Services: []application.Service{
			application.NewService(videoPreview),
			application.NewService(imagePreview),
			application.NewService(videoCompose),
			application.NewService(fileManager),
		},
		Assets: application.AssetOptions{
			Handler: assetHandler,
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:          "Device Mockup",
		Width:          500,
		Height:         400,
		EnableFileDrop: true,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
		},
		URL: "/",
	}).OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {

		ctx := event.Context()
		files := ctx.DroppedFiles()
		target := ctx.DropTargetDetails()

		if len(files) > 0 {

			emitKey := "image-file-dropped"
			emitData := map[string]any{
				"elementID": target.ElementID,
				"path":      files[0],
			}

			application.Get().Event.Emit(emitKey, emitData)
		}
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
