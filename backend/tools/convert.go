package tools

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"wails-device-mockup/backend/utility"
)

// ComposeVideoWithFrame 將 inputVideo 疊進 frameImage 的透明螢幕區域，
// 產生指定尺寸的 MP4 輸出檔。
//
// 目前參數依照指定的 iPhone 外框設定：
// - 外框畫布：2760 x 1350
// - 螢幕位置：(200, 85)
// - 螢幕尺寸：2320 x 1180
// - 輸出尺寸：1280 x 720
func ComposeVideoWithFrame(ctx context.Context, inputVideo string, inputVideoWidth int, inputVideoHeight int, frameImage string, frameImageWidth int, frameImageHeight int, outputVideo string) error {

	// 檢查系統是否可找到 ffmpeg。
	// Wails 開發環境可透過 PATH 找到；正式包版時，
	// 視需求可改成你 app bundle 中 ffmpeg 的絕對路徑。
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		ffmpegPath = "/opt/homebrew/bin/ffmpeg"
	}

	// 確保輸出資料夾存在。
	if err := os.MkdirAll(filepath.Dir(outputVideo), 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	filterComplex, err := buildFilterComplex(ctx, inputVideo, inputVideoWidth, inputVideoHeight, frameImage, frameImageWidth, frameImageHeight)
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(
		ctx,
		ffmpegPath,
		"-y", // 若輸出檔已存在，直接覆寫
		"-i", inputVideo,
		"-loop", "1",
		"-i", frameImage,
		"-filter_complex", filterComplex,
		"-map", "[out]",
		"-map", "0:a?", // 有音訊時才映射；沒有音訊不報錯
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-c:a", "aac",
		"-shortest",
		outputVideo,
	)

	// ffmpeg 將大部分資訊與錯誤寫至 stderr；
	// CombinedOutput 可取得 stdout + stderr，便於在 Wails 前端顯示原因。
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"ffmpeg compose failed: %w\n%s",
			err,
			string(output),
		)
	}

	return nil
}

func buildFilterComplex(ctx context.Context, inputVideo string, inputVideoWidth int, inputVideoHeight int, frameImage string, frameImageWidth int, frameImageHeight int) (string, error) {

	picture, _, err := utility.ImageInformation(frameImage)

	if err != nil {
		return "", err
	}

	_, fps, _, videoHeight, err := utility.VideoInformation(ctx, inputVideo)

	if err != nil {
		return "", err
	}

	// 將前端預覽中量到的 video 尺寸，換算回原始 PNG 的座標系
	//
	// 例如：
	// 2760 / 1000 * 770 = 2125.2 → 2125
	// 1350 / 492 * 433 = 1188.1 → 1188
	scaleVideoWidth := picture.Width * inputVideoWidth / frameImageWidth
	scaleVideoHeight := picture.Height * inputVideoHeight / frameImageHeight

	// 將換算後的影片區域置中於原始外框 PNG
	//
	// 例如：
	// (2760 - 2125) / 2 = 317.5 → 317
	// (1350 - 1188) / 2 = 81
	overlayX := (picture.Width - scaleVideoWidth) / 2
	overlayY := (picture.Height - scaleVideoHeight) / 2

	finalVideoWidth := videoHeight * picture.Width / picture.Height

	// 使用 fmt.Sprintf 把所有計算結果放進 FFmpeg filter graph
	filterComplex := fmt.Sprintf(
		"color=c=#1c1c1c:s=%dx%d:r=%s[bg];"+
			"[0:v]scale=%d:%d:force_original_aspect_ratio=increase,"+
			"crop=%d:%d,"+
			"setsar=1[screen];"+
			"[bg][screen]overlay=%d:%d[base];"+
			"[base][1:v]overlay=0:0:shortest=1,"+
			"scale=%d:%d,"+
			"setsar=1[out]",

		picture.Width,
		picture.Height,
		fps,

		// 原影片縮放與裁切後的螢幕區域
		scaleVideoWidth,
		scaleVideoHeight,
		scaleVideoWidth,
		scaleVideoHeight,

		// 影片放到外框畫布的位置
		overlayX,
		overlayY,

		finalVideoWidth,
		videoHeight,
	)

	return filterComplex, nil
}
