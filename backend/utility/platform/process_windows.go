//go:build windows

//go:dev windows

package platform

import (
	"fmt"
	"os/exec"
	"runtime"
	"syscall"
)

// 設定子程序的執行屬性，主要用於在 Windows 上隱藏子程序的控制台視窗
// 參數 cmd:
//   - 要執行的 *exec.Cmd 物件
//
// 回傳值:
//   - 設定後的 *exec.Cmd 物件（通常就是傳入的同一個指標）
func HideWindow(cmd *exec.Cmd) *exec.Cmd {

	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}

	return cmd
}

// 回傳 FFmpeg 執行檔的標準名稱
func FFMpegDefaultName() string {
	return "ffmpeg.exe"
}

// 回傳 FFprobe 執行檔的標準名稱
func FFProbeDefaultName() string {
	return "ffprobe.exe"
}

// 回傳 Windows 系統中 FFmpeg 執行檔的預設安裝路徑
func FFMpegDefaultPath() string {
	return fmt.Sprintf("C:/Program Files/ffmpeg/bin/%s", FFMpegDefaultName())
}

// 回傳 Windows 系統中 FFprobe 執行檔的預設安裝路驗
func FFProbeDefaultPath() string {
	return fmt.Sprintf("C:/Program Files/ffmpeg/bin/%s", FFProbeDefaultName())
}
