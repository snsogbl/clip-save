<template>
  <div class="color-content">
    <div class="color-item-detail">
      <div class="color-display-area">
        <div 
          class="color-circle-large" 
          :style="{ backgroundColor: currentColor }"
          @click="openColorPicker"
          title="点击选择颜色"
        ></div>
        <!-- 隐藏的调色板 -->
        <input 
          ref="colorPickerRef" 
          type="color" 
          :value="hexValue || '#000000'"
          @input="onColorChange"
          class="color-picker-hidden"
        />
        <!-- 隐藏的临时元素用于解析颜色 -->
        <div ref="colorParserRef" :style="{ color: currentColor }" class="color-parser-hidden"></div>
        <div class="color-values-container">
          <div class="color-value-item" @click="copyToClipboard(currentColor)" title="点击复制">
            <span class="color-value-label">原始值:</span>
            <span class="color-value">{{ currentColor }}</span>
          </div>
          <div v-if="rgbValue" class="color-value-item" @click="copyToClipboard(`rgb(${rgbValue.r}, ${rgbValue.g}, ${rgbValue.b})`)" title="点击复制">
            <span class="color-value-label">RGB:</span>
            <span class="color-value">
              rgb({{ rgbValue.r }}, {{ rgbValue.g }}, {{ rgbValue.b }})
            </span>
          </div>
          <div v-if="rgbValue" class="color-value-item" @click="copyToClipboard(hexValue)" title="点击复制">
            <span class="color-value-label">HEX:</span>
            <span class="color-value">
              {{ hexValue }}
            </span>
          </div>
          <div v-if="rgbValue?.a !== undefined" class="color-value-item" @click="copyToClipboard(`${(rgbValue.a * 100).toFixed(0)}%`)" title="点击复制">
            <span class="color-value-label">透明度:</span>
            <span class="color-value">
              {{ (rgbValue.a * 100).toFixed(0) }}%
            </span>
          </div>
        </div>
        <!-- 复制成功提示 -->
        <div v-if="showCopyTip" class="copy-tip">已复制</div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref, watch } from "vue";

interface Props {
  color: string;
}

const props = defineProps<Props>();

// 定义 emit
const emit = defineEmits<{
  (e: 'colorChange', color: string): void
}>();

// 使用 Vue ref 引用模板中的元素
const colorParserRef = ref<HTMLDivElement | null>(null);
const colorPickerRef = ref<HTMLInputElement | null>(null);

// 当前显示的颜色（内部可变状态）
const currentColor = ref(props.color);

// 复制提示状态
const showCopyTip = ref(false);
let copyTipTimer: number | null = null;

//  props.color 变化，同步到内部状态
watch(() => props.color, (newColor) => {
  currentColor.value = newColor;
});

// 解析颜色值为 RGB
const rgbValue = computed(() => {
  if (!currentColor.value) return null;

  const color = currentColor.value.trim().toLowerCase();

  // 解析 HEX 格式 (#fff, #ffffff)
  if (color.startsWith('#')) {
    let hex = color.slice(1);
    
    // 处理 3 位 HEX (#fff -> #ffffff)
    if (hex.length === 3) {
      hex = hex.split('').map(char => char + char).join('');
    }
    
    if (hex.length === 6) {
      const r = parseInt(hex.slice(0, 2), 16);
      const g = parseInt(hex.slice(2, 4), 16);
      const b = parseInt(hex.slice(4, 6), 16);
      return { r, g, b };
    }
  }

  // 解析 RGB/RGBA 格式
  const rgbMatch = color.match(/rgba?\((\d+),\s*(\d+),\s*(\d+)(?:,\s*([\d.]+))?\)/);
  if (rgbMatch) {
    return {
      r: parseInt(rgbMatch[1]),
      g: parseInt(rgbMatch[2]),
      b: parseInt(rgbMatch[3]),
      a: rgbMatch[4] ? parseFloat(rgbMatch[4]) : undefined,
    };
  }

  // 解析 HSL/HSLA 格式
  const hslMatch = color.match(/hsla?\((\d+),\s*(\d+)%,\s*(\d+)%(?:,\s*([\d.]+))?\)/);
  if (hslMatch) {
    const h = parseInt(hslMatch[1]) / 360;
    const s = parseInt(hslMatch[2]) / 100;
    const l = parseInt(hslMatch[3]) / 100;
    const a = hslMatch[4] ? parseFloat(hslMatch[4]) : undefined;
    
    // HSL 转 RGB
    const hslToRgb = (h: number, s: number, l: number) => {
      let r, g, b;
      
      if (s === 0) {
        r = g = b = l; // 灰度
      } else {
        const hue2rgb = (p: number, q: number, t: number) => {
          if (t < 0) t += 1;
          if (t > 1) t -= 1;
          if (t < 1/6) return p + (q - p) * 6 * t;
          if (t < 1/2) return q;
          if (t < 2/3) return p + (q - p) * (2/3 - t) * 6;
          return p;
        };
        
        const q = l < 0.5 ? l * (1 + s) : l + s - l * s;
        const p = 2 * l - q;
        r = hue2rgb(p, q, h + 1/3);
        g = hue2rgb(p, q, h);
        b = hue2rgb(p, q, h - 1/3);
      }
      
      return {
        r: Math.round(r * 255),
        g: Math.round(g * 255),
        b: Math.round(b * 255)
      };
    };
    
    const rgb = hslToRgb(h, s, l);
    return { ...rgb, a };
  }

  // 对于命名颜色，使用 getComputedStyle 作为备选方案
  if (colorParserRef.value) {
    const computedColor = window.getComputedStyle(colorParserRef.value).color;
    const computedRgbMatch = computedColor.match(/rgba?\((\d+),\s*(\d+),\s*(\d+)(?:,\s*([\d.]+))?\)/);
    if (computedRgbMatch) {
      return {
        r: parseInt(computedRgbMatch[1]),
        g: parseInt(computedRgbMatch[2]),
        b: parseInt(computedRgbMatch[3]),
        a: computedRgbMatch[4] ? parseFloat(computedRgbMatch[4]) : undefined,
      };
    }
  }

  return null;
});

