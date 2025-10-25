<template>
  <div class="url-content">
    <div class="url-display">{{ url }}</div>
    <div class="button-group">
      <button class="url-open-btn" @click="handleOpenUrl">
        <el-icon :size="14" style="margin-right: 4px">
          <Link />
        </el-icon>
        {{ $t('components.url.openInBrowser') }}
      </button>
      <button
        class="url-open-btn"
        @click="handleGenerateQR"
        :disabled="isGenerating"
      >
        {{ isGenerating ? $t('components.url.generating') : $t('components.url.generateQR') }}
      </button>
    </div>
    <!-- 二维码显示区域 -->
    <div v-if="qrCodeData" class="qr-code-section">
      <div class="qr-title">{{ $t('components.url.generatedQR') }}</div>
      <div class="qr-container">
        <el-image
          :src="`data:image/png;base64,${qrCodeData}`"
          :alt="$t('components.url.generatedQR')"
          class="qr-image"
          fit="scale-down"
          preview-teleported
          :preview-src-list="[`data:image/png;base64,${qrCodeData}`]"
        />
        <div class="qr-actions">
          <button class="qr-action-btn" @click="saveQRCode">{{ $t('components.url.saveQR') }}</button>
          <button class="qr-action-btn secondary" @click="copyQRCode">
            {{ $t('components.url.copyQR') }}
          </button>
        </div>
      </div>
    </div>
    <!-- URL 参数表格 -->
    <div v-if="urlParams.length > 0" class="url-params-section">
      <div class="params-title">{{ $t('components.url.urlParams') }}</div>
      <table class="params-table">
        <thead>
          <tr>
            <th class="param-key-header">{{ $t('components.url.key') }}</th>
            <th class="param-value-header">{{ $t('components.url.value') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="(param, index) in urlParams"
            :key="index"
            class="param-row"
          >
            <td class="param-key">{{ param.key }}</td>
            <td class="param-value">{{ param.value || "" }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref, watch } from "vue";
import { useI18n } from 'vue-i18n';
import { Link } from "@element-plus/icons-vue";
import { GenerateQRCode, SaveImagePNG, CopyImageToClipboard } from "../../../../wailsjs/go/main/App";
import { ElMessage } from "element-plus";

const { t } = useI18n();

interface Props {
  url: string;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  openUrl: [url: string];
}>();

// 响应式数据
const isGenerating = ref(false);
const qrCodeData = ref("");

// 解析 URL 参数
const urlParams = computed(() => {
  try {
    const urlObj = new URL(props.url);
    const params: Array<{ key: string; value: string }> = [];

    // 解析查询参数
    urlObj.searchParams.forEach((value, key) => {
      params.push({ key, value });
    });

    return params;
  } catch (error) {
    // 如果 URL 格式不正确，尝试手动解析
    const questionMarkIndex = props.url.indexOf("?");
    if (questionMarkIndex === -1) {
      return [];
    }

    const queryString = props.url.substring(questionMarkIndex + 1);
    const params: Array<{ key: string; value: string }> = [];

    const pairs = queryString.split("&");
    for (const pair of pairs) {
      const equalIndex = pair.indexOf("=");
      if (equalIndex !== -1) {
        const key = decodeURIComponent(pair.substring(0, equalIndex));
        const value = decodeURIComponent(pair.substring(equalIndex + 1));
        params.push({ key, value });
      } else if (pair) {
        // 只有键没有值的情况
        params.push({ key: decodeURIComponent(pair), value: "" });
      }
    }

    return params;
  }
});

function handleOpenUrl() {
  emit("openUrl", props.url);
}

// 监听URL变化
watch(
  () => props.url,
  (newUrl, oldUrl) => {
    if (newUrl !== oldUrl) {
      console.log("URL changed:", { oldUrl, newUrl });

      // 当URL变化时，清空之前生成的二维码
      if (qrCodeData.value) {
        qrCodeData.value = "";
        console.log("Cleared previous QR code due to URL change");
      }

      // 重置生成状态
      if (isGenerating.value) {
        isGenerating.value = false;
      }
    }
  },
  { immediate: true }
);

// 生成二维码
async function handleGenerateQR() {
  if (isGenerating.value) return;

  isGenerating.value = true;
  try {
    const qrData = await GenerateQRCode(props.url, 256);
    qrCodeData.value = qrData;
    ElMessage.success(t('components.url.qrGenerated'));
  } catch (error) {
    console.error("生成二维码失败:", error);
    ElMessage.error(t('components.url.qrGenerateFailed'));
  } finally {
    isGenerating.value = false;
  }
}

// 保存二维码
async function saveQRCode() {
  if (!qrCodeData.value) return;

  try {
    const ts = new Date();
    const pad = (n: number) => n.toString().padStart(2, "0");
    const suggested = `qrcode-${ts.getFullYear()}${pad(ts.getMonth() + 1)}${pad(
      ts.getDate()
    )}-${pad(ts.getHours())}${pad(ts.getMinutes())}${pad(ts.getSeconds())}.png`;

    // 在保存前抑制窗口隐藏
    (window as any).__suppressHideWindow = true;

    const savePath = await SaveImagePNG(qrCodeData.value, suggested);
    if (savePath) {
      ElMessage.success(t('components.url.qrSaved'));
    }
  } catch (error) {
    console.error("保存二维码失败:", error);
    ElMessage.error(t('components.url.qrSaveFailed'));
  } finally {
    // 恢复隐藏行为
    (window as any).__suppressHideWindow = false;
  }
}

// 复制二维码到剪贴板
async function copyQRCode() {
  if (!qrCodeData.value) return;

  try {
    // 使用后端API复制图片到剪贴板
    await CopyImageToClipboard(qrCodeData.value);
    ElMessage.success(t('components.url.qrCopied'));
  } catch (error) {
    console.error("复制二维码失败:", error);
    ElMessage.error(t('components.url.qrCopyFailed'));
  }
}
</script>

<style scoped>
.url-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.url-display {
  font-size: 15px;
  color: #2196f3;
  background-color: #f8f8f8;
  border-radius: 6px;
  word-break: break-all;
  line-height: 1.6;
  padding: 12px;
}

.button-group {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.url-open-btn {
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

.url-open-btn:hover {
  background-color: #2196f3;
  color: #ffffff;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(33, 150, 243, 0.3);
}

.url-open-btn:active {
  transform: translateY(0);
  box-shadow: 0 1px 3px rgba(33, 150, 243, 0.2);
}

.url-open-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.url-open-btn:disabled:hover {
  background-color: #e3f2fd;
  color: #2196f3;
  transform: none;
  box-shadow: none;
}

.url-params-section {
  margin-top: 16px;
  background-color: #ffffff;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 12px;
}

.params-title {
  font-size: 14px;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 8px;
}

.params-table {
  width: 100%;
  border-collapse: collapse;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  overflow: hidden;
  font-size: 13px;
}

.params-table thead {
  background-color: #f8f8f8;
}

.param-key-header,
.param-value-header {
  padding: 8px 12px;
  text-align: left;
  font-weight: 600;
  color: #1a1a1a;
  border-bottom: 1px solid #e0e0e0;
}

.param-key-header {
  width: 40%;
  border-right: 1px solid #e0e0e0;
}

.param-value-header {
  width: 60%;
}

.param-row {
  background-color: #ffffff;
  transition: background-color 0.2s ease;
}

.param-row:hover {
  background-color: #f8f8f8;
}

.param-row:nth-child(even) {
  background-color: #fafafa;
}

.param-row:nth-child(even):hover {
  background-color: #f0f0f0;
}

.param-key,
.param-value {
  padding: 8px 12px;
  border-bottom: 1px solid #e0e0e0;
  word-break: break-all;
  line-height: 1.4;
}

.param-key {
  font-weight: 500;
  color: #1a1a1a;
  border-right: 1px solid #e0e0e0;
}

.param-value {
  color: #6d6d70;
}

.param-row:last-child .param-key,
.param-row:last-child .param-value {
  border-bottom: none;
}

/* 二维码相关样式 */
.qr-code-section {
  margin-top: 16px;
  background-color: #ffffff;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 16px;
}

.qr-title {
  font-size: 14px;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 12px;
}

.qr-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.qr-image {
  width: 200px;
  height: 200px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  background-color: #ffffff;
}

.qr-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  justify-content: center;
}

.qr-action-btn {
  padding: 6px 12px;
  border: 1px solid #4caf50;
  border-radius: 6px;
  background-color: #4caf50;
  color: #ffffff;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.qr-action-btn:hover {
  background-color: #45a049;
  border-color: #45a049;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(76, 175, 80, 0.3);
}

.qr-action-btn.secondary {
  background-color: #ffffff;
  color: #4caf50;
  border-color: #4caf50;
}

.qr-action-btn.secondary:hover {
  background-color: #4caf50;
  color: #ffffff;
}
</style>
