package backend

import (
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"wails-device-mockup/backend/tools"
)

type VideoPreviewService struct {
	mutex     sync.RWMutex
	videoPath string
}

// 設定目前允許 WebView 播放的影片；回傳固定的虛擬 URL；前端不可自行指定本機檔案路徑
//
// 參數:
//   - path: 影片檔案路徑（相對或絕對路徑）
//
// 回傳:
//   - url: 前端可用的虛擬 URL（例如 "/media/video/input.mp4?v=123"）
//   - error: 若路徑無效或檔案不存在則回傳錯誤
func (service *VideoPreviewService) SetVideoPath(path string) (string, error) {

	absolutePath, info, err := tools.ValidatePath(path, tools.VideoMimeMap)
	if err != nil {
		return "", err
	}

	prefix := tools.MediaPrefix["video-input"]

	service.mutex.Lock()
	service.videoPath = absolutePath
	service.mutex.Unlock()

	return tools.MediaURL(prefix, absolutePath, info), nil
}

// 處理影片播放請求
//
// 參數:
//   - writer: HTTP 回應寫入器
//   - request: HTTP 請求物件
//   - service: 影片預覽服務
func (service *VideoPreviewService) serveCurrentVideo(writer http.ResponseWriter, request *http.Request) {

	service.mutex.RLock()
	path := service.videoPath
	service.mutex.RUnlock()

	file, err := tools.OpenFile(path)
	if err != nil {
		http.NotFound(writer, request)
		return
	}
	defer file.Close()

	contentType, ok := tools.VideoContentType(strings.ToLower(filepath.Ext(path)))

	if !ok {
		http.Error(writer, "不支援的影片格式", http.StatusUnsupportedMediaType)
		return
	}

	tools.ServeFile(writer, request, file, contentType)
}
