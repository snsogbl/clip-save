<template>
  <!-- 设置页面 Drawer -->
  <el-drawer
    :model-value="showSetting"
    @update:model-value="showSetting = $event"
    :title="$t('settings.title')"
    direction="rtl"
    :size="showRightPanel ? '600px' : '100%'"
    @close="handleSettingBack"
    class="settings-drawer"
    destroy-on-close
  >
    <SettingView />
  </el-drawer>

  <!-- 脚本选择器 -->
  <ScriptSelector
    v-if="currentItem"
    :model-value="showScriptSelector"
    :show-right-panel="showRightPanel"
    @update:model-value="showScriptSelector = $event"
    :item="currentItem"
  />

  <div
    ref="containerRef"
    class="clipboard-container"
    style="--wails-draggable: no-drag"
  >
    <!-- 极简模式 -->
    <MinimalModeView
      v-if="!showRightPanel"
      :items="items"
      :current-item="currentItem"
      :loading="loading"
      :left-tab="leftTab"
      :search-keyword="searchKeyword"
      :is-always-on-top="isAlwaysOnTop"
      :is-command-pressed="isCommandPressed"
      @toggle-always-on-top="toggleAlwaysOnTop"
      @show-setting="showSetting = true"
      @switch-tab="switchLeftTab"
      @update:search-keyword="searchKeyword = $event"
      @search-keydown="handleSearchKeydown"
      @search-change="onSearchChange"
      @select-item="selectItem"
      @double-click="handleDoubleClick"
      @copy-item="copyItem"
      @delete-item="deleteItem"
      @collect-item="collectItem"
      @send-item="sendItem"
      ref="minimalModeRef"
    />

    <!-- 正常模式 -->
    <NormalModeView
      v-else
      :items="items"
      :current-item="currentItem"
      :loading="loading"
      :left-tab="leftTab"
      :search-keyword="searchKeyword"
      @update:search-keyword="searchKeyword = $event"
      :filter-type="filterType"
      :is-always-on-top="isAlwaysOnTop"
      :is-command-pressed="isCommandPressed"
      @toggle-always-on-top="toggleAlwaysOnTop"
      @show-setting="showSetting = true"
      @change-language="changeLanguage"
      @filter-change="handleFilterChange"
      @search-keydown="handleSearchKeydown"
      @search-change="onSearchChange"
      @switch-tab="switchLeftTab"
      @select-item="selectItem"
      @double-click="handleDoubleClick"
      ref="normalModeRef"
    >
      <template #content-area>
        <div class="content-area" ref="contentAreaRef">
          <ClipboardTitleView
            v-if="currentItem"
            :item="currentItem"
            @copy="copyItem"
            @delete="deleteItem"
            @collect="collectItem"
            @run-script="handleRunScript"
          />
          <div class="content-display">
            <div v-if="!currentItem" class="welcome-text">
              {{ $t("main.welcome") }}
            </div>
            <!-- 图片内容展示 -->
            <ClipboardImageView
              v-else-if="
                currentItem.ContentType === 'Image' && currentItem.ImageData
              "
              ref="imageViewRef"
              :imageData="
                Array.isArray(currentItem.ImageData)
                  ? currentItem.ImageData.map((b) =>
                      String.fromCharCode(b)
                    ).join('')
                  : String(currentItem.ImageData)
              "
              :ocrText="currentItem.OCRText"
              :itemId="currentItem.ID"
              :playing="isSayPlaying"
            />
            <!-- 文件内容展示 -->
            <ClipboardFileView
              v-else-if="currentItem.ContentType === 'File'"
              :files="parseFileInfo(currentItem)"
              @open-file="openInFinder"
            />
            <!-- URL 内容展示 -->
            <ClipboardUrlView
              v-else-if="currentItem.ContentType === 'URL'"
              :url="currentItem.Content"
              @open-url="openURL"
            />
            <!-- 颜色内容展示 -->
            <ClipboardColorView
              v-else-if="currentItem.ContentType === 'Color'"
              :color="currentItem.Content"
            />
            <!-- JSON 内容展示/编辑 -->
            <ClipboardJsonView
              ref="jsonEditorRef"
              v-else-if="currentItem.ContentType === 'JSON'"
              :text="currentItem?.Content || '{}'"
            />
            <!-- 文本内容展示 -->
            <ClipboardTextView
              v-else
              ref="textEditorRef"
              :text="currentItem?.Content || '空内容'"
              :playing="isSayPlaying"
            />
          </div>
        </div>
        <!-- 脚本执行结果 -->
        <ScriptResultView
          v-if="currentItem && scriptResults[currentItem.ID]"
          :item-id="currentItem.ID"
          :result="scriptResults[currentItem.ID]"
        />
      </template>

      <template #info-panel>
        <div v-if="currentItem" class="info-panel">
          <el-descriptions title="">
            <el-descriptions-item :label="$t('main.source')">
              {{ currentItem.Source }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('main.contentType')">
              {{ currentItem.ContentType }}
            </el-descriptions-item>
            <template v-if="currentItem.ContentType === 'File'">
              <el-descriptions-item :label="$t('main.fileCount')">
                {{ currentItem.WordCount }}
              </el-descriptions-item>
            </template>
            <template v-else>
              <el-descriptions-item :label="$t('main.charCount')">
                {{ currentItem.CharCount }}
              </el-descriptions-item>
              <el-descriptions-item :label="$t('main.wordCount')">
                {{ currentItem.WordCount }}
              </el-descriptions-item>
            </template>
            <el-descriptions-item :label="$t('main.createTime')">
              {{ new Date(currentItem.Timestamp).toLocaleString("zh-CN") }}
            </el-descriptions-item>
          </el-descriptions>
        </div>
      </template>
    </NormalModeView>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, onUnmounted, nextTick, computed, watch } from "vue";
