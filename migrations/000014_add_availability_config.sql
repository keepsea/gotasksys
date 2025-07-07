-- 000014_add_availability_config.sql

-- 1. 创建系统配置表，用于存储键值对配置
CREATE TABLE system_configs (
    config_key VARCHAR(255) PRIMARY KEY,
    config_value VARCHAR(255) NOT NULL,
    description TEXT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 插入一条默认的全局每日工时配置
INSERT INTO
    system_configs (
        config_key,
        config_value,
        description
    )
VALUES (
        'global_daily_work_hours',
        '8.0',
        '全局默认的每日可用工作时长（小时）'
    );

-- 2. 创建法定节假日表
CREATE TABLE holidays (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    holiday_date DATE NOT NULL UNIQUE,
    description VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 3. 为users表增加个人每日可用工时字段
ALTER TABLE users ADD COLUMN daily_capacity_hours FLOAT;