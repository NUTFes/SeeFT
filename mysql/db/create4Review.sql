-- レビューのテーブル
CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    task_id INTEGER NOT NULL,
    staffing_rating INTEGER NOT NULL, -- 人数は適切でしたか
    manual_rating INTEGER NOT NULL,   -- マニュアルわかりやすかったですか
    comment VARCHAR(255) NOT NULL DEFAULT '', -- 他にあれば教えてください
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE,
    FOREIGN KEY (task_id) REFERENCES tasks (id) ON UPDATE CASCADE
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_reviews_timestamp') THEN
        CREATE TRIGGER update_reviews_timestamp
        BEFORE UPDATE ON reviews
        FOR EACH ROW
        EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;
