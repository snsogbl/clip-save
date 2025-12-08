<template>
  <div class="image-content">
    <el-image
      :src="`data:image/png;base64,${imageData}`"
      :alt="$t('components.image.clipboardImage')"
      class="content-image"
      fit="scale-down"
      preview-teleported
      :preview-src-list="[`data:image/png;base64,${imageData}`]"
    />
    <!-- 二维码识别结果 -->
    <div v-if="qrCodeResult" class="qr-result">
      <div class="qr-result-header">
        <span class="qr-result-title">{{
          $t("components.image.qrContent")
        }}</span>
        <el-button class="me-button" round @click="copyQRResult">{{
          $t("components.image.copy")
        }}</el-button>
      </div>
      <div class="qr-result-content">{{ qrCodeResult }}</div>
    </div>
    <div class="button-group">
      <el-button class="me-button" round @click="handleSave">
        {{ $t("components.image.saveToLocal") }}
      </el-button>
      <el-button
        v-if="isQRCode"
        class="me-button"
        round
        @click="handleQRCode"
        :disabled="isRecognizing"
      >
        {{
          isRecognizing
            ? $t("components.image.recognizing")
            : $t("components.image.recognizeQR")
        }}
      </el-button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, watch, onMounted, onUnmounted } from "vue";
import { useI18n } from "vue-i18n";
import {
  SaveImagePNG,
  DetectQRCode,
  RecognizeQRCode,
} from "../../../../wailsjs/go/main/App";
import { ElMessage } from "element-plus";

const { t } = useI18n();

interface Props {
  imageData: string;
}

const props = defineProps<Props>();

// 响应式数据
const isQRCode = ref(false);
const isRecognizing = ref(false);
const qrCodeResult = ref("");
const isDetecting = ref(false);
let detectTimeout: number | null = null;

// 检测图片中是否包含二维码
async function detectQRCode() {
  if (!props.imageData || isDetecting.value) {
    return;
  }

  isDetecting.value = true;
  try {
    const hasQRCode = await DetectQRCode(props.imageData);
    isQRCode.value = hasQRCode;
    // 如果检测到二维码，清空之前的结果
    if (!hasQRCode) {
      qrCodeResult.value = "";
    }
  } catch (error) {
    console.error("检测二维码失败:", error);
    isQRCode.value = false;
  } finally {
    isDetecting.value = false;
  }
}

// 识别二维码内容
async function handleQRCode() {
  if (isRecognizing.value) return;

  isRecognizing.value = true;
  try {
    const result = await RecognizeQRCode(props.imageData);
    qrCodeResult.value = result;
    ElMessage.success(t("components.image.qrGenerated"));
  } catch (error) {
    console.error("识别二维码失败:", error);
    ElMessage.error(t("components.image.qrGenerateFailed"));
  } finally {
    isRecognizing.value = false;
  }
}

// 复制二维码结果
function copyQRResult() {
  if (qrCodeResult.value) {
    navigator.clipboard
      .writeText(qrCodeResult.value)
      .then(() => {
        ElMessage.success(t("message.copySuccess"));
      })
      .catch(() => {
        ElMessage.error(t("components.image.qrCopyFailed"));
      });
  }
}

function handleSave() {
  try {
    const ts = new Date();
    const pad = (n: number) => n.toString().padStart(2, "0");
    const suggested = `clipboard-${ts.getFullYear()}${pad(
      ts.getMonth() + 1
    )}${pad(ts.getDate())}-${pad(ts.getHours())}${pad(ts.getMinutes())}${pad(
      ts.getSeconds()
    )}.png`;
    // 在保存前抑制窗口隐藏（避免保存对话框导致的 blur 触发）
    (window as any).__suppressHideWindow = true;
    SaveImagePNG(props.imageData, suggested)
      .then((savePath) => {
        if (savePath) {
          ElMessage.success(t("components.image.qrSaved"));
        }
      })
      .catch((e) => {
        console.error("保存图片失败", e);
        ElMessage.error(t("components.image.qrSaveFailed"));
      })
      .finally(() => {
        // 恢复隐藏行为
        (window as any).__suppressHideWindow = false;
      });
  } catch (e) {
    console.error("保存图片失败", e);
    ElMessage.error(t("components.image.qrSaveFailed"));
  }
}

// 监听 imageData 变化，当图片数据变化时重新检测二维码
watch(
  () => props.imageData,
  (newImageData) => {
    if (newImageData) {
      // 重置状态
      isQRCode.value = false;
      qrCodeResult.value = "";

      // 清除之前的定时器
      if (detectTimeout) {
        clearTimeout(detectTimeout);
      }

      // 防抖检测，避免频繁调用
      detectTimeout = setTimeout(() => {
        detectQRCode();
      }, 300);
    }
  },
  { immediate: true }
);

// 组件挂载时检测二维码（作为备用）
onMounted(() => {
  if (props.imageData) {
    detectQRCode();
  }
});

// 组件卸载时清理定时器
onUnmounted(() => {
  if (detectTimeout) {
    clearTimeout(detectTimeout);
  }
});
</script>

<style scoped>
.image-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.content-image {
  max-width: 100%;
  height: 480px;
  /* height: auto; */
  cursor: pointer;
  border-radius: 8px;
  transition: transform 0.2s ease;
  margin: 0 auto;
}

.content-image:hover {
  transform: scale(1.02);
}

.button-group {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

/* 模态框样式 */
.image-modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.8);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  backdrop-filter: blur(4px);
}

.modal-content {
  position: relative;
  max-width: 90%;
  max-height: 90%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal-image {
  max-width: 100%;
  max-height: 100%;
  border-radius: 8px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  transition: transform 0.2s ease;
  transform-origin: center center;
}

/* 缩放控制样式 */
.zoom-controls {
  position: absolute;
  bottom: -50px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 8px;
  background: rgba(0, 0, 0, 0.8);
  padding: 8px 16px;
  border-radius: 20px;
  backdrop-filter: blur(4px);
}

/* 二维码识别结果样式 */
.qr-result {
  background-color: #f5f5f5;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 12px;
  margin-top: 8px;
}

.qr-result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.qr-result-title {
  font-weight: 600;
  color: #333;
  font-size: 14px;
}

.qr-result-content {
  background-color: #ffffff;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  padding: 8px;
  font-size: 13px;
  line-height: 1.4;
  word-break: break-all;
  white-space: pre-wrap;
  max-height: 200px;
  overflow-y: auto;
}
</style>
