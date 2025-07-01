-- 000006_add_dates_to_periodic_tasks.sql
ALTER TABLE periodic_tasks
ADD COLUMN start_date TIMESTAMPTZ,
ADD COLUMN end_date TIMESTAMPTZ;