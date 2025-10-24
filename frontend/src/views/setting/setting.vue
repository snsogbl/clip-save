<template>
  <div class="setting-container" style="--wails-draggable: no-drag">
    <!-- 顶部导航栏 -->
    <div class="setting-header">
      <el-button @click="$emit('back')" text>
        <el-icon :size="20" style="margin-right: 8px">
          <ArrowLeft />
        </el-icon>
        返回
      </el-button>
      <h2>设置</h2>
      <div style="width: 80px"></div>
    </div>

    <!-- 设置内容 -->
    <div class="setting-content">
      <div class="setting-section">
        <h3>安全设置</h3>
        <div class="setting-item">
          <div class="setting-item-left">
            <el-icon :size="20" class="setting-icon">
              <Lock />
            </el-icon>
            <div class="setting-item-info">
              <div class="setting-item-title">应用密码</div>
              <div class="setting-item-desc">
                设置密码后，每次打开应用需要输入密码
              </div>
            </div>
          </div>
          <el-button @click="showPasswordDialog = true" size="small">
            {{ settings.password ? "修改密码" : "设置密码" }}
          </el-button>
        </div>

        <div class="setting-item" v-if="settings.password">
          <div class="setting-item-left">
            <el-icon :size="20" class="setting-icon">
              <Key />
            </el-icon>
            <div class="setting-item-info">
              <div class="setting-item-title">移除密码</div>
              <div class="setting-item-desc">移除密码后可直接打开应用</div>
            </div>
          </div>
          <el-button @click="removePassword" size="small" type="danger">
            移除密码
          </el-button>
          <el-button @click="lockPassword" size="small"> 锁定 </el-button>
        </div>
      </div>

      <div class="setting-section">
        <h3>通用设置</h3>
        <div class="setting-item">
          <div class="setting-item-left">
            <el-icon :size="20" class="setting-icon">
              <Clock />
            </el-icon>
            <div class="setting-item-info">
              <div class="setting-item-title">自动清理历史</div>
              <div class="setting-item-desc">
                自动删除超过指定天数的剪贴板历史
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
              <div class="setting-item-title">保留天数</div>
              <div class="setting-item-desc">历史记录保留的天数</div>
            </div>
          </div>
          <el-input-number
            v-model="settings.retentionDays"
            :min="1"
            :max="365"
            size="small"
          />
        </div>

        <div class="setting-item">
          <div class="setting-item-left">
            <el-icon :size="20" class="setting-icon">
              <Operation />
            </el-icon>
            <div class="setting-item-info">
              <div class="setting-item-title">全局快捷键</div>
              <div class="setting-item-desc">
                按下快捷键唤起应用窗口，例如：Control+v, Command+Shift+C
              </div>
              <!-- <div class="setting-item-tip" v-if="isHotkeyModified">
                <el-icon :size="14" style="color: #f56c6c; margin-right: 4px">
                  <Warning />
                </el-icon>
                <span style="color: #f56c6c; font-size: 12px">
                  修改快捷键后需要重启应用才能生效
                </span>
              </div> -->
            </div>
          </div>
          <el-input
            v-model="settings.hotkey"
            placeholder="Control+V"
            size="small"
            style="width: 150px"
            @keydown="captureHotkey"
          />
        </div>

        <div class="setting-item">
          <div class="setting-item-left">
            <el-icon :size="20" class="setting-icon">
              <Delete />
            </el-icon>
            <div class="setting-item-info">
              <div class="setting-item-title">全部清除</div>
              <div class="setting-item-desc">
                清除所有剪贴板历史记录，此操作不可恢复
              </div>
            </div>
          </div>
          <el-button @click="clearAllItems" size="small" type="danger">
            清除全部
          </el-button>
        </div>
      </div>

      <div class="setting-section">
        <h3>界面设置</h3>
        <div class="setting-item">
          <div class="setting-item-left">
            <el-icon :size="20" class="setting-icon">
              <List />
            </el-icon>
            <div class="setting-item-info">
              <div class="setting-item-title">每页显示数量</div>
              <div class="setting-item-desc">列表中每次加载的记录数量</div>
            </div>
          </div>
          <el-input-number
            v-model="settings.pageSize"
            :min="10"
            :max="200"
            :step="10"
            size="small"
          />
        </div>
      </div>

      <div class="setting-section">
        <h3>关于</h3>
        <div class="about-info">
          <div class="about-item">
            <span class="about-label">应用名称：</span>
            <span class="about-value">剪存</span>
          </div>
          <div class="about-item">
            <span class="about-label">版本号：</span>
            <span class="about-value">1.0.3</span>
          </div>
          <div class="about-item">
            <span class="about-label">描述：</span>
            <span class="about-value">一个优雅的剪贴板历史管理工具</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 密码设置对话框 -->
    <el-dialog
      v-model="showPasswordDialog"
      title="设置应用密码"
      width="400px"
      :close-on-click-modal="false"
    >
      <el-form @submit.prevent="savePassword">
        <el-form-item label="新密码" required>
          <el-input
            v-model="newPassword"
            type="password"
            placeholder="请输入新密码"
            show-password
          />
        </el-form-item>
        <el-form-item label="确认密码" required>
          <el-input
            v-model="confirmPassword"
            type="password"
            placeholder="请再次输入密码"
            show-password
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showPasswordDialog = false">取消</el-button>
        <el-button type="primary" @click="savePassword">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, watch, computed } from "vue";
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
import {
  ClearAllItems,
  ClearItemsOlderThanDays,
  GetAppSettings,
  SaveAppSettings,
  RestartRegisterHotkey,
} from "../../../wailsjs/go/main/App";

