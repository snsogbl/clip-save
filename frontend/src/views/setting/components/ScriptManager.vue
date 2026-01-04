<template>
  <el-dialog
    v-model="visible"
    :title="$t('settings.scripts.title')"
    width="80%"
    :close-on-click-modal="false"
  >
    <div class="script-manager">
      <div class="script-manager-header">
        <el-button
          class="me-button"
          size="small"
          @click="showOnlineScriptList = true"
        >
          {{ $t("settings.scripts.scriptMarket") }}
        </el-button>
        <el-button class="me-button" size="small" @click="handleFindScripts">
          {{ $t("settings.scripts.findScripts") }}
        </el-button>
        <el-button size="small" type="primary" @click="handleNewScript">
          {{ $t("settings.scripts.newScript") }}
        </el-button>
      </div>

      <el-table
        ref="tableRef"
        :data="scripts"
        height="66vh"
        style="width: 100%"
        v-loading="loading"
        row-key="ID"
      >
        <el-table-column width="60" fixed="left">
          <template #default>
            <div class="drag-handle-cell">
              <el-icon class="drag-handle-icon" style="cursor: move; color: #909399">
                <Rank />
              </el-icon>
            </div>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('settings.scripts.order')"
          width="66"
          fixed="left"
        >
          <template #default="{ $index }">
            <span style="font-size: 12px">{{ $index + 1 }}</span>
          </template>
        </el-table-column>
        <el-table-column
          prop="Name"
          :label="$t('settings.scripts.name')"
          width="120"
          show-overflow-tooltip
        />
        <el-table-column
          prop="Description"
          :label="$t('settings.scripts.description')"
          show-overflow-tooltip
        />
        <el-table-column
          prop="Trigger"
          :label="$t('settings.scripts.trigger')"
          width="150"
        >
          <template #default="{ row }">
            <span v-if="row.Trigger === 'after_save'">{{
              $t("settings.scripts.triggerAfterSave")
            }}</span>
            <span v-else-if="row.Trigger === 'manual'">{{
              $t("settings.scripts.triggerManual")
            }}</span>
            <span v-else>{{ row.Trigger }}</span>
          </template>
        </el-table-column>
        <el-table-column
          prop="Enabled"
          :label="$t('settings.scripts.enabled')"
          width="100"
        >
          <template #default="{ row }">
            <el-switch
              v-model="row.Enabled"
              @change="handleEnabledChange(row)"
              :loading="updatingEnabledMap.get(row.ID) || false"
            />
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('common.actions')"
          width="134"
          fixed="right"
        >
          <template #default="{ row }">
            <el-button size="small" @click="handleEditScript(row.ID)">
              {{ $t("common.edit") }}
            </el-button>
            <el-button
              size="small"
              type="danger"
              @click="handleDeleteScript(row.ID, row.Name)"
            >
              {{ $t("common.delete") }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 脚本编辑器 -->
    <ScriptEditor
      v-model="showScriptEditor"
      :script-id="editingScriptId"
      @saved="handleScriptSaved"
    />

    <!-- 在线脚本列表 -->
    <OnlineScriptList
      v-model="showOnlineScriptList"
      @installed="handleScriptSaved"
    />
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, onMounted, onUnmounted } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { Rank } from "@element-plus/icons-vue";
import Sortable from "sortablejs";
import { useI18n } from "vue-i18n";
import {
  GetAllUserScripts,
  DeleteUserScript,
  UpdateUserScriptOrder,
  OpenURL,
  GetUserScriptByID,
  SaveUserScript,
} from "../../../../wailsjs/go/main/App";
import { common } from "../../../../wailsjs/go/models";
import ScriptEditor from "./ScriptEditor.vue";
import OnlineScriptList from "./OnlineScriptList.vue";

const { t } = useI18n();

const props = defineProps<{
  modelValue: boolean;
}>();

const emit = defineEmits<{
  "update:modelValue": [value: boolean];
}>();