import { EventsOn, WindowGetSize } from "../../../wailsjs/runtime/runtime";
import { useI18n } from "vue-i18n";
import { useCommandNumberShortcut } from "../../composables/useCommandNumberShortcut";
import {
  SearchClipboardItems,
  CopyToClipboard,
  DeleteClipboardItem,
  OpenFileInFinder,
  OpenURL,
  ClearItemsOlderThanDays,
  GetAppSettings,
  ToggleFavorite,
  HideWindowAndQuit,
  SetLanguage,
  AutoPasteCurrentItem,
  GetClipboardItemByID,
  AutoPasteCurrentItemToPreviousApp,
  SetWindowAlwaysOnTop,
  ActivatePreviousApp,
  StopSay,
  IsSayPlaying,
} from "../../../wailsjs/go/main/App";

// 组件导入
import MinimalModeView from "./components/MinimalModeView.vue";
import NormalModeView from "./components/NormalModeView.vue";
import ClipboardUrlView from "./components/clipboardUrlView.vue";
import ClipboardColorView from "./components/clipboardColorView.vue";
import ClipboardFileView from "./components/clipboardFileView.vue";
import ClipboardTextView from "./components/clipboardTextView.vue";
import ClipboardImageView from "./components/clipboardImageView.vue";
import ClipboardJsonView from "./components/clipboardJsonView.vue";
import ClipboardTitleView from "./components/clipboardTitleView.vue";
import ScriptResultView from "./components/ScriptResultView.vue";
import ScriptSelector from "./components/ScriptSelector.vue";
import SettingView from "../setting/setting.vue";
import { ElMessageBox, ElMessage } from "element-plus";
import { common } from "../../../wailsjs/go/models";

const { t, locale } = useI18n();

// 使用 Wails 生成的类型
type ClipboardItem = common.ClipboardItem;

interface FileInfo {
  name: string;
  path: string;
  size: number;
  is_dir: boolean;
  exists: boolean;
  extension: string;
}

