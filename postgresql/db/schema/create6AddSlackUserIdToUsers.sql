ALTER TABLE users ADD COLUMN slack_user_id VARCHAR(255);
ALTER TABLE users ADD CONSTRAINT users_slack_user_id_key UNIQUE (slack_user_id);