const visible = ref(false);
const loading = ref(false);
const scripts = ref<common.UserScript[]>([]);
const showScriptEditor = ref(false);
const editingScriptId = ref<string | undefined>();
const showOnlineScriptList = ref(false);
const tableRef = ref<any>(null);
let sortableInstance: Sortable | null = null;
let scrollContainer: HTMLElement | null = null;
let mouseMoveHandler: ((e: MouseEvent) => void) | null = null;
let rafId: number | null = null;
// 跟踪每个脚本的更新状态
const updatingEnabledMap = ref<Map<string, boolean>>(new Map());
const isFirstLoad = ref(true);
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
  if (!val) {
    // 对话框关闭时销毁 Sortable 实例
    if (sortableInstance) {
      sortableInstance.destroy();
      sortableInstance = null;
    }
    // 对话框关闭时取消待执行的更新
    if (updateDebounceTimer) {
      clearTimeout(updateDebounceTimer);
      updateDebounceTimer = null;
    }
    // 移除鼠标移动监听
    if (mouseMoveHandler) {
      document.removeEventListener("mousemove", mouseMoveHandler, { capture: true } as any);
      mouseMoveHandler = null;
    }
    scrollContainer = null;
  }
});

onUnmounted(() => {
  if (sortableInstance) {
    sortableInstance.destroy();
    sortableInstance = null;
  }
  // 清理防抖定时器
  if (updateDebounceTimer) {
    clearTimeout(updateDebounceTimer);
    updateDebounceTimer = null;
  }
  // 移除鼠标移动监听
  if (mouseMoveHandler) {
    document.removeEventListener("mousemove", mouseMoveHandler, { capture: true } as any);
    mouseMoveHandler = null;
  }
  scrollContainer = null;
});

async function loadScripts() {
  loading.value = true;
  try {
    scripts.value = await GetAllUserScripts();
    if (isFirstLoad.value && (!scripts.value || scripts.value.length === 0)) {
      showOnlineScriptList.value = true;
      isFirstLoad.value = false;
    }
    // 加载完成后初始化拖拽
    await nextTick();
    initSortable();
  } catch (error: any) {
    ElMessage.error(`加载脚本失败: ${error.message || error}`);
  } finally {
    loading.value = false;
  }
}

