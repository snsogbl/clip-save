<template>
  <div class="json-content">
    <JsonEditorVue
      v-model="model"
      :mode="Mode.text"
      :mainMenuBar="true"
      :navigationBar="true"
      :statusBar="true"
      :askToFormat="true"
      :readOnly="false"
    />
  </div>
</template>

<script lang="ts" setup>
import { ref, watch, onMounted } from "vue";
import JsonEditorVue from "json-editor-vue";
import { Mode } from "vanilla-jsoneditor";
import { ClipboardSetText } from "../../../../wailsjs/runtime/runtime";

const props = defineProps<{
  text: string;
}>();

const model = ref<any>({});

watch(
  () => props.text,
  (val) => {
    try {
      model.value = val ? JSON.parse(val) : {};
    } catch {
      model.value = {};
    }
  },
  { immediate: true }
);

async function copyEdited() {
  try {
    const str = JSON.stringify(model.value, null, 2);
    await ClipboardSetText(str);
  } catch (e) {}
}

onMounted(() => {});
</script>

<style scoped>
.json-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.json-actions {
  margin-top: 8px;
}
</style>
