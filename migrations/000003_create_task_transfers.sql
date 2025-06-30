-- 000003_create_task_transfers.sql
CREATE TABLE task_transfers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    task_id BIGINT NOT NULL REFERENCES tasks (id),
    from_user_id UUID NOT NULL REFERENCES users (id),
    to_user_id UUID NOT NULL REFERENCES users (id),
    effort_spent_by_initiator INT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- 'pending', 'accepted', 'rejected'
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);