// 响应式数据
const items = ref<ClipboardItem[]>([]);
const currentItem = ref<ClipboardItem | null>(null);
const containerRef = ref<HTMLElement | null>(null);
const contentAreaRef = ref<HTMLElement | null>(null);
const textEditorRef = ref<InstanceType<typeof ClipboardTextView> | null>(null);
const imageViewRef = ref<InstanceType<typeof ClipboardImageView> | null>(null);
const jsonEditorRef = ref<InstanceType<typeof ClipboardJsonView> | null>(null);
const minimalModeRef = ref<InstanceType<typeof MinimalModeView> | null>(null);
const normalModeRef = ref<InstanceType<typeof NormalModeView> | null>(null);

const searchKeyword = ref("");
const filterType = ref("");
const loading = ref(false);
const showSetting = ref(false);
const showScriptSelector = ref(false);
const leftTab = ref<"all" | "fav">("all");
const isSayPlaying = ref(false);

// 脚本执行结果存储
interface ScriptExecutionResult {
  error?: string;
  returnValue?: any;
  timestamp: number;
  scriptName?: string;
  status?: "executing" | "completed" | "error" | "streaming";
}

const scriptResults = ref<Record<string, ScriptExecutionResult>>({});
const executingScripts = ref<Record<string, string>>({});
// 流式内容存储（用于实时更新）
const streamContents = ref<Record<string, string>>({});

// 窗口大小相关
const windowWidth = ref(1280);
const showRightPanel = computed(() => windowWidth.value >= 800);

// 窗口置顶状态
const isAlwaysOnTop = ref(false);

// 使用 Command+数字键快捷键
const { isCommandPressed } = useCommandNumberShortcut({
  enabled: () => !showScriptSelector.value,
  itemCount: () => items.value.length,
  onSelect: (index: number) => {
    if (items.value[index]) {
      autoPasteCurrentItem(items.value[index]);
    }
  },
});

// 缓存的设置数据
let cachedSettings: {
  pageSize: number;
  autoClean: boolean;
  retentionDays: number;
  doubleClickPaste?: boolean;
} | null = null;

// 定时器和事件清理
let autoCleanInterval: ReturnType<typeof setInterval> | null = null;
let resizeObserver: ResizeObserver | null = null;
const eventCleanupFunctions: (() => void)[] = [];
const isFirstWatch = ref(true);

// 监听showRightPanel变化
watch(showRightPanel, (newVal) => {
  if (isFirstWatch.value) {
    isFirstWatch.value = false;
    return;
  }
  loadItems();
});

// 获取当前搜索输入框引用
const getSearchInputRef = () => {
  if (showRightPanel.value) {
    return normalModeRef.value?.searchInputRef;
  } else {
    return minimalModeRef.value?.searchInputRef;
  }
};

// 获取当前列表容器引用
const getItemListRef = () => {
  if (showRightPanel.value) {
    return normalModeRef.value?.itemListRef;
  } else {
    return minimalModeRef.value?.itemListRef;
  }
};

// 切换窗口置顶状态
function toggleAlwaysOnTop() {
  isAlwaysOnTop.value = !isAlwaysOnTop.value;
  SetWindowAlwaysOnTop(isAlwaysOnTop.value);
  ElMessage.success({
    placement: showRightPanel.value ? "top-right" : "top",
    message: isAlwaysOnTop.value
      ? t("message.windowTopped")
      : t("message.windowUnTopped"),
  });
}

