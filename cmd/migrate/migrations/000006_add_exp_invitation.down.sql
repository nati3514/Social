-- +goose Down
-- +begin
ALTER TABLE user_invitations 
    DROP COLUMN IF EXISTS expires_at;
-- +end
