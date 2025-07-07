-- 000013_create_leaves_table.sql
CREATE TABLE leaves (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    reason TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'approved', -- V1.0 简化，暂无审批流程
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 为 user_id 和日期创建索引以优化查询
CREATE INDEX idx_leaves_user_id ON leaves (user_id);

CREATE INDEX idx_leaves_start_date_end_date ON leaves (start_date, end_date);