<template>
  <div class="hotkey-display-container">
    <div 
      v-for="(key, index) in parsedKeys" 
      :key="index"
      class="hotkey-key"
    >
      {{ key }}
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue';

interface Props {
  hotkey: string;
}

const props = defineProps<Props>();

const parsedKeys = computed(() => {
  if (!props.hotkey) return [];
  
  const keys = props.hotkey.split('+').map(key => key.trim());
  
  return keys.map(key => {
    // 转换修饰键显示
    switch (key.toLowerCase()) {
      case 'control':
        return '⌃';
      case 'command':
      case 'meta':
        return '⌘';
      case 'shift':
        return '⇧';
      case 'alt':
      case 'option':
        return '⌥';
      default:
        return key.toUpperCase();
    }
  });
});
</script>

<style scoped>
.hotkey-display-container {
  display: flex;
  align-items: center;
  gap: 4px;
}

.hotkey-key {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 24px;
  height: 24px;
  padding: 0 6px;
  background-color: #f0f0f0;
  border: 1px solid #d0d0d0;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
  color: #333;
  line-height: 1;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.hotkey-key:first-child {
  margin-left: 0;
}

.hotkey-key:last-child {
  margin-right: 0;
}
</style>
