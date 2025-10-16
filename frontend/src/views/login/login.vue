<template>
  <div class="login-container">
    <button class="close-btn" @click="hideApp" aria-label="关闭">×</button>
    <div class="login-box">
      <div class="logo-section">
        <div class="app-icon">🔒</div>
        <h1 class="app-name">剪存</h1>
        <p class="app-subtitle">请输入密码解锁</p>
      </div>

      <el-form @submit.prevent="handleLogin" class="login-form">
        <el-form-item>
          <el-input
            v-model="password"
            type="password"
            placeholder="请输入密码"
            size="large"
            show-password
            @keyup.enter="handleLogin"
            autofocus
          >
            <template #prefix>
              <el-icon><Lock /></el-icon>
            </template>
          </el-input>
        </el-form-item>

        <el-button
          type="primary"
          size="large"
          @click="handleLogin"
          :loading="loading"
          class="login-btn"
        >
          解锁
        </el-button>

        <div v-if="errorMessage" class="error-message">
          {{ errorMessage }}
        </div>
      </el-form>

      <div class="login-footer">
        <p class="hint-text">忘记密码？请删除数据库文件重置应用</p>
        <p class="hint-text">数据库位置: ~/.clipsave/clipboard.db</p>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { Lock } from '@element-plus/icons-vue';
import { Hide } from '../../../wailsjs/runtime/runtime';

// 定义事件
const emit = defineEmits(['unlock']);

const password = ref('');
const loading = ref(false);
const errorMessage = ref('');

// 处理登录
async function handleLogin() {
  if (!password.value) {
    errorMessage.value = '请输入密码';
    return;
  }

  loading.value = true;
  errorMessage.value = '';

  try {
    // 发送密码给父组件验证
    emit('unlock', password.value);
  } catch (error) {
    errorMessage.value = '验证失败，请重试';
  } finally {
    loading.value = false;
  }
}

function hideApp() {
  Hide();
}
</script>

<style scoped>
.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  position: relative;
}

.login-box {
  background: #ffffff;
  border-radius: 20px;
  padding: 48px 40px;
  width: 400px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.logo-section {
  text-align: center;
  margin-bottom: 40px;
}

.app-icon {
  font-size: 64px;
  margin-bottom: 16px;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
}

.app-name {
  margin: 0 0 8px 0;
  font-size: 32px;
  font-weight: 700;
  color: #1a1a1a;
  letter-spacing: 2px;
}

.app-subtitle {
  margin: 0;
  font-size: 14px;
  color: #8e8e93;
}

.login-form {
  margin-bottom: 24px;
}

.login-btn {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 600;
  margin-top: 8px;
}

.error-message {
  color: #ff3b30;
  font-size: 13px;
  text-align: center;
  margin-top: 16px;
  padding: 8px;
  background-color: #fff5f5;
  border-radius: 8px;
}

.login-footer {
  text-align: center;
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid #f0f0f0;
}

.hint-text {
  margin: 4px 0;
  font-size: 12px;
  color: #8e8e93;
}

.close-btn {
  position: absolute;
  top: 16px;
  left: 16px;
  width: 28px;
  height: 28px;
  border: none;
  background: transparent;
  color: #ffffff;
  font-size: 20px;
  line-height: 28px;
  text-align: center;
  border-radius: 6px;
  cursor: pointer;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.15);
  color: #ffffff;
}
</style>

