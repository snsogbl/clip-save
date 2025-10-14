<template>
  <div class="url-content">
    <div class="url-display">{{ url }}</div>
    <button class="url-open-btn" @click="handleOpenUrl">
      <el-icon :size="14" style="margin-right: 4px">
        <Link />
      </el-icon>
      在浏览器中打开
    </button>

    <!-- URL 参数表格 -->
    <div v-if="urlParams.length > 0" class="url-params-section">
      <div class="params-title">URL 参数</div>
      <table class="params-table">
        <thead>
          <tr>
            <th class="param-key-header">键</th>
            <th class="param-value-header">值</th>
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
import { computed } from "vue";
import { Link } from "@element-plus/icons-vue";

interface Props {
  url: string;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  openUrl: [url: string];
}>();

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
</script>

<style scoped>
.url-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.url-display {
  font-size: 14px;
  color: #2196f3;
  background-color: #f8f8f8;
  border-radius: 6px;
  word-break: break-all;
  font-family: "SF Mono", Monaco, Consolas, monospace;
  line-height: 1.6;
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
  font-family: "SF Mono", Monaco, Consolas, monospace;
}

.param-value {
  color: #6d6d70;
  font-family: "SF Mono", Monaco, Consolas, monospace;
}

.param-row:last-child .param-key,
.param-row:last-child .param-value {
  border-bottom: none;
}
</style>

