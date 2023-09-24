CREATE TABLE IF NOT EXISTS bureaus (
    id SERIAL PRIMARY KEY,
    bureau VARCHAR(255) NOT NULL,
    color VARCHAR(255) NOT NULL DEFAULT 'fffafa',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_bureaus_timestamp') THEN
        CREATE TRIGGER update_bureaus_timestamp
        BEFORE UPDATE ON bureaus
        FOR EACH ROW
        EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;
