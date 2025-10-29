CREATE TABLE IF NOT EXISTS comments (
	id BIGSERIAL PRIMARY KEY,
	content TEXT NOT NULL,
	user_id BIGINT NOT NULL,
	post_id BIGINT NOT NULL,
	created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);