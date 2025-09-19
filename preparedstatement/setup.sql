-- 既存テーブル削除
DROP TABLE IF EXISTS users;

-- ユーザーテーブル作成
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name   TEXT NOT NULL,
    email  TEXT NOT NULL UNIQUE,
    status TEXT NOT NULL
);

-- インデックス作成
-- 1. email による検索用
CREATE INDEX idx_users_email ON users(email);

-- 2. status による検索用
CREATE INDEX idx_users_status ON users(status);

-- 3. 複合インデックス（email + status）
CREATE INDEX idx_users_email_status ON users(email, status);

-- 4. 部分インデックス（active ユーザー限定）
CREATE INDEX idx_users_email_active ON users(email) WHERE status = 'active';
