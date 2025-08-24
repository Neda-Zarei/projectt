-- Migration: Update to new plan system structure
-- Date: 2025-01-XX

-- Drop existing tables if they exist to start fresh
DROP TABLE IF EXISTS plan_histories CASCADE;
DROP VIEW IF EXISTS expired_user_plans;
DROP VIEW IF EXISTS expiring_user_plans;
DROP FUNCTION IF EXISTS expire_user_plans();

-- Create limitations table
CREATE TABLE IF NOT EXISTS limitations (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL UNIQUE
);

-- Create plans table with new structure
DROP TABLE IF EXISTS plans CASCADE;
CREATE TABLE plans (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    custom BOOLEAN NOT NULL DEFAULT true,
    payg BOOLEAN NOT NULL DEFAULT false
);

-- Create prices table
CREATE TABLE IF NOT EXISTS prices (
    plan_id INTEGER NOT NULL,
    month INTEGER NOT NULL,
    price INTEGER NOT NULL DEFAULT 1000000,
    PRIMARY KEY (plan_id, month),
    FOREIGN KEY (plan_id) REFERENCES plans(id) ON DELETE CASCADE
);

-- Create plan_limitations junction table
CREATE TABLE IF NOT EXISTS plan_limitations (
    plan_id INTEGER NOT NULL,
    limitation_id INTEGER NOT NULL,
    value INTEGER DEFAULT 1,
    PRIMARY KEY (plan_id, limitation_id),
    FOREIGN KEY (plan_id) REFERENCES plans(id) ON DELETE CASCADE,
    FOREIGN KEY (limitation_id) REFERENCES limitations(id) ON DELETE CASCADE
);

-- Update users table structure
ALTER TABLE users DROP COLUMN IF EXISTS external_id;
ALTER TABLE users DROP COLUMN IF EXISTS first_name;
ALTER TABLE users DROP COLUMN IF EXISTS last_name;
ALTER TABLE users DROP COLUMN IF EXISTS date_of_birth;
ALTER TABLE users DROP COLUMN IF EXISTS is_active;

-- Add new columns to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS oauth_id VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS name VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS password VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS company_name VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS job_title VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS active BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE users ADD COLUMN IF NOT EXISTS subscribe_news BOOLEAN NOT NULL DEFAULT true;
ALTER TABLE users ADD COLUMN IF NOT EXISTS subscribe_notifications BOOLEAN NOT NULL DEFAULT true;

-- Update email column to be NOT NULL
ALTER TABLE users ALTER COLUMN email SET NOT NULL;

-- Update user_plans table structure
DROP TABLE IF EXISTS user_plans CASCADE;
CREATE TABLE user_plans (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    plan_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    ex_time TIMESTAMP WITH TIME ZONE NOT NULL,
    FOREIGN KEY (plan_id) REFERENCES plans(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_user_plans_user_id ON user_plans(user_id);
CREATE INDEX IF NOT EXISTS idx_user_plans_plan_id ON user_plans(plan_id);
CREATE INDEX IF NOT EXISTS idx_user_plans_ex_time ON user_plans(ex_time);
CREATE INDEX IF NOT EXISTS idx_user_plans_deleted_at ON user_plans(deleted_at);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_oauth_id ON users(oauth_id);

-- Create updated function for plan expiration
CREATE OR REPLACE FUNCTION expire_user_plans()
RETURNS void AS $$
BEGIN
    -- Soft delete expired plans
    UPDATE user_plans 
    SET deleted_at = NOW()
    WHERE deleted_at IS NULL 
    AND ex_time <= NOW();
END;
$$ LANGUAGE plpgsql;

-- Create updated views for expired and expiring plans
CREATE OR REPLACE VIEW expired_user_plans AS
SELECT up.*, p.title as plan_title, u.name as user_name, u.email as user_email
FROM user_plans up
JOIN plans p ON up.plan_id = p.id
JOIN users u ON up.user_id = u.id
WHERE up.deleted_at IS NOT NULL;

CREATE OR REPLACE VIEW expiring_user_plans AS
SELECT up.*, p.title as plan_title, u.name as user_name, u.email as user_email
FROM user_plans up
JOIN plans p ON up.plan_id = p.id
JOIN users u ON up.user_id = u.id
WHERE up.deleted_at IS NULL
AND up.ex_time <= NOW() + INTERVAL '7 days'
AND up.ex_time > NOW();

-- Insert some sample data for testing
INSERT INTO limitations (title) VALUES 
    ('API Calls per Month'),
    ('Storage GB'),
    ('Concurrent Users'),
    ('Premium Support')
ON CONFLICT (title) DO NOTHING;

INSERT INTO plans (title, custom, payg) VALUES 
    ('Basic Plan', true, false),
    ('Pro Plan', true, false),
    ('Enterprise Plan', true, false),
    ('Pay as You Go', false, true)
ON CONFLICT DO NOTHING;