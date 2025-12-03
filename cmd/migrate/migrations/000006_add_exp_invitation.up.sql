-- +goose Up
-- +begin
ALTER TABLE user_invitations 
    ADD COLUMN IF NOT EXISTS expires_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() + INTERVAL '7 days');
-- +end

-- +begin
COMMENT ON COLUMN user_invitations.expires_at IS 'When the invitation expires and becomes invalid';
-- +end
