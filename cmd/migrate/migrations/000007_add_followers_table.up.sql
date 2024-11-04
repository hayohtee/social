CREATE TABLE IF NOT EXISTS followers(
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    follower_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY(user_id, follower_id)
); 