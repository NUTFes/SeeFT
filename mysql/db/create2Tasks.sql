CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    task VARCHAR(255) ,
    place VARCHAR(255) ,
    url VARCHAR(255) ,
    superviser VARCHAR(255) ,
    color VARCHAR(255) DEFAULT 'fffafa' ,
    notes VARCHAR(255),
    year_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
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

