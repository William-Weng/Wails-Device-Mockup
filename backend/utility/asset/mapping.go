package asset

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"wails-device-mockup/backend/utility/constant"
)

// 驗證給定的檔案路徑是否有效，並回傳絕對路徑與檔案資訊
//
// 參數:
//   - path: 圖片檔案路徑（相對或絕對路徑）
//   - supportMineMap: 支援的副檔名與 MIME type 對應表（例如 imageMimeMap、videoMimeMap）
//
// 回傳:
//   - absolutePath: 驗證通過的絕對路徑
//   - info: 檔案的 os.FileInfo 資訊
//   - error: 若路徑無效、檔案不存在、是資料夾或格式不支援則回傳錯誤
func ValidatePath(path string, supportMineMap map[string]string) (string, os.FileInfo, error) {

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "", nil, fmt.Errorf("無法取得檔案絕對路徑：%w", err)
	}

	info, err := os.Stat(absolutePath)
	if err != nil {
		return "", nil, fmt.Errorf("找不到檔案：%w", err)
	}

	if info.IsDir() {
		return "", nil, fmt.Errorf("請選擇檔案，不是資料夾")
	}

	extension := strings.ToLower(filepath.Ext(absolutePath))

	if _, ok := supportMineMap[extension]; !ok {
		return "", nil, fmt.Errorf("不支援的檔案格式：%s", filepath.Ext(absolutePath))
	}

	return absolutePath, info, nil
}

// 根據基礎路徑、檔案路徑和檔案資訊，組合出前端可用的虛擬 URL；URL 包含版本參數（?v=修改時間），避免瀏覽器快取舊檔案
//
// 參數:
//   - basePath: 虛擬路徑前綴（例如 "/media/image/screenshot"）
//   - filePath: 檔案的絕對路徑，用於取出副檔名
//   - info: 檔案的 os.FileInfo，用於取得修改時間
//
// 回傳:
//   - 組合後的 URL（例如 "/media/image/screenshot.png?v=1696234567890123456"）
func MediaURL(basePath string, filePath string, info os.FileInfo) string {
	extension := strings.ToLower(filepath.Ext(filePath))
	return fmt.Sprintf("%s%s?v=%d", basePath, extension, info.ModTime().UnixNano())
}

// 以指定的 Content-Type 提供檔案內容
//
// 參數:
//   - writer: HTTP 回應寫入器
//   - request: HTTP 請求物件
//   - file: 要提供的檔案
//   - contentType: 檔案的 MIME type
func ServeFile(writer http.ResponseWriter, request *http.Request, file *os.File, contentType string) {

	info, err := file.Stat()
	if err != nil {
		http.Error(writer, "無法讀取檔案資訊", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", contentType)
	writer.Header().Set("Cache-Control", "no-store")

	http.ServeContent(writer, request, info.Name(), info.ModTime(), file)
}

// 根據副檔名回傳對應的圖片 MIME type
//
// 參數:
//   - extension: 檔案副檔名（例如 ".png"、".jpg"）
//
// 回傳:
//   - contentType: 對應的 MIME type（例如 "image/png"）
//   - ok: 若副檔名不在支援列表中為 false
func ImageContentType(extension string) (string, bool) {
	contentType, ok := constant.ImageMimeMap[extension]
	return contentType, ok
}

// 根據副檔名回傳對應的影片 MIME type
//
// 參數:
//   - extension: 檔案副檔名（例如 ".mp4"、".webm"）
//
// 回傳:
//   - contentType: 對應的 MIME type（例如 "video/mp4"）
//   - ok: 若副檔名不在支援列表中為 false
func VideoContentType(extension string) (string, bool) {
	contentType, ok := constant.VideoMimeMap[extension]
	return contentType, ok
}
