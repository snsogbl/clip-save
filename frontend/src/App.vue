<script lang="ts" setup>
import { ref, onMounted } from "vue";
import { EventsEmit } from "../wailsjs/runtime/runtime";
import ClipboardHistory from "./views/clipboardHistory/clipboardHistory.vue";
import Login from "./views/login/login.vue";
import {
  GetAppSettings,
  VerifyPassword,
  HideWindow,
} from "../wailsjs/go/main/App";
import { ElMessage } from "element-plus";
import { useI18n } from 'vue-i18n';
import "highlight.js/styles/github.css";

const { t } = useI18n();

const isLocked = ref(true);
const isLoading = ref(true);

// 检查是否设置了密码
async function checkPassword() {
  try {
    const settings = await GetAppSettings();
    if (settings) {
      const parsed = JSON.parse(settings);
      // 如果没有设置密码或密码为空，直接解锁
      if (!parsed.password || parsed.password === "") {
        isLocked.value = false;
        console.log("📖 未设置密码，直接进入应用");
      } else {
        console.log("🔒 应用已锁定，需要密码");
      }
    } else {
      // 没有设置，直接解锁
      isLocked.value = false;
    }
  } catch (error) {
    console.error("检查密码失败:", error);
    // 出错时直接解锁，避免用户被锁在外面
    isLocked.value = false;
  } finally {
    isLoading.value = false;
  }
}

// 验证密码
async function handleUnlock(password: string) {
  try {
    const isValid = await VerifyPassword(password);
    if (isValid) {
      isLocked.value = false;
      ElMessage.success(t('login.unlockSuccess'));
      console.log("✅ 密码验证成功");
    } else {
      ElMessage.error(t('login.passwordError'));
      console.log("❌ 密码验证失败");
    }
  } catch (error) {
    ElMessage.error(t('login.verifyError', [error]));
    console.error("验证密码失败:", error);
  }
}

const addKeyListener = () => {
  document.addEventListener("keydown", (event) => {
    // 当图片预览(ElImage Viewer)打开时，按 Esc 不隐藏窗口
    const hasImagePreview = !!document.querySelector('.el-image-viewer__wrapper');
    const hasDialog = !!document.querySelector('.el-overlay');
    const shouldSuppress = (window as any).__suppressHideWindow || hasImagePreview || hasDialog;
    if ((event.key === 'Escape' || event.keyCode === 27) && !shouldSuppress) {
      HideWindow();
    }
    if ((event.metaKey || event.ctrlKey) && event.key === "w") {
      event.preventDefault();
      HideWindow();
    }

    // 拦截 ⌘+↑ / ⌘+↓，避免列表滚动，并触发上一条/下一条
    if ((event.metaKey || event.ctrlKey) && (event.key === "ArrowUp" || event.key === "ArrowDown")) {
      event.preventDefault();
      if (event.key === "ArrowUp") {
        EventsEmit("nav.prev");
      } else {
        EventsEmit("nav.next");
      }
    }
  });
  // window.addEventListener("blur", (event) => {
  //   // 当有系统对话框（如保存文件）弹出时，不要自动隐藏
  //   // 使用全局标记进行抑制
  //   const shouldSuppress = (window as any).__suppressHideWindow;
  //   if (shouldSuppress) return;
  //   HideWindow();
  // });
};

onMounted(() => {
  checkPassword();
  addKeyListener();
});
</script>

<template>
  <div style="--wails-draggable: drag;">
    <!-- <div style="width: 100px;height: 100px;background-color: antiquewhite;"></div> -->
    <div v-if="isLoading" class="loading-screen">
      <div class="loading-spinner"></div>
    </div>
    <Login v-else-if="isLocked" @unlock="handleUnlock" />
    <ClipboardHistory v-else />
  </div>
</template>

<style>
@import "/src/assets/sass/iconfont.css";
@import "/src/assets/sass/theme.css";

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

#app {
  height: 100vh;
  overflow: hidden;
  /* border-radius: 8px; */
  background-color: #fff;
  /* background-color: rgba(255, 255, 255, 1);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px); */
}

.loading-screen {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.loading-spinner {
  width: 50px;
  height: 50px;
  border: 4px solid rgba(255, 255, 255, 0.3);
  border-top-color: #ffffff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
