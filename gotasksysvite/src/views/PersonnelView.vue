<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { usePersonnelStore } from '@/stores/personnel'
import { storeToRefs } from 'pinia'
import { ElRow, ElCol, ElCard, ElProgress, ElTag, ElDivider, ElEmpty } from 'element-plus'

const personnelStore = usePersonnelStore()
const { personnelList, isLoading } = storeToRefs(personnelStore)

// 定义状态灯颜色映射
const statusLightColors = {
  idle: '#67C23A',       // 绿色 - 空闲
  normal: '#409EFF',     // 蓝色 - 正常
  busy: '#E6A23C',         // 橙色 - 繁忙
  overloaded: '#F56C6C', // 红色 - 超负荷
}

// 定义进度条颜色规则
const progressColors = [
  { color: '#409EFF', percentage: 75 },
  { color: '#E6A23C', percentage: 100 },
  { color: '#F56C6C', percentage: 120 }, // 超过100%时显示红色
]

// 计算进度条百分比 (假设8小时为100%)
const calculatePercentage = (load: number) => {
  return Math.min((load / 8) * 100, 120) // 最大显示120%，给超负荷留出空间
}

onMounted(() => {
  personnelStore.fetchPersonnelStatus()
})
</script>

<template>
  <div>
    <h1>人员看板</h1>
    <p>实时查看团队每位成员的工作负载和进行中的任务。</p>
    
    <el-row :gutter="20">
      <el-col :span="24" :sm="12" :md="8" :lg="6" v-for="person in personnelList" :key="person.user_id">
        <el-card class="person-card">
          <div class="person-header">
            <div class="status-light" :style="{ backgroundColor: statusLightColors[person.status_light] }"></div>
            <span class="person-name">{{ person.real_name }}</span>
            <el-tag size="small" type="info">{{ person.role }}</el-tag>
            <el-tag v-if="person.has_overdue_task" type="danger" size="small" style="margin-left: auto;">有超期</el-tag>
          </div>
          
          <div class="load-info">
            <span>当前负载: {{ person.current_load }} 小时</span>
            <el-progress :percentage="calculatePercentage(person.current_load)" :color="progressColors" />
          </div>

          <el-divider />

          <div class="task-list">
            <strong>进行中的任务:</strong>
            <div v-if="person.active_tasks?.length">
              <ul>
                <li v-for="task in person.active_tasks" :key="task.id">{{ task.title }}</li>
              </ul>
            </div>
            <el-empty v-else description="暂无进行中任务" :image-size="50" />
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
.person-card {
  margin-bottom: 20px;
}
.person-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 16px;
}
.status-light {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  flex-shrink: 0;
}
.person-name {
  font-weight: bold;
  font-size: 16px;
}
.load-info {
  font-size: 14px;
  color: #606266;
}
.load-info .el-progress {
  margin-top: 8px;
}
.task-list {
  font-size: 14px;
}
.task-list ul {
  padding-left: 20px;
  margin: 8px 0 0;
}
.task-list li {
  margin-bottom: 4px;
}
</style>