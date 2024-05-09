CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    task VARCHAR(255) ,
    place_id INTEGER NOT NULL,
    url VARCHAR(255) ,
    superviser_id INTEGER,
    color VARCHAR(255) DEFAULT 'fffafa' ,
    remark VARCHAR(255),
    year_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (place_id) REFERENCES places (id) ON UPDATE CASCADE,
    FOREIGN KEY (superviser_id) REFERENCES users (id) ON UPDATE CASCADE,
    FOREIGN KEY (year_id) REFERENCES years (id) ON UPDATE CASCADE
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_tasks_timestamp') THEN
        CREATE TRIGGER update_tasks_timestamp
        BEFORE UPDATE ON tasks
        FOR EACH ROW
        EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;

