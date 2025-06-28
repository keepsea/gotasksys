-- 000001_init_schema.sql

-- 开启 pgcrypto 扩展, 用于生成 UUID
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- 用户表 (Users)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    real_name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL, -- 'system_admin', 'manager', 'member'
    avatar VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 任务类型配置表 (Task Types)
-- 我们将通过系统后台管理，所以也需要一张表
CREATE TABLE task_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    name VARCHAR(255) UNIQUE NOT NULL,
    is_enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 任务表 (Tasks)
CREATE TABLE tasks (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    task_type_id UUID REFERENCES task_types (id),
    status VARCHAR(50) NOT NULL,
    priority VARCHAR(50) NOT NULL,
    effort INT,
    original_effort INT,
    creator_id UUID REFERENCES users (id),
    reviewer_id UUID REFERENCES users (id),
    assignee_id UUID REFERENCES users (id),
    parent_task_id BIGINT REFERENCES tasks (id),
    evaluation JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    approved_at TIMESTAMPTZ,
    claimed_at TIMESTAMPTZ,
    due_date TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 创建索引以提高查询性能
CREATE INDEX idx_tasks_status ON tasks (status);

CREATE INDEX idx_tasks_assignee_id ON tasks (assignee_id);

CREATE INDEX idx_tasks_priority ON tasks (priority);

CREATE INDEX idx_tasks_creator_id ON tasks (creator_id);