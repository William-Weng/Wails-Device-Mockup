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
func CreateScreenMaskOnTemp(frameImagePath string, innerColor color.RGBA, outterColor color.RGBA) (string, error) {

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

	if err := writeScreenMaskWithBorder(frame, region, maskPath, innerColor, outterColor); err != nil {
		return "", err
	}

	return maskPath, nil
}

// 將 DetectScreenRegion 找到的螢幕透明區域輸出為灰階 PNG mask
//
// 輸出圖片與 frame 擁有相同的尺寸：
//   - 白色（255）：允許影片顯示的螢幕像素
//   - 黑色（0）：遮蔽影片的區域，例如外框與圓角外側
func writeScreenMaskWithBorder(frame image.Image, region constant.Region, outputPath string, innerColor color.RGBA, outterColor color.RGBA) error {

	bounds := frame.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width <= 0 || height <= 0 {
		return fmt.Errorf("圖片尺寸無效")
	}

	mask := image.NewRGBA(bounds)
	maxIndex := width * height
	pix := mask.Pix

	// 1. 初始化：將整張圖片塗成「不透明黑色」
	// Go 的 NewRGBA 預設是全透明，所以要把背景填滿黑色 (R:0, G:0, B:0, A:255)
	for i := 0; i < len(pix); i += 4 {
		pix[i+0] = outterColor.R
		pix[i+1] = outterColor.G
		pix[i+2] = outterColor.B
		pix[i+3] = outterColor.A
	}

	// 2. 建立快速查表，標記哪些位置屬於區域內
	inRegion := make([]bool, maxIndex)
	for _, index := range region.Pixels {
		if index >= 0 && index < maxIndex {
			inRegion[index] = true
		}
	}

	// 3. 處理 region 區域（挖空變透明 或 畫上外框）
	for _, index := range region.Pixels {
		if index < 0 || index >= maxIndex {
			continue
		}

		base := index * 4
		pix[base] = innerColor.R
		pix[base+1] = innerColor.G
		pix[base+2] = innerColor.B
		pix[base+3] = innerColor.A
	}

	// 建立目錄與檔案
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
