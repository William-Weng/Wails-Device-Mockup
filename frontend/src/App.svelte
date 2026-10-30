<script lang="ts">
    import { onMount } from "svelte";
    import { Events } from "@wailsio/runtime";
    import { dialog } from "./utility/dialog";

    import { SetFramePath, ImageInformation } from "../bindings/wails-device-mockup/backend/imagepreviewservice";
    import { SetVideoPath } from "../bindings/wails-device-mockup/backend/videopreviewservice";
    import { Compose, Cancel } from "../bindings/wails-device-mockup/backend/videocomposeservice";
    import { ParseFilePath } from "../bindings/wails-device-mockup/backend/filemanagerservice";

    const supportedImageTypes = new Set<ImageType>([".png", ".webp"]);
    const supportedVideoTypes = new Set<VideoType>([".mp4", ".mov", ".ts", ".mkv", ".wmv"]);

    let framePath = "";
    let videoPath = "";
    let fileInfo: string[] = [];

    let isPortrait = $state(false);
    let isConverting = $state(false);
    let isCancelling = $state(false);
    let showControls = $state(true);
    let covertingText = $state("0%");
    let aspectRatio = $state(2760/1350);
    let frameSrc = $state("");
    let videoSrc = $state("");

    let frameWidth = $state(0);
    let frameHeight = $state(0);

    let videoWidthScale = $state(75);
    let videoHeight = $state(0);
    let videoWidth = $derived(Math.round((frameWidth * videoWidthScale) / 100));

    onMount(() => {
        const unsubscribeDrop = Events.On("image-file-dropped", (event) => {
            void imageFileDroppedAction(event.data as DroppedData);
        });

        const unsubscribeCombine = Events.On("combine-progress", (event) => {
          void parseCombineProgress(event.data as CombineProgress);
        });

        return () => {
          unsubscribeDrop();
          unsubscribeCombine();
        }
    });

    /**
     * 處理檔案拖放後的非同步核心邏輯
     */
    async function imageFileDroppedAction(data: DroppedData): Promise<void> {

        if (isConverting) { return; }
        fileInfo = await ParseFilePath(data.path);
        const ext = fileInfo[2];

        if (isImageType(ext)) {
            framePath = data.path;
            frameSrc = await SetFramePath(framePath);
            try {
              const info = await ImageInformation(framePath) as FrameInformation
              isPortrait = (info.width < info.height)
              aspectRatio = info.width / info.height
            } catch (err) {
              aspectRatio = 2760 / 1350
            }

            return;
        }

        if (isVideoType(ext)) {
            videoPath = data.path;
            videoSrc = await SetVideoPath(videoPath);
            return;
        }

        await dialog("warning", "錯誤", `不支援 ${ext} 格式`);
    }

    // 觸發合成影片的按鈕點擊事件
    async function composeVideo() {

        if (frameSrc == "") { return; }
        if (videoSrc == "") { return; }

        const timestamp = Date.now();
        const outputPath = `${fileInfo[0]}/${fileInfo[1]}-output-${timestamp}.mp4`;

        isConverting = true;

        try {
            await Compose(videoPath, videoWidth, videoHeight, framePath, frameWidth, frameHeight, outputPath);
            reset()
            dialog("info", "影片合成完成", `${outputPath}`);
        } catch (err) {
            reset()
            dialog("info", "影片合成失敗", `${err}`);
        }
    }

    // 觸發取消影片合成的按鈕點擊事件
    async function cancelComposeVideo() {

      if (!isConverting || isCancelling) { return; }
      isCancelling = true;

      try {
          await Cancel();
          dialog("info", "取消影片合成", "目前沒有正在進行的合成任務");
      } catch (err) {
          dialog("warning", "影片取消合成失敗", `${err}`);
      } finally {
        isConverting = false;
        isCancelling = false;
      }
    }

    // 顯示轉換進度
    async function parseCombineProgress(data: CombineProgress): Promise<void> {
      const percent = ((data.progress / data.duration) * 100).toFixed(2);
      covertingText = `${percent}%`;
    }

    // 檢查是否為支援的圖片副檔名
    function isImageType(extension: string): extension is VideoType {
        return supportedImageTypes.has(extension.toLocaleLowerCase() as ImageType);
    }

    // 檢查是否為支援的影片副檔名
    function isVideoType(extension: string): extension is VideoType {
        return supportedVideoTypes.has(extension.toLocaleLowerCase() as VideoType);
    }

    /**
     * 監聽全域鍵盤按下事件；當使用者按下 F 鍵時，切換 slider-field 與 submit-container 的顯示／隱藏
     * @param {KeyboardEvent} event - 瀏覽器傳入的鍵盤事件
     */
    function handleKeydown(event: any) {

        // 避免使用者正在輸入文字時誤觸
        const target = event.target;
        const isTyping = target instanceof HTMLInputElement || target instanceof HTMLTextAreaElement || target.isContentEditable;

        if (isTyping) return;

        // event.key 在按下 Shift + F 時會是 'F'
        if (event.key.toLowerCase() === 'f') {
          event.preventDefault();
          showControls = !showControls;
        }
    }

    // 回復原設定值
    function reset() {
      isConverting = false;
      covertingText = "0%"
    }
</script>

<!-- 將 keydown 事件監聽器掛到 window；因此即使焦點不在特定按鈕或輸入框上，按下 F 也能觸發；注意：<svelte:window> 必須放在元件最上層，不能放在 {#if} 或其他元素內。 -->
<svelte:window onkeydown={handleKeydown} />

<main class:preview-only={!showControls}>
    <section class="video-frame" class:is-portrait={isPortrait} style:--aspect-ratio={aspectRatio}  data-file-drop-target>
        <!-- 影片層 -->
        <div class="layer layer-video">
            <video controls src={videoSrc} style="width: {videoWidth}px;" bind:clientHeight={videoHeight}>
                <track kind="captions" />
            </video>
        </div>
        <!-- 外框層 -->
        <div class="layer layer-frame" bind:clientWidth={frameWidth} bind:clientHeight={frameHeight}>
            {#if isPortrait}
                <div class="spacer"></div>
            {/if}

            <img src={frameSrc} alt="手機外框" />

            {#if isPortrait}
                <div class="spacer"></div>
            {/if}
        </div>
    </section>
    <!-- 控制層 -->
    {#if showControls}
        <section class="slider-field">
            <div class="slider-header">
                <label for="videoWidthScale">寬度比例</label>
                <output for="videoWidthScale">{videoWidthScale}%</output>
            </div>
            <input type="range" min="10" max="120" step="0.1" disabled={isConverting} bind:value={videoWidthScale}/>
        </section>
        <section class="submit-container">
            <input class="submit-action" class:cancel-action={isConverting} type="button" value={isConverting ? covertingText : "開始轉換"} disabled={isCancelling} onclick={isConverting ? cancelComposeVideo : composeVideo}/>
        </section>
    {/if}
</main>

<style>
    @import "./less/main.less"
</style>
