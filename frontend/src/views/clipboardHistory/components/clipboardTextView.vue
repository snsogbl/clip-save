<template>
  <div class="text-content">
    <!-- 原始文本显示 -->
    <pre><code>{{ text }}</code></pre>

    <!-- 解码按钮 -->
    <div v-if="showDecodeButtons" class="decode-buttons">
      <button
        v-if="needsURIDecoding"
        class="decode-btn"
        @click="toggleURIDecode"
      >
        <el-icon :size="14" style="margin-right: 4px">
          <Link />
        </el-icon>
        解码 URI
      </button>
      <button
        v-if="needsUnicodeDecoding"
        class="decode-btn"
        @click="toggleUnicodeDecode"
      >
        <el-icon :size="14" style="margin-right: 4px">
          <Document />
        </el-icon>
        解码 Unicode
      </button>
    </div>

    <!-- 解码后的文本显示区域 -->
    <div v-if="hasDecodedText" class="decoded-section">
      <div class="section-title">解码后文本</div>
      <pre class="content-text decoded">{{ decodedText }}</pre>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, watch, nextTick, onMounted } from "vue";
import { Link, Document } from "@element-plus/icons-vue";
import hljs from "highlight.js";

const props = defineProps<{
  text: string;
}>();

const decodedText = ref("");


// 检测是否需要URI解码
const needsURIDecoding = computed(() => {
  // 检测是否包含URL编码字符（%XX格式）
  const uriPattern = /%[0-9A-Fa-f]{2}/;
  return uriPattern.test(props.text);
});

// 检测是否需要Unicode解码
const needsUnicodeDecoding = computed(() => {
  // 检测是否包含Unicode转义序列（\uXXXX格式）
  const unicodePattern = /\\u[0-9A-Fa-f]{4}/;
  return unicodePattern.test(props.text);
});

// 是否显示解码按钮
const showDecodeButtons = computed(() => {
  return needsURIDecoding.value || needsUnicodeDecoding.value;
});

// 是否有解码后的文本
const hasDecodedText = computed(() => {
  return decodedText.value !== "";
});

// Unicode解码函数
function decodeUnicode(str: string): string {
  return str.replace(/\\u([0-9A-Fa-f]{4})/g, (match, hex) => {
    return String.fromCharCode(parseInt(hex, 16));
  });
}

// 切换URI解码
function toggleURIDecode() {
  try {
    decodedText.value = decodeURIComponent(props.text);
    highlightCode(); // 解码后重新高亮
  } catch (e) {
    console.error("URI解码失败:", e);
    decodedText.value = "解码失败：" + e;
  }
}

// 切换Unicode解码
function toggleUnicodeDecode() {
  try {
    decodedText.value = decodeUnicode(props.text);
  } catch (e) {
    console.error("Unicode解码失败:", e);
    decodedText.value = "解码失败：" + e;
  }
}

// 高亮代码块
const highlightCode = () => {
  nextTick(() => {
    document.querySelectorAll("pre code").forEach((el) => {
      const result = hljs.highlightAuto(el.textContent || '');
      el.innerHTML = result.value;
      el.className = `hljs ${result.language || ''}`;
    });
  });
};

// 当文本变化时，清空解码文本并重新高亮
watch(
  () => props.text,
  () => {
    decodedText.value = "";
    highlightCode();
  }
);

// 组件挂载时进行高亮
onMounted(() => {
  highlightCode();
});
</script>

<style scoped>
.text-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.content-text {
  white-space: pre-wrap;
  word-wrap: break-word;
  font-size: 15px;
  line-height: 1.7;
  margin: 0;
  color: #1a1a1a;
  letter-spacing: 0.01em;
}

.decode-buttons {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.decode-btn {
  padding: 8px 16px;
  border: 1px solid #2196f3;
  border-radius: 8px;
  background-color: #e3f2fd;
  color: #2196f3;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.decode-btn:hover {
  background-color: #2196f3;
  color: #ffffff;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(33, 150, 243, 0.3);
}

.decode-btn:active {
  transform: translateY(0);
  box-shadow: 0 1px 3px rgba(33, 150, 243, 0.2);
}

.decoded-section {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #e0e0e0;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 8px;
}

.content-text.decoded {
  background-color: #e8f5e9;
  border: 1px solid #4caf50;
  padding: 12px 16px;
  border-radius: 6px;
}

/* highlight.js 字体大小配置 */
pre code {
  font-size: 14px;
  line-height: 1.5;
}

.hljs {
  font-size: 14px;
  line-height: 1.5;
}
</style>
