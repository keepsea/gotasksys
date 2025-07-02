-- 000011_create_system_avatars.sql
CREATE TABLE system_avatars (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    url VARCHAR(255) NOT NULL UNIQUE,
    description VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 为头像库添加一些初始数据，方便我们测试
INSERT INTO
    system_avatars (url, description)
VALUES (
        'https://zos.alipayobjects.com/rmsportal/jkjgkEfvpUPVyRjUImniV.png',
        '默认头像1'
    ),
    (
        'https://zos.alipayobjects.com/rmsportal/ODTLcjxAfvqbxHnVXCYX.png',
        '默认头像2'
    ),
    (
        'https://zos.alipayobjects.com/rmsportal/ARTUuuJTTMkEGtKuJINY.png',
        '默认头像3'
    );