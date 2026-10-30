//go:build windows

//go:dev windows

package platform

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"syscall"
)

// 設定子程序的執行屬性，主要用於在 Windows 上隱藏子程序的控制台視窗
// 參數 cmd:
//   - 要執行的 *exec.Cmd 物件
func HideWindow(cmd *exec.Cmd) {

	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
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

// 設定子程序的行程群組
func ConfigureProcessGroup(cmd *exec.Cmd) {}

// 強制終止 cmd 及其所有子程序
func KillProcessTree(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return nil
	}

	killCmd := exec.Command("taskkill", "/PID", itoa(cmd.Process.Pid), "/T", "/F")
	HideWindow(killCmd)

	return killCmd.Run()
}

// itoa 將整數轉成字串
func itoa(n int) string {
	return strconv.Itoa(n)
}
