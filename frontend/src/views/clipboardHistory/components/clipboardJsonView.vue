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
import { CopyTextToClipboard } from "../../../../wailsjs/go/main/App";
import { ElMessage } from "element-plus";
import { useI18n } from "vue-i18n";
const { t } = useI18n();

const props = defineProps<{
  text: string;
}>();

const model = ref<any>({});

watch(
  () => props.text,
  (val) => {
    try {
      model.value = val ? val : {};
    } catch {
      model.value = {};
    }
  },
  { immediate: true }
);

async function copyEdited() {
  try {
    await CopyTextToClipboard(model.value);
    ElMessage.success(t("message.copySuccess"));
    console.log("已复制到剪贴板");
  } catch (error) {
    console.error("复制失败:", error);
    ElMessage.error(t("message.copyError", [error]));
  }
}

defineExpose({
  copyEdited
});
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
