/**
 * 從 Wails「image-file-dropped」事件取得的圖片拖放資料
 * - slotId：圖片被拖放到的區域 ID
 * - path：圖片在本機檔案系統中的完整路徑
 */
type DroppedData = {
  elementID: string;
  path: string;
};
