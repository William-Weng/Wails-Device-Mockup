package backend

import (
	"context"
	"time"
	"wails-device-mockup/backend/tools"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type VideoComposeService struct {
	App *application.App
}

// Compose 供 Wails 前端呼叫。
// Wails 會將 Go error 傳回前端 Promise rejection。
func (service *VideoComposeService) Compose(inputVideo string, inputVideoWidth int, inputVideoHeight int, frameImage string, frameImageWidth int, frameImageHeight int, outputVideo string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	return tools.ComposeVideoWithFrame(service.App, ctx, inputVideo, inputVideoWidth, inputVideoHeight, frameImage, frameImageWidth, frameImageHeight, outputVideo)
}
