CREATE TABLE IF NOT EXISTS grades (
    id SERIAL PRIMARY KEY,
    grade VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_grades_timestamp') THEN
        CREATE TRIGGER update_grades_timestamp
        BEFORE UPDATE ON grades
        FOR EACH ROW
        EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;
