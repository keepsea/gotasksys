<script setup lang="ts">
import { RouterView, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { ref } from 'vue'
import { Menu as IconMenu, ArrowDown } from '@element-plus/icons-vue'

const authStore = useAuthStore()
const router = useRouter()

const handleLogout = () => {
  authStore.logout()
}

// 模拟用户信息，我们稍后会从API获取
const currentUser = ref({
  realName: '管理员' 
})
</script>

<template>
  <el-container class="main-container">
    <el-aside width="200px">
      <el-menu default-active="/" class="el-menu-vertical-demo" router>
        <div class="logo">任务看板系统</div>
        <el-menu-item index="/">
          <el-icon><icon-menu /></el-icon>
          <span>首页/驾驶舱</span>
        </el-menu-item>
        <el-menu-item index="/tasks">
          <el-icon><icon-menu /></el-icon>
          <span>任务看板</span>
        </el-menu-item>
        <el-menu-item index="/personnel">
          <el-icon><icon-menu /></el-icon>
          <span>人员看板</span>
        </el-menu-item>
        <el-menu-item index="/admin">
          <el-icon><icon-menu /></el-icon>
          <span>系统管理</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    
    <el-container direction="vertical"> <el-header class="header">
        <div class="header-right">
          <el-dropdown>
            <span class="el-dropdown-link">
              {{ currentUser.realName }}
              <el-icon class="el-icon--right"><arrow-down /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="handleLogout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      <el-main>
        <RouterView />
      </el-main>
    
    </el-container>
  
  </el-container>
</template>

<style scoped>
.main-container, .el-container {
  height: 100vh;
}

.el-aside {
  background-color: #304156;
  box-shadow: 2px 0 6px rgba(0,21,41,.35);
  z-index: 10;
}

.logo {
  height: 60px;
  line-height: 60px;
  text-align: center;
  font-size: 20px;
  font-weight: bold;
  color: #fff;
  background-color: #2b2f3a;
}

.el-menu {
  height: calc(100% - 60px); /* 减去logo的高度 */
  background-color: #304156;
  border-right: none;
}

.el-menu-item {
  color: #bfcbd9;
}

.el-menu-item:hover {
  background-color: #263445 !important; /* !important 强制覆盖默认hover样式 */
}

.el-menu-item.is-active {
  color: #409eff !important;
  background-color: #263445 !important;
}

.header {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  background-color: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, .08);
  padding: 0 20px;
}

.el-dropdown-link {
  cursor: pointer;
  display: flex;
  align-items: center;
}

.el-main {
  background-color: #f0f2f5;
  padding: 20px;
}
</style>