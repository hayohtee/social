ALTER TABLE posts
ADD COLUMN IF NOT EXISTS version INT NOT NULL DEFAULT 1;