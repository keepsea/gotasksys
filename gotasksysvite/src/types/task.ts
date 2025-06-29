// src/types/task.ts
// 这个文件用于存放所有与Task相关的、可复用的TypeScript类型定义

export interface CreateTaskInput {
    title: string;
    description: string;
    priority: string;
    effort: number;
    task_type_id: string;
}