<script setup lang="ts">
import { onMounted, computed, ref } from 'vue'
import { useTaskStore, type TaskStatus } from '@/stores/task' // 确保导入TaskStatus类型
import { storeToRefs } from 'pinia'
import { ElCard, ElDialog, ElButton } from 'element-plus'
import draggable from 'vuedraggable'
import TaskForm from '@/components/TaskForm.vue'
import type { CreateTaskInput } from '@/types/task'

const taskStore = useTaskStore()
const { tasksByStatus, isLoading } = storeToRefs(taskStore)

const isDialogVisible = ref(false)

// --- 关键的类型修正 ---
// 1. 定义一个接口来描述列的结构
interface StatusColumn {
  key: TaskStatus;
  title: string;
}

// 2. 明确告诉TypeScript，statusColumns是一个 StatusColumn 类型的数组
const statusColumns = computed<StatusColumn[]>(() => [
  { key: 'pending_review', title: '待审核' },
  { key: 'in_pool', title: '任务池' },
  { key: 'in_progress', title: '进行中' },
  { key: 'pending_evaluation', title: '待评价' },
  { key: 'completed', title: '已完成' },
])
// --------------------

const handleTaskDrop = (event: any) => {
  const taskId = event.item.dataset.taskId
  const newStatus = event.to.dataset.status as TaskStatus
  if (!newStatus || !taskId) return
  taskStore.moveTask(Number(taskId), newStatus)
}

const handleCreateTask = async (formData: CreateTaskInput) => {
  await taskStore.createTask(formData)
  isDialogVisible.value = false
}

onMounted(() => {
  taskStore.fetchTasks()
})
</script>

<template>
  <div class="task-board-container">
    <div class="board-header">
      <h1>任务看板</h1>
      <el-button type="primary" @click="isDialogVisible = true">创建新任务</el-button>
    </div>

    <div class="board-columns-wrapper">
      <div v-for="column in statusColumns" :key="column.key" class="board-column">
        <h3 class="column-title">
          {{ column.title }}
          <span class="task-count">{{ tasksByStatus[column.key]?.length || 0 }}</span>
        </h3>
        <draggable
          v-model="tasksByStatus[column.key]"
          class="task-list"
          group="tasks"
          itemKey="id"
          @end="handleTaskDrop"
          :data-status="column.key"
        >
          <template #item="{ element: task }">
            <el-card :data-task-id="task.id" class="task-card shadow-never">
              <p>{{ task.title }}</p>
              <div class="task-meta">
                <span>优先级: {{ task.priority }}</span>
              </div>
            </el-card>
          </template>
        </draggable>
        <div v-if="!tasksByStatus[column.key]?.length && !isLoading" class="no-task">
          无任务
        </div>
      </div>
    </div>
    <el-dialog v-model="isDialogVisible" title="创建新任务" width="50%">
      <TaskForm @submit="handleCreateTask" />
    </el-dialog>
  </div>
</template>

<style scoped>
.board-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  flex-shrink: 0;
}
h1 {
  margin: 0;
}
.task-board-container {
  display: flex;
  flex-direction: column;
  height: 100%;
}
.board-columns-wrapper {
  display: flex;
  flex-grow: 1;
  gap: 16px;
  min-height: 0;
}
.board-column {
  flex: 1 1 0;
  min-width: 280px;
  background-color: #f4f5f7;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
}
.column-title {
  padding: 12px 16px;
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  flex-shrink: 0;
  border-bottom: 1px solid #ebeef5;
  color: #303133;
}
.task-count {
  color: #909399;
  font-weight: normal;
  font-size: 14px;
  margin-left: 8px;
}
.task-list {
  flex-grow: 1;
  overflow-y: auto;
  padding: 8px;
  min-height: 50px;
}
.task-list::-webkit-scrollbar {
    width: 6px;
}
.task-list::-webkit-scrollbar-thumb {
    background: #ccc;
    border-radius: 3px;
}
.task-list::-webkit-scrollbar-thumb:hover {
    background: #aaa;
}
.task-card {
  margin-bottom: 10px;
  cursor: grab;
  border: 1px solid #e9e9eb;
}
.task-card:hover {
  border-color: #c0c4cc;
}
.task-meta {
  font-size: 12px;
  color: #888;
  margin-top: 10px;
}
.no-task {
  text-align: center;
  color: #999;
  padding: 20px;
  font-size: 14px;
}
</style>