// 从数据库获取设置（带缓存）
async function getSettings(forceRefresh = false) {
  if (cachedSettings && !forceRefresh) {
    return cachedSettings;
  }

  try {
    const savedSettings = await GetAppSettings();
    if (savedSettings) {
      const parsed = JSON.parse(savedSettings);
      cachedSettings = {
        pageSize: parsed.pageSize || 50,
        autoClean: parsed.autoClean !== undefined ? parsed.autoClean : true,
        retentionDays: parsed.retentionDays || 30,
        doubleClickPaste:
          parsed.doubleClickPaste !== undefined
            ? parsed.doubleClickPaste
            : true,
      };
      return cachedSettings;
    }
  } catch (e) {
    console.error("❌ 读取设置失败:", e);
  }

  cachedSettings = {
    pageSize: 50,
    autoClean: true,
    retentionDays: 30,
    doubleClickPaste: true,
  };
  return cachedSettings;
}

// 加载剪贴板项目
async function loadItems() {
  try {
    loading.value = true;
    const settings = await getSettings();
    const pageSize = settings?.pageSize || 50;

    const result = await SearchClipboardItems(
      leftTab.value === "fav",
      searchKeyword.value,
      filterType.value,
      pageSize,
      !showRightPanel.value
    );
    items.value = result || [];

    if (items.value.length > 0) {
      selectItem(items.value[0]);
    } else {
      currentItem.value = null;
    }
  } catch (error) {
    console.error("加载剪贴板项目失败:", error);
  } finally {
    loading.value = false;
  }
}

// 静默检查更新
async function checkForUpdates() {
  try {
    const settings = await getSettings();
    const pageSize = settings?.pageSize || 50;

    const result = await SearchClipboardItems(
      leftTab.value === "fav",
      searchKeyword.value,
      filterType.value,
      pageSize,
      !showRightPanel.value
    );
    const newItems = result || [];

    if (
      newItems.length !== items.value.length ||
      (newItems.length > 0 &&
        items.value.length > 0 &&
        newItems[0].ID !== items.value[0].ID)
    ) {
      items.value = newItems;

      if (!currentItem.value && newItems.length > 0) {
        selectItem(newItems[0]);
      }
    }
  } catch (error) {
    console.error("检查更新失败:", error);
  }
}

// 选择项目
async function selectItem(item: ClipboardItem) {
  if (
    currentItem.value?.ContentType === "Image" &&
    currentItem.value.ImageData &&
    showRightPanel.value
  ) {
    if (currentItem.value.ID !== item.ID) {
      currentItem.value.ImageData = null as any;
    }
  }

  if (currentItem.value && currentItem.value.ID !== item.ID) {
    delete scriptResults.value[currentItem.value.ID];
  }

  if (item.ContentType === "Image" && !item.ImageData && showRightPanel.value) {
    try {
      const fullItem = await GetClipboardItemByID(item.ID);
      if (fullItem) {
        currentItem.value = fullItem;
      } else {
        currentItem.value = item;
      }
    } catch (error) {
      console.error("加载图片数据失败:", error);
      currentItem.value = item;
    }
  } else {
    currentItem.value = item;
  }

  await nextTick();
  const container = getItemListRef();
  if (!container) return;
  const activeEl = container.querySelector(
    ".list-item.active, .minimal-list-item.active, .normal-list-item.active"
  ) as HTMLElement | null;
  if (activeEl) {
    activeEl.scrollIntoView({ block: "nearest" });
  }

  if (contentAreaRef.value) {
    contentAreaRef.value.scrollTo({ top: 0, behavior: "smooth" });
  }
}

// 发送项目
async function sendItem(item: ClipboardItem) {
  if (currentItem.value?.ID !== item.ID) {
    await selectItem(item);
    await nextTick();
  }

  await copyItem(item.ID);

  if (isAlwaysOnTop.value) {
    AutoPasteCurrentItemToPreviousApp();
  } else {
    ActivatePreviousApp();
    AutoPasteCurrentItem();
  }
}

// 处理双击事件
async function handleDoubleClick(item: ClipboardItem) {
  const settings = await getSettings();
  if (settings?.doubleClickPaste === false) {
    return;
  }
  autoPasteCurrentItem(item);
}

