<template>
  <div class="setting-container" style="--wails-draggable: no-drag">
    <!-- 顶部导航栏 -->
    <div class="setting-header">
      <el-button @click="$emit('back')" text>
        <el-icon :size="20" style="margin-right: 8px">
          <ArrowLeft />
        </el-icon>
        {{ $t('settings.back') }}
      </el-button>
      <h2>{{ $t('settings.title') }}</h2>
      <div style="width: 80px"></div>
    </div>

    <!-- 设置内容 -->
    <div class="setting-content">
      <div class="setting-section">
        <h3>{{ $t('settings.security') }}</h3>
        <div class="setting-item">
          <div class="setting-item-left">
            <el-icon :size="20" class="setting-icon">
              <Lock />
            </el-icon>
            <div class="setting-item-info">
              <div class="setting-item-title">{{ $t('settings.appPassword') }}</div>
              <div class="setting-item-desc">
                {{ $t('settings.passwordDesc') }}
              </div>
            </div>
          </div>
          <el-button @click="showPasswordDialog = true">
            {{ settings.password ? $t('settings.changePassword') : $t('settings.setPassword') }}
          </el-button>
        </div>

        <div class="setting-item" v-if="settings.password">
          <div class="setting-item-left">
            <el-icon :size="20" class="setting-icon">
              <Key />
            </el-icon>
            <div class="setting-item-info">
              <div class="setting-item-title">{{ $t('settings.removePassword') }}</div>
              <div class="setting-item-desc">{{ $t('settings.removePasswordDesc') }}</div>
            </div>
          </div>
          <el-button @click="removePassword" type="danger">
            {{ $t('settings.removePassword') }}
          </el-button>
          <el-button @click="lockPassword">{{ $t('settings.lock') }}</el-button>
        </div>
      </div>

      <div class="setting-section">
        <h3>{{ $t('settings.general') }}</h3>
        <div class="setting-item">
          <div class="setting-item-left">
            <el-icon :size="20" class="setting-icon">
              <Clock />
            </el-icon>
            <div class="setting-item-info">
              <div class="setting-item-title">{{ $t('settings.autoClean') }}</div>
              <div class="setting-item-desc">
                {{ $t('settings.autoCleanDesc') }}
              </div>
            </div>
          </div>
          <el-switch v-model="settings.autoClean" />
        </div>

        <div class="setting-item" v-if="settings.autoClean">
          <div class="setting-item-left">
            <el-icon :size="20" class="setting-icon">
              <Calendar />
            </el-icon>
            <div class="setting-item-info">
              <div class="setting-item-title">{{ $t('settings.retentionDays') }}</div>
              <div class="setting-item-desc">{{ $t('settings.retentionDaysDesc') }}</div>
            </div>
          </div>
          <el-input-number
            v-model="settings.retentionDays"
            :min="1"
            :max="365"
          />
        </div>

        <!-- 全局快捷键设置 -->
        <div class="setting-item">
          <div class="setting-item-left">
            <el-icon :size="20" class="setting-icon">
              <Operation />
            </el-icon>
            <div class="setting-item-info">
              <div class="setting-item-title">{{ $t('settings.hotkey') }}</div>
              <div class="setting-item-desc">{{ $t('settings.hotkeyDesc', [settings.hotkey]) }}</div>
            </div>
          </div>
          <div class="hotkey-input-area">
            <div
              class="hotkey-display"
              v-if="isRecording && currentRecordingHotkey"
            >
              <hotkey-display :hotkey="currentRecordingHotkey" />
            </div>
            <div
              class="hotkey-display"
              v-else-if="settings.hotkey && !isRecording"
            >
              <hotkey-display :hotkey="settings.hotkey" />
            </div>
            <div class="hotkey-placeholder" v-else-if="isRecording">
              {{ $t('settings.recordingPlaceholder') }}
            </div>
            <div class="hotkey-placeholder" v-else>{{ $t('settings.recordPlaceholder') }}</div>
            <el-button
              @click="startRecording"
              :disabled="isRecording"
              size="small"
              type="primary"
              style="margin-left: 12px"
            >
              {{ isRecording ? $t('settings.recording') : $t('settings.record') }}
            </el-button>
          </div>
        </div>

        <div class="setting-item">
          <div class="setting-item-left">
            <el-icon :size="20" class="setting-icon">
              <Delete />
            </el-icon>
            <div class="setting-item-info">
              <div class="setting-item-title">{{ $t('settings.clearAll') }}</div>
              <div class="setting-item-desc">
                {{ $t('settings.clearAllDesc') }}
              </div>
            </div>
          </div>
          <el-button @click="clearAllItems" type="danger">
            {{ $t('settings.clearAllButton') }}
          </el-button>
        </div>
      </div>

      <div class="setting-section">
        <h3>{{ $t('settings.interface') }}</h3>
        <div class="setting-item">
          <div class="setting-item-left">
            <el-icon :size="20" class="setting-icon">
              <List />
            </el-icon>
            <div class="setting-item-info">
              <div class="setting-item-title">{{ $t('settings.pageSize') }}</div>
              <div class="setting-item-desc">{{ $t('settings.pageSizeDesc') }}</div>
            </div>
          </div>
          <el-input-number
            v-model="settings.pageSize"
            :min="10"
            :max="200"
            :step="10"
          />
        </div>
        
        <!-- 语言设置 -->
        <div class="setting-item">
          <div class="setting-item-left">
            <el-icon :size="20" class="setting-icon">
              <Operation />
            </el-icon>
            <div class="setting-item-info">
              <div class="setting-item-title">{{ $t('settings.language') }}</div>
              <div class="setting-item-desc">{{ $t('settings.languageDesc') }}</div>
            </div>
          </div>
          <el-select style="width: 120px;" v-model="currentLanguage" @change="changeLanguage">
            <el-option label="中文" value="zh-CN" />
            <el-option label="English" value="en-US" />
            <el-option label="Français" value="fr-FR" />
            <el-option label="العربية" value="ar-SA" />
          </el-select>
        </div>
      </div>

      <div class="setting-section">
        <h3>{{ $t('settings.about') }}</h3>
        <div class="about-info">
          <div class="about-item">
            <span class="about-label">{{ $t('settings.appName') }}</span>
            <span class="about-value">{{ $t('app.name') }}</span>
          </div>
          <div class="about-item">
            <span class="about-label">{{ $t('settings.version') }}</span>
            <span class="about-value">{{ $t('app.version') }}</span>
          </div>
          <div class="about-item">
            <span class="about-label">{{ $t('settings.description') }}</span>
            <span class="about-value">{{ $t('app.description') }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 密码设置对话框 -->
    <el-dialog
      v-model="showPasswordDialog"
      :title="$t('passwordDialog.title')"
      width="400px"
      :close-on-click-modal="false"
    >
      <el-form @submit.prevent="savePassword">
        <el-form-item :label="$t('passwordDialog.newPassword')" required>
          <el-input
            v-model="newPassword"
            type="password"
            :placeholder="$t('passwordDialog.newPlaceholder')"
            show-password
          />
        </el-form-item>
        <el-form-item :label="$t('passwordDialog.confirmPassword')" required>
          <el-input
            v-model="confirmPassword"
            type="password"
            :placeholder="$t('passwordDialog.confirmPlaceholder')"
            show-password
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showPasswordDialog = false">{{ $t('passwordDialog.cancel') }}</el-button>
        <el-button type="primary" @click="savePassword">{{ $t('passwordDialog.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, onUnmounted, watch, computed } from "vue";
import { ElLoading, ElMessage, ElMessageBox } from "element-plus";
import {
  ArrowLeft,
  Clock,
  Calendar,
  List,
  Lock,
  Key,
  Delete,
  Operation,
  Warning,
} from "@element-plus/icons-vue";
import HotkeyDisplay from "./components/HotkeyDisplay.vue";
import { useHotkey } from "../../composables/useHotkey";
import { useI18n } from 'vue-i18n';
import {
  ClearAllItems,
  ClearItemsOlderThanDays,
  GetAppSettings,
  SaveAppSettings,
  RestartRegisterHotkey,
  GetCurrentLanguage,
  SetLanguage,
} from "../../../wailsjs/go/main/App";

const { t, locale } = useI18n();

// 定义事件
const emit = defineEmits(["back"]);

// 设置数据
const settings = ref({
  autoClean: true,
  retentionDays: 30,
  pageSize: 100,
  password: "", // 加密后的密码
  hotkey: "Command+Option+c", // 全局快捷键
});

// 当前语言
const currentLanguage = ref('zh-CN');

// 原始快捷键值，用于比较是否有修改
const originalHotkey = ref("");

// 快捷键重启状态
const isHotkeyRestarting = ref(false);

// 快捷键启用状态
const hotkeyEnabled = ref(true);

// 使用快捷键组合式函数
const {
  isRecording,
  currentRecordingHotkey,
  startRecording,
  stopRecording,
  cleanup: cleanupHotkey,
} = useHotkey(settings);

// 计算属性：判断快捷键是否被修改
const isHotkeyModified = computed(() => {
  return originalHotkey.value && settings.value.hotkey !== originalHotkey.value;
});

// 重启快捷键的函数
const restartHotkey = async () => {
  if (isHotkeyRestarting.value) {
    console.log("快捷键正在重启中，跳过重复调用");
    return;
  }

  isHotkeyRestarting.value = true;

  try {
    await RestartRegisterHotkey();
    ElMessage.success("快捷键已更新");
    originalHotkey.value = settings.value.hotkey;
  } catch (error) {
    console.error("重启快捷键失败:", error);
    ElMessage.error("快捷键更新失败，请重试");
  } finally {
    isHotkeyRestarting.value = false;
  }
};

watch(isHotkeyModified, () => {
  if (isHotkeyModified.value) {
    const loading = ElLoading.service({
      lock: true,
      text: "设置中...",
      // background: "rgba(0, 0, 0, 0.7)",
    });
    // 使用较短的延迟，因为后端已经优化了同步机制
    setTimeout(() => {
      restartHotkey();
      loading.close();
    }, 500);
  }
});

// 密码对话框
const showPasswordDialog = ref(false);
const newPassword = ref("");
const confirmPassword = ref("");

// 加载设置（从数据库）
async function loadSettings() {
  try {
    const savedSettings = await GetAppSettings();
    if (savedSettings) {
      const parsed = JSON.parse(savedSettings);
      settings.value = { ...settings.value, ...parsed };
      // 保存原始快捷键值用于比较
      originalHotkey.value = settings.value.hotkey;
      // 初始化快捷键启用状态
      hotkeyEnabled.value = !!settings.value.hotkey;
      console.log("✅ 已从数据库加载设置:", settings.value);
    } else {
      // 数据库应该已经有默认设置，如果没有则使用代码中的默认值
      console.log("⚠️ 数据库中无设置，使用代码默认值");
      await autoSaveSettings(); // 保存默认设置到数据库
      // 保存原始快捷键值用于比较
      originalHotkey.value = settings.value.hotkey;
    }
    
    // 加载当前语言
    try {
      const lang = await GetCurrentLanguage();
      currentLanguage.value = lang;
      locale.value = lang as any;
    } catch (e) {
      console.error("❌ 获取当前语言失败:", e);
    }
  } catch (e) {
    console.error("❌ 加载设置失败:", e);
  }
}

// 切换语言
async function changeLanguage(lang: string) {
  try {
    await SetLanguage(lang);
    locale.value = lang as any;
    currentLanguage.value = lang;
    ElMessage.success(t('message.settingsSaved'));
  } catch (error) {
    console.error("切换语言失败:", error);
    ElMessage.error(t('message.settingsError'));
  }
}

// 自动保存设置（到数据库）
async function autoSaveSettings() {
  try {
    await SaveAppSettings(JSON.stringify(settings.value));
    console.log("💾 设置已自动保存到数据库:", settings.value);
  } catch (e) {
    console.error("❌ 保存设置失败:", e);
  }
}

// 手动保存设置（显示成功消息）
async function saveSettings() {
  try {
    await SaveAppSettings(JSON.stringify(settings.value));
    ElMessage.success(t('message.settingsSaved'));
    console.log("✅ 设置已手动保存到数据库:", settings.value);
  } catch (e) {
    console.error("❌ 保存设置失败:", e);
    ElMessage.error(t('message.settingsError'));
  }
}

// 保存并返回
function saveAndBack() {
  saveSettings();
  setTimeout(() => {
    emit("back");
  }, 500); // 延迟返回，让用户看到保存成功的提示
}

// 立即手动清理
async function manualCleanNow() {
  if (!settings.value.autoClean) {
    ElMessage.warning("请先启用自动清理功能");
    return;
  }

  const retentionDays = settings.value.retentionDays || 30;

  try {
    ElMessage.info(`正在清理超过 ${retentionDays} 天的记录...`);
    console.log(`🗑️ 手动清理: 删除超过 ${retentionDays} 天的记录`);

    await ClearItemsOlderThanDays(retentionDays);

    ElMessage.success("清理完成！");
    console.log("✅ 手动清理完成");
  } catch (error) {
    console.error("❌ 清理失败:", error);
    ElMessage.error("清理失败: " + error);
  }
}

// 保存密码
async function savePassword() {
  if (!newPassword.value) {
    ElMessage.warning(t('passwordDialog.passwordRequired'));
    return;
  }

  if (newPassword.value !== confirmPassword.value) {
    ElMessage.error(t('passwordDialog.passwordMismatch'));
    return;
  }

  if (newPassword.value.length < 4) {
    ElMessage.warning(t('passwordDialog.passwordTooShort'));
    return;
  }

  try {
    const hashedPassword = await hashPassword(newPassword.value);
    settings.value.password = hashedPassword;

    await autoSaveSettings();

    ElMessage.success(t('passwordDialog.success'));
    showPasswordDialog.value = false;
    newPassword.value = "";
    confirmPassword.value = "";
  } catch (error) {
    console.error("设置密码失败:", error);
    ElMessage.error(t('passwordDialog.error'));
  }
}

// 移除密码
async function removePassword() {
  try {
    await ElMessageBox.confirm(
      t('message.removePasswordConfirm'),
      t('message.removePasswordTitle'),
      {
        confirmButtonText: t('passwordDialog.confirm'),
        cancelButtonText: t('passwordDialog.cancel'),
        type: "warning",
      }
    );

    settings.value.password = "";
    await autoSaveSettings();
    ElMessage.success(t('message.removePasswordSuccess'));
  } catch (error) {
    // 用户取消
  }
}

// 锁定重启应用
async function lockPassword() {
  window.location.reload();
}

// 清除所有剪贴板历史
async function clearAllItems() {
  try {
    await ElMessageBox.confirm(
      t('message.clearConfirm'),
      t('message.clearConfirmTitle'),
      {
        confirmButtonText: t('message.clearConfirmBtn'),
        cancelButtonText: t('message.clearCancelBtn'),
        type: "warning",
      }
    );

    ElMessage.info(t('message.clearProcessing'));
    console.log("🗑️ 开始清除所有剪贴板记录");

    await ClearAllItems();

    ElMessage.success(t('message.clearSuccess'));
    console.log("✅ 清除所有记录完成");
  } catch (error) {
    if (error === "cancel") {
      // 用户取消操作
      return;
    }
    console.error("❌ 清除失败:", error);
    ElMessage.error(t('message.clearError', [error]));
  }
}

async function hashPassword(password: string): Promise<string> {
  const encoder = new TextEncoder();
  const data = encoder.encode(password);
  const hashBuffer = await crypto.subtle.digest("SHA-256", data);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  const hashHex = hashArray
    .map((b) => b.toString(16).padStart(2, "0"))
    .join("");
  return hashHex;
}


//设置变化，自动保存
watch(
  settings,
  () => {
    autoSaveSettings();
  },
  { deep: true }
);

onMounted(() => {
  loadSettings();
});

// 组件卸载时清理快捷键相关资源
onUnmounted(() => {
  cleanupHotkey();
});
</script>

<style scoped>
.setting-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #fafafa;
}

.setting-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 68px;
  background-color: #ffffff;
  border-bottom: 1px solid #e0e0e0;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.setting-header h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #1a1a1a;
}

.setting-content {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
  max-width: 800px;
  margin: 0 auto;
  width: 100%;
}

.setting-section {
  background-color: #ffffff;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.setting-section h3 {
  margin: 0 0 20px 0;
  font-size: 18px;
  font-weight: 600;
  color: #1a1a1a;
}

.setting-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 0;
  border-bottom: 1px solid #f0f0f0;
}

.setting-item:last-child {
  border-bottom: none;
}

.setting-item-left {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  flex: 1;
}

.setting-icon {
  color: #666;
  margin-top: 2px;
}

.setting-item-info {
  flex: 1;
}

.setting-item-title {
  font-size: 16px;
  font-weight: 500;
  color: #000;
  margin-bottom: 4px;
}

.setting-item-desc {
  font-size: 14px;
  color: #333;
}

.setting-item-tip {
  display: flex;
  align-items: center;
  margin-top: 4px;
}

.about-info {
  padding: 8px 0;
}

.about-item {
  display: flex;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.about-item:last-child {
  border-bottom: none;
}

.about-label {
  font-weight: 600;
  color: #000;
  min-width: 100px;
  font-size: 16px;
}

.about-value {
  color: #333;
  font-size: 16px;
}

.setting-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
  padding: 24px 0 12px;
}

.auto-save-tip {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #8e8e93;
  font-size: 14px;
  padding: 0 0 24px;
}

/* 快捷键设置样式 */
.hotkey-input-area {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 200px;
  justify-content: end;
}

.hotkey-display {
  margin: 0;
}

.hotkey-placeholder {
  font-size: 12px;
  color: #999;
  font-style: italic;
  min-width: 120px;
}
</style>
