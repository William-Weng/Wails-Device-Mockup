package constant

import "image"

type VideoInfo struct {
	Streams []struct {
		CodecName string `json:"codec_name"`
		Width     int    `json:"width"`
		Height    int    `json:"height"`
		FrameRate string `json:"r_frame_rate"`
	} `json:"streams"`
}

// Region 表示一塊相連透明區域的搜尋結果。Rect 的座標相對於圖片左上角；Max 不包含在矩形內
type Region struct {
	Rect        image.Rectangle
	Area        int
	TouchesEdge bool
	Pixels      []int
}
