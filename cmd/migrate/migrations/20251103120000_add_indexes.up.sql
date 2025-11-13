-- Create the extension for trigram search with superuser privileges
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'pg_trgm') THEN
        CREATE EXTENSION IF NOT EXISTS pg_trgm;
    END IF;
END $$;

-- Create GIN indexes for text search
CREATE INDEX IF NOT EXISTS idx_comments_content ON comments USING gin (content gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_posts_title ON posts USING gin (title gin_trgm_ops);
-- For array types, we need to use the array_ops operator class
CREATE INDEX IF NOT EXISTS idx_posts_tags ON posts USING gin (tags array_ops);
CREATE INDEX IF NOT EXISTS idx_users_username ON users USING gin (username gin_trgm_ops);
-- Regular B-tree indexes
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts (user_id);
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments (post_id);