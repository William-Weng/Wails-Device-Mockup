//go:build !windows

//go:dev !windows

package platform

import (
	"fmt"
	"os/exec"
	"syscall"
)

// 在非 Windows 平台不修改子程序設定
func HideWindow(cmd *exec.Cmd) {}

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

func ConfigureProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
}

func KillProcessTree(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return nil
	}

	// Negative PID kills the whole process group.
	return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
}
