// src/model/task.ts
export interface Task {
    id: number
    title: string
    description: string
    task_type_id: string
    status: string
    priority: string
    effort: number
    original_effort: number
    creator_id: string
    reviewer_id: string | null
    assignee_id: string | null
    parent_task_id: number | null
    evaluation: any | null // 简化处理
    created_at: string
    approved_at: string | null
    claimed_at: string | null
    due_date: string | null
    completed_at: string | null
    updated_at: string
}