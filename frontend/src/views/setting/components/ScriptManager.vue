<template>
  <el-dialog
    v-model="visible"
    :title="$t('settings.scripts.title')"
    width="80%"
    :close-on-click-modal="false"
  >
    <div class="script-manager">
      <div class="script-manager-header">
        <el-button class="me-button" size="small" @click="handleFindScripts">
          {{ $t('settings.scripts.findScripts') }}
        </el-button>
        <el-button size="small" type="primary" @click="handleNewScript">
          {{ $t('settings.scripts.newScript') }}
        </el-button>
      </div>

      <el-table :data="scripts" style="width: 100%" v-loading="loading" row-key="ID">
        <el-table-column :label="$t('settings.scripts.order')" width="86" fixed="left">
          <template #default="{ row }">
            <el-input
              v-model="row.SortOrder"
              @change="handleOrderChange(row)"
              style="width: 100%"
            />
          </template>
        </el-table-column>
        <el-table-column prop="Name" :label="$t('settings.scripts.name')" width="120" show-overflow-tooltip/>
        <el-table-column prop="Description" :label="$t('settings.scripts.description')" show-overflow-tooltip/>
        <el-table-column prop="Trigger" :label="$t('settings.scripts.trigger')" width="150">
          <template #default="{ row }">
            <span v-if="row.Trigger === 'after_save'">{{ $t('settings.scripts.triggerAfterSave') }}</span>
            <span v-else-if="row.Trigger === 'manual'">{{ $t('settings.scripts.triggerManual') }}</span>
            <span v-else>{{ row.Trigger }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="Enabled" :label="$t('settings.scripts.enabled')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.Enabled ? 'success' : 'info'">
              {{ row.Enabled ? $t('common.enabled') : $t('common.disabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('common.actions')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEditScript(row.ID)">
              {{ $t('common.edit') }}
            </el-button>
            <el-button size="small" type="danger" @click="handleDeleteScript(row.ID, row.Name)">
              {{ $t('common.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 脚本编辑器 -->
    <ScriptEditor
      v-model="showScriptEditor"
      :script-id="editingScriptId"
      @saved="handleScriptSaved"
    />
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Link } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import { GetAllUserScripts, DeleteUserScript, UpdateUserScriptOrder, OpenURL } from '../../../../wailsjs/go/main/App'
import { common } from '../../../../wailsjs/go/models'
import ScriptEditor from './ScriptEditor.vue'

const { t } = useI18n()

const props = defineProps<{
  modelValue: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

const visible = ref(false)
const loading = ref(false)
const scripts = ref<common.UserScript[]>([])
const showScriptEditor = ref(false)
const editingScriptId = ref<string | undefined>()

watch(
  () => props.modelValue,
  (val) => {
    visible.value = val
    if (val) {
      loadScripts()
    }
  }
)

watch(visible, (val) => {
  emit('update:modelValue', val)
})

async function loadScripts() {
  loading.value = true
  try {
    scripts.value = await GetAllUserScripts()
  } catch (error: any) {
    ElMessage.error(`加载脚本失败: ${error.message || error}`)
  } finally {
    loading.value = false
  }
}

function handleNewScript() {
  editingScriptId.value = undefined
  showScriptEditor.value = true
}

async function handleFindScripts() {
  try {
    await OpenURL('https://github.com/snsogbl/clip-save/tree/main/scriptingExample')
  } catch (error: any) {
    ElMessage.error(`打开链接失败: ${error.message || error}`)
  }
}

function handleEditScript(id: string) {
  editingScriptId.value = id
  showScriptEditor.value = true
}

async function handleDeleteScript(id: string, name: string) {
  try {
    await ElMessageBox.confirm(
      t('settings.scripts.deleteConfirm', { name }) || `确定要删除脚本 "${name}" 吗？`,
      t('settings.scripts.deleteTitle') || '删除脚本',
      {
        confirmButtonText: t('common.delete') || '删除',
        cancelButtonText: t('common.cancel') || '取消',
        type: 'warning',
      }
    )

    await DeleteUserScript(id)
    ElMessage.success(t('settings.scripts.deleteSuccess') || '脚本已删除')
    await loadScripts()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(`${t('settings.scripts.deleteError') || '删除脚本失败'}: ${error.message || error}`)
    }
  }
}

function handleScriptSaved() {
  loadScripts()
}

let saveTimer: ReturnType<typeof setTimeout> | null = null

async function handleOrderChange(row: common.UserScript) {
  // changedScript.SortOrder 检测不是数字，则不保存
  if (isNaN(parseInt(row.SortOrder.toString()))) {
    ElMessage.error(t('settings.scripts.orderInvalid') || '顺序无效')
    return
  }
  await saveOrder(row)
}

async function saveOrder(changedScript: common.UserScript) {
  try {
    const sortOrder = parseInt(changedScript.SortOrder.toString()) || 0
    await UpdateUserScriptOrder(changedScript.ID, sortOrder)
    ElMessage.success(t('settings.scripts.orderUpdated') || '顺序已更新')
  } catch (error: any) {
    ElMessage.error(`${t('settings.scripts.orderUpdateError') || '更新顺序失败'}: ${error.message || error}`)
    // 重新加载以恢复原始顺序
    await loadScripts()
  }
}
</script>

<style scoped>
.script-manager {
  min-height: 400px;
}

.script-manager-header {
  margin-bottom: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>

