package tools

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"wails-device-mockup/backend/utility"
	"wails-device-mockup/backend/utility/platform"

	"github.com/wailsapp/wails/v3/pkg/application"
)

var timeRegexp = regexp.MustCompile(`time=(\d+:\d{2}:\d{2}(?:\.\d+)?)`)

// ComposeVideoWithFrame 將 inputVideo 疊進 frameImage 的透明螢幕區域，生成一張mask圖片，然後產生指定尺寸的 MP4 輸出檔
//
// 目前參數依照指定的 iPhone 外框設定：
// - 外框畫布：2760 x 1350
// - 螢幕位置：(200, 85)
// - 螢幕尺寸：2320 x 1180
// - 輸出尺寸：1280 x 720
func ComposeVideoWithFrame(app *application.App, ctx context.Context, inputVideo string, inputVideoWidth int, inputVideoHeight int, frameImage string, frameImageWidth int, frameImageHeight int, outputVideo string) error {

	videoDuration, err := detectMediaDuration(ctx, inputVideo)
	if err != nil {
		return err
	}

	ffmpegPath, err := exec.LookPath(platform.FFMpegDefaultName())
	if err != nil {
		ffmpegPath = platform.FFMpegDefaultPath()
	}

	maskImage, err := CreateScreenMaskOnTemp(frameImage)
	if err != nil {
		return err
	}

	// 確保輸出資料夾存在
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
		"-loop", "1",
		"-i", maskImage,
		"-filter_complex", filterComplex,
		"-map", "[out]",
		"-map", "0:a?", // 有音訊時才映射；沒有音訊不報錯
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-c:a", "aac",
		"-shortest",
		outputVideo,
	)

	cmd = platform.HideWindow(cmd)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return nil
		}
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}

	// if err := cmd.Run(); err != nil {
	// 	return cmd.Process.Kill()
	// }

	err = utility.RunCommandStream(
		cmd,
		func(source string, text string) {

			if source != "stderr" {
				return
			}

			time, err := parseTimeText(text)
			if err != nil {
				return
			}

			second, err := timeToSeconds(time)
			if err != nil {
				return
			}

			app.Event.Emit("combine-progress", map[string]float64{
				"duration": videoDuration,
				"progress": second,
			})
		},
	)

	if err != nil {
		return fmt.Errorf("ffmpeg compose failed: %w\n", err)
	}

	return nil
}

// 使用 ffprobe 取得媒體檔案的總時長
//   - context：Context
//   - inputPath：要讀取資訊的影音檔案路徑
//
// 回傳:
//   - duration：媒體總長度，單位為秒，例如 125.45 代表 2 分 5.45 秒。
//   - err：ffprobe 執行失敗、輸出無法轉為秒數，或取得非正時長時的錯誤。
func detectMediaDuration(context context.Context, inputPath string) (duration float64, err error) {

	ffprobePath, err := exec.LookPath(platform.FFProbeDefaultPath())
	if err != nil {
		ffprobePath = platform.FFMpegDefaultPath()
	}

	cmd := exec.CommandContext(
		context,
		ffprobePath,
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		inputPath,
	)

	cmd = platform.HideWindow(cmd)
	output, err := cmd.Output()

	if err != nil {
		return 0, fmt.Errorf("ffprobe 執行失敗：%w", err)
	}

	duration, err = strconv.ParseFloat(strings.TrimSpace(string(output)), 64)

	if err != nil {
		return 0, fmt.Errorf("無法解析影片長度：%w", err)
	}

	if duration <= 0 {
		return 0, errors.New("影片長度必須大於 0")
	}

	return duration, nil
}

