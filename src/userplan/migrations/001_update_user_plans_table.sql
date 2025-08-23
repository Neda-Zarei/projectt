-- Migration: Update user_plans table for expiration logic
-- Date: 2024-01-XX

-- Add new expires_at column
ALTER TABLE user_plans ADD COLUMN IF NOT EXISTS expires_at TIMESTAMP WITH TIME ZONE;

-- Update existing records to set expires_at based on end_date
UPDATE user_plans 
SET expires_at = DATE_TRUNC('day', end_date + INTERVAL '1 day')
WHERE expires_at IS NULL AND end_date IS NOT NULL;

-- Set expires_at to current time for records without end_date
UPDATE user_plans 
SET expires_at = DATE_TRUNC('day', NOW())
WHERE expires_at IS NULL;

-- Make expires_at NOT NULL
ALTER TABLE user_plans ALTER COLUMN expires_at SET NOT NULL;

-- Add index on expires_at for efficient expiration queries
CREATE INDEX IF NOT EXISTS idx_user_plans_expires_at ON user_plans(expires_at);

-- Drop old columns that are no longer needed
ALTER TABLE user_plans DROP COLUMN IF EXISTS end_date;
ALTER TABLE user_plans DROP COLUMN IF EXISTS is_auto_renew;
ALTER TABLE user_plans DROP COLUMN IF EXISTS last_renewal_at;
ALTER TABLE user_plans DROP COLUMN IF EXISTS next_renewal_at;
ALTER TABLE user_plans DROP COLUMN IF EXISTS payment_gateway;
ALTER TABLE user_plans DROP COLUMN IF EXISTS transaction_id;

-- Create a function to handle plan expiration
CREATE OR REPLACE FUNCTION expire_user_plans()
RETURNS void AS $$
BEGIN
    -- Update expired plans
    UPDATE user_plans 
    SET status = 'expired'
    WHERE status = 'active' 
    AND expires_at <= NOW();
END;
$$ LANGUAGE plpgsql;

-- Create a view for expired plans
CREATE OR REPLACE VIEW expired_user_plans AS
SELECT * FROM user_plans 
WHERE status = 'active' 
AND expires_at <= NOW();

-- Create a view for expiring plans (for notifications)
CREATE OR REPLACE VIEW expiring_user_plans AS
SELECT * FROM user_plans 
WHERE status = 'active' 
AND expires_at <= NOW() + INTERVAL '7 days'
AND expires_at > NOW();
