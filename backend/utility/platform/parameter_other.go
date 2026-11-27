//go:build !darwin

//go:dev !darwin

package platform

// CombineParameter 組合「影片 + 靜態圖 + 遮罩」合成輸出 MP4 的 FFmpeg 參數。
//
// 參數：
//   - inputVideo：原始輸入影片路徑
//   - frameImage：要合成到影片上的靜態圖片路徑
//   - maskImage：遮罩圖片路徑，供 filterComplex 使用
//   - filterComplex：FFmpeg 的 filter_complex 字串，負責實際合成邏輯
//   - outputVideo：輸出 MP4 檔案路徑
//
// 回傳：
//   - 可直接傳給 exec.CommandContext(ffmpegPath, args...) 的參數切片
func CombineParameter(inputVideo string, frameImage string, maskImage string, filterComplex string, outputVideo string) []string {

	return []string{
		"-y",
		"-hide_banner",

		"-i", inputVideo,

		"-loop", "1",
		"-i", frameImage,

		"-loop", "1",
		"-i", maskImage,

		"-filter_complex", filterComplex,

		"-map", "[out]",
		"-map", "0:a?",

		// 一般 CPU H.264 編碼
		"-c:v", "libx264",
		"-crf", "23",
		"-preset", "medium",
		"-pix_fmt", "yuv420p",

		"-c:a", "aac",

		// 讓中斷後的 MP4 仍可播放
		"-f", "mp4",
		"-movflags", "+frag_keyframe+empty_moov+default_base_moof",

		"-shortest",

		outputVideo,
	}
}