// 初始化 Sortable
function initSortable() {
  if (!tableRef.value) return;

  // 销毁旧的实例
  if (sortableInstance) {
    sortableInstance.destroy();
    sortableInstance = null;
  }

  // 获取表格的 tbody 元素和滚动容器
  const tbody = tableRef.value.$el.querySelector(
    ".el-table__body-wrapper tbody"
  );
  if (!tbody) return;

  // 获取滚动容器（表格的 body wrapper）
  // Element Plus 表格的滚动容器是 .el-scrollbar__wrap 或 .el-table__body-wrapper
  scrollContainer = tableRef.value.$el.querySelector(
    ".el-scrollbar__wrap"
  ) as HTMLElement || tableRef.value.$el.querySelector(
    ".el-table__body-wrapper"
  ) as HTMLElement;

  if (!scrollContainer) {
    console.warn("无法找到滚动容器，尝试查找所有可能的滚动元素");
    // 尝试查找所有可能的滚动容器
    const possibleContainers = tableRef.value.$el.querySelectorAll(
      ".el-table__body-wrapper, .el-scrollbar__wrap, .el-table__body"
    );
    if (possibleContainers.length > 0) {
      // 找到第一个有滚动条的元素
      for (let i = 0; i < possibleContainers.length; i++) {
        const el = possibleContainers[i] as HTMLElement;
        if (el.scrollHeight > el.clientHeight) {
          scrollContainer = el;
          break;
        }
      }
    }
  }

  if (!scrollContainer) {
    console.error("无法找到滚动容器");
    return;
  }

  console.log("滚动容器找到:", scrollContainer, {
    scrollHeight: scrollContainer.scrollHeight,
    clientHeight: scrollContainer.clientHeight,
    scrollTop: scrollContainer.scrollTop
  });

  // 创建鼠标移动处理函数
  const container = scrollContainer; // 保存引用，避免闭包问题
  mouseMoveHandler = (e: MouseEvent) => {
    if (!container) return;
    
    // 取消之前的动画帧
    if (rafId !== null) {
      cancelAnimationFrame(rafId);
    }
    
    // 使用 requestAnimationFrame 确保流畅滚动
    rafId = requestAnimationFrame(() => {
      if (!container) {
        rafId = null;
        return;
      }
      
      const rect = container.getBoundingClientRect();
      const mouseY = e.clientY;
      const scrollThreshold = 80; // 滚动触发区域
      const scrollSpeed = 30; // 滚动速度（增大）
      
      // 检查鼠标是否在滚动容器内或附近
      const isNearContainer = mouseY >= rect.top - 50 && mouseY <= rect.bottom + 50;
      if (!isNearContainer) {
        rafId = null;
        return;
      }
      
      // 检查是否可以滚动
      const canScroll = container.scrollHeight > container.clientHeight;
      if (!canScroll) {
        rafId = null;
        return;
      }
      
      // 检查是否接近顶部
      const distanceFromTop = mouseY - rect.top;
      if (distanceFromTop < scrollThreshold && distanceFromTop > -50) {
        const newScrollTop = container.scrollTop - scrollSpeed;
        container.scrollTop = Math.max(0, newScrollTop);
      }
      // 检查是否接近底部
      const distanceFromBottom = rect.bottom - mouseY;
      if (distanceFromBottom < scrollThreshold && distanceFromBottom > -50) {
        const maxScroll = container.scrollHeight - container.clientHeight;
        const newScrollTop = container.scrollTop + scrollSpeed;
        container.scrollTop = Math.min(maxScroll, newScrollTop);
      }
      
      rafId = null;
    });
  };

  sortableInstance = Sortable.create(tbody, {
    handle: ".drag-handle-cell", // 指定拖拽手柄（整个单元格）
    animation: 150,
    ghostClass: "sortable-ghost",
    chosenClass: "sortable-chosen",
    dragClass: "sortable-drag",
    forceFallback: false,
    // 滚动配置（作为备选）
    scroll: scrollContainer || true,
    scrollSensitivity: 80,
    scrollSpeed: 20,
    bubbleScroll: true,
    // 开始拖动时添加鼠标移动监听
    onStart: () => {
      if (mouseMoveHandler && scrollContainer) {
        // 使用 capture 模式确保能捕获到事件
        document.addEventListener("mousemove", mouseMoveHandler, { capture: true, passive: true });
      }
    },
    // 结束拖动时移除鼠标移动监听
    onEnd: async (evt) => {
      // 移除鼠标移动监听
      if (mouseMoveHandler) {
        document.removeEventListener("mousemove", mouseMoveHandler, { capture: true } as any);
      }
      // 取消待执行的动画帧
      if (rafId !== null) {
        cancelAnimationFrame(rafId);
        rafId = null;
      }
      
      const { oldIndex, newIndex } = evt;
      if (
        oldIndex === undefined ||
        newIndex === undefined ||
        oldIndex === newIndex
      ) {
        return;
      }

      // 检查索引是否有效
      if (
        oldIndex < 0 ||
        oldIndex >= scripts.value.length ||
        newIndex < 0 ||
        newIndex >= scripts.value.length
      ) {
        console.error("拖拽索引无效:", {
          oldIndex,
          newIndex,
          length: scripts.value.length,
        });
        return;
      }

      // 更新本地数组顺序
      const movedScript = scripts.value.splice(oldIndex, 1)[0];
      if (!movedScript) {
        console.error("无法获取被移动的脚本");
        return;
      }
      scripts.value.splice(newIndex, 0, movedScript);

      // 批量更新所有脚本的顺序
      await updateAllScriptsOrder();
    },
  });
}

// 防抖定时器
let updateDebounceTimer: ReturnType<typeof setTimeout> | null = null;