// 将 RGB 转换为 HEX
const hexValue = computed(() => {
  if (!rgbValue.value) return "";

  const toHex = (n: number) => {
    const hex = Math.round(n).toString(16).padStart(2, "0");
    return hex.toUpperCase();
  };

  return `#${toHex(rgbValue.value.r)}${toHex(rgbValue.value.g)}${toHex(
    rgbValue.value.b
  )}`;
});

// 打开调色板
const openColorPicker = () => {
  if (colorPickerRef.value) {
    colorPickerRef.value.click();
  }
};

// 处理颜色变化
const onColorChange = (event: Event) => {
  const target = event.target as HTMLInputElement;
  const newColor = target.value;
  
  // 立即更新UI
  currentColor.value = newColor;
  
  // 通知父组件颜色变化
  emit('colorChange', newColor);
};

// 复制到剪贴板
const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text);
    
    // 显示复制成功提示
    showCopyTip.value = true;
    
    // 清除之前的定时器
    if (copyTipTimer) {
      clearTimeout(copyTipTimer);
    }
    
    // 1.5秒后隐藏提示
    copyTipTimer = setTimeout(() => {
      showCopyTip.value = false;
    }, 1500);
  } catch (err) {
    console.error('复制失败:', err);
  }
};
</script>

<style scoped>
/* 颜色显示样式 */
.color-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.color-item-detail {
  padding: 20px;
  background-color: #ffffff;
  border-radius: 12px;
  border: 1px solid #e0e0e0;
  transition: all 0.2s ease;
}

.color-item-detail:hover {
  border-color: #2196f3;
  box-shadow: 0 2px 8px rgba(33, 150, 243, 0.1);
}

.color-display-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 24px;
  border-radius: 8px;
}

.color-circle-large {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  border: 3px solid #e0e0e0;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transition: transform 0.2s ease;
  cursor: pointer;
}

.color-circle-large:hover {
  transform: scale(1.05);
  border-color: #2196f3;
}

.color-values-container {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.color-value-item {
  display: flex;
  align-items: center;
  gap: 12px;
  background-color: #ffffff;
  padding: 12px 16px;
  border-radius: 6px;
  border: 1px solid #e0e0e0;
  cursor: pointer;
  transition: all 0.2s ease;
}

.color-value-item:hover {
  background-color: #f5f5f5;
  border-color: #2196f3;
  transform: translateX(2px);
}

.color-value-label {
  font-size: 14px;
  font-weight: 600;
  color: #6d6d70;
  min-width: 70px;
}

.color-value {
  font-size: 15px;
  font-weight: 600;
  color: #1a1a1a;
  font-family: "SF Mono", Monaco, Consolas, monospace;
  flex: 1;
}

/* 隐藏的调色板 */
.color-picker-hidden {
  position: absolute;
  width: 0;
  height: 0;
  opacity: 0;
  pointer-events: none;
}

/* 隐藏的颜色解析元素 */
.color-parser-hidden {
  position: absolute;
  visibility: hidden;
  pointer-events: none;
}

/* 复制成功提示 */
.copy-tip {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background-color: rgba(0, 0, 0, 0.75);
  color: #ffffff;
  padding: 12px 24px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  z-index: 9999;
  animation: fadeInOut 1.5s ease;
  pointer-events: none;
}

@keyframes fadeInOut {
  0% {
    opacity: 0;
    transform: translate(-50%, -50%) scale(0.9);
  }
  20% {
    opacity: 1;
    transform: translate(-50%, -50%) scale(1);
  }
  80% {
    opacity: 1;
    transform: translate(-50%, -50%) scale(1);
  }
  100% {
    opacity: 0;
    transform: translate(-50%, -50%) scale(0.9);
  }
}
</style>

