package backend

import (
	"context"
	"time"
	"wails-device-mockup/backend/tools"
)

type VideoComposeService struct{}

// Compose 供 Wails 前端呼叫。
// Wails 會將 Go error 傳回前端 Promise rejection。
func (service *VideoComposeService) Compose(inputVideo string, inputVideoWidth int, inputVideoHeight int, frameImage string, frameImageWidth int, frameImageHeight int, outputVideo string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	return tools.ComposeVideoWithFrame(ctx, inputVideo, inputVideoWidth, inputVideoHeight, frameImage, frameImageWidth, frameImageHeight, outputVideo)
}
