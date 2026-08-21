ALTER TABLE users ADD COLUMN IF NOT EXISTS slack_user_id VARCHAR(255);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'users_slack_user_id_key'
    ) THEN
        ALTER TABLE users ADD CONSTRAINT users_slack_user_id_key UNIQUE (slack_user_id);
    END IF;
END $$;
