<template>
  <el-dialog
    v-model="visible"
    width="70%"
    :close-on-click-modal="false"
    :close-on-press-escape="true"
    @close="handleClose"
  >
    <template #header>
      <div class="dialog-header">
        <span>{{ t("scripts.selectScript") }}</span>
        <el-button
          type="primary"
          link
          size="small"
          @click="handleManageScripts"
        >
          <el-icon><Setting /></el-icon>
          {{ t("settings.scripts.manage") }}
        </el-button>
      </div>
    </template>
    <div v-loading="loading" class="script-selector">
      <div v-if="scripts.length === 0" class="empty-state">
        {{ t("scripts.noManualScripts") }}
      </div>
      <div v-else class="script-list">
        <div
          v-for="script in scripts"
          :key="script.ID"
          class="script-item"
          @click="handleSelectScript(script)"
        >
          <div class="script-info">
            <div class="script-name">{{ script.Name }}</div>
          </div>
          <el-icon><ArrowRight /></el-icon>
        </div>
      </div>
    </div>

    <!-- 脚本管理器 -->
    <ScriptManager v-model="showScriptManager" @close="loadScripts" />
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";
import { ElMessage } from "element-plus";
import { ArrowRight, Setting } from "@element-plus/icons-vue";
import { useI18n } from "vue-i18n";
import { GetEnabledUserScriptsByTrigger } from "../../../../wailsjs/go/main/App";
import {
  executeScriptInBrowser,
  shouldTriggerScript,
} from "../../../scripts/executor";
import { common } from "../../../../wailsjs/go/models";
import ScriptManager from "../../setting/components/ScriptManager.vue";

const { t } = useI18n();

const props = defineProps<{
  modelValue: boolean;
  item: common.ClipboardItem | null;
}>();

const emit = defineEmits<{
  "update:modelValue": [value: boolean];
  "script-executed": [itemId: string, result: any];
}>();

const visible = ref(false);
const loading = ref(false);
const scripts = ref<common.UserScript[]>([]);
const showScriptManager = ref(false);

watch(
  () => props.modelValue,
  (val) => {
    visible.value = val;
    if (val) {
      loadScripts();
    }
  }
);

watch(visible, (val) => {
  emit("update:modelValue", val);
});

// 处理对话框关闭
function handleClose() {
  visible.value = false;
}

// 处理管理脚本
function handleManageScripts() {
  showScriptManager.value = true;
}

async function loadScripts() {
  if (!props.item) return;

  loading.value = true;
  try {
    scripts.value = await GetEnabledUserScriptsByTrigger("manual") || [];
  } catch (error: any) {
    ElMessage.error(`加载脚本失败: ${error.message || error}`);
  } finally {
    loading.value = false;
  }
}

async function handleSelectScript(script: common.UserScript) {
  if (!props.item) return;

  // 检查脚本过滤条件
  if (!shouldTriggerScript(script, props.item)) {
    ElMessage.warning(t("scripts.filterNotMatch", { name: script.Name }));
    return;
  }

  try {
    const result = await executeScriptInBrowser(script, props.item);

    // 发送脚本执行结果事件
    emit("script-executed", props.item.ID, {
      ...result,
      scriptName: script.Name,
    });

    ElMessage.success(
      `${script.Name} ${
        result.error ? t("scripts.executeError") : t("scripts.executeSuccess")
      }`
    );
    visible.value = false;
  } catch (error: any) {
    const errorResult = {
      error: error.message || String(error),
      timestamp: Date.now(),
      scriptName: script.Name,
    };
    emit("script-executed", props.item.ID, errorResult);
    ElMessage.error(`${t("scripts.executeError")}: ${error.message || error}`);
  }
}
</script>

<style scoped>
.script-selector {
  /* min-height: 200px; */
}

.empty-state {
  text-align: center;
  padding: 40px;
  color: #999;
}

.script-list {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  gap: 4px;
}

.script-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
  width: calc(100% / 6 - 6px);
}

.script-item:hover {
  border-color: #484848;
  background-color: #f5f7fa;
}

.script-info {
  flex: 1;
}

.script-name {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.script-desc {
  font-size: 12px;
  color: #909399;
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  margin-top: -4px;
}
</style>
