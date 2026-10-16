package backend

import "wails-device-mockup/backend/utility"

type FileManagerService struct{}

// 將給定的檔案路徑拆解為三個部分：目錄路徑、檔名（不含副檔名）、副檔名
//
// 參數：
//   - path：檔案的完整路徑，例如 "/Users/ios/Desktop/古本屋.jpg" 或 "data/config.yaml"
//
// 回傳值：
//   - "/Users/ios/Desktop/古本屋.jpg" → ("/Users/ios/Desktop", "古本屋", ".jpg")
//   - "/Users/ios/Desktop/.gitignore" → ("/Users/ios/Desktop", ".gitignore", "")
func (service *FileManagerService) ParseFilePath(path string) (rootdir string, filename string, extension string) {
	return utility.ParseFilePath(path)
}
