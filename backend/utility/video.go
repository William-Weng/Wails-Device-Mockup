package utility

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"wails-device-mockup/backend/utility/constant"
)

// 使用 ffprobe 取得影片第一條視訊流的基本資訊
//
// 它只解析影片 metadata，不會解碼完整影片，因此適合在執行 FFmpeg 合成、轉檔或加外框前，先確認影片規格
//
// 參數：
//   - ctx：用來控制 ffprobe 的生命週期；例如呼叫端可設定 timeout，當 context 被取消或逾時時，ffprobe 子程序也會被終止
//   - path：影片檔案的完整路徑或相對路徑，例如 "input.mp4"
//
// 回傳值：
//   - codec：視訊編碼名稱，例如 "h264"、"hevc"、"vp9"、"av1"
//   - fps：FFmpeg 回傳的影格率字串，例如 "30/1"、"30000/1001"；可直接用於 FFmpeg filter 的 r 參數，例如 color=...:r=30/1
//   - width：影片寬度，單位為像素
//   - height：影片高度，單位為像素
//   - err：ffprobe 執行失敗、JSON 解析失敗，或找不到視訊流時的錯誤
func VideoInformation(ctx context.Context, path string) (codec, fps string, width, height int, err error) {

	ffprobePath, err := exec.LookPath("ffprobe")
	if err != nil {
		ffprobePath = "/opt/homebrew/bin/ffprobe"
	}

	cmd := exec.CommandContext(ctx, ffprobePath,
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=codec_name,width,height,r_frame_rate",
		"-of", "json",
		path,
	)

	stdout, err := cmd.Output()
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("ffprobe failed: %w", err)
	}

	var result constant.VideoInfo
	if err := json.Unmarshal(stdout, &result); err != nil {
		return "", "", 0, 0, fmt.Errorf("parse ffprobe output: %w", err)
	}

	if len(result.Streams) == 0 {
		return "", "", 0, 0, fmt.Errorf("no video stream found: %s", path)
	}

	stream := result.Streams[0]
	return stream.CodecName, stream.FrameRate, stream.Width, stream.Height, nil
}
