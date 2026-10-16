/**
 * 從 Wails「image-file-dropped」事件取得的圖片拖放資料
 * - slotId：圖片被拖放到的區域 ID
 * - path：圖片在本機檔案系統中的完整路徑
 */
type DroppedData = {
  elementID: string;
  path: string;
};

/**
 * 合成或處理進度的型別定義
 * - duration：影片總時長（秒）
 * - progress：當前整理的影片長度（秒）
 */
type CombineProgress = {
  duration: number;
  progress: number;
};

/**
 * 圖片框的詳細資訊
 * - width：圖片框的寬度（像素, pixel）
 * - height：圖片框的高度（像素, pixel
 * - extension：檔案的副檔名（例如：'jpg', 'png', 'mp4'）
 */
type FrameInformation = {
  width: number;
  height: number;
  extension: string;
};
