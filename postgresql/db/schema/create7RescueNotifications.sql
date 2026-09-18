-- レスキューの対応状況・返答が変わったことを送信者にSlack DMで知らせるための未送信キュー。
-- レスキューは種類ごとに別テーブル(question/shorthanded/trouble_rescues)でIDも別採番のため、
-- rescue_type と rescue_id の組で対象を指す(外部キーは張れない)
CREATE TABLE IF NOT EXISTS rescue_notifications (
    id SERIAL PRIMARY KEY,
    rescue_type VARCHAR(50) NOT NULL, -- "trouble", "question", "shorthanded"
    rescue_id INTEGER NOT NULL,
    user_id INTEGER,
    kind VARCHAR(50) NOT NULL, -- "inProgress", "done", "response"
    status VARCHAR(255) NOT NULL,
    response VARCHAR(255),
    is_sent BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS rescue_notifications_unsent_idx
    ON rescue_notifications (created_at)
    WHERE is_sent = false;
