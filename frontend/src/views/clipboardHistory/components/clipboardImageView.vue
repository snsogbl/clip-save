<template>
  <div class="image-content">
    <el-image
      :src="`data:image/png;base64,${imageData}`"
      alt="剪贴板图片"
      class="content-image"
      preview-teleported
      :preview-src-list="[`data:image/png;base64,${imageData}`]"
    />
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";

interface Props {
  imageData: string;
}

const props = defineProps<Props>();
</script>

<style scoped>
.image-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.content-image {
  max-width: 100%;
  height: auto;
  cursor: pointer;
  border-radius: 8px;
  transition: transform 0.2s ease;
}

.content-image:hover {
  transform: scale(1.02);
}

/* 模态框样式 */
.image-modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.8);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  backdrop-filter: blur(4px);
}

.modal-content {
  position: relative;
  max-width: 90%;
  max-height: 90%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal-image {
  max-width: 100%;
  max-height: 100%;
  border-radius: 8px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  transition: transform 0.2s ease;
  transform-origin: center center;
}

.close-btn {
  position: absolute;
  top: -40px;
  right: -10px;
  background: rgba(255, 255, 255, 0.9);
  border: none;
  border-radius: 50%;
  width: 32px;
  height: 32px;
  font-size: 20px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s ease;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 1);
}

/* 缩放控制样式 */
.zoom-controls {
  position: absolute;
  bottom: -50px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 8px;
  background: rgba(0, 0, 0, 0.8);
  padding: 8px 16px;
  border-radius: 20px;
  backdrop-filter: blur(4px);
}

.zoom-btn {
  background: rgba(255, 255, 255, 0.2);
  border: none;
  border-radius: 50%;
  width: 32px;
  height: 32px;
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  font-size: 12px;
}

.zoom-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.3);
  transform: scale(1.1);
}

.zoom-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.zoom-btn:last-child {
  width: auto;
  padding: 0 12px;
  border-radius: 16px;
  font-size: 12px;
}

.zoom-percentage {
  color: white;
  font-size: 12px;
  min-width: 40px;
  text-align: center;
  font-weight: 500;
}
</style>
