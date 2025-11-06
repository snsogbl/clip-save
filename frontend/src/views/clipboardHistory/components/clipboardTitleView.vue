<template>
  <div v-if="item" class="clipboard-title-view">
    <div class="title-content">
      <div class="title-main">{{ displayContent.length > 30 ? displayContent.substring(0, 27) + "..." : displayContent }}</div>
      <div class="title-source">{{ $t('main.source') }} {{ item.Source || $t('main.unknown') }}</div>
    </div>
    <div class="title-actions">
      <el-button class="me-button" @click="handleCopy" round>
        <el-icon><DocumentCopy /></el-icon>
        <span>{{ $t('main.copy') }}</span>
      </el-button>
      <el-button
        class="me-button"
        :class="{ active: item?.IsFavorite === 1 }"
        @click="handleCollect"
        round
      >
        <el-icon><Star /></el-icon>
      <span>{{ item?.IsFavorite === 1 ? $t('main.unfavorite') : $t('main.favorite') }}</span>
      </el-button>
      <el-button type="danger" @click="handleDelete" round>
        <el-icon><Delete /></el-icon>
        <span>{{ $t('main.delete') }}</span>
      </el-button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed } from "vue";
import { DocumentCopy, Delete, Star } from "@element-plus/icons-vue";

interface ClipboardItem {
  ID: string;
  Content: string;
  ContentType: string;
  Source: string;
  IsFavorite: number;
}

const props = defineProps<{
  item: ClipboardItem | null;
}>();

const emit = defineEmits<{
  copy: [id: string];
  delete: [id: string];
  collect: [id: string];
}>();

const displayContent = computed(() => {
  if (!props.item) return "";
  const content = props.item.Content || "";
  // 如果内容太长，截断并显示...
  if (content.length > 50) {
    return content.substring(0, 47) + "...";
  }
  return content;
});

function handleCopy() {
  if (props.item) {
    emit("copy", props.item.ID);
  }
}

function handleDelete() {
  if (props.item) {
    emit("delete", props.item.ID);
  }
}

function handleCollect() {
  if (props.item) {
    emit("collect", props.item.ID);
  }
}
</script>

<style scoped>
.clipboard-title-view {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background-color: #fff;
  border-bottom: 1px solid #f0f0f0;
}

.title-content {
  flex: 1;
  min-width: 0;
}

.title-main {
  font-size: 18px;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 8px;
  line-height: 1.5;
  word-break: break-word;
}

.title-source {
  font-size: 13px;
  color: #999;
  line-height: 1.4;
}

.title-actions {
  display: flex;
  margin-left: 20px;
  flex-shrink: 0;
}
</style>
