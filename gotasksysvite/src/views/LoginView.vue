<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth' // 导入我们的auth store

const username = ref('admin') // 预填用户名方便测试
const password = ref('Password123!@#') // 预填密码方便测试

const authStore = useAuthStore() // 获取auth store的实例

const handleLogin = async () => {
  if (!username.value || !password.value) {
    alert('请输入用户名和密码');
    return;
  }
  // 调用store中的login动作
  await authStore.login(username.value, password.value)
}
</script>

<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <div class="card-header">
          <span>任务看板系统登录</span>
        </div>
      </template> <el-form @submit.prevent="handleLogin">
        <el-form-item label="用户名">
          <el-input v-model="username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" native-type="submit" style="width: 100%;">登 录</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f0f2f5;
}
.login-card {
  width: 400px;
}
.card-header {
  text-align: center;
  font-size: 20px;
  font-weight: bold;
}
</style>