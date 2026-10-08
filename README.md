# [Wails-Device-Mockup](https://v3.wails.io/features/drag-and-drop/files/)

![Go](https://img.shields.io/badge/Go-1.27.1-00ADD8?logo=go&logoColor=white)
![Node.js](https://img.shields.io/badge/Node.js-20.19.2-339933?logo=nodedotjs&logoColor=white)
![Wails](https://img.shields.io/badge/Wails-v3.0.0--beta.26-DF0000?logo=wails&logoColor=white)
![LICENSE](https://img.shields.io/github/license/William-Weng/Wails-Device-Mockup?style=flat&label=LICENSE&color=yellow)
![Tag](https://img.shields.io/github/v/tag/William-Weng/Wails-Device-Mockup?style=flat&label=Tag)
![Stars](https://img.shields.io/github/stars/William-Weng/Wails-Device-Mockup?style=flat&label=Stars)

使用 Wails、Go、Svelte 與 TypeScript 建立的裝置 Mockup 工具，可將影片合成至裝置外框中，快速產生適合展示、簡報與產品宣傳使用的裝置預覽影片。

## [操作介面](https://v3.wails.io/zh-tw/reference/cli/)

https://github.com/user-attachments/assets/12e90b7a-07fb-4eb2-ab8f-faddcd435319

## [功能特色](https://v3.wails.io/zh-tw/concepts/build-system/)

- 將影片合成至裝置 Mockup 畫面中。
- 支援影片預覽與輸出。
- 使用 FFmpeg 處理影片合成。
- 提供簡單直覺的圖形化操作介面。
- 支援 macOS App 使用情境。
- 前端使用 Svelte 與 TypeScript。
- 後端使用 Go 與 Wails。
- 可擴充不同裝置尺寸、外框與影片比例。
- 透過本機處理影片，不需要上傳至遠端伺服器。

## 安裝 FFmpeg

### macOS

使用 Homebrew：

```bash
brew install ffmpeg
```

確認安裝成功：

```bash
ffmpeg -version
```

確認 FFmpeg 的位置：

```bash
which ffmpeg
```

常見位置如下：

```text
/opt/homebrew/bin/ffmpeg
```

或：

```text
/usr/local/bin/ffmpeg
```

### Windows

可以使用 `winget` 安裝：

```powershell
winget install Gyan.FFmpeg
```

安裝完成後，重新開啟終端機並確認：

```powershell
ffmpeg -version
```

### Linux

以 Debian 或 Ubuntu 為例：

```bash
sudo apt update
sudo apt install ffmpeg
```

確認安裝成功：

```bash
ffmpeg -version
```

## [建置指令整理](https://v3.wails.io/zh-tw/guides/build/building/)

| 目的 | 指令 |
| --- | --- |
| 建立新專案 | `wails3 init -n <專案名稱> -t <前端框架>` |
| 產生 bindings | `wails3 generate bindings` |
| 更新 Windows 與 macOS 專用圖示檔案 | `wails3 generate icons -input build/appicon.png -windowsfilename build/windows/icon.ico -macfilename build/darwin/icons.icns` |
| 更新建置資源 | `wails3 update build-assets -config build/config.yml -dir build` |
| 拉取（下載）用於跨平台交叉編譯的 Docker 映像檔 | `wails3 task setup:docker` |
| 建置 Windows x64 | `wails3 build GOOS=windows GOARCH=amd64` |
| 打包 macOS arm64 | `wails3 package GOOS=darwin GOARCH=arm64` |
| 打包 Linux arm64 | `wails3 build GOOS=linux GOARCH=arm64` |
