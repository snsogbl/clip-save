<template>
  <div class="text-content">
    <!-- 原始文本显示 -->
    <pre class="content-text"><code>{{ text }}</code></pre>

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
        {{ $t("components.text.decodeUri") }}
      </button>
      <button
        v-if="needsUnicodeDecoding"
        class="decode-btn"
        @click="toggleUnicodeDecode"
      >
        <el-icon :size="14" style="margin-right: 4px">
          <Document />
        </el-icon>
        {{ $t("components.text.decodeUnicode") }}
      </button>
    </div>

    <!-- 解码后的文本显示区域 -->
    <div v-if="hasDecodedText" class="decoded-section">
      <div class="section-title">{{ $t("components.text.decodedText") }}</div>
      <pre class="content-text decoded">{{ decodedText }}</pre>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, watch, nextTick, onMounted } from "vue";
import { useI18n } from "vue-i18n";
import { Link, Document } from "@element-plus/icons-vue";
import hljs from "highlight.js";

const { t } = useI18n();

const props = defineProps<{
  text: string;
}>();

const decodedText = ref("");
const languageType = ref("");

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
  } catch (e) {
    console.error("URI解码失败:", e);
    decodedText.value = t("components.text.decodeFailed", [e]);
  }
}

// 切换Unicode解码
function toggleUnicodeDecode() {
  try {
    decodedText.value = decodeUnicode(props.text);
  } catch (e) {
    console.error("Unicode解码失败:", e);
    decodedText.value = t("components.text.decodeFailed", [e]);
  }
}

const checkIsCode = () => {
  const text = props.text || "";
  if (!text) return false;

  // 大文本或超长行数直接跳过高亮，避免卡顿
  if (text.length > 50000) return false; // ~50KB 阈值
  const lines = text.split(/\r?\n/);
  if (lines.length > 2000) return false;

  // 代码特征正则（多语言通用 + 常见场景）
  const indicators: RegExp[] = [
    /function\s+\w+/m, // JS/TS
    /\b(class|interface|enum|struct)\b/m, // 多语言
    /\b(import|export)\b/m, // JS/TS/ESM
    /\b(let|const|var)\b/m, // JS/TS
    /\b(if|else|for|while|switch|case|return)\b[^{;]*[({;]/m, // 控制流
    /=>/m, // 箭头函数
    /#include\b|using\s+namespace\b|template\s*</m, // C/C++
    /SELECT\s+.+\s+FROM\b|INSERT\s+INTO\b|UPDATE\b|DELETE\s+FROM\b/i, // SQL
    /\/\/|\/\*|\*\//m, // 注释
    /^\s*#!/m, // shebang
    /^\s*<\w+[^>]*>.*<\/\w+>\s*$/m, // HTML/XML 单行标签
  ];

  let score = 0;
  for (const re of indicators) {
    if (re.test(text)) {
      score++;
      if (score >= 2) break;
    }
  }

  // JSON 判定（常见复制粘贴）
  const isJsonLike =
    /^\s*[\[{][\s\S]*[\]}]\s*$/.test(text) && /"\s*[\w$-]+\s*"\s*:/m.test(text);
  if (isJsonLike) score += 2;

  // 符号密度和缩进行比例（代码一般符号更多、缩进更多）
  const codeSymbolCount = (text.match(/[{}();=<>[\]]/g) || []).length;
  const symbolDensity = codeSymbolCount / Math.max(text.length, 1);
  const indentedLines = lines.filter((l) => /^\s{2,}|\t/.test(l)).length;
  if (symbolDensity > 0.02) score++;
  if (indentedLines / Math.max(lines.length, 1) > 0.2) score++;

  return score >= 2;
};

// 高亮代码块
const highlightCode = () => {
  nextTick(() => {
    if (checkIsCode()) {
      document.querySelectorAll("pre code").forEach((el) => {
        const testResult = hljs.highlightAuto(
          el.textContent.slice(0, 100) || ""
        );
        if (testResult.language) {
          const result = hljs.highlightAuto(el.textContent || "");
          languageType.value = result.language || "";
          if (result.language) {
            el.innerHTML = result.value;
            el.className = `hljs ${result.language || ""}`;
          }
        }
      });
    }
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
  font-size: 15px;
  line-height: 1.3;
}

.hljs {
  font-size: 15px;
  line-height: 1.3;
}
</style>
