CREATE TABLE IF NOT EXISTS dates (
    id SERIAL PRIMARY KEY,
    year_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    date VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (year_id) REFERENCES years (id) ON UPDATE CASCADE
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_dates_timestamp') THEN
        CREATE TRIGGER update_dates_timestamp
        BEFORE UPDATE ON dates
        FOR EACH ROW
        EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;

