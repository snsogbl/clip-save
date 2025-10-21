<template>
  <div class="image-content">
    <el-image
      :src="`data:image/png;base64,${imageData}`"
      alt="剪贴板图片"
      class="content-image"
      fit="scale-down"
      preview-teleported
      :preview-src-list="[`data:image/png;base64,${imageData}`]"
    />
    <!-- 二维码识别结果 -->
    <div v-if="qrCodeResult" class="qr-result">
      <div class="qr-result-header">
        <span class="qr-result-title">二维码内容：</span>
        <button class="copy-btn" @click="copyQRResult">复制</button>
      </div>
      <div class="qr-result-content">{{ qrCodeResult }}</div>
    </div>
    <div class="button-group">
      <button class="save-btn" @click="handleSave">保存到本地</button>
      <button
        v-if="isQRCode"
        class="save-btn qr-btn"
        @click="handleQRCode"
        :disabled="isRecognizing"
      >
        {{ isRecognizing ? "识别中..." : "识别二维码" }}
      </button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, watch, onMounted, onUnmounted } from "vue";
import {
  SaveImagePNG,
  DetectQRCode,
  RecognizeQRCode,
} from "../../../../wailsjs/go/main/App";
import { ElMessage } from "element-plus";

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
    ElMessage.success("二维码识别成功");
  } catch (error) {
    console.error("识别二维码失败:", error);
    ElMessage.error("识别二维码失败");
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
        ElMessage.success("已复制到剪贴板");
      })
      .catch(() => {
        ElMessage.error("复制失败");
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
          ElMessage.success("图片已保存");
        }
      })
      .catch((e) => {
        console.error("保存图片失败", e);
        ElMessage.error("保存失败");
      })
      .finally(() => {
        // 恢复隐藏行为
        (window as any).__suppressHideWindow = false;
      });
  } catch (e) {
    console.error("保存图片失败", e);
    ElMessage.error("保存失败");
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

.close-btn {
  position: absolute;
  top: -40px;
  right: -10px;
  background: rgba(255, 255, 255, 0.9);
  border: none;
  border-radius: 50%;
  width: 32px;
  height: 32px;
  font-size: 20px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s ease;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 1);
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

.zoom-btn {
  background: rgba(255, 255, 255, 0.2);
  border: none;
  border-radius: 50%;
  width: 32px;
  height: 32px;
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  font-size: 12px;
}

.zoom-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.3);
  transform: scale(1.1);
}

.zoom-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.zoom-btn:last-child {
  width: auto;
  padding: 0 12px;
  border-radius: 16px;
  font-size: 12px;
}

.zoom-percentage {
  color: white;
  font-size: 12px;
  min-width: 40px;
  text-align: center;
  font-weight: 500;
}

.save-btn {
  width: fit-content;
  padding: 6px 12px;
  border: 1px solid #2196f3;
  border-radius: 6px;
  background-color: #e3f2fd;
  color: #2196f3;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.save-btn:hover {
  background-color: #2196f3;
  color: #ffffff;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(33, 150, 243, 0.3);
}

.save-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.save-btn:disabled:hover {
  background-color: #e3f2fd;
  color: #2196f3;
  transform: none;
  box-shadow: none;
}

.qr-btn {
  background-color: #4caf50;
  border-color: #4caf50;
  color: #ffffff;
}

.qr-btn:hover:not(:disabled) {
  background-color: #45a049;
  border-color: #45a049;
  box-shadow: 0 2px 8px rgba(76, 175, 80, 0.3);
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

.copy-btn {
  padding: 4px 8px;
  border: 1px solid #2196f3;
  border-radius: 4px;
  background-color: #ffffff;
  color: #2196f3;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.copy-btn:hover {
  background-color: #2196f3;
  color: #ffffff;
}

.qr-result-content {
  background-color: #ffffff;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  padding: 8px;
  font-family: monospace;
  font-size: 13px;
  line-height: 1.4;
  word-break: break-all;
  white-space: pre-wrap;
  max-height: 200px;
  overflow-y: auto;
}
</style>
