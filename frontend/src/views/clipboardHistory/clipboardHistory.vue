<template>
  <div class="clipboard-container" style="--wails-draggable: no-drag">
    <!-- 设置页面 -->
    <SettingView v-if="showSetting" @back="showSetting = false" />

    <!-- 剪贴板历史主页面 -->
    <template v-else>
      <!-- 顶部工具栏 -->
      <div class="toolbar">
        <el-input
          v-model="searchKeyword"
          type="text"
          class="search-input"
          :placeholder="$t('main.searchPlaceholder')"
          @input="onSearchChange"
          clearable
          style="--wails-draggable: no-drag"
        />
        <el-select
          v-model="filterType"
          class="filter-select"
          @change="onSearchChange"
          :placeholder="$t('main.filterAll')"
        >
          <el-option :label="$t('main.filterAll')" value="" />
          <el-option :label="$t('main.filterText')" value="Text" />
          <el-option :label="$t('main.filterImage')" value="Image" />
          <el-option :label="$t('main.filterFile')" value="File" />
          <el-option :label="$t('main.filterUrl')" value="URL" />
          <el-option :label="$t('main.filterColor')" value="Color" />
        </el-select>
        <el-button class="setting-btn" circle @click="showSetting = true">
          <el-icon :size="20">
            <Setting />
          </el-icon>
        </el-button>
      </div>

      <!-- 主内容区域 -->
      <div class="main-content">
        <!-- 左侧列表 -->
        <div class="left-panel">
          <div class="panel-header">
            <el-tabs v-model="leftTab" class="tabs" @tab-click="switchLeftTab">
              <el-tab-pane :label="$t('main.listTitle')" name="all">
              </el-tab-pane>
              <el-tab-pane :label="$t('main.favorite')" name="fav">
              </el-tab-pane>
            </el-tabs>
          </div>
          <div class="item-list">
            <div v-if="loading" class="loading">{{ $t("main.loading") }}</div>
            <div v-else-if="items.length === 0" class="empty-state">
              <div class="empty-icon">📋</div>
              <div class="empty-text">{{ $t("main.emptyState") }}</div>
            </div>
            <div
              v-else
              v-for="item in items"
              :key="item.ID"
              class="list-item"
              :class="{ active: currentItem?.ID === item.ID }"
              @click="selectItem(item)"
            >
              <div class="item-header">
                <el-icon class="item-icon" :size="18">
                  <Document v-if="item.ContentType === 'Text'" />
                  <Link v-else-if="item.ContentType === 'URL'" />
                  <Folder v-else-if="item.ContentType === 'File'" />
                  <Brush v-else-if="item.ContentType === 'Color'" />
                  <Picture v-else-if="item.ContentType === 'Image'" />
                  <Document v-else />
                </el-icon>
                <span class="item-content">{{ getPreview(item) }}</span>
                <el-icon
                  v-if="item.IsFavorite === 1"
                  :size="16"
                  style="color: #f5a623"
                >
                  <Star />
                </el-icon>
                <div
                  v-if="item.ContentType === 'Color'"
                  class="color-circle-small"
                  :style="{ backgroundColor: item.Content }"
                ></div>
              </div>
              <div class="item-footer">
                <span class="item-type" style="width: 40px">{{
                  item.ContentType
                }}</span>
                <span class="item-time">{{ formatTime(item.Timestamp) }}</span>
              </div>
            </div>
          </div>
          <div class="panel-footer">
            <strong>{{ $t("main.clipboardHistory") }}</strong>
          </div>
        </div>

        <!-- 右侧详情 -->
        <div class="right-panel" style="--wails-draggable: no-drag">
          <div class="content-area">
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
              <!-- 文本内容展示 -->
              <ClipboardTextView
                v-else
                :text="currentItem?.Content || '空内容'"
              />
            </div>

            <div v-if="currentItem" class="info-panel">
              <div class="info-row">
                <span class="info-label">{{ $t("main.source") }}</span>
                <span class="info-value">{{ currentItem.Source }}</span>
              </div>
              <div class="info-row">
                <span class="info-label">{{ $t("main.contentType") }}</span>
                <span class="info-value">{{ currentItem.ContentType }}</span>
              </div>
              <template v-if="currentItem.ContentType !== 'File'">
                <div class="info-row">
                  <span class="info-label">{{ $t("main.charCount") }}</span>
                  <span class="info-value">{{ currentItem.CharCount }}</span>
                </div>
                <div class="info-row">
                  <span class="info-label">{{ $t("main.wordCount") }}</span>
                  <span class="info-value">{{ currentItem.WordCount }}</span>
                </div>
              </template>
              <template v-if="currentItem.ContentType === 'File'">
                <div class="info-row">
                  <span class="info-label">{{ $t("main.fileCount") }}</span>
                  <span class="info-value">{{ currentItem.WordCount }}</span>
                </div>
              </template>
              <div class="info-row">
                <span class="info-label">{{ $t("main.createTime") }}</span>
                <span class="info-value">{{
                  new Date(currentItem.Timestamp).toLocaleString("zh-CN")
                }}</span>
              </div>
            </div>
          </div>

          <div v-if="currentItem" class="actions-bar">
            <button class="action-btn" @click="copyItem(currentItem.ID)">
              <el-icon :size="16" style="margin-right: 6px">
                <DocumentCopy />
              </el-icon>
              {{ $t("main.copy") }}
            </button>
            <button
              class="action-btn delete"
              @click="deleteItem(currentItem.ID)"
            >
              <el-icon :size="16" style="margin-right: 6px">
                <Delete />
              </el-icon>
              {{ $t("main.delete") }}
            </button>
            <button
              class="action-btn delete"
              @click="collectItem(currentItem.ID)"
            >
              <el-icon :size="16" style="margin-right: 6px">
                <Star />
              </el-icon>
              {{
                currentItem.IsFavorite === 1
                  ? $t("main.unfavorite")
                  : $t("main.favorite")
              }}
            </button>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, nextTick } from "vue";
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
} from "../../../wailsjs/go/main/App";