// 自动粘贴当前项目
async function autoPasteCurrentItem(item: ClipboardItem) {
  if (!currentItem.value) return;
  if (currentItem.value?.ID !== item.ID) {
    await selectItem(item);
    await nextTick();
  }

  await copyItem(item.ID);

  if (isAlwaysOnTop.value) {
    AutoPasteCurrentItemToPreviousApp();
  } else {
    // 是否mac系统
    const isMac = navigator.platform.toUpperCase().indexOf('MAC') >= 0
    if (isMac){
      HideWindowAndQuit();
    }else{
      ActivatePreviousApp();
    }
    AutoPasteCurrentItem();
  }
}

// 复制项目
async function copyItem(id: string) {
  if (currentItem.value?.ContentType === "JSON") {
    jsonEditorRef.value?.copyEdited();
  } else {
    try {
      await CopyToClipboard(id);
      ElMessage.success(t("message.copySuccess"));
    } catch (error) {
      console.error("复制失败:", error);
      ElMessage.error(t("message.copyError", [error]));
    }
  }
}

// 删除项目
async function deleteItem(id: string) {
  ElMessageBox.confirm(
    t("message.deleteConfirm"),
    t("message.deleteConfirmTitle"),
    {
      confirmButtonText: t("message.deleteConfirmBtn"),
      cancelButtonText: t("message.deleteCancelBtn"),
      type: "warning",
    }
  ).then(async () => {
    try {
      await DeleteClipboardItem(id);
      const index = items.value.findIndex((item) => item.ID === id);
      items.value.splice(index, 1);
      currentItem.value = items.value[index] || items.value[index - 1] || null;
      ElMessage.success(t("message.deleteSuccess"));
    } catch (error) {
      console.error("删除失败:", error);
      ElMessage.error(t("message.deleteError", [error]));
    }
  });
}

// 收藏
async function collectItem(id: string) {
  try {
    const newVal = await ToggleFavorite(id);
    if (currentItem.value && currentItem.value.ID === id) {
      currentItem.value.IsFavorite = newVal;
    }

    const index = items.value.findIndex((i) => i.ID === id);
    if (index !== -1) {
      if (leftTab.value === "fav" && newVal === 0) {
        const isCurrent = currentItem.value?.ID === id;
        const nextItem =
          items.value[index + 1] || items.value[index - 1] || null;
        items.value.splice(index, 1);
        if (isCurrent) {
          if (nextItem) {
            selectItem(nextItem);
          } else {
            currentItem.value = null;
          }
        }
      } else {
        items.value[index].IsFavorite = newVal;
      }
    }
    ElMessage.success(
      newVal === 1 ? t("message.favoriteAdded") : t("message.favoriteRemoved")
    );
  } catch (error) {
    console.error("收藏失败:", error);
    ElMessage.error(t("message.favoriteError"));
  }
}

// 切换左侧标签页
async function switchLeftTab(tab: "all" | "fav") {
  if (leftTab.value === tab) return;
  leftTab.value = tab;
  await loadItems();
  await nextTick();
  getItemListRef()?.focus();
}

// 处理运行脚本按钮点击
function handleRunScript() {
  if (currentItem.value) {
    showScriptSelector.value = true;
  }
}

// 搜索和过滤变化时重新加载
const onSearchChange = () => {
  loadItems();
};

// 处理过滤器变化
const handleFilterChange = (type: string) => {
  filterType.value = type;
  onSearchChange();
};

