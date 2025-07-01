-- 000005_create_periodic_tasks.sql
CREATE TABLE periodic_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    -- Cron表达式，用于定义周期，例如 '0 9 * * 1' 代表每周一上午9点
    cron_expression VARCHAR(255) NOT NULL,
    -- 自动生成的任务的默认属性
    default_assignee_id UUID REFERENCES users (id),
    default_effort INT NOT NULL,
    default_priority VARCHAR(50) NOT NULL,
    default_task_type_id UUID REFERENCES task_types (id),
    is_active BOOLEAN NOT NULL DEFAULT TRUE, -- 用于控制启停
    created_by_id UUID NOT NULL REFERENCES users (id), -- 记录此计划由谁创建
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);