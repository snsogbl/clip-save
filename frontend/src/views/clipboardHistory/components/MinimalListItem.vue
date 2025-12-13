<template>
  <div
    class="minimal-list-item"
    :class="{ active: currentItem?.ID === item.ID }"
    @click="$emit('select', item)"
    @dblclick="$emit('double-click', item)"
  >
    <!-- 数字标签（按住 Command 时显示前 9 个） -->
    <div v-if="isCommandPressed && index < 9" class="quick-access-badge">
      {{ index + 1 }}
    </div>

    <!-- 图片类型：直接显示缩略图 -->
    <div
      v-if="item.ContentType === 'Image' && item.ImageData"
      class="minimal-image-wrapper"
    >
      <div class="minimal-image-container">
        <img :src="getImageDataURL(item)" class="minimal-image-thumbnail" />
        <div v-if="item.IsFavorite === 1" class="minimal-favorite-badge">
          <el-icon :size="14"><Star /></el-icon>
        </div>
      </div>
    </div>

    <!-- 其他类型：显示文本内容 -->
    <div v-else class="minimal-text-content">
      <span class="item-content">{{ item.Content }}</span>
      <div style="display: flex; margin-left: auto; align-items: center">
        <el-icon v-if="item.IsFavorite === 1" :size="16" style="color: #f5a623">
          <Star />
        </el-icon>
        <div
          v-if="item.ContentType === 'Color'"
          class="color-circle-small"
          :style="{ backgroundColor: item.Content }"
        ></div>
      </div>
    </div>

    <!-- 操作按钮（悬停显示） -->
    <div class="minimal-actions" @click.stop>
      <el-icon
        class="iconfont action-icon icon-copy icon-fasong"
        :size="18"
        @click.stop="$emit('send', item)"
        :title="$t('main.send')"
      >
      </el-icon>
      <el-icon
        class="action-icon icon-copy"
        :size="18"
        @click.stop="$emit('copy', item.ID)"
        :title="$t('main.copy')"
      >
        <CopyDocument />
      </el-icon>
      <el-icon
        class="action-icon icon-delete"
        :size="18"
        @click.stop="$emit('delete', item.ID)"
        :title="$t('main.delete')"
      >
        <Delete />
      </el-icon>
      <el-icon
        class="action-icon icon-favorite"
        :size="18"
        :class="{ active: item.IsFavorite === 1 }"
        @click.stop="$emit('collect', item.ID)"
        :title="$t('main.favorite')"
      >
        <Star />
      </el-icon>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useI18n } from "vue-i18n";
import { Star, CopyDocument, Delete } from "@element-plus/icons-vue";
import { common } from "../../../../wailsjs/go/models";

type ClipboardItem = common.ClipboardItem;

const { t } = useI18n();

defineProps<{
  item: ClipboardItem;
  index: number;
  currentItem: ClipboardItem | null;
  isCommandPressed: boolean;
}>();

defineEmits<{
  select: [item: ClipboardItem];
  "double-click": [item: ClipboardItem];
  copy: [id: string];
  delete: [id: string];
  collect: [id: string];
  send: [item: ClipboardItem];
}>();

// 获取图片的 data URL
function getImageDataURL(item: ClipboardItem): string {
  if (!item.ImageData) return "";
  const imageData = Array.isArray(item.ImageData)
    ? item.ImageData.map((b) => String.fromCharCode(b)).join("")
    : String(item.ImageData);
  return `data:image/png;base64,${imageData}`;
}
</script>

<style scoped>
@import "./styles/minimal-list-item.css";
</style>
