package constant

// 支援的圖片副檔名與 MIME type 對應表
var ImageMimeMap = map[string]string{
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".webp": "image/webp",
	".gif":  "image/gif",
}

// 支援的影片副檔名與 MIME type 對應表
var VideoMimeMap = map[string]string{
	".mp4":  "video/mp4",
	".mov":  "video/quicktime",
	".m4v":  "video/x-m4v",
	".webm": "video/webm",
}

// 定義前端使用的虛擬路徑前綴，避免直接暴露本機檔案路徑
var MediaPrefix = map[string]string{
	"video-input": "/media/video/input",
	"image-frame": "/media/image/frame",
}
