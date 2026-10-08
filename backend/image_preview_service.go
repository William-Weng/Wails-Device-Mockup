package backend

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"wails-device-mockup/backend/tools"

	"sync"
)

type ImagePreviewService struct {
	mutex          sync.RWMutex
	screenshotPath string
	framePath      string
}

// 設定目前要在前端預覽的外框圖片，回傳的 URL 可直接指定給 <img src>
//
// 參數:
//   - path: 圖片檔案路徑（相對或絕對路徑）
//
// 回傳:
//   - url: 前端可用的虛擬 URL（例如 "/media/image/frame.png?v=123"）
//   - error: 若路徑無效或檔案不存在則回傳錯誤
func (service *ImagePreviewService) SetFramePath(path string) (string, error) {
	prefix := tools.MediaPrefix["image-frame"]
	return service.setImagePath(path, prefix, tools.ImageTypeFrame)
}

// 是 ImagePreviewService 的內部共用函式，負責驗證、保存指定類型的圖片路徑，並回傳前端可用的虛擬媒體 URL
//
// 參數:
//   - path: 圖片檔案路徑，可為相對或絕對路徑
//   - prefix: 虛擬 URL 前綴，通常取自 mediaPrefix
//   - imageType: 圖片用途類型，例如 tools.ImageTypeScreenshot 或 tools.ImageTypeFrame
//
// 回傳:
//   - url: 前端可指定給 <img src> 的虛擬 URL
//   - error: 路徑無效、檔案不存在、格式不支援或圖片類型不支援時的錯誤
func (service *ImagePreviewService) setImagePath(path string, prefix string, imageType tools.ImageType) (string, error) {

	absolutePath, info, err := tools.ValidatePath(path, tools.ImageMimeMap)
	if err != nil {
		return "", err
	}

	service.mutex.Lock()
	defer service.mutex.Unlock()

	switch imageType {
	case tools.ImageTypeFrame:
		service.framePath = absolutePath
	default:
		return "", fmt.Errorf("不支援的圖片類型: %v", imageType)
	}

	return tools.MediaURL(prefix, absolutePath, info), nil
}

// serveCurrentImage 處理目前指定類型的圖片預覽 HTTP 請求；此方法僅供 Go 端的 Asset Handler 使用，不應暴露為 Wails 可由前端呼叫的公開方法
//
// 參數:
//   - writer: HTTP 回應寫入器，用於寫入 HTTP status、header 與圖片內容
//   - request: WebView 對虛擬媒體 URL 發出的 HTTP 請求
//   - imageType: 要回傳的圖片類型；可為 tools.ImageTypeScreenshot 或 tools.ImageTypeFrame
func (service *ImagePreviewService) serveCurrentImage(writer http.ResponseWriter, request *http.Request, imageType tools.ImageType) {

	service.mutex.RLock()

	var path string
	var validType bool

	switch imageType {
	case tools.ImageTypeFrame:
		path = service.framePath
		validType = true
	}

	service.mutex.RUnlock()

	if !validType {
		http.Error(writer, "不支援的圖片類型", http.StatusBadRequest)
		return
	}

	file, err := tools.OpenFile(path)
	if err != nil {
		http.NotFound(writer, request)
		return
	}
	defer file.Close()

	contentType, ok := tools.ImageContentType(strings.ToLower(filepath.Ext(path)))
	if !ok {
		http.Error(writer, "不支援的圖片格式", http.StatusUnsupportedMediaType)
		return
	}

	tools.ServeFile(writer, request, file, contentType)
}
