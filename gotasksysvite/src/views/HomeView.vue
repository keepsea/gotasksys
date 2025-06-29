<script setup lang="ts">
import { onMounted } from 'vue'
import { useDashboardStore } from '@/stores/dashboard'
import { storeToRefs } from 'pinia'
import { ElRow, ElCol, ElCard, ElStatistic, ElSkeleton } from 'element-plus'

// 获取dashboard store的实例
const dashboardStore = useDashboardStore()

// 使用storeToRefs来保持响应性
const { summaryData, isLoading } = storeToRefs(dashboardStore)

// onMounted是一个“生命周期钩子”，它会在组件被加载到页面上后立即执行
onMounted(() => {
  dashboardStore.fetchSummary()
})
</script>

<template>
  <div>
    <h1>驾驶舱</h1>
    <p>欢迎回来！这是当前团队任务的整体概览。</p>
    
    <el-skeleton :loading="isLoading" animated>
      <template #template>
        <el-row :gutter="20">
          <el-col :span="8"><el-skeleton-item variant="p" style="width: 100%; height: 120px;" /></el-col>
          <el-col :span="8"><el-skeleton-item variant="p" style="width: 100%; height: 120px;" /></el-col>
          <el-col :span="8"><el-skeleton-item variant="p" style="width: 100%; height: 120px;" /></el-col>
        </el-row>
      </template>
      <template #default>
        <el-row v-if="summaryData" :gutter="20">
          <el-col :span="8">
            <el-card>
              <el-statistic title="待审批任务数" :value="summaryData.pending_review_count" />
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card>
              <el-statistic title="任务池任务数" :value="summaryData.in_pool_count" />
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card>
              <el-statistic title="进行中任务总数" :value="summaryData.in_progress_count" />
            </el-card>
          </el-col>
        </el-row>
      </template>
    </el-skeleton>
  </div>
</template>

<style scoped>
h1 {
  margin-bottom: 20px;
}
.el-card {
  text-align: center;
}
</style>