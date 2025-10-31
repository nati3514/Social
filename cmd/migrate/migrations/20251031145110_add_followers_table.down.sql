-- Drop indexes first
DROP INDEX IF EXISTS idx_followers_follower_id;
DROP INDEX IF EXISTS idx_followers_user_id;

-- Then drop the table
DROP TABLE IF EXISTS followers;