-- 000004_make_task_type_nullable.sql
ALTER TABLE tasks ALTER COLUMN task_type_id DROP NOT NULL;