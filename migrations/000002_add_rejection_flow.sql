-- 000002_add_rejection_flow.sql
ALTER TABLE tasks ADD COLUMN rejection_reason TEXT;

-- 注意：下面的命令是针对PostgreSQL的，用于给status字段增加一个新的枚举值。
-- 如果是其他数据库，可能语法不同。
-- 为了简化，我们也可以不在数据库层面做强校验，而是在代码层面保证status值的合法性。
-- 此处我们暂时不修改字段类型，仅在代码中新增逻辑。