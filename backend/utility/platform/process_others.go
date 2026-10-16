//go:build !windows

//go:dev !windows

package platform

import (
	"fmt"
	"os/exec"
)

// 在非 Windows 平台不修改子程序設定
func HideWindow(cmd *exec.Cmd) *exec.Cmd {
	return cmd
}

// 回傳 FFmpeg 執行檔的標準名稱
func FFMpegDefaultName() string {
	return "ffmpeg"
}

// 回傳 FFprobe 執行檔的標準名稱
func FFProbeDefaultName() string {
	return "ffprobe"
}

// 回傳 macOS 環境下 FFmpeg 執行檔的預設路徑
func FFMpegDefaultPath() string {
	return fmt.Sprintf("/opt/homebrew/bin/%s", FFMpegDefaultName())
}

// 回傳 macOS 環境下 FFprobe 執行檔的預設路徑
func FFProbeDefaultPath() string {
	return fmt.Sprintf("/opt/homebrew/bin/%s", FFProbeDefaultName())
}
