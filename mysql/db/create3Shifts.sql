CREATE TABLE IF NOT EXISTS shifts (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    year_id INTEGER NOT NULL,
    date_id INTEGER NOT NULL,
    time_id INTEGER NOT NULL,
    weather_id INTEGER NOT NULL,
    is_attendance BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (task_id) REFERENCES tasks (id) ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE,
    FOREIGN KEY (year_id) REFERENCES years (id) ON UPDATE CASCADE,
    FOREIGN KEY (date_id) REFERENCES dates (id) ON UPDATE CASCADE,
    FOREIGN KEY (time_id) REFERENCES times (id) ON UPDATE CASCADE,
    FOREIGN KEY (weather_id) REFERENCES weathers (id) ON UPDATE CASCADE
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_shifts_timestamp') THEN
        CREATE TRIGGER update_shifts_timestamp
        BEFORE UPDATE ON shifts
        FOR EACH ROW
        EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;

