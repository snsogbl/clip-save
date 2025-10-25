import { ref, computed } from 'vue';
import { ElMessage } from 'element-plus';

export interface HotkeySettings {
  hotkey: string;
}

export function useHotkey(settings: { value: HotkeySettings }) {
  // 快捷键录制状态
  const isRecording = ref(false);
  
  // 当前录制的快捷键（实时显示）
  const currentRecordingHotkey = ref("");
  
  // 录制超时定时器
  let recordingTimeout: any = null;

  // 开始录制快捷键
  function startRecording() {
    recordingTimeout && clearTimeout(recordingTimeout);
    isRecording.value = true;
    currentRecordingHotkey.value = "";

    // 添加全局键盘事件监听
    document.addEventListener("keydown", handleKeyDown, true);

    // 5秒后自动停止录制
    recordingTimeout = setTimeout(() => {
      clearTimeout(recordingTimeout);
      recordingTimeout = null;
      if (isRecording.value) {
        stopRecording();
      }
    }, 5000);
  }

  // 停止录制快捷键
  function stopRecording() {
    isRecording.value = false;
    currentRecordingHotkey.value = "";
    document.removeEventListener("keydown", handleKeyDown, true);
  }

  // 处理键盘按下事件
  function handleKeyDown(event: KeyboardEvent) {
    if (!isRecording.value) return;

    event.preventDefault();
    event.stopPropagation();

    const modifiers: string[] = [];
    const keyMap: { [key: string]: string } = {
      Control: "Control",
      Meta: "Command",
      Shift: "Shift",
      Alt: "Option",
    };

    console.log("event", event);

    // 收集修饰键
    if (event.ctrlKey) modifiers.push("Control");
    if (event.metaKey) modifiers.push("Command");
    if (event.shiftKey) modifiers.push("Shift");
    if (event.altKey) modifiers.push("Option");

    // 获取主键
    let key = event.key;

    // 跳过单独的修饰键，但更新显示
    if (keyMap[key]) {
      // 只显示修饰键
      if (modifiers.length > 0) {
        currentRecordingHotkey.value = modifiers.join("+");
        console.log("currentRecordingHotkey", currentRecordingHotkey.value);
      }
      return;
    }

    // 处理特殊字符问题（如 Option+Command+V 会变成 √）
    // 使用 event.code 来获取物理按键而不是字符
    const physicalKey = event.code;
    let displayKey = key;

    // 将物理按键代码转换为显示键
    if (physicalKey.startsWith("Key")) {
      // 字母键：KeyA -> A
      displayKey = physicalKey.replace("Key", "");
    } else if (physicalKey.startsWith("Digit")) {
      // 数字键：Digit1 -> 1
      displayKey = physicalKey.replace("Digit", "");
    } else {
      // 其他特殊键的映射
      const specialKeyMap: { [key: string]: string } = {
        Space: "Space",
        Enter: "Enter",
        Return: "Return",
        Tab: "Tab",
        Escape: "Escape",
        Backspace: "Delete",
        Delete: "Delete",
        ArrowUp: "Up",
        ArrowDown: "Down",
        ArrowLeft: "Left",
        ArrowRight: "Right",
        F1: "F1",
        F2: "F2",
        F3: "F3",
        F4: "F4",
        F5: "F5",
        F6: "F6",
        F7: "F7",
        F8: "F8",
        F9: "F9",
        F10: "F10",
        F11: "F11",
        F12: "F12",
      };

      if (specialKeyMap[physicalKey]) {
        displayKey = specialKeyMap[physicalKey];
      } else if (physicalKey.startsWith("F")) {
        // F键：F1, F2, etc.
        displayKey = physicalKey;
      }
    }

    // 将字母转为大写
    if (displayKey.length === 1 && /[a-zA-Z]/.test(displayKey)) {
      displayKey = displayKey.toUpperCase();
    }

    // 构建快捷键字符串并实时显示
    if (modifiers.length > 0) {
      const newHotkey = [...modifiers, displayKey].join("+");
      currentRecordingHotkey.value = newHotkey;

      // 验证快捷键
      const validation = validateHotkey(newHotkey);
      if (!validation.valid) {
        ElMessage.warning(validation.message);
        stopRecording();
        return;
      }

      // 检查冲突
      if (checkHotkeyConflict(newHotkey)) {
        ElMessage.warning("此快捷键与系统快捷键冲突，请选择其他组合");
        stopRecording();
        return;
      }

      // 设置快捷键并停止录制
      settings.value.hotkey = newHotkey;

      // 停止录制
      stopRecording();
    } else {
      // 如果没有修饰键，只显示主键
      currentRecordingHotkey.value = displayKey;
    }
  }

  // 验证快捷键格式
  function validateHotkey(hotkey: string): { valid: boolean; message?: string } {
    if (!hotkey) {
      return { valid: false, message: "快捷键不能为空" };
    }

    const parts = hotkey.split("+");
    if (parts.length < 2) {
      return { valid: false, message: "快捷键至少需要两个键" };
    }

    // 检查是否包含修饰键
    const modifiers = parts.slice(0, -1);
    const mainKey = parts[parts.length - 1];

    const validModifiers = [
      "Control",
      "Command",
      "Meta",
      "Shift",
      "Alt",
      "Option",
    ];
    const hasValidModifier = modifiers.some((mod) =>
      validModifiers.includes(mod)
    );

    if (!hasValidModifier) {
      return {
        valid: false,
        message: "快捷键必须包含修饰键（Control、Command、Shift、Alt）",
      };
    }

    // 检查主键是否有效
    const validKeys =
      /^[A-Z0-9]$|^(F[1-9]|F1[0-2])$|^(Space|Enter|Return|Tab|Escape|Delete|Up|Down|Left|Right)$/i;
    if (!validKeys.test(mainKey)) {
      return { valid: false, message: "不支持的主键" };
    }

    return { valid: true };
  }

  // 检查快捷键冲突
  function checkHotkeyConflict(hotkey: string): boolean {
    // 常见的系统快捷键
    const systemHotkeys = [
      "Command+Q",
      "Command+W",
      "Command+C",
      "Command+V",
      "Command+X",
      "Command+Z",
      "Command+A",
      "Command+S",
      "Command+N",
      "Command+O",
      "Command+P",
      "Command+F",
      "Command+H",
      "Command+M",
      "Command+T",
      "Command+R",
      "Command+D",
      "Command+L",
      "Control+C",
      "Control+X",
      "Control+Z",
      "Control+A",
      "Control+S",
      "Control+N",
      "Control+O",
      "Control+P",
      "Control+F",
      "Control+H",
      "Control+R",
      "Control+T",
      "Control+W",
      "Control+Q",
      "Control+Shift+Z",
      "Control+Y",
      "Command+Shift+3",
      "Command+Shift+4",
      "Command+Shift+5",
      "Command+Space",
      "Control+Space",
      "Command+Tab",
      "Control+Tab",
    ];

    return systemHotkeys.includes(hotkey);
  }

  // 清理函数
  function cleanup() {
    stopRecording();
    if (recordingTimeout) {
      clearTimeout(recordingTimeout);
      recordingTimeout = null;
    }
  }

  return {
    // 状态
    isRecording,
    currentRecordingHotkey,
    
    // 方法
    startRecording,
    stopRecording,
    cleanup,
    
    // 工具函数
    validateHotkey,
    checkHotkeyConflict,
  };
}