// buildFilterComplex 建立 FFmpeg 的複雜濾鏡參數（-filter_complex）
// 參數說明：
//   - ctx: 上下文，用於控制影片資訊讀取的生命週期
//   - inputVideo: 原始輸入影片的路徑（對應 FFmpeg 的 [0:v]）
//   - inputVideoWidth/Height: 前端預覽畫面上，使用者調整後的影片寬高
//   - frameImage: 外框圖片路徑（對應 FFmpeg 的 [1:v]，如手機外框，通常包含一個主體框與一個透明鏤空區）
//   - frameImageWidth/Height: 前端預覽畫面上，外框圖片的寬高
func buildFilterComplex(ctx context.Context, inputVideo string, inputVideoWidth int, inputVideoHeight int, frameImage string, frameImageWidth int, frameImageHeight int) (string, error) {

	picture, _, err := utility.ImageInformation(frameImage)
	if err != nil {
		return "", err
	}

	_, fps, _, videoHeight, err := utility.VideoInformation(ctx, inputVideo)
	if err != nil {
		return "", err
	}

	// 將前端預覽中的影片尺寸，換算回原始 frame PNG 座標
	scaleVideoWidth := picture.Width * inputVideoWidth / frameImageWidth
	scaleVideoHeight := picture.Height * inputVideoHeight / frameImageHeight

	// 將影片矩形置中在完整 frame 畫布
	overlayX := (picture.Width - scaleVideoWidth) / 2
	overlayY := (picture.Height - scaleVideoHeight) / 2

	// 保持既有輸出比例：輸出高度等於原影片高度
	// finalVideoWidth := videoHeight * picture.Width / picture.Height
	finalVideoWidth := -2

	filterComplex := fmt.Sprintf(
		// 先建立最底層背景
		"color=c=#000000:s=%[1]dx%[2]d:r=%[3]s,format=rgba[bg];"+

			// 將影片縮放成 cover，並裁成前端預覽實際顯示大小
			"[0:v]scale=%[4]d:%[5]d:force_original_aspect_ratio=increase,"+
			"crop=%[4]d:%[5]d,"+
			"setsar=1,format=rgba[screen];"+

			// 建立完整尺寸的透明影片圖層，並將影片放到中間
			"color=c=black@0.0:s=%[1]dx%[2]d:r=%[3]s,format=rgba[canvas];"+
			"[canvas][screen]overlay=%[6]d:%[7]d:format=auto[videoPlaced];"+

			// mask：白色保留畫面、黑色設為透明；將 mask 灰階亮度帶入 videoPlaced 的 alpha channel
			"[2:v]format=gray[screenMask];"+
			"[videoPlaced][screenMask]alphamerge[maskedVideo];"+

			// 疊圖順序：背景 → 已 mask 的影片 → 手機外框
			"[bg][maskedVideo]overlay=0:0:format=auto[withVideo];"+
			"[withVideo][1:v]overlay=0:0:format=auto,"+

			// 輸出尺寸與相容格式
			"scale=%[8]d:%[9]d:flags=lanczos,"+
			"setsar=1,format=yuv420p[out]",

		picture.Width,    // 1
		picture.Height,   // 2
		fps,              // 3
		scaleVideoWidth,  // 4
		scaleVideoHeight, // 5
		overlayX,         // 6
		overlayY,         // 7
		finalVideoWidth,  // 8
		videoHeight,      // 9
	)

	return filterComplex, nil
}

// 從 FFmpeg 的進度輸出文字中，擷取 time= 後面的時間字串，例如 "00:00:37.59"
//
// 參數 line：FFmpeg 輸出的一行文字，例如：frame= 1126 fps= 86 ... time=00:00:37.59 bitrate=...
//
// 回傳：
//   - 成功：時間字串，例如 "00:00:37.59"
//   - 失敗：空字串與錯誤
func parseTimeText(line string) (string, error) {

	match := timeRegexp.FindStringSubmatch(line)
	if len(match) >= 2 {
		return match[1], nil
	}

	return "", fmt.Errorf("not match")
}

// 將 FFmpeg 的時間字串轉為總秒數
//
// 支援格式："HH:MM:SS" 或 "HH:MM:SS.xxx"；例如："00:00:37.59" 會回傳 37.59。
//
// 參數 time：FFmpeg 時間字串
//
// 回傳：
//   - 成功：總秒數
//   - 失敗：0 與錯誤
func timeToSeconds(timeText string) (float64, error) {

	parts := strings.Split(timeText, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid ffmpeg time: %q", timeText)
	}

	h, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}

	m, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}

	s, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return 0, err
	}

	return float64(h*3600+m*60) + s, nil
}