const { t } = useI18n();
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
} from "@element-plus/icons-vue";
import ClipboardUrlView from "./components/clipboardUrlView.vue";
import ClipboardColorView from "./components/clipboardColorView.vue";
import ClipboardFileView from "./components/clipboardFileView.vue";
import ClipboardTextView from "./components/clipboardTextView.vue";
import ClipboardImageView from "./components/clipboardImageView.vue";
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
const searchKeyword = ref("");
const filterType = ref("");
const loading = ref(false);
const showSetting = ref(false);
const leftTab = ref<"all" | "fav">("all");

// 从数据库获取设置
async function getSettings() {
  try {
    const savedSettings = await GetAppSettings();
    if (savedSettings) {
      return JSON.parse(savedSettings);
    }
  } catch (e) {
    console.error("❌ 读取设置失败:", e);
  }
  // 返回默认值（数据库初始化时应该已经创建了默认设置）
  return { pageSize: 100, autoClean: true, retentionDays: 30 };
}

// 加载剪贴板项目
async function loadItems() {
  try {
    loading.value = true;
    const settings = await getSettings();
    const pageSize = settings.pageSize || 100;
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
    const settings = await getSettings();
    const pageSize = settings.pageSize || 100;

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
function selectItem(item: ClipboardItem) {
  currentItem.value = item;
}

// 复制项目
async function copyItem(id: string) {
  try {
    await CopyToClipboard(id);
    ElMessage.success(t("message.copySuccess"));
    console.log("已复制到剪贴板");
  } catch (error) {
    console.error("复制失败:", error);
    ElMessage.error(t("message.copyError", [error]));
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
      currentItem.value = null;
      await loadItems();
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

function switchLeftTab(tab: "all" | "fav") {
  leftTab.value = tab;
  loadItems();
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

// 自动清理超过指定天数的历史记录
async function autoCleanOldItems() {
  const settings = await getSettings();

  if (!settings.autoClean) {
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
  loadItems();

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
});

function hideApp() {
  setTimeout(() => {
    HideWindow();
  }, 100);
}
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
  padding: 20px;
  border-bottom: 1px solid #e0e0e0;
}

.filter-select {
  width: 140px;
}

.setting-btn {
  border: 1px solid #e0e0e0;
  color: #666;
  transition: all 0.2s ease;
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
}

.left-panel {
  width: 380px;
  border-right: 1px solid #e0e0e0;
  display: flex;
  flex-direction: column;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.04);
}

.panel-header {
  padding: 10px 20px 0px 20px;
  /* border-bottom: 1px solid #f0f0f0; */
}

.panel-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #1a1a1a;
  letter-spacing: -0.02em;
}

.item-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.loading,
.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: #86868b;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.list-item {
  padding: 14px 20px;
  margin: 0 12px 2px;
  /* border-radius: 10px; */
  cursor: pointer;
  transition: all 0.2s ease;
  border-bottom: 1px solid #bebebe;
}

.list-item:hover {
  background-color: #f8f8f8;
}

.list-item.active {
  background-color: #f8f8f8;
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
  overflow: auto;
}

.content-area {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.content-display {
  margin-bottom: 24px;
  padding: 14px;
  background-color: #fff;
  border-radius: 16px;
  border: 1px solid #e8e8e8;
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
  padding: 20px;
  background-color: rgba(255, 255, 255, 0.5);
  border-radius: 12px;
  border: 1px solid #e8e8e8;
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
  color: #1a1a1a;
  min-width: 90px;
  font-size: 14px;
}

.info-value {
  color: #333;
  font-size: 14px;
}

.actions-bar {
  display: flex;
  gap: 16px;
  padding: 10px;
  border-top: 1px solid #e0e0e0;
  background-color: transparent;
}

.action-btn {
  padding: 12px 24px;
  border: 1px solid #d1d1d6;
  border-radius: 10px;
  background-color: transparent;
  color: #1a1a1a;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
}

.action-btn:hover {
  background-color: #f2f2f7;
  border-color: #007aff;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.action-btn.delete:hover {
  background-color: #fff5f5;
  border-color: #ff3b30;
  color: #ff3b30;
}

.action-btn:active {
  transform: translateY(0);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
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

:deep(.el-tabs__item) {
  font-size: 18px;
  font-weight: 600;
}
</style>
