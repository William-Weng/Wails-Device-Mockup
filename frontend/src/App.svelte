<script lang="ts">
    import { onMount } from "svelte";
    import { Events } from "@wailsio/runtime";
    import { dialog } from "./utility/dialog";

    import { SetFramePath, ImageInformation } from "../bindings/wails-device-mockup/backend/imagepreviewservice";
    import { SetVideoPath } from "../bindings/wails-device-mockup/backend/videopreviewservice";
    import { Compose } from "../bindings/wails-device-mockup/backend/videocomposeservice";
    import { ParseFilePath } from "../bindings/wails-device-mockup/backend/filemanagerservice";

    const supportedImageTypes = new Set<ImageType>([".png", ".webp"]);
    const supportedVideoTypes = new Set<VideoType>([".mp4", ".mov", ".ts", ".mkv", ".wmv"]);

    let framePath = "";
    let videoPath = "";
    let outputPath = "";

    let isPortrait = $state(false);
    let isCoverting = $state(false);
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

        if (isCoverting) { return; }

        const info = await ParseFilePath(data.path);
        const ext = info[2];

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
            const timestamp = Date.now();
            videoPath = data.path;
            outputPath = `${info[0]}/${info[1]}-output-${timestamp}.mp4`;
            videoSrc = await SetVideoPath(videoPath);
            return;
        }

        await dialog("warning", "錯誤", `不支援 ${ext} 格式`);
    }

    // 觸發合成影片的按鈕點擊事件
    async function composeVideo() {

        if (frameSrc == "") { return; }
        if (videoSrc == "") { return; }

        isCoverting = true;

        try {
            await Compose(videoPath, videoWidth, videoHeight, framePath, frameWidth, frameHeight, outputPath);
            reset()
            dialog("info", "影片合成完成", `${outputPath}`);
        } catch (err) {
            reset()
            dialog("info", "影片合成失敗", `${err}`);
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

    // 回復原設定值
    function reset() {
      isCoverting = false;
      covertingText = "0%"
    }
</script>

<main>
    <!-- 使用 Svelte 5 的 class 綁定語法 -->
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

    <section class="slider-field">
        <div class="slider-header">
            <label for="videoWidthScale">寬度比例</label>
            <output for="videoWidthScale">{videoWidthScale}%</output>
        </div>
        <input type="range" min="10" max="120" step="0.1" disabled={isCoverting} bind:value={videoWidthScale}/>
    </section>
    <section class="submit-container">
        <input class="submit-action" type="button" value={isCoverting ? covertingText : "開始轉換"} disabled={isCoverting} onclick={composeVideo}/>
    </section>
</main>
