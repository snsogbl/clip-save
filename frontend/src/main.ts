import {createApp} from 'vue'
import App from './App.vue'
import './style.css';
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
// 在 Element Plus CSS 之后导入 theme.css，确保自定义变量覆盖默认值
import './assets/sass/theme.css'
import i18n from './i18n'
import { initScriptExecutor } from './scripts/executor'

// 初始化脚本执行器
initScriptExecutor()

createApp(App).use(ElementPlus).use(i18n).mount('#app')
