<template>
  <div
    class="normal-list-item"
    :class="{ active: currentItem?.ID === item.ID }"
    @click="$emit('select', item)"
    @dblclick="$emit('double-click', item)"
  >
    <!-- 数字标签（按住 Command 时显示前 9 个） -->
    <div v-if="isCommandPressed && index < 9" class="quick-access-badge">
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
      <span class="item-content">{{ item.Content }}</span>
      <div
        v-if="item.ContentType === 'Color'"
        class="color-circle-small"
        :style="{ backgroundColor: item.Content }"
      ></div>
      <el-icon v-if="item.IsFavorite === 1" :size="16" style="color: #f5a623">
        <Star />
      </el-icon>
    </div>
    <div class="item-footer">
      <span class="item-type" style="width: 40px">{{ item.ContentType }}</span>
      <span class="item-time">{{ formatTime(item.Timestamp) }}</span>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {
  Document,
  Link,
  Folder,
  Brush,
  Picture,
  Star,
} from "@element-plus/icons-vue";
import { common } from "../../../../wailsjs/go/models";

type ClipboardItem = common.ClipboardItem;

defineProps<{
  item: ClipboardItem;
  index: number;
  currentItem: ClipboardItem | null;
  isCommandPressed: boolean;
}>();

defineEmits<{
  select: [item: ClipboardItem];
  "double-click": [item: ClipboardItem];
}>();

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
</script>

<style scoped>
@import "./styles/normal-list-item.css";
</style>