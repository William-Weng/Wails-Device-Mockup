package utility

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"strconv"
	"strings"

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

// HexToRGBA converts "#ff0000", "#f00", "#ff000080" to application.RGBA.
func HexToRGBA(hex string) (color.RGBA, error) {

	hex = strings.TrimPrefix(strings.TrimSpace(hex), "#")

	switch len(hex) {
	case 3:
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]}) + "ff"
	case 6:
		hex += "ff"
	case 8:
	default:
		return color.RGBA{}, fmt.Errorf("invalid hex color: %q", hex)
	}

	value, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("invalid hex color: %q", hex)
	}

	return color.RGBA{R: uint8(value >> 24), G: uint8(value >> 16), B: uint8(value >> 8), A: uint8(value)}, nil
}
