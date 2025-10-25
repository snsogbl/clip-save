import {createApp} from 'vue'
import App from './App.vue'
import './style.css';
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import i18n from './i18n'

createApp(App).use(ElementPlus).use(i18n).mount('#app')
