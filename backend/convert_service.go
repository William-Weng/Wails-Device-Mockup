package backend

import (
	"context"
	"errors"
	"image/color"
	"sync"
	"time"
	"wails-device-mockup/backend/tools"
	"wails-device-mockup/backend/utility"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// VideoComposeService 負責管理「合成影片」任務，並確保同一時間可以安全地啟動、取消、清理 FFmpeg 任務
type VideoComposeService struct {
	App           *application.App
	mutex         sync.Mutex         // 保護 cancel 與 currentTaskID
	cancel        context.CancelFunc // 是目前正在執行任務的 context.CancelFunc
	currentTaskID string             // 用來識別「目前儲存的 cancel」屬於哪一次 Compose 任務
}

// ServiceShutdown 是 Wails v3 的 service lifecycle hook
//
// 當 Wails App 關閉時會自動呼叫，例如使用者按下視窗 X、Cmd+Q 或系統登出；這裡取消目前任務的 context，讓 exec.CommandContext 終止還在執行的 FFmpeg
func (service *VideoComposeService) ServiceShutdown() error {
	return service.Cancel()
}

// Compose 供 Wails 前端呼叫，開始將 inputVideo 與 frameImage 合成為 outputVideo
func (service *VideoComposeService) Compose(inputVideo string, inputVideoWidth int, inputVideoHeight int, frameImage string, frameImageWidth int, frameImageHeight int, outputVideo string) error {

	service.mutex.Lock()
	if service.cancel != nil {
		service.cancel()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	taskID := time.Now().Format("20060102150405.000000000")

	service.cancel = cancel
	service.currentTaskID = taskID
	service.mutex.Unlock()

	defer func() {

		service.mutex.Lock()

		if service.currentTaskID == taskID {
			service.cancel = nil
			service.currentTaskID = ""
		}
		service.mutex.Unlock()

		cancel()
	}()

	innerColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	outterColor := color.RGBA{R: 0, G: 0, B: 0, A: 255}

	return tools.ComposeVideoWithFrame(service.App, ctx, innerColor, outterColor, inputVideo, inputVideoWidth, inputVideoHeight, frameImage, frameImageWidth, frameImageHeight, outputVideo)
}

// CreateScreenMaskOnTemp 讀取手機外框圖，偵測螢幕透明區域，並將灰階 mask 固定輸出到系統暫存資料夾的 screen-mask.png
func (service *VideoComposeService) CreateScreenMaskOnTemp(frameImagePath string, maskColor string) (string, error) {

	innerColor := color.RGBA{R: 0, G: 0, B: 0, A: 0}

	outterColor, err := utility.HexToRGBA(maskColor)
	if err != nil {
		outterColor = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	}

	return tools.CreateScreenMaskOnTemp(frameImagePath, innerColor, outterColor)
}

// Cancel 供 Wails 前端呼叫，主動停止目前正在執行的 FFmpeg 合成任務
//
// 它會取消目前任務的 context；若 tools.ComposeVideoWithFrame 使用 exec.CommandContext(ctx, ...)，FFmpeg 就會被 Go 自動終止
func (service *VideoComposeService) Cancel() error {

	service.mutex.Lock()
	defer service.mutex.Unlock()

	if service.cancel != nil {
		service.cancel()
		return nil
	}

	return errors.New("no running compose task")
}
