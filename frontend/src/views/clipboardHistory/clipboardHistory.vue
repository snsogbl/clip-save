<template>
  <!-- 设置页面 Drawer -->
  <el-drawer
    v-model="showSetting"
    :title="$t('settings.title')"
    direction="rtl"
    size="600px"
    @close="handleSettingBack"
    class="settings-drawer"
    destroy-on-close
  >
    <SettingView />
  </el-drawer>
  <div class="clipboard-container" style="--wails-draggable: no-drag">
    <!-- 顶部工具栏 -->
    <div class="toolbar" style="--wails-draggable: drag">
      <div class="toolbar-left">
        <div class="title-bg">
          <el-icon :size="20" class="iconfont icon-shandian"> </el-icon>
          <span class="toolbar-left-text">
            {{ $t("app.title") }}
          </span>
        </div>
      </div>
      <div class="toolbar-right">
        <el-dropdown placement="bottom">
          <el-icon :size="20" class="iconfont icon-duoyuyan"> </el-icon>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="changeLanguage('zh-CN')"
                >中文</el-dropdown-item
              >
              <el-dropdown-item @click="changeLanguage('en-US')"
                >English</el-dropdown-item
              >
              <el-dropdown-item @click="changeLanguage('fr-FR')"
                >Français</el-dropdown-item
              >
              <el-dropdown-item @click="changeLanguage('ar-SA')"
                >العربية</el-dropdown-item
              >
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-select
          v-model="filterType"
          class="filter-select"
          @change="onSearchChange"
          size="large"
          :placeholder="$t('main.filterAll')"
        >
          <el-option :label="$t('main.filterAll')" value="" />
          <el-option :label="$t('main.filterText')" value="Text" />
          <el-option :label="$t('main.filterImage')" value="Image" />
          <el-option :label="$t('main.filterFile')" value="File" />
          <el-option :label="$t('main.filterUrl')" value="URL" />
          <el-option :label="$t('main.filterColor')" value="Color" />
          <el-option :label="$t('main.filterJSON')" value="JSON" />
        </el-select>
        <el-input
          ref="searchInputRef"
          v-model="searchKeyword"
          type="text"
          class="search-input"
          :prefix-icon="Search"
          :placeholder="$t('main.searchPlaceholder')"
          @keydown="handleSearchKeydown"
          @input="onSearchChange"
          clearable
          style="--wails-draggable: no-drag"
        />
        <el-button class="me-button" circle @click="showSetting = true">
          <el-icon :size="20">
            <Setting />
          </el-icon>
        </el-button>
      </div>
    </div>

    <!-- 主内容区域 -->
    <div class="main-content">
      <!-- 左侧列表 -->
      <div class="left-panel">
        <div class="tab-buttons">
          <el-button
            round
            class="me-button"
            :class="{ active: leftTab === 'all' }"
            @click="switchLeftTab('all')"
          >
            <el-icon :size="20" class="iconfont icon-liebiao"> </el-icon>
            <span>{{ $t("main.listTitle") }}</span>
          </el-button>
          <el-button
            round
            class="me-button"
            :class="{ active: leftTab === 'fav' }"
            @click="switchLeftTab('fav')"
          >
            <el-icon><Star /></el-icon>
            <span>{{ $t("main.favorite") }}</span>
          </el-button>
        </div>
        <div class="item-list" ref="itemListRef" tabindex="-1">
          <div v-if="loading" class="loading">{{ $t("main.loading") }}</div>
          <div v-else-if="items.length === 0" class="empty-state">
            <el-icon :size="48" class="iconfont icon-kongyemian"> </el-icon>
            <div class="empty-text">{{ $t("main.emptyState") }}</div>
          </div>
          <div
            v-else
            v-for="(item, index) in items"
            :key="item.ID"
            class="list-item"
            :class="{ active: currentItem?.ID === item.ID }"
            @click="selectItem(item)"
            @dblclick="handleDoubleClick(item)"
          >
            <!-- 数字标签（按住 Command 时显示前 9 个） -->
            <div
              v-if="isCommandPressed && index < 9"
              class="quick-access-badge"
            >
              {{ index + 1 }}
            </div>
            <div class="item-header">
              <el-icon class="item-icon" :size="18">
                <Document v-if="item.ContentType === 'Text'" />
                <Link v-else-if="item.ContentType === 'URL'" />
                <Folder v-else-if="item.ContentType === 'File'" />
                <Brush v-else-if="item.ContentType === 'Color'" />
                <Picture v-else-if="item.ContentType === 'Image'" />
                <Document v-else-if="item.ContentType === 'JSON'" />
                <Document v-else />
              </el-icon>
              <span class="item-content">{{ getPreview(item) }}</span>
              <div
                v-if="item.ContentType === 'Color'"
                class="color-circle-small"
                :style="{ backgroundColor: item.Content }"
              ></div>
              <el-icon
                v-if="item.IsFavorite === 1"
                :size="16"
                style="color: #f5a623"
              >
                <Star />
              </el-icon>
            </div>
            <div class="item-footer">
              <span class="item-type" style="width: 40px">{{
                item.ContentType
              }}</span>
              <span class="item-time">{{ formatTime(item.Timestamp) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧详情 -->
      <div class="right-panel" style="--wails-draggable: no-drag">
        <div class="content-area" ref="contentAreaRef">
          <ClipboardTitleView
            v-if="currentItem"
            :item="currentItem"
            @copy="copyItem"
            @delete="deleteItem"
            @collect="collectItem"
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
              :imageData="currentItem.ImageData"
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
            />
          </div>
        </div>
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
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, onUnmounted, nextTick } from "vue";
import { EventsOn } from "../../../wailsjs/runtime/runtime";
import { useI18n } from "vue-i18n";
import {
  SearchClipboardItems,
  CopyToClipboard,
  DeleteClipboardItem,
  OpenFileInFinder,
  OpenURL,
  ClearItemsOlderThanDays,
  GetAppSettings,
  HideWindow,
  ToggleFavorite,
  HideWindowAndQuit,
  SetLanguage,
  AutoPasteCurrentItem
} from "../../../wailsjs/go/main/App";

const { t, locale } = useI18n();
import {
  Document,
  Link,
  Folder,
  Brush,
  Picture,
  DocumentCopy,
  Delete,
  Setting,
  Star,
  Search,
  List,
} from "@element-plus/icons-vue";
import ClipboardUrlView from "./components/clipboardUrlView.vue";
import ClipboardColorView from "./components/clipboardColorView.vue";
import ClipboardFileView from "./components/clipboardFileView.vue";
import ClipboardTextView from "./components/clipboardTextView.vue";
import ClipboardImageView from "./components/clipboardImageView.vue";
import ClipboardJsonView from "./components/clipboardJsonView.vue";
import ClipboardTitleView from "./components/clipboardTitleView.vue";
import SettingView from "../setting/setting.vue";
import { ElMessageBox, ElMessage } from "element-plus";

interface ClipboardItem {
  ID: string;
  Content: string;
  ContentType: string;
  ImageData: any; // Go []byte 会被序列化为 base64 字符串
  FilePaths: string; // JSON 数组格式
  FileInfo: string; // JSON 格式
  Timestamp: string;
  Source: string;
  CharCount: number;
  WordCount: number;
  IsFavorite: number;
}

interface FileInfo {
  name: string;
  path: string;
  size: number;
  is_dir: boolean;
  exists: boolean;
  extension: string;
}

const items = ref<ClipboardItem[]>([]);
const currentItem = ref<ClipboardItem | null>(null);
const itemListRef = ref<HTMLElement | null>(null);
const contentAreaRef = ref<HTMLElement | null>(null);
const searchInputRef = ref<HTMLInputElement | null>(null);
const textEditorRef = ref<InstanceType<typeof ClipboardTextView> | null>(null);
const searchKeyword = ref("");
const filterType = ref("");
const loading = ref(false);
const showSetting = ref(false);
const leftTab = ref<"all" | "fav">("all");
const jsonEditorRef = ref<InstanceType<typeof ClipboardJsonView> | null>(null);
const isCommandPressed = ref(false);

// 缓存的设置数据，避免频繁查询数据库
let cachedSettings: {
  pageSize: number;
  autoClean: boolean;
  retentionDays: number;
} | null = null;

// 从数据库获取设置（带缓存）
async function getSettings(forceRefresh = false) {
  // 如果已有缓存且不需要强制刷新，直接返回缓存
  if (cachedSettings && !forceRefresh) {
    return cachedSettings;
  }

  try {
    const savedSettings = await GetAppSettings();
    if (savedSettings) {
      cachedSettings = JSON.parse(savedSettings);
      return cachedSettings;
    }
  } catch (e) {
    console.error("❌ 读取设置失败:", e);
  }
  // 返回默认值（数据库初始化时应该已经创建了默认设置）
  cachedSettings = { pageSize: 100, autoClean: true, retentionDays: 30 };
  return cachedSettings;
}

// 加载剪贴板项目
async function loadItems() {
  try {
    loading.value = true;
    const settings = await getSettings();
    const pageSize = settings?.pageSize || 100;
    console.log("📊 使用页面大小:", pageSize);

    const result = await SearchClipboardItems(
      leftTab.value === "fav",
      searchKeyword.value,
      filterType.value,
      pageSize
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

// 静默检查更新（不显示加载状态）
async function checkForUpdates() {
  try {
    // 使用缓存的设置，避免频繁查询数据库
    const settings = await getSettings();
    const pageSize = settings?.pageSize || 100;

    const result = await SearchClipboardItems(
      leftTab.value === "fav",
      searchKeyword.value,
      filterType.value,
      pageSize
    );
    const newItems = result || [];

    // 只在数据真正变化时才更新（比较第一个项目的ID和总数）
    if (
      newItems.length !== items.value.length ||
      (newItems.length > 0 &&
        items.value.length > 0 &&
        newItems[0].ID !== items.value[0].ID)
    ) {
      items.value = newItems;

      // 如果没有选中项，自动选中第一项
      if (!currentItem.value && newItems.length > 0) {
        selectItem(newItems[0]);
      }

      console.log("✨ 检测到剪贴板更新");
    }
  } catch (error) {
    console.error("检查更新失败:", error);
  }
}

// 选择项目
async function selectItem(item: ClipboardItem) {
  currentItem.value = item;
  await nextTick();
  // 确保当前选中项进入可视区域
  const container = itemListRef.value;
  if (!container) return;
  const activeEl = container.querySelector(
    ".list-item.active"
  ) as HTMLElement | null;
  if (activeEl) {
    activeEl.scrollIntoView({ block: "nearest" });
  }
  // 将内容区域滚动到顶部
  if (contentAreaRef.value) {
    contentAreaRef.value.scrollTo({ top: 0, behavior: "smooth" });
  }
}

// 处理双击事件
async function handleDoubleClick(item: ClipboardItem) {
  // 如果双击的项目不是当前选中的，先选中它
  if (currentItem.value?.ID !== item.ID) {
    await selectItem(item);
    // 等待 DOM 更新，特别是 JSON 编辑器
    await nextTick();
  }
  // 复制当前项
  await copyItem(item.ID);
  HideWindowAndQuit();
  AutoPasteCurrentItem();
}

// 复制项目
async function copyItem(id: string) {
  if (currentItem.value?.ContentType === "JSON") {
    jsonEditorRef.value?.copyEdited();
  } else {
    try {
      await CopyToClipboard(id);
      ElMessage.success(t("message.copySuccess"));
      console.log("已复制到剪贴板");
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
    // 就地更新左侧 items
    const index = items.value.findIndex((i) => i.ID === id);
    if (index !== -1) {
      // 在收藏页签下，取消收藏需要从列表移除
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
        // 其他情况仅更新该项的收藏状态
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

async function switchLeftTab(tab: "all" | "fav") {
  if (leftTab.value === tab) return;
  leftTab.value = tab;
  await loadItems();
  await nextTick();
  itemListRef.value?.focus();
}
// 格式化时间
function formatTime(timestamp: string): string {
  const date = new Date(timestamp);
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  const minutes = Math.floor(diff / 60000);
  const hours = Math.floor(diff / 3600000);
  const days = Math.floor(diff / 86400000);

  if (minutes < 60) {
    return `${minutes}分钟前`;
  } else if (hours < 24) {
    return `${hours}小时前`;
  } else if (days < 7) {
    return `${days}天前`;
  } else {
    return date.toLocaleString("zh-CN", {
      month: "2-digit",
      day: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
    });
  }
}

// 获取预览文本
function getPreview(item: ClipboardItem): string {
  let preview = item.Content || "空内容";
  if (preview.length > 30) {
    preview = preview.substring(0, 27) + "...";
  }
  return preview;
}

// 搜索和过滤变化时重新加载
const onSearchChange = () => {
  loadItems();
};

// 处理搜索框键盘按下事件
function handleSearchKeydown(event: KeyboardEvent) {
  // 检测 Cmd+Enter 或 Ctrl+Enter
  if ((event.metaKey || event.ctrlKey) && event.key === "Enter") {
    event.preventDefault();
    event.stopPropagation();
    // 直接执行复制并退出功能
    if (currentItem.value) {
      handleDoubleClick(currentItem.value);
    }
    return;
  }
  // 检测 Cmd+Left 或 Ctrl+Left（切换到列表页签）
  if ((event.metaKey || event.ctrlKey) && event.key === "ArrowLeft") {
    event.preventDefault();
    event.stopPropagation();
    switchLeftTab("all").then(() => {
      // 切换后恢复搜索框焦点
      nextTick(() => {
        searchInputRef.value?.focus();
      });
    });
    return;
  }
  // 检测 Cmd+Right 或 Ctrl+Right（切换到收藏页签）
  if ((event.metaKey || event.ctrlKey) && event.key === "ArrowRight") {
    event.preventDefault();
    event.stopPropagation();
    switchLeftTab("fav").then(() => {
      // 切换后恢复搜索框焦点
      nextTick(() => {
        searchInputRef.value?.focus();
      });
    });
    return;
  }
}

// 处理全局键盘事件（用于 Command+数字键快速粘贴）
function handleGlobalKeydown(event: KeyboardEvent) {
  // 检测 Command/Ctrl 键按下
  if (event.metaKey || event.ctrlKey) {
    if (!isCommandPressed.value) {
      isCommandPressed.value = true;
    }
    
    // 检测 Command+数字键（1-9）
    const numKey = parseInt(event.key);
    if (!isNaN(numKey) && numKey >= 1 && numKey <= 9) {
      event.preventDefault();
      event.stopPropagation();
      // 快速粘贴对应索引的项目（索引从 0 开始，所以减 1）
      const index = numKey - 1;
      if (items.value[index]) {
        handleDoubleClick(items.value[index]);
      }
      // 重置状态
      isCommandPressed.value = false;
      return;
    }
  } else {
    // 非 Command 键按下时，如果之前是按下的状态，检查是否是 Command 键本身
    if (event.key !== "Meta" && event.key !== "Control" && isCommandPressed.value) {
      // 如果按下的不是 Command 键，说明 Command 已经松开
      isCommandPressed.value = false;
    }
  }
}

// 处理全局键盘松开事件
function handleGlobalKeyup(event: KeyboardEvent) {
  // Command/Ctrl 键松开
  if (event.key === "Meta" || event.key === "Control" || event.key === "MetaLeft" || event.key === "MetaRight" || event.key === "ControlLeft" || event.key === "ControlRight") {
    isCommandPressed.value = false;
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
    console.log("已在 Finder 中打开文件");
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
    console.log("已在浏览器中打开 URL");
  } catch (error) {
    console.error("在浏览器中打开 URL 失败:", error);
    ElMessage.error(t("message.openUrlError", [error]));
  }
}

// 处理设置页面返回
async function handleSettingBack() {
  console.log("handleSettingBack");
  showSetting.value = false;
  await getSettings(true);
}

// 自动清理超过指定天数的历史记录
async function autoCleanOldItems() {
  // 使用缓存的设置，避免频繁查询数据库
  const settings = await getSettings();

  if (!settings?.autoClean) {
    return; // 未启用自动清理
  }

  const retentionDays = settings.retentionDays || 30;

  try {
    console.log(`🗑️ 执行自动清理: 删除超过 ${retentionDays} 天的记录`);
    await ClearItemsOlderThanDays(retentionDays);
    console.log(`✅ 自动清理完成`);
  } catch (error) {
    console.error("❌ 自动清理失败:", error);
  }
}

// 初始化和定时刷新
onMounted(() => {
  // 初始化设置缓存
  getSettings().then(() => {
    loadItems();
  });

  // 每1秒静默检查更新（不会导致闪烁）
  setInterval(() => {
    checkForUpdates();
  }, 1000);

  // 启动时执行一次自动清理
  autoCleanOldItems();

  // 每小时执行一次自动清理
  setInterval(() => {
    autoCleanOldItems();
  }, 60 * 60 * 1000); // 1小时 = 60分钟 * 60秒 * 1000毫秒

  // 监听全局键盘事件（用于 Command+数字键快速粘贴）
  window.addEventListener("keydown", handleGlobalKeydown);
  window.addEventListener("keyup", handleGlobalKeyup);

  // 监听窗口显示事件：从后台切换到前台时，选中第一个列表项
  EventsOn("window.show", () => {
    setTimeout(() => {
      checkForUpdates();
      if (items.value.length > 0) {
        selectItem(items.value[0]);
      }
      searchInputRef.value?.focus();
    }, 100);
  });

  // 监听菜单事件：上一条/下一条
  EventsOn("nav.prev", () => {
    if (items.value.length === 0) return;
    if (!currentItem.value) {
      selectItem(items.value[0]);
      return;
    }
    const idx = items.value.findIndex((i) => i.ID === currentItem.value!.ID);
    const nextIdx = Math.max(0, idx - 1);
    selectItem(items.value[nextIdx]);
  });
  EventsOn("nav.next", () => {
    if (items.value.length === 0) return;
    if (!currentItem.value) {
      selectItem(items.value[0]);
      return;
    }
    const idx = items.value.findIndex((i) => i.ID === currentItem.value!.ID);
    const nextIdx = Math.min(items.value.length - 1, idx + 1);
    selectItem(items.value[nextIdx]);
  });

  EventsOn("nav.switch", (tab: "all" | "fav") => {
    switchLeftTab(tab);
  });
  EventsOn("nav.setting", () => {
    showSetting.value = true;
  });
  EventsOn("copy.current", () => {
    copyItem(currentItem.value!.ID);
  });
  EventsOn("delete.current", () => {
    deleteItem(currentItem.value!.ID);
  });
  EventsOn("collect.current", () => {
    collectItem(currentItem.value!.ID);
  });
  EventsOn("search.item", () => {
    searchInputRef.value?.focus();
  });
  EventsOn("translate.current", () => {
    textEditorRef.value?.translateText();
  });
});

function changeLanguage(lang: string) {
  SetLanguage(lang);
  locale.value = lang as any;
}

// 组件卸载时清理事件监听器
onUnmounted(() => {
  window.removeEventListener("keydown", handleGlobalKeydown);
  window.removeEventListener("keyup", handleGlobalKeyup);
});
</script>

<style scoped>
.clipboard-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
}

.toolbar {
  display: flex;
  gap: 12px;
  padding: 8px 14px;
  border-bottom: 1px solid #e0e0e0;
  /* align-items: center; */

  .toolbar-left {
    display: flex;
    align-items: center;
    gap: 8px;

    .title-bg {
      margin-left: 60px;
      .toolbar-left-text {
        font-size: 16px;
        font-weight: 600;
      }
    }
  }
  .toolbar-right {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-left: auto;
    .search-input {
      width: 300px;
    }
  }
}
.search-input :deep(.el-input__wrapper) {
  border-radius: 20px;
}

.filter-select {
  width: 80px;
  color: #000;
}

.filter-select :deep(.el-select__wrapper) {
  border: none;
  box-shadow: none;
  color: #000;
  padding: 0 !important;
}
.filter-select :deep(.el-select__wrapper):hover {
  border: none;
  box-shadow: none;
}
.filter-select :deep(.el-select__placeholder.is-transparent) {
  color: #000;
}
.filter-select :deep(.el-select__placeholder) {
  color: #000;
  text-align: right;
}
.filter-select :deep(.el-select__caret) {
  color: #000;
}

.setting-btn {
  border: 1px solid #e0e0e0;
  color: #666;
  transition: all 0.2s ease;
  width: 30px !important;
  height: 30px !important;
}

.setting-btn:hover {
  background-color: #f8f8f8;
  border-color: #007aff;
  color: #007aff;
  transform: scale(1.05);
}

.setting-btn:active {
  transform: scale(0.98);
}

.layout-btn {
  color: #888;
}

.main-content {
  display: flex;
  flex: 1;
  overflow: hidden;
  background-color: #fafafa;
}

.left-panel {
  width: 280px;
  display: flex;
  flex-direction: column;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.04);
  background-color: #fff;
  margin: 12px;
  border-radius: 12px;
}

.item-list {
  flex: 1;
  overflow-y: auto;
}

/* 去除程序化聚焦后的蓝色边框 */
.item-list:focus {
  outline: none;
  box-shadow: none;
}

.loading,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: 40px 20px;
  color: #86868b;
  gap: 10px;
}

.list-item {
  padding: 10px;
  margin: 0 12px 10px;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s ease;
  border: 1px solid #e8e8e8;
  position: relative;
}

.list-item.active {
  border: 1px solid #999;
  background-color: #fafafa;
}

.list-item:hover {
  background-color: #fafafa;
}

.item-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.item-icon {
  color: #666;
  display: flex;
  align-items: center;
}

.item-content {
  flex: 1;
  font-size: 15px;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: #2c2c2e;
  text-align: left;
  line-height: 1.4;
}

.item-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 13px;
  color: #8e8e93;
  margin-top: 6px;
}

.item-type {
  background-color: #f2f2f7;
  color: #6d6d70;
  padding: 2px 8px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 500;
  min-width: 40px;
  text-align: center;
}

.panel-footer {
  padding: 16px;
  border-top: 1px solid #f0f0f0;
  color: #000;
  text-align: center;
  font-size: 16px;
  font-weight: 600;
  margin-top: auto;
  border-radius: 0 0 0 0;
}

.right-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  background-color: transparent;
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

.content-image {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
}

.info-panel {
  margin: 12px 20px 12px 8px;
  padding: 8px 12px 0px 12px;
  background-color: #fff;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.info-row {
  display: flex;
  gap: 12px;
  margin-bottom: 12px;
}

.info-row:last-child {
  margin-bottom: 0;
}

.info-label {
  font-weight: 600;
  color: #333;
  min-width: 90px;
  font-size: 14px;
}

.info-value {
  color: #1a1a1a;
  font-size: 14px;
}

/* 颜色显示样式 - 仅保留小圆圈样式（列表中使用） */
.color-circle-small {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: 2px solid #e0e0e0;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
  flex-shrink: 0;
  margin-left: auto;
}

.tab-buttons {
  padding: 16px 16px 12px 20px;
  display: inline-flex;
  gap: 4px;
}

.quick-access-badge {
  position: absolute;
  top: 0px;
  right: 0px;
  width: 16px;
  height: 16px;
  background: rgba(153, 153, 153, 0.6);
  color: #fff;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: 500;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
  z-index: 10;
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
</style>