// 定义事件
const emit = defineEmits(["back"]);

// 设置数据
const settings = ref({
  autoClean: true,
  retentionDays: 30,
  pageSize: 100,
  password: "", // 加密后的密码
  hotkey: "Control+v", // 全局快捷键
});

// 原始快捷键值，用于比较是否有修改
const originalHotkey = ref("");

// 快捷键重启状态
const isHotkeyRestarting = ref(false);

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
      console.log("✅ 已从数据库加载设置:", settings.value);
    } else {
      // 数据库应该已经有默认设置，如果没有则使用代码中的默认值
      console.log("⚠️ 数据库中无设置，使用代码默认值");
      await autoSaveSettings(); // 保存默认设置到数据库
      // 保存原始快捷键值用于比较
      originalHotkey.value = settings.value.hotkey;
    }
  } catch (e) {
    console.error("❌ 加载设置失败:", e);
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
    ElMessage.success("设置已保存");
    console.log("✅ 设置已手动保存到数据库:", settings.value);
  } catch (e) {
    console.error("❌ 保存设置失败:", e);
    ElMessage.error("保存设置失败");
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
    ElMessage.warning("请输入密码");
    return;
  }

  if (newPassword.value !== confirmPassword.value) {
    ElMessage.error("两次输入的密码不一致");
    return;
  }

  if (newPassword.value.length < 4) {
    ElMessage.warning("密码长度至少4位");
    return;
  }

  try {
    const hashedPassword = await hashPassword(newPassword.value);
    settings.value.password = hashedPassword;

    await autoSaveSettings();

    ElMessage.success("密码设置成功！下次启动应用需要输入密码");
    showPasswordDialog.value = false;
    newPassword.value = "";
    confirmPassword.value = "";
  } catch (error) {
    console.error("设置密码失败:", error);
    ElMessage.error("设置密码失败");
  }
}

// 移除密码
async function removePassword() {
  try {
    await ElMessageBox.confirm(
      "移除密码后，将不再需要密码即可打开应用。确定要移除密码吗？",
      "确认移除",
      {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      }
    );

    settings.value.password = "";
    await autoSaveSettings();
    ElMessage.success("密码已移除");
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
      "确定要清除所有剪贴板历史记录吗？此操作不可恢复！",
      "确认清除",
      {
        confirmButtonText: "确定清除",
        cancelButtonText: "取消",
        type: "warning",
      }
    );

    ElMessage.info("正在清除所有记录...");
    console.log("🗑️ 开始清除所有剪贴板记录");

    await ClearAllItems();

    ElMessage.success("已成功清除所有记录！");
    console.log("✅ 清除所有记录完成");

    // 刷新页面以更新显示
    setTimeout(() => {
      emit("back");
    }, 1000);
  } catch (error) {
    if (error === "cancel") {
      // 用户取消操作
      return;
    }
    console.error("❌ 清除失败:", error);
    ElMessage.error("清除失败: " + error);
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

// 捕获快捷键输入
function captureHotkey(event: KeyboardEvent) {
  event.preventDefault();

  const modifiers: string[] = [];
  const keyMap: { [key: string]: string } = {
    Control: "Control",
    Meta: "Command",
    Shift: "Shift",
    Alt: "Alt",
  };

  // 收集修饰键
  if (event.ctrlKey) modifiers.push("Control");
  if (event.metaKey) modifiers.push("Command");
  if (event.shiftKey) modifiers.push("Shift");
  if (event.altKey) modifiers.push("Alt");

  // 获取主键
  let key = event.key;

  // 跳过单独的修饰键
  if (keyMap[key]) {
    return;
  }

  // 将字母转为大写
  if (key.length === 1) {
    key = key.toUpperCase();
  }

  // 构建快捷键字符串
  if (modifiers.length > 0) {
    settings.value.hotkey = [...modifiers, key].join("+");
  }
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
  padding: 20px 24px;
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
  font-size: 15px;
  font-weight: 500;
  color: #1a1a1a;
  margin-bottom: 4px;
}

.setting-item-desc {
  font-size: 13px;
  color: #8e8e93;
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
  color: #1a1a1a;
  min-width: 100px;
  font-size: 14px;
}

.about-value {
  color: #6d6d70;
  font-size: 14px;
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
  font-size: 13px;
  padding: 0 0 24px;
}
</style>
