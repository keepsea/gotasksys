// src/stores/taskType.ts
import { ref } from 'vue'
import { defineStore } from 'pinia'
import apiClient from '@/services/api'

// 定义任务类型的结构
export interface TaskType {
    id: string;
    name: string;
    is_enabled: boolean;
}

export const useTaskTypeStore = defineStore('taskType', () => {
    const taskTypes = ref<TaskType[]>([])
    const isLoading = ref(false)

    async function fetchTaskTypes() {
        // 如果已经加载过，就不再重复加载
        if (taskTypes.value.length > 0) return;

        isLoading.value = true
        try {
            // 注意：这里我们调用的是新的、非管理员的接口
            const response = await apiClient.get<TaskType[]>('/task-types')
            // 只保留启用的类型
            taskTypes.value = response.data.filter(t => t.is_enabled)
        } catch (error) {
            console.error('Failed to fetch task types:', error)
        } finally {
            isLoading.value = false
        }
    }

    return { taskTypes, isLoading, fetchTaskTypes }
})