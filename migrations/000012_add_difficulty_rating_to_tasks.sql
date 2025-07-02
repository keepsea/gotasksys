-- 000012_add_difficulty_rating_to_tasks.sql
ALTER TABLE tasks ADD COLUMN difficulty_rating JSONB;