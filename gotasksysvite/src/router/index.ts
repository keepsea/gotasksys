// src/router/index.ts

import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

import MainLayout from '@/layouts/MainLayout.vue'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import TasksView from '../views/TasksView.vue'
import PersonnelView from '../views/PersonnelView.vue' // <-- 1. 导入新页面

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView
    },
    {
      path: '/',
      component: MainLayout,
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'home',
          component: HomeView
        },
        {
          path: 'tasks', // 注意：嵌套路由的path不需要开头的'/'
          name: 'tasks',
          component: TasksView
        },
        // --- 2. 新增人员看板的路由 ---
        {
          path: 'personnel',
          name: 'personnel',
          component: PersonnelView
        }
        // -----------------------------
      ]
    }
  ]
})

// 导航守卫保持不变
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();
  const isAuthenticated = !!authStore.token;
  const requiresAuth = to.meta.requiresAuth;

  if (requiresAuth && !isAuthenticated) {
    next({ name: 'login' });
  } else if (to.name === 'login' && isAuthenticated) {
    next({ name: 'home' });
  } else {
    next();
  }
});

export default router