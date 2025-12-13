<template>
  <div class="normal-mode-container">
    <!-- 顶部工具栏 -->
    <div class="toolbar" style="--wails-draggable: drag">
      <div class="toolbar-left">
        <div class="title-bg" :style="{ marginLeft: isMacOS ? '60px' : '0px' }">
          <el-icon :size="20" class="iconfont icon-shandian"> </el-icon>
          <span class="toolbar-left-text">
            {{ $t("app.title") }}
          </span>
        </div>
      </div>
      <div class="toolbar-right">
        <el-icon
          :size="18"
          class="iconfont icon-chuangti-zhiding"
          :style="{
            color: isAlwaysOnTop ? '#000' : '#bebebe',
            transform: isAlwaysOnTop ? 'rotate(0deg)' : 'rotate(45deg)',
          }"
          @click="$emit('toggle-always-on-top')"
          title="置顶"
        ></el-icon>
        <el-dropdown placement="bottom">
          <el-icon :size="20" class="iconfont icon-duoyuyan"> </el-icon>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="$emit('change-language', 'zh-CN')"
                >中文</el-dropdown-item
              >
              <el-dropdown-item @click="$emit('change-language', 'en-US')"
                >English</el-dropdown-item
              >
              <el-dropdown-item @click="$emit('change-language', 'fr-FR')"
                >Français</el-dropdown-item
              >
              <el-dropdown-item @click="$emit('change-language', 'ar-SA')"
                >العربية</el-dropdown-item
              >
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-select
          :model-value="filterType"
          class="filter-select"
          @change="$emit('filter-change', $event)"
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
          :model-value="searchKeyword"
          type="text"
          class="search-input"
          :prefix-icon="Search"
          :placeholder="$t('main.searchPlaceholder')"
          @keydown="$emit('search-keydown', $event)"
          @input="$emit('search-change', $event)"
          clearable
          style="--wails-draggable: no-drag"
        />
        <el-button class="me-button" circle @click="$emit('show-setting')">
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
        <div class="tab-buttons" style="--wails-draggable: drag">
          <el-button
            round
            class="me-button"
            :class="{ active: leftTab === 'all' }"
            @click="$emit('switch-tab', 'all')"
          >
            <el-icon :size="20" class="iconfont icon-liebiao"> </el-icon>
            <span>{{ $t("main.listTitle") }}</span>
          </el-button>
          <el-button
            round
            class="me-button"
            :class="{ active: leftTab === 'fav' }"
            @click="$emit('switch-tab', 'fav')"
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
          <NormalListItem
            v-else
            v-for="(item, index) in items"
            :key="item.ID"
            :item="item"
            :index="index"
            :current-item="currentItem"
            :is-command-pressed="isCommandPressed"
            @select="$emit('select-item', item)"
            @double-click="$emit('double-click', item)"
          />
        </div>
      </div>

      <!-- 右侧详情 -->
      <div class="right-panel" style="--wails-draggable: no-drag">
        <slot name="content-area" />
        <slot name="info-panel" />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { useI18n } from "vue-i18n";
import { Setting, Search, Star } from "@element-plus/icons-vue";
import NormalListItem from "./NormalListItem.vue";
import { common } from "../../../../wailsjs/go/models";

type ClipboardItem = common.ClipboardItem;

const { t } = useI18n();

defineProps<{
  items: ClipboardItem[];
  currentItem: ClipboardItem | null;
  loading: boolean;
  leftTab: "all" | "fav";
  searchKeyword: string;
  filterType: string;
  isAlwaysOnTop: boolean;
  isCommandPressed: boolean;
}>();

defineEmits<{
  "toggle-always-on-top": [];
  "show-setting": [];
  "change-language": [lang: string];
  "filter-change": [type: string];
  "search-keydown": [event: KeyboardEvent];
  "search-change": [event: Event];
  "switch-tab": [tab: "all" | "fav"];
  "select-item": [item: ClipboardItem];
  "double-click": [item: ClipboardItem];
}>();

const searchInputRef = ref<HTMLInputElement | null>(null);
const itemListRef = ref<HTMLElement | null>(null);

// 检测是否为 macOS
const isMacOS = ref(navigator.platform.toUpperCase().indexOf('MAC') >= 0);

defineExpose({
  searchInputRef,
  itemListRef,
});
</script>

<style scoped>
@import "./styles/normal-mode.css";
</style>
