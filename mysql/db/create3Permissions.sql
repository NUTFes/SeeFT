CREATE TABLE IF NOT EXISTS permissions (
    user_id INTEGER PRIMARY KEY,
    allow_shift BOOLEAN NOT NULL DEFAULT FALSE,
    allow_task BOOLEAN NOT NULL DEFAULT FALSE,
    allow_user BOOLEAN NOT NULL DEFAULT FALSE,
    allow_property BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_permissions_timestamp') THEN
        CREATE TRIGGER update_permissions_timestamp
        BEFORE UPDATE ON permissions
        FOR EACH ROW
        EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;

