// src/stores/dashboard.ts

import { ref } from 'vue'
import { defineStore } from 'pinia'
import apiClient from '@/services/api'

// 定义统计数据的结构
export interface DashboardSummary {
    pending_review_count: number
    in_pool_count: number
    in_progress_count: number
}

export const useDashboardStore = defineStore('dashboard', () => {
    // state: 用于存放统计数据
    const summaryData = ref<DashboardSummary | null>(null)
    const isLoading = ref(false)

    // action: 用于调用API获取数据
    async function fetchSummary() {
        isLoading.value = true
        try {
            const response = await apiClient.get<DashboardSummary>('/dashboard/summary')
            summaryData.value = response.data
        } catch (error) {
            console.error('Failed to fetch dashboard summary:', error)
            // 在这里可以添加错误提示
        } finally {
            isLoading.value = false
        }
    }

    return { summaryData, isLoading, fetchSummary }
})