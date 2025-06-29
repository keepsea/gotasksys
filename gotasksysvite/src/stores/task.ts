// src/stores/task.ts (重构版)

import { ref } from 'vue'
import { defineStore } from 'pinia'
import apiClient from '@/services/api'
import type { Task } from '@/model/task'
import type { CreateTaskInput } from '@/types/task' // 我们将很快创建这个类型文件

// 定义状态列的类型
export type TaskStatus = 'pending_review' | 'in_pool' | 'in_progress' | 'pending_evaluation' | 'completed'
type GroupedTasks = Record<TaskStatus, Task[]>

export const useTaskStore = defineStore('task', () => {
    // --- 状态变更 ---
    // 将 tasksByStatus 从 computed 改为 ref，使其可被直接修改
    const tasksByStatus = ref<GroupedTasks>({
        pending_review: [],
        in_pool: [],
        in_progress: [],
        pending_evaluation: [],
        completed: [],
    })
    const isLoading = ref(false)

    // --- Actions ---
    // === 使用下面这段全新的 fetchTasks 函数 ===
    async function fetchTasks() {
        isLoading.value = true;
        try {
            // 1. 从后端API获取原始的任务列表
            const response = await apiClient.get<Task[]>('/tasks');
            const allTasks = response.data;
            // --- 探针 1: 检查我们是否收到了原始数据 ---
            console.log("[taskStore] 1. 从API收到的原始数据:", allTasks);
            // ------------------------------------------
            // 2. 在内存中创建一个全新的、空的、已分组的对象
            const newGroupedTasks: GroupedTasks = {
                pending_review: [],
                in_pool: [],
                in_progress: [],
                pending_evaluation: [],
                completed: [],
            };

            // 3. 遍历原始列表，将每个任务放入新分组对象的对应数组中
            for (const task of allTasks) {
                // 确保任务状态是已知的，防止意外错误
                if (newGroupedTasks[task.status as TaskStatus]) {
                    newGroupedTasks[task.status as TaskStatus].push(task);
                }
            }

            // 4. 最关键的一步：用这个全新的、已经分类好的对象，去完整替换掉旧的看板数据
            tasksByStatus.value = newGroupedTasks;
            // --- 探针 2: 检查store里的数据是否已正确分组更新 ---
            console.log("[taskStore] 2. 已分组并更新到store的数据:", tasksByStatus.value);
            // ----------------------------------------------------

        } catch (error) {
            console.error('Failed to fetch tasks:', error);
        } finally {
            isLoading.value = false;
        }
    }

    // 拖拽成功后，调用这个函数来通知后端
    async function moveTask(taskId: number, newStatus: TaskStatus) {
        // 找出对应的API端点
        const apiEndpoints: Partial<Record<TaskStatus, string>> = {
            in_pool: 'approve',
            in_progress: 'claim',
            pending_evaluation: 'complete',
        }
        const endpoint = apiEndpoints[newStatus]

        if (!endpoint) {
            console.error(`No API endpoint defined for status: ${newStatus}`)
            await fetchTasks() // 如果是无效的移动，则刷新整个看板来撤销前端的移动
            return
        }

        try {
            // 只调用API，不再手动修改前端状态，因为vuedraggable已经帮我们完成了
            await apiClient.post(`/tasks/${taskId}/${endpoint}`)
        } catch (error) {
            console.error(`Failed to move task ${taskId} to ${newStatus}:`, error)
            alert(`操作失败: ${error}`)
            // 如果API调用失败，则刷新整个看板，让卡片回到原来的位置
            await fetchTasks()
        }
    }
    // --- 新增的创建任务动作 ---
    async function createTask(taskData: CreateTaskInput) {
        try {
            // 1. 调用API，并接收后端返回的新创建的任务对象
            const response = await apiClient.post<Task>('/tasks', taskData);
            const newTask = response.data;

            // 2. 直接将这个新任务，添加到"待审核"列的数组中
            if (tasksByStatus.value.pending_review) {
                tasksByStatus.value.pending_review.push(newTask);
            } else {
                // 如果因为某些原因该列不存在，则刷新整个列表作为后备方案
                await fetchTasks();
            }

        } catch (error) {
            console.error('Failed to create task:', error)
            alert(`创建任务失败: ${error}`);
            // 如果创建失败，也刷新一下列表以确保数据一致性
            await fetchTasks();
        }
    }
    return {
        isLoading,
        tasksByStatus, // 现在导出的是一个ref
        fetchTasks,
        moveTask,
        createTask
    }
})