// 处理搜索框键盘按下事件
function handleSearchKeydown(event: KeyboardEvent) {
  if (event.shiftKey) return;

  if ((event.metaKey || event.ctrlKey) && event.key === "Backspace") {
    event.preventDefault();
    event.stopPropagation();
    if (currentItem.value) {
      deleteItem(currentItem.value.ID);
    }
    return;
  }

  if ((event.metaKey || event.ctrlKey) && event.key === "Enter") {
    event.preventDefault();
    event.stopPropagation();
    if (currentItem.value) {
      autoPasteCurrentItem(currentItem.value);
    }
    return;
  }

  if ((event.metaKey || event.ctrlKey) && event.key === "ArrowLeft") {
    event.preventDefault();
    event.stopPropagation();
    switchLeftTab("all").then(() => {
      nextTick(() => {
        getSearchInputRef()?.focus();
      });
    });
    return;
  }

  if ((event.metaKey || event.ctrlKey) && event.key === "ArrowRight") {
    event.preventDefault();
    event.stopPropagation();
    switchLeftTab("fav").then(() => {
      nextTick(() => {
        getSearchInputRef()?.focus();
      });
    });
    return;
  }
}

// 解析文件信息
function parseFileInfo(item: ClipboardItem): FileInfo[] {
  if (!item.FileInfo) return [];
  try {
    return JSON.parse(item.FileInfo);
  } catch (e) {
    console.error("解析文件信息失败:", e);
    return [];
  }
}

// 在 Finder 中打开文件
async function openInFinder(filePath: string) {
  try {
    await OpenFileInFinder(filePath);
    ElMessage.success(t("message.openFinderSuccess"));
  } catch (error) {
    console.error("在 Finder 中打开文件失败:", error);
    ElMessage.error(t("message.openFinderError", [error]));
  }
}

// 在浏览器中打开 URL
async function openURL(url: string) {
  try {
    await OpenURL(url);
    ElMessage.success(t("message.openUrlSuccess"));
  } catch (error) {
    console.error("在浏览器中打开 URL 失败:", error);
    ElMessage.error(t("message.openUrlError", [error]));
  }
}

// 处理设置页面返回
async function handleSettingBack() {
  showSetting.value = false;
  await getSettings(true);
}

// 自动清理超过指定天数的历史记录
async function autoCleanOldItems() {
  const settings = await getSettings();

  if (!settings?.autoClean) {
    return;
  }

  const retentionDays = settings.retentionDays || 30;

  try {
    await ClearItemsOlderThanDays(retentionDays);
  } catch (error) {
    console.error("❌ 自动清理失败:", error);
  }
}

// 检查并更新窗口大小
async function updateWindowSize() {
  try {
    const size = await WindowGetSize();
    if (size) {
      windowWidth.value = size.w;
    }
  } catch (error) {
    console.error("获取窗口大小失败:", error);
  }
}

// 切换语言
function changeLanguage(lang: string) {
  SetLanguage(lang);
  locale.value = lang as any;
}

  // 初始化和定时刷新
