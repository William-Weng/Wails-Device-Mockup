package utility

import (
	"os/exec"
	"wails-device-mockup/backend/utility/platform"
)

// 設定子程序的執行屬性，主要用於在 Windows 上隱藏子程序的控制台視窗
// 參數 cmd:
//   - 要執行的 *exec.Cmd 物件
func HideWindow(cmd *exec.Cmd) {
	platform.HideWindow(cmd)
}
