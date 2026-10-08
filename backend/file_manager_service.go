package backend

import "wails-device-mockup/backend/utility"

type FileManagerService struct{}

func (service *FileManagerService) ParseFilePath(path string) (rootdir string, filename string, extension string) {
	return utility.ParseFilePath(path)
}
