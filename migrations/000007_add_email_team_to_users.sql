-- 000007_add_email_team_to_users.sql
ALTER TABLE users
ADD COLUMN email VARCHAR(255) UNIQUE,
ADD COLUMN team VARCHAR(255);