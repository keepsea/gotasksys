// src/main.ts (带有“记忆恢复”功能的版本)

import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

import App from './App.vue'
import router from './router'

// --- 新增的导入 ---
import { useAuthStore } from './stores/auth'
import apiClient from './services/api'
// -----------------

const app = createApp(App)

app.use(createPinia())

// === 新增：应用启动时，检查并恢复登录状态 ===
// 这个操作必须在 app.use(router) 之前，以确保路由守卫能获取到正确的登录状态
const authStore = useAuthStore()
if (authStore.token) {
    apiClient.defaults.headers.common['Authorization'] = `Bearer ${authStore.token}`
}
// ===========================================

app.use(router)
app.use(ElementPlus)

app.mount('#app')