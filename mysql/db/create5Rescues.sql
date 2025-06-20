-- トラブルレスキューテーブル
CREATE TABLE IF NOT EXISTS trouble_rescues (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    task_id INTEGER NOT NULL,
    place VARCHAR(255),
    detail VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL DEFAULT 'todo',
    response VARCHAR(255),
    time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE,
    FOREIGN KEY (task_id) REFERENCES tasks (id) ON UPDATE CASCADE
);

-- 質問レスキューテーブル
CREATE TABLE IF NOT EXISTS question_rescues (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    question VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL DEFAULT 'todo',
    response VARCHAR(255),
    time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE
);

-- 人手不足レスキューテーブル
CREATE TABLE IF NOT EXISTS shorthanded_rescues (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    task_id INTEGER NOT NULL,
    missing_number INTEGER NOT NULL,
    place VARCHAR(255),
    status VARCHAR(255) NOT NULL DEFAULT 'todo',
    response VARCHAR(255),
    time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE,
    FOREIGN KEY (task_id) REFERENCES tasks (id) ON UPDATE CASCADE
);

-- トリガー設定
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_trouble_rescues_timestamp') THEN
        CREATE TRIGGER update_trouble_rescues_timestamp
        BEFORE UPDATE ON trouble_rescues
        FOR EACH ROW
        EXECUTE FUNCTION update_updated_at_column();
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_question_rescues_timestamp') THEN
        CREATE TRIGGER update_question_rescues_timestamp
        BEFORE UPDATE ON question_rescues
        FOR EACH ROW
        EXECUTE FUNCTION update_updated_at_column();
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_shorthanded_rescues_timestamp') THEN
        CREATE TRIGGER update_shorthanded_rescues_timestamp
        BEFORE UPDATE ON shorthanded_rescues
        FOR EACH ROW
        EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;
