-- 000009_create_system_avatars.sql
CREATE TABLE system_avatars (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    url VARCHAR(255) NOT NULL UNIQUE,
    description VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE system_avatars
ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT now();