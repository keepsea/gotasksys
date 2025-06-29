// src/stores/auth.ts

import { defineStore } from 'pinia';
import { ref } from 'vue';
import apiClient from '@/services/api'; // 导入我们创建的API客户端
import router from '@/router';

export const useAuthStore = defineStore('auth', () => {
    const token = ref(localStorage.getItem('token') || '');

    // 登录动作
    async function login(username: string, password: string) {
        try {
            // 1. 调用后端的登录API
            const response = await apiClient.post('/login', {
                username,
                password,
            });

            // 2. 从响应中获取token
            const newToken = response.data.token;
            token.value = newToken;

            // 3. 将token存入浏览器的localStorage，以便刷新页面后状态不丢失
            localStorage.setItem('token', newToken);

            // 4. 设置axios的默认请求头，之后的所有请求都会自动带上token
            apiClient.defaults.headers.common['Authorization'] = `Bearer ${newToken}`;

            // 5. 登录成功后，跳转到首页
            router.push('/');

        } catch (error) {
            // 简单处理错误，实际项目中可以更复杂
            console.error('Login failed:', error);
            alert('登录失败，请检查用户名和密码！');
        }
    }

    // 登出动作
    function logout() {
        token.value = '';
        localStorage.removeItem('token');
        delete apiClient.defaults.headers.common['Authorization'];
        router.push('/login');
    }

    return { token, login, logout };
});