// 批量更新所有脚本的顺序（使用防抖优化，只更新需要更新的脚本）
async function updateAllScriptsOrder() {
  // 清除之前的定时器
  if (updateDebounceTimer) {
    clearTimeout(updateDebounceTimer);
  }

  // 设置新的防抖定时器（300ms 防抖，避免频繁更新）
  updateDebounceTimer = setTimeout(async () => {
    // 检查对话框是否仍然打开
    if (!visible.value) {
      updateDebounceTimer = null;
      return;
    }

    try {
      // 构建需要更新的列表（只包含顺序改变的脚本）
      const updates: Array<{ id: string; order: number }> = [];
      scripts.value.forEach((script, index) => {
        const newSortOrder = index + 1;
        // 确保 SortOrder 是数字类型进行比较
        const currentOrder = Number(script.SortOrder) || 0;
        if (currentOrder !== newSortOrder) {
          script.SortOrder = newSortOrder;
          updates.push({ id: script.ID, order: newSortOrder });
        }
      });

      // 如果没有需要更新的，直接返回
      if (updates.length === 0) {
        updateDebounceTimer = null;
        return;
      }

      // 并行更新所有需要更新的脚本（使用 Promise.all 提高性能）
      const updatePromises = updates.map(({ id, order }) =>
        UpdateUserScriptOrder(id, order)
      );
      await Promise.all(updatePromises);

      // 再次检查对话框是否仍然打开（避免在关闭后显示消息）
      if (visible.value) {
        ElMessage.success(t("settings.scripts.orderUpdated") || "顺序已更新");
      }
    } catch (error: any) {
      // 再次检查对话框是否仍然打开
      if (visible.value) {
        ElMessage.error(
          `${t("settings.scripts.orderUpdateError") || "更新顺序失败"}: ${
            error.message || error
          }`
        );
        // 重新加载以恢复原始顺序
        await loadScripts();
      }
    } finally {
      updateDebounceTimer = null;
    }
  }, 300); // 300ms 防抖
}

function handleNewScript() {
  editingScriptId.value = undefined;
  showScriptEditor.value = true;
}

function handleFindScripts() {
  OpenURL(
    "https://github.com/snsogbl/clip-save/blob/main/scriptingExample/README.md"
  );
}

function handleEditScript(id: string) {
  editingScriptId.value = id;
  showScriptEditor.value = true;
}

async function handleDeleteScript(id: string, name: string) {
  try {
    await ElMessageBox.confirm(
      t("settings.scripts.deleteConfirm", { name }) ||
        `确定要删除脚本 "${name}" 吗？`,
      t("settings.scripts.deleteTitle") || "删除脚本",
      {
        confirmButtonText: t("common.delete") || "删除",
        cancelButtonText: t("common.cancel") || "取消",
        type: "warning",
      }
    );

    await DeleteUserScript(id);
    ElMessage.success(t("settings.scripts.deleteSuccess") || "脚本已删除");
    await loadScripts();
  } catch (error: any) {
    if (error !== "cancel") {
      ElMessage.error(
        `${t("settings.scripts.deleteError") || "删除脚本失败"}: ${
          error.message || error
        }`
      );
    }
  }
}

function handleScriptSaved() {
  loadScripts();
}

async function handleEnabledChange(row: common.UserScript) {
  // 设置更新状态，防止重复点击
  updatingEnabledMap.value.set(row.ID, true);

  try {
    // 获取完整的脚本数据
    const fullScript = await GetUserScriptByID(row.ID);
    if (!fullScript) {
      throw new Error("脚本不存在");
    }

    // 更新启用状态
    fullScript.Enabled = row.Enabled;

    // 保存脚本
    const scriptJSON = JSON.stringify(fullScript);
    await SaveUserScript(scriptJSON);

    ElMessage.success(
      row.Enabled
        ? t("settings.scripts.enabledSuccess") || "脚本已启用"
        : t("settings.scripts.disabledSuccess") || "脚本已禁用"
    );
  } catch (error: any) {
    // 恢复原状态
    row.Enabled = !row.Enabled;
    ElMessage.error(
      `${t("settings.scripts.enabledUpdateError") || "更新启用状态失败"}: ${
        error.message || error
      }`
    );
  } finally {
    updatingEnabledMap.value.set(row.ID, false);
  }
}
</script>

<style scoped>
.script-manager {
  min-height: 72vh;
}

.script-manager-header {
  margin-bottom: 10px;
  display: flex;
  justify-content: flex-end;
}

.drag-handle-cell {
  cursor: move;
  user-select: none;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  padding: 8px;
  transition: background-color 0.2s;
}

.drag-handle-cell:hover {
  background-color: #f5f7fa;
}

.drag-handle-icon {
  transition: color 0.2s;
}

.drag-handle-cell:hover .drag-handle-icon {
  color: #409eff !important;
}

/* Sortable 拖拽样式 */
:deep(.sortable-ghost) {
  opacity: 0.4;
  background-color: #f0f0f0;
}

:deep(.sortable-chosen) {
  background-color: #e6f7ff;
}

:deep(.sortable-drag) {
  opacity: 0.8;
}
</style>
