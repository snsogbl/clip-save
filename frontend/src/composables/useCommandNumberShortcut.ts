import { ref, onMounted, onUnmounted, type Ref } from "vue";

interface UseCommandNumberShortcutOptions {
  /**
   * 是否启用快捷键（通常是一个 ref，比如 visible.value）
   */
  enabled: Ref<boolean> | (() => boolean);
  /**
   * 可选择的项目数量（最多显示 9 个）
   */
  itemCount: Ref<number> | (() => number);
  /**
   * 当按下 Command/Ctrl + 数字键时的回调函数
   * @param index 选中的索引（0-8）
   */
  onSelect: (index: number) => void | Promise<void>;
  /**
   * 是否使用捕获阶段捕获事件（默认 false）
   */
  useCapture?: boolean;
  /**
   * 是否阻止事件传播（默认 false）
   */
  stopPropagation?: boolean;
  /**
   * 是否阻止其他监听器处理（默认 false）
   */
  stopImmediatePropagation?: boolean;
}

/**
 * Command/Ctrl + 数字键（1-9）快捷键 composable
 * 用于快速选择列表中的项目
 */
export function useCommandNumberShortcut(
  options: UseCommandNumberShortcutOptions
) {
  const {
    enabled,
    itemCount,
    onSelect,
    useCapture = false,
    stopPropagation = false,
    stopImmediatePropagation = false,
  } = options;

  const isCommandPressed = ref(false);

  // 窗口可见性变化处理函数
  const handleVisibilityChange = () => {
    if (document.visibilityState === "hidden") {
      isCommandPressed.value = false;
    }
  };

  // 处理全局键盘事件（用于 Command+数字键快速选择）
  function handleGlobalKeydown(event: KeyboardEvent) {
    // 检查是否启用
    const isEnabled =
      typeof enabled === "function" ? enabled() : enabled.value;
    if (!isEnabled) return;

    // 检测 Command/Ctrl 键按下
    if (event.metaKey || event.ctrlKey) {
      // 只有在窗口可见时才显示标签
      if (!isCommandPressed.value && document.visibilityState === "visible") {
        isCommandPressed.value = true;
      }

      // 检测 Command+数字键（1-9）
      const numKey = parseInt(event.key);
      if (!isNaN(numKey) && numKey >= 1 && numKey <= 9) {
        // 阻止默认行为和传播
        event.preventDefault();
        if (stopPropagation) {
          event.stopPropagation();
        }
        if (stopImmediatePropagation) {
          event.stopImmediatePropagation();
        }

        // 快速选择对应索引的项目（索引从 0 开始，所以减 1）
        const count = typeof itemCount === "function" ? itemCount() : itemCount.value;
        const index = numKey - 1;
        if (index < count) {
          onSelect(index);
        }
        // 重置状态
        isCommandPressed.value = false;
        return;
      }
    } else {
      // 非 Command 键按下时，如果之前是按下的状态，检查是否是 Command 键本身
      if (
        event.key !== "Meta" &&
        event.key !== "Control" &&
        isCommandPressed.value
      ) {
        // 如果按下的不是 Command 键，说明 Command 已经松开
        isCommandPressed.value = false;
      }
    }
  }

  // 处理全局键盘松开事件
  function handleGlobalKeyup(event: KeyboardEvent) {
    // 检查是否启用
    const isEnabled =
      typeof enabled === "function" ? enabled() : enabled.value;
    if (!isEnabled) return;

    // Command/Ctrl 键松开
    if (
      event.key === "Meta" ||
      event.key === "Control" ||
      event.key === "MetaLeft" ||
      event.key === "MetaRight" ||
      event.key === "ControlLeft" ||
      event.key === "ControlRight"
    ) {
      isCommandPressed.value = false;
    }
  }

  onMounted(() => {
    // 监听全局键盘事件
    window.addEventListener("keydown", handleGlobalKeydown, useCapture);
    window.addEventListener("keyup", handleGlobalKeyup, useCapture);
    // 监听窗口可见性变化，隐藏窗口时重置状态
    document.addEventListener("visibilitychange", handleVisibilityChange);
  });

  onUnmounted(() => {
    // 清理 DOM 事件监听器
    window.removeEventListener("keydown", handleGlobalKeydown, useCapture);
    window.removeEventListener("keyup", handleGlobalKeyup, useCapture);
    document.removeEventListener("visibilitychange", handleVisibilityChange);
  });

  return {
    isCommandPressed,
  };
}

