<template>
  <el-dialog
    v-model="visible"
    :title="
      isEdit
        ? $t('settings.scripts.editScript')
        : $t('settings.scripts.newScript')
    "
    width="80%"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <el-form :model="form" label-width="120px" label-position="left">
      <el-form-item :label="$t('settings.scripts.name')" required>
        <el-input
          v-model="form.name"
          :placeholder="$t('settings.scripts.namePlaceholder')"
        />
      </el-form-item>

      <el-form-item :label="$t('settings.scripts.description')">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="2"
          :placeholder="$t('settings.scripts.descriptionPlaceholder')"
        />
      </el-form-item>

      <el-form-item :label="$t('settings.scripts.enabled')">
        <el-switch v-model="form.enabled" />
      </el-form-item>

      <el-form-item :label="$t('settings.scripts.trigger')" required>
        <el-select
          v-model="form.trigger"
          :placeholder="$t('settings.scripts.triggerPlaceholder')"
        >
          <el-option
            :label="$t('settings.scripts.triggerManual')"
            value="manual"
          />
          <el-option
            :label="$t('settings.scripts.triggerAfterSave')"
            value="after_save"
          />
        </el-select>
      </el-form-item>

      <el-form-item :label="$t('settings.scripts.contentTypes')">
        <el-select
          v-model="form.contentType"
          multiple
          :placeholder="$t('settings.scripts.contentTypesPlaceholder')"
          style="width: 100%"
        >
          <el-option label="Text" value="Text" />
          <el-option label="Image" value="Image" />
          <el-option label="File" value="File" />
          <el-option label="URL" value="URL" />
          <el-option label="Color" value="Color" />
          <el-option label="JSON" value="JSON" />
        </el-select>
        <div class="form-item-hint">
          {{ $t("settings.scripts.contentTypesHint") }}
        </div>
      </el-form-item>

      <el-form-item :label="$t('settings.scripts.keywords')">
        <el-input-tag
          v-model="form.keywords"
          draggable
          :placeholder="$t('settings.scripts.keywordsPlaceholder')"
          delimiter=","
        />
        <div class="form-item-hint">
          {{ $t("settings.scripts.keywordsHint") }}
        </div>
      </el-form-item>

      <el-form-item :label="$t('settings.scripts.script')" required>
        <div class="script-editor-container">
          <textarea
            ref="scriptEditor"
            v-model="form.script"
            class="script-editor"
            :placeholder="$t('settings.scripts.scriptPlaceholder')"
          />
        </div>
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="handleClose">{{ $t("common.cancel") }}</el-button>
      <el-button type="primary" @click="handleSave" :loading="saving">
        {{ $t("common.save") }}
      </el-button>
    </template>

  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";
import { ElMessage } from "element-plus";
import {
  GetUserScriptByID,
  SaveUserScript,
  GetClipboardItems,
} from "../../../../wailsjs/go/main/App";

interface Script {
  ID?: string;
  name: string;
  enabled: boolean;
  trigger: string;
  contentType: string[];
  keywords: string[];
  script: string;
  description: string;
  sortOrder?: number;
}

const props = defineProps<{
  modelValue: boolean;
  scriptId?: string;
}>();

const emit = defineEmits<{
  "update:modelValue": [value: boolean];
  saved: [];
}>();

const visible = ref(false);
const isEdit = ref(false);
const saving = ref(false);
const scriptEditor = ref<HTMLTextAreaElement | null>(null);

const form = ref<Script>({
  name: "",
  enabled: true,
  trigger: "manual",
  contentType: [],
  keywords: [],
  script: "",
  description: "",
});

watch(
  () => props.modelValue,
  (val) => {
    visible.value = val;
    if (val) {
      if (props.scriptId) {
        loadScript();
      } else {
        resetForm();
      }
    }
  }
);

watch(visible, (val) => {
  emit("update:modelValue", val);
});

async function loadScript() {
  if (!props.scriptId) {
    resetForm();
    return;
  }

  try {
    const script = await GetUserScriptByID(props.scriptId);
    if (script) {
      form.value = {
        ID: script.ID,
        name: script.Name || "",
        enabled: script.Enabled ?? true,
        trigger: script.Trigger || "manual",
        contentType: script.ContentType || [],
        keywords: script.Keywords || [],
        script: script.Script || "",
        description: script.Description || "",
        sortOrder: script.SortOrder || 0,
      };
      isEdit.value = true;
    } else {
      ElMessage.error("脚本不存在");
      handleClose();
    }
  } catch (error: any) {
    ElMessage.error(`加载脚本失败: ${error.message || error}`);
    handleClose();
  }
}

function resetForm() {
  form.value = {
    name: "",
    enabled: true,
    trigger: "manual",
    contentType: [],
    keywords: [],
    script: `// 示例脚本JavaScript代码
// item.Content - 剪贴板内容
// item.ContentType - 内容类型 (Text, Image, File, URL, Color, JSON)
// item.Timestamp - 时间戳
// item.Source - 来源应用

if (item.ContentType === 'Text') {
  // 返回结果
  return item.Content + ' 执行完成';
} else {
  return '不支持的类型';
}
`,
    description: "",
  };
  isEdit.value = false;
}

function handleClose() {
  visible.value = false;
  resetForm();
}

async function handleSave() {
  if (!form.value.name.trim()) {
    ElMessage.warning("请输入脚本名称");
    return;
  }

  if (!form.value.script.trim()) {
    ElMessage.warning("请输入脚本代码");
    return;
  }

  saving.value = true;

  try {
    const scriptData = {
      ID: form.value.ID || "",
      Name: form.value.name,
      Enabled: form.value.enabled,
      Trigger: form.value.trigger,
      ContentType: form.value.contentType,
      Keywords: form.value.keywords,
      Script: form.value.script,
      Description: form.value.description,
      SortOrder: form.value.sortOrder || 0,
    };

    await SaveUserScript(JSON.stringify(scriptData));
    ElMessage.success(isEdit.value ? "脚本已更新" : "脚本已创建");
    emit("saved");
    handleClose();
  } catch (error: any) {
    ElMessage.error(`保存脚本失败: ${error.message || error}`);
  } finally {
    saving.value = false;
  }
}

</script>

<style scoped>
.form-item-hint {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.keywords-tags {
  margin-top: 8px;
}

.script-editor-container {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  overflow: hidden;
  flex: 1;
}

.script-editor-toolbar {
  padding: 8px;
  background: #f5f7fa;
  border-bottom: 1px solid #dcdfe6;
  display: flex;
  gap: 8px;
}

.script-editor {
  width: 100%;
  min-height: 300px;
  padding: 12px;
  font-family: "Monaco", "Menlo", "Ubuntu Mono", "Consolas", monospace;
  font-size: 13px;
  line-height: 1.5;
  border: none;
  outline: none;
  resize: vertical;
  background: #fff;
}

.test-result {
  margin-top: 16px;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 4px;
}

.test-result h4 {
  margin-top: 0;
  margin-bottom: 12px;
}

.test-success {
  color: #67c23a;
}

.test-error {
  color: #f56c6c;
}

.test-result pre {
  background: #fff;
  padding: 8px;
  border-radius: 4px;
  overflow-x: auto;
  margin: 8px 0;
}

.test-result ul {
  margin: 8px 0;
  padding-left: 20px;
}
</style>
