<template>
  <div class="minimal-mode-container">
    <!-- 极简模式顶部设置按钮 -->
    <div
      class="minimal-setting-buttons"
      style="--wails-draggable: drag"
      v-if="isMacOS"
    >
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
      <el-icon
        :size="18"
        class="iconfont setting-icon"
        @click="$emit('show-setting')"
      >
        <Setting />
      </el-icon>
    </div>

    <!-- 极简模式标签页和搜索 -->
    <div class="minimal-tab-buttons" :style="{ paddingTop: isMacOS ? '0' : '8px' }">
      <el-button
        round
        class="me-button"
        size="small"
        :class="{ active: leftTab === 'all' }"
        @click="$emit('switch-tab', 'all')"
      >
        <span>{{ $t("main.listTitle") }}</span>
      </el-button>
      <el-button
        round
        class="me-button"
        size="small"
        :class="{ active: leftTab === 'fav' }"
        @click="$emit('switch-tab', 'fav')"
      >
        <span>{{ $t("main.favorite") }}</span>
      </el-button>
      <el-input
        ref="searchInputRef"
        :model-value="searchKeyword"
        @update:model-value="$emit('update:search-keyword', $event)"
        type="text"
        class="search-input minimal-mode"
        size="small"
        :prefix-icon="Search"
        :placeholder="$t('main.searchPlaceholder')"
        @keydown="$emit('search-keydown', $event)"
        @input="$emit('search-change', $event)"
        clearable
        style="--wails-draggable: no-drag"
      />
      <template v-if="!isMacOS">
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
        <el-icon
          :size="18"
          class="iconfont setting-icon"
          @click="$emit('show-setting')"
        >
          <Setting />
        </el-icon>
      </template>
    </div>

    <!-- 极简模式列表 -->
    <div class="minimal-item-list" ref="itemListRef" tabindex="-1">
      <div v-if="loading" class="loading">{{ $t("main.loading") }}</div>
      <div v-else-if="items.length === 0" class="empty-state">
        <el-icon :size="48" class="iconfont icon-kongyemian"> </el-icon>
        <div class="empty-text">{{ $t("main.emptyState") }}</div>
      </div>
      <MinimalListItem
        v-else
        v-for="(item, index) in items"
        :key="item.ID"
        :item="item"
        :index="index"
        :current-item="currentItem"
        :is-command-pressed="isCommandPressed"
        @select="$emit('select-item', item)"
        @double-click="$emit('double-click', item)"
        @copy="$emit('copy-item', $event)"
        @send="$emit('send-item', $event)"
        @delete="$emit('delete-item', $event)"
        @collect="$emit('collect-item', $event)"
      />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { useI18n } from "vue-i18n";
import { Setting, Search } from "@element-plus/icons-vue";
import MinimalListItem from "./MinimalListItem.vue";
import { common } from "../../../../wailsjs/go/models";

type ClipboardItem = common.ClipboardItem;

const { t } = useI18n();

// 检测是否为 macOS
const isMacOS = ref(navigator.platform.toUpperCase().indexOf("MAC") >= 0);

defineProps<{
  items: ClipboardItem[];
  currentItem: ClipboardItem | null;
  loading: boolean;
  leftTab: "all" | "fav";
  searchKeyword: string;
  isAlwaysOnTop: boolean;
  isCommandPressed: boolean;
}>();

defineEmits<{
  "toggle-always-on-top": [];
  "show-setting": [];
  "switch-tab": [tab: "all" | "fav"];
  "update:search-keyword": [value: string];
  "search-keydown": [event: KeyboardEvent];
  "search-change": [event: Event];
  "select-item": [item: ClipboardItem];
  "double-click": [item: ClipboardItem];
  "copy-item": [id: string];
  "delete-item": [id: string];
  "collect-item": [id: string];
  "send-item": [id: ClipboardItem];
}>();

const searchInputRef = ref<HTMLInputElement | null>(null);
const itemListRef = ref<HTMLElement | null>(null);

defineExpose({
  searchInputRef,
  itemListRef,
});
</script>

<style scoped>
@import "./styles/minimal-mode.css";
</style>
