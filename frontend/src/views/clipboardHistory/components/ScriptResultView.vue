<template>
  <div v-if="lastResult" class="script-results">
    <div class="result-header">
      <el-tag :type="lastResult.error ? 'danger' : 'success'" size="small">
        {{ lastResult.error ? $t('scripts.executeError') : $t('scripts.executeSuccess') }}
      </el-tag>
      <span class="result-time">{{ formatTime(lastResult.timestamp) }}</span>
    </div>
    
    <!-- 错误信息 -->
    <div v-if="lastResult.error" class="result-error">
      <el-icon><Warning /></el-icon>
      <span>{{ lastResult.error }}</span>
    </div>
    
    <!-- 显示脚本返回值 -->
    <div v-else-if="lastResult.returnValue !== undefined" class="result-return">
      <div class="result-content">
        <pre class="return-value">{{ formatReturnValue(lastResult.returnValue) }}</pre>
      </div>
    </div>
    
    <!-- 如果没有返回值也没有错误，显示空状态 -->
    <div v-else class="result-empty">
      {{ $t('scripts.noReturnValue') }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Warning } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

interface ScriptExecutionResult {
  error?: string
  returnValue?: any
  timestamp: number
  scriptName?: string
}

const props = defineProps<{
  itemId: string
  results: ScriptExecutionResult[]
}>()

// 只显示最后一次执行的结果
const lastResult = computed(() => {
  if (!props.results || props.results.length === 0) {
    return null
  }
  return props.results[props.results.length - 1]
})

// 格式化时间
function formatTime(timestamp: number): string {
  const date = new Date(timestamp)
  return date.toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  })
}

// 格式化返回值
function formatReturnValue(value: any): string {
  if (value === null) {
    return 'null'
  }
  if (value === undefined) {
    return 'undefined'
  }
  if (typeof value === 'string') {
    return value
  }
  if (typeof value === 'object') {
    try {
      return JSON.stringify(value, null, 2)
    } catch {
      return String(value)
    }
  }
  return String(value)
}
</script>

<style scoped>
.script-results {
  margin: 12px 20px 0px 8px;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  background-color: #fff;
  padding: 12px;
}

.result-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.result-time {
  font-size: 12px;
  color: #999;
}

.result-error {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px;
  background-color: #fef0f0;
  border-left: 3px solid #f56c6c;
  border-radius: 4px;
  color: #f56c6c;
  margin-top: 8px;
}

.result-return {
  margin-top: 12px;
}

.result-content {
  background-color: #f5f7fa;
  border-radius: 4px;
  padding: 12px;
  max-height: 150px;
  overflow-y: auto;
}

.return-value {
  margin: 0;
  font-size: 13px;
  color: #303133;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
}

.result-empty {
  margin-top: 12px;
  padding: 12px;
  text-align: center;
  color: #909399;
  font-size: 13px;
}
</style>

