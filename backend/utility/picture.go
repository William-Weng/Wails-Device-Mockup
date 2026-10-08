package utility

import (
	"image"
	_ "image/png"

	_ "golang.org/x/image/webp"

	"os"
)

// 會讀取指定路徑的圖片檔案，並在不完整解碼圖片像素的情況下，取得圖片的基本資訊
//
// 回傳值：
//   - image.Config：包含圖片寬、高與 ColorModel（色彩模型）
//   - string：Go 辨識出的圖片格式，例如 "png"、"webp"
//   - error：開檔失敗、格式不支援或圖片標頭解析失敗時的錯誤
func ImageInformation(path string) (image.Config, string, error) {

	file, err := os.Open(path)
	if err != nil {
		return image.Config{}, "", err
	}
	defer file.Close()

	config, format, err := image.DecodeConfig(file)
	if err != nil {
		return image.Config{}, "", err
	}

	return config, format, nil
}
