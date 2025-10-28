-- 用户表 (users)
-- 基于PostgreSQL的认证服务用户表结构

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,                                            -- 用户ID，自增主键
    username VARCHAR(50) UNIQUE,                                         -- 用户名（可选，唯一）
    nick_name VARCHAR(100),                                              -- 昵称
    password_hash VARCHAR(255) NOT NULL,                                 -- 密码哈希（bcrypt加密）
    phone VARCHAR(20) UNIQUE,                                            -- 手机号（可选，唯一）
    email VARCHAR(100) UNIQUE,                                           -- 邮箱地址（可选，唯一）
    status SMALLINT DEFAULT 0,                                           -- 用户状态：0-未激活，1-正常，2-禁用，3-注销，9-已删除
    last_login_at BIGINT,                                                -- 最后登录时间（Unix时间戳）
    creator VARCHAR(100) DEFAULT 'system',                               -- 创建者
    updater VARCHAR(100) DEFAULT 'system',                               -- 更新者
    create_time BIGINT DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)::BIGINT, -- 创建时间（Unix时间戳）
    update_time BIGINT DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)::BIGINT, -- 更新时间（Unix时间戳）
    delete_time BIGINT DEFAULT 0                                         -- 删除时间（Unix时间戳，用于软删除）
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
CREATE INDEX IF NOT EXISTS idx_users_create_time ON users(create_time);
CREATE INDEX IF NOT EXISTS idx_users_last_login_at ON users(last_login_at);
CREATE INDEX IF NOT EXISTS idx_users_delete_time ON users(delete_time);

-- 添加约束
ALTER TABLE users ADD CONSTRAINT chk_users_status
    CHECK (status IN (0, 1, 2, 3, 9));  -- 限制状态值范围

-- 添加注释
COMMENT ON TABLE users IS '用户表 - 存储用户基本信息';
COMMENT ON COLUMN users.id IS '用户ID，自增主键';
COMMENT ON COLUMN users.username IS '用户名，可选，唯一索引';
COMMENT ON COLUMN users.nick_name IS '昵称，显示用';
COMMENT ON COLUMN users.password_hash IS '密码哈希，使用bcrypt加密';
COMMENT ON COLUMN users.phone IS '手机号，可选，唯一索引';
COMMENT ON COLUMN users.email IS '邮箱地址，可选，唯一索引';
COMMENT ON COLUMN users.status IS '用户状态：0-未激活，1-正常，2-禁用，3-注销，9-已删除';
COMMENT ON COLUMN users.last_login_at IS '最后登录时间，Unix时间戳';
COMMENT ON COLUMN users.creator IS '创建者';
COMMENT ON COLUMN users.updater IS '更新者';
COMMENT ON COLUMN users.create_time IS '创建时间，Unix时间戳';
COMMENT ON COLUMN users.update_time IS '更新时间，Unix时间戳';
COMMENT ON COLUMN users.delete_time IS '删除时间，Unix时间戳，用于软删除';

-- 创建更新时间触发器函数
CREATE OR REPLACE FUNCTION update_updated_time_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.update_time = EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)::BIGINT;
    NEW.updater = COALESCE(current_setting('app.current_user', true), 'system');
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 创建用户表的更新时间触发器
DROP TRIGGER IF EXISTS update_users_update_time ON users;
CREATE TRIGGER update_users_update_time
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_time_column();

-- 插入默认系统用户（可选）
INSERT INTO users (username, nick_name, password_hash, status, creator, updater)
VALUES ('admin', '系统管理员', '$2a$12$placeholder_hash_bcrypt', 1, 'system', 'system')
ON CONFLICT (username) DO NOTHING;