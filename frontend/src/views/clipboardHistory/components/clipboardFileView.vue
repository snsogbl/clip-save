<template>
  <div class="file-content">
    <div
      v-for="(fileInfo, index) in files"
      :key="index"
      class="file-item-detail"
    >
      <div class="file-item-header">
        <el-icon :size="32" class="file-icon-large">
          <Folder v-if="fileInfo.is_dir" />
          <Document v-else />
        </el-icon>
        <div class="file-item-info">
          <div class="file-name">{{ fileInfo.name }}</div>
          <div class="file-meta">
            <span v-if="!fileInfo.is_dir">{{
              formatFileSize(fileInfo.size)
            }}</span>
            <span v-if="!fileInfo.exists" class="file-not-exists"
              >{{ $t('components.file.fileNotExists') }}</span
            >
          </div>
        </div>
      </div>
      <div class="file-path">{{ fileInfo.path }}</div>
      <el-button
        round
        v-if="fileInfo.exists"
        class="me-button"
        @click="handleOpenFile(fileInfo.path)"
      >
        <el-icon :size="14" style="margin-right: 4px">
          <FolderOpened />
        </el-icon>
        {{ $t('components.file.openInFinder') }}
      </el-button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useI18n } from 'vue-i18n';
import { Document, Folder, FolderOpened } from "@element-plus/icons-vue";

const { t } = useI18n();

interface FileInfo {
  name: string;
  path: string;
  size: number;
  is_dir: boolean;
  exists: boolean;
  extension: string;
}

interface Props {
  files: FileInfo[];
}

defineProps<Props>();

const emit = defineEmits<{
  openFile: [path: string];
}>();

// 格式化文件大小
function formatFileSize(size: number): string {
  const units = ["B", "KB", "MB", "GB", "TB"];
  let unitIndex = 0;
  let fileSize = size;

  while (fileSize >= 1024 && unitIndex < units.length - 1) {
    fileSize /= 1024;
    unitIndex++;
  }

  return `${fileSize.toFixed(1)} ${units[unitIndex]}`;
}

function handleOpenFile(path: string) {
  (window as any).__suppressHideWindow = true;
  emit("openFile", path);
  setTimeout(() => {
    (window as any).__suppressHideWindow = false;
  }, 200);
}
</script>

<style scoped>
/* 文件展示样式 */
.file-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.file-item-detail {
  padding: 20px;
  background-color: #ffffff;
  border-radius: 12px;
  border: 1px solid #e0e0e0;
  transition: all 0.2s ease;
}

.file-item-detail:hover {
  border-color: #2196f3;
  box-shadow: 0 2px 8px rgba(33, 150, 243, 0.1);
}

.file-item-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.file-icon-large {
  color: #2196f3;
  flex-shrink: 0;
}

.file-item-info {
  flex: 1;
  min-width: 0;
}

.file-name {
  font-size: 16px;
  font-weight: 600;
  color: #1a1a1a;
  word-wrap: break-word;
  margin-bottom: 4px;
}

.file-meta {
  font-size: 13px;
  color: #8e8e93;
  display: flex;
  gap: 8px;
  align-items: center;
}

.file-not-exists {
  color: #ff3b30;
  font-weight: 500;
}

.file-path {
  font-size: 12px;
  color: #6d6d70;
  background-color: #f8f8f8;
  padding: 8px 12px;
  border-radius: 6px;
  word-break: break-all;
  margin-bottom: 12px;
}
</style>