onMounted(async () => {
  getSettings().then(() => {
    loadItems();
    autoCleanOldItems();
    nextTick(() => {
      isFirstWatch.value = false;
    });
  });

  // 监听剪贴板更新事件
  eventCleanupFunctions.push(
    EventsOn("clipboard.updated", () => {
      checkForUpdates();
    })
  );

  // 统一监听播放状态事件
  eventCleanupFunctions.push(
    EventsOn("say.started", () => {
      isSayPlaying.value = true;
    })
  );

  eventCleanupFunctions.push(
    EventsOn("say.finished", () => {
      isSayPlaying.value = false;
    })
  );

  eventCleanupFunctions.push(
    EventsOn("say.stopped", () => {
      isSayPlaying.value = false;
    })
  );

  // 每小时执行一次自动清理
  autoCleanInterval = setInterval(() => {
    autoCleanOldItems();
  }, 60 * 60 * 1000);

  // 初始化窗口大小
  updateWindowSize();

  // 使用 ResizeObserver 监听容器大小变化
  nextTick(() => {
    if (containerRef.value) {
      resizeObserver = new ResizeObserver((entries) => {
        for (const entry of entries) {
          const width = entry.contentRect.width;
          windowWidth.value = width;
        }
      });
      resizeObserver.observe(containerRef.value);

      eventCleanupFunctions.push(() => {
        if (resizeObserver) {
          resizeObserver.disconnect();
          resizeObserver = null;
        }
      });
    }
  });

  // 监听各种事件
  eventCleanupFunctions.push(
    EventsOn("window.show", () => {
      isCommandPressed.value = false;
      setTimeout(() => {
        checkForUpdates();
        updateWindowSize();
        if (items.value.length > 0) {
          selectItem(items.value[0]);
        }
        getSearchInputRef()?.focus();
      }, 100);
    })
  );

  // 其他事件监听...
  eventCleanupFunctions.push(
    EventsOn("nav.prev", () => {
      if (items.value.length === 0) return;
      if (!currentItem.value) {
        selectItem(items.value[0]);
        return;
      }
      const idx = items.value.findIndex((i) => i.ID === currentItem.value!.ID);
      const nextIdx = Math.max(0, idx - 1);
      selectItem(items.value[nextIdx]);
    })
  );

  eventCleanupFunctions.push(
    EventsOn("nav.next", () => {
      if (items.value.length === 0) return;
      if (!currentItem.value) {
        selectItem(items.value[0]);
        return;
      }
      const idx = items.value.findIndex((i) => i.ID === currentItem.value!.ID);
      const nextIdx = Math.min(items.value.length - 1, idx + 1);
      selectItem(items.value[nextIdx]);
    })
  );

  eventCleanupFunctions.push(
    EventsOn("nav.switch", (tab: "all" | "fav") => {
      switchLeftTab(tab);
    })
  );

  eventCleanupFunctions.push(
    EventsOn("nav.setting", () => {
      showSetting.value = true;
    })
  );

  eventCleanupFunctions.push(
    EventsOn("nav.runScript", () => {
      if (currentItem.value) {
        showScriptSelector.value = true;
      }
    })
  );

  eventCleanupFunctions.push(
    EventsOn("copy.current", () => {
      if (currentItem.value) {
        copyItem(currentItem.value.ID);
      }
    })
  );

  eventCleanupFunctions.push(
    EventsOn("delete.current", () => {
      if (currentItem.value) {
        deleteItem(currentItem.value.ID);
      }
    })
  );

  eventCleanupFunctions.push(
    EventsOn("collect.current", () => {
      if (currentItem.value) {
        collectItem(currentItem.value.ID);
      }
    })
  );

  eventCleanupFunctions.push(
    EventsOn("search.item", () => {
      getSearchInputRef()?.focus();
    })
  );

  eventCleanupFunctions.push(
    EventsOn("enter.item", () => {
      autoPasteCurrentItem(currentItem.value || {} as ClipboardItem);
    })
  );

  eventCleanupFunctions.push(
    EventsOn("translate.current", () => {
      textEditorRef.value?.translateText();
    })
  );

  eventCleanupFunctions.push(
    EventsOn("play.current", () => {
      if (!currentItem.value) return;
      
      // 根据内容类型播放
      if (currentItem.value.ContentType === "Image") {
        // 图片类型：播放 OCR 文字
        imageViewRef.value?.playOCRText();
      } else if (currentItem.value.ContentType === "Text") {
        // 文本类型：播放文本内容
        textEditorRef.value?.playText();
      } else {
        ElMessage.error(t("components.text.playEmpty"));
      }
    })
  );

  // 监听脚本执行事件
  eventCleanupFunctions.push(
    EventsOn(
      "script.executing",
      (data: { itemId: string; scriptName: string; scriptId: string }) => {
        const { itemId, scriptName, scriptId } = data;
        executingScripts.value[itemId] = scriptId;
        scriptResults.value[itemId] = {
          scriptName,
          timestamp: Date.now(),
          status: "executing",
        };
      }
    )
  );

  eventCleanupFunctions.push(
    EventsOn(
      "script.executed",
      (data: {
        itemId: string;
        scriptName: string;
        result: ScriptExecutionResult;
      }) => {
        const { itemId, result } = data;
        delete executingScripts.value[itemId];
        // 如果有流式内容，使用流式内容作为返回值
        if (streamContents.value[itemId]) {
          scriptResults.value[itemId] = {
            ...result,
            returnValue: streamContents.value[itemId],
            status: (result.error ? "error" : "completed") as
              | "error"
              | "completed",
            timestamp: result.timestamp || Date.now(),
          };
          // 清理流式内容
          delete streamContents.value[itemId];
        } else {
        scriptResults.value[itemId] = {
          ...result,
          status: (result.error ? "error" : "completed") as
            | "error"
            | "completed",
          timestamp: result.timestamp || Date.now(),
        };
      }
      }
    )
  );

  // 监听浏览器原生流式事件（性能最好，无 IPC 开销）
  const handleStreamStart = (e: CustomEvent) => {
    const { itemId } = e.detail;
    // 初始化流式内容
    streamContents.value[itemId] = "";
    // 更新状态为流式
    if (scriptResults.value[itemId]) {
      scriptResults.value[itemId].status = "streaming";
    }
  };

  const handleStreamChunk = (e: CustomEvent) => {
    const { itemId, chunk } = e.detail;
    // 实时追加数据块
    if (streamContents.value[itemId] !== undefined) {
      streamContents.value[itemId] += chunk;
      // 实时更新 scriptResults 的 returnValue，触发视图更新
      if (scriptResults.value[itemId]) {
        scriptResults.value[itemId].returnValue = streamContents.value[itemId];
        scriptResults.value[itemId].status = "streaming";
      }
    }
  };

  const handleStreamEnd = (e: CustomEvent) => {
    const { itemId } = e.detail;
    // 流结束，状态改为 completed
    if (scriptResults.value[itemId] && streamContents.value[itemId]) {
      scriptResults.value[itemId].status = "completed";
      scriptResults.value[itemId].returnValue = streamContents.value[itemId];
    }
    // 清理流式内容（保留在 scriptResults 中）
    delete streamContents.value[itemId];
  };

  // 添加浏览器事件监听器
  window.addEventListener("script-stream-start", handleStreamStart as EventListener);
  window.addEventListener("script-stream-chunk", handleStreamChunk as EventListener);
  window.addEventListener("script-stream-end", handleStreamEnd as EventListener);

  // 在清理函数中移除事件监听器
  eventCleanupFunctions.push(() => {
    window.removeEventListener("script-stream-start", handleStreamStart as EventListener);
    window.removeEventListener("script-stream-chunk", handleStreamChunk as EventListener);
    window.removeEventListener("script-stream-end", handleStreamEnd as EventListener);
  });
});

// 组件卸载时清理
onUnmounted(() => {
  if (autoCleanInterval) {
    clearInterval(autoCleanInterval);
    autoCleanInterval = null;
  }

  eventCleanupFunctions.forEach((cleanup) => cleanup());
  eventCleanupFunctions.length = 0;

  if (
    currentItem.value?.ContentType === "Image" &&
    currentItem.value.ImageData
  ) {
    currentItem.value.ImageData = null as any;
  }
  currentItem.value = null;
  items.value = [];
});
</script>

<style scoped>
.clipboard-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
}

.content-area {
  margin: 12px 20px 0px 8px;
  border-radius: 16px;
  overflow-y: auto;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.content-display {
  padding: 14px;
  background-color: #fff;
}

.welcome-text {
  color: #86868b;
  text-align: center;
  padding: 40px 20px;
  font-size: 16px;
}

.info-panel {
  margin: 12px 20px 12px 8px;
  padding: 8px 12px 0px 12px;
  background-color: #fff;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}
</style>

<style>
.el-drawer {
  background-color: #fafafa !important;
}
.el-drawer__body {
  background-color: #fafafa !important;
  padding: 0 !important;
}
.el-drawer__header{
  margin-bottom: 10px !important;
}
</style>
