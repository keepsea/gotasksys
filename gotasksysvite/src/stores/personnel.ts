// src/stores/personnel.ts
import { ref } from 'vue'
import { defineStore } from 'pinia'
import apiClient from '@/services/api'

// 定义从后端获取到的单个人员状态的数据结构
export interface PersonnelStatus {
    user_id: string;
    real_name: string;
    role: string;
    current_load: number;
    status_light: 'idle' | 'normal' | 'busy' | 'overloaded';
    active_tasks: {
        id: number;
        title: string;
    }[];
    has_overdue_task: boolean;
}

export const usePersonnelStore = defineStore('personnel', () => {
    const personnelList = ref<PersonnelStatus[]>([])
    const isLoading = ref(false)

    async function fetchPersonnelStatus() {
        isLoading.value = true
        try {
            const response = await apiClient.get<PersonnelStatus[]>('/personnel/status')
            personnelList.value = response.data
        } catch (error) {
            console.error('Failed to fetch personnel status:', error)
        } finally {
            isLoading.value = false
        }
    }

    return { personnelList, isLoading, fetchPersonnelStatus }
})