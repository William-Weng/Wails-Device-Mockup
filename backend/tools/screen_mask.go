package tools

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
	"path/filepath"
	"wails-device-mockup/backend/utility/constant"
)

// CreateScreenMaskOnTemp 讀取手機外框圖，偵測螢幕透明區域，並將灰階 mask 固定輸出到系統暫存資料夾的 screen-mask.png
func CreateScreenMaskOnTemp(frameImagePath string) (string, error) {

	input, err := os.Open(frameImagePath)
	if err != nil {
		return "", fmt.Errorf("open frame image: %w", err)
	}
	defer input.Close()

	frame, _, err := image.Decode(input)
	if err != nil {
		return "", fmt.Errorf("decode frame image: %w", err)
	}

	region, err := DetectScreenRegion(frame)
	if err != nil {
		return "", fmt.Errorf("detect screen region: %w", err)
	}

	maskPath := filepath.Join(os.TempDir(), "screen-mask.png")

	if err := writeScreenMask(frame, region, maskPath); err != nil {
		return "", err
	}

	return maskPath, nil
}

// 將 DetectScreenRegion 找到的螢幕透明區域輸出為灰階 PNG mask
//
// 輸出圖片與 frame 擁有相同的尺寸：
//   - 白色（255）：允許影片顯示的螢幕像素
//   - 黑色（0）：遮蔽影片的區域，例如外框與圓角外側
func writeScreenMask(frame image.Image, region constant.Region, outputPath string) error {

	bounds := frame.Bounds()
	width := bounds.Dx()

	if width <= 0 || bounds.Dy() <= 0 {
		return fmt.Errorf("圖片尺寸無效")
	}

	mask := image.NewGray(bounds)

	for _, index := range region.Pixels {
		x := index % width
		y := index / width

		// Pixels 使用相對於 frame.Bounds().Min 的 index，所以要補回 bounds.Min
		mask.SetGray(bounds.Min.X+x, bounds.Min.Y+y, color.Gray{Y: 255})
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return fmt.Errorf("create mask directory: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create mask file: %w", err)
	}
	defer file.Close()

	if err := png.Encode(file, mask); err != nil {
		return fmt.Errorf("encode mask PNG: %w", err)
	}

	return nil
}
