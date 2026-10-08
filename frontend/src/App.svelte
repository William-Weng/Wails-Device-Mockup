<script lang="ts">
    import { onMount } from "svelte";
    import { Events } from "@wailsio/runtime";
    import { dialog } from "./utility/dialog";

    import { SetFramePath } from "../bindings/wails-device-mockup/backend/imagepreviewservice";
    import { SetVideoPath } from "../bindings/wails-device-mockup/backend/videopreviewservice";
    import { Compose } from "../bindings/wails-device-mockup/backend/videocomposeservice";
    import { ParseFilePath } from "../bindings/wails-device-mockup/backend/filemanagerservice";

    const supportedImageTypes = new Set<ImageType>([".png", ".webp"]);
    const supportedVideoTypes = new Set<VideoType>([".mp4", ".mov"]);

    let framePath = "";
    let videoPath = "";
    let outputPath = "";

    let isCoverting = $state(false);
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

        return () => unsubscribeDrop();
    });

    async function imageFileDroppedAction(data: DroppedData): Promise<void> {
        const info = await ParseFilePath(data.path);
        const ext = info[2];

        if (isImageType(ext)) {
            framePath = data.path;
            frameSrc = await SetFramePath(framePath);
            return;
        }

        if (isVideoType(ext)) {
            videoPath = data.path;
            outputPath = `${info[0]}/${info[1]}-output.mp4`;
            videoSrc = await SetVideoPath(videoPath);
            return;
        }

        await dialog("warning", "錯誤", `不支援 ${ext} 格式`);
    }

    async function composeVideo() {
        if (frameSrc == "") {
            return;
        }
        if (videoSrc == "") {
            return;
        }

        isCoverting = true;

        try {
            await Compose(
                videoPath,
                videoWidth,
                videoHeight,
                framePath,
                frameWidth,
                frameHeight,
                outputPath,
            );

            isCoverting = false;
            dialog("info", "影片合成完成", `${outputPath}`);
        } catch (err) {
            isCoverting = false;
            dialog("info", "影片合成失敗", `${err}`);
        }
    }

    function isImageType(extension: string): extension is VideoType {
        return supportedImageTypes.has(extension as ImageType);
    }

    function isVideoType(extension: string): extension is VideoType {
        return supportedVideoTypes.has(extension as VideoType);
    }
</script>

<main>
    <section class="video-frame" data-file-drop-target>
        <div class="layer layer-video">
            <video
                controls
                src={videoSrc}
                style="width: {videoWidth}px;"
                bind:clientHeight={videoHeight}
            >
                <track kind="captions" />
            </video>
        </div>
        <div
            class="layer layer-frame"
            bind:clientWidth={frameWidth}
            bind:clientHeight={frameHeight}
        >
            <img src={frameSrc} alt="" />
        </div>
    </section>
    <section class="slider-field">
        <div class="slider-header">
            <label for="videoWidthScale">寬度比例</label>
            <output for="videoWidthScale">{videoWidthScale}%</output>
        </div>
        <input
            type="range"
            min="10"
            max="100"
            step="1"
            bind:value={videoWidthScale}
        />
    </section>
    <section class="submit-container">
        <input
            class="submit-action"
            type="button"
            value={isCoverting ? "轉換中…" : "開始轉換"}
            disabled={isCoverting}
            onclick={composeVideo}
        />
    </section>
</main>
