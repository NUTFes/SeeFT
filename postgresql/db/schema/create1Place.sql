CREATE TABLE IF NOT EXISTS places (
    id SERIAL PRIMARY KEY,
    place VARCHAR(255) NOT NULL,
    remark VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

DO $$
BEGIN
    -- 旧手順で構築したDBには誤った名前(update_grades_timestamp)で
    -- placesのトリガーが作られているため、存在する場合は正しい名前に揃える
    IF EXISTS (
        SELECT 1 FROM pg_trigger t
        JOIN pg_class c ON c.oid = t.tgrelid
        WHERE c.relname = 'places' AND t.tgname = 'update_grades_timestamp'
    ) THEN
        ALTER TRIGGER update_grades_timestamp ON places RENAME TO update_places_timestamp;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger t
        JOIN pg_class c ON c.oid = t.tgrelid
        WHERE c.relname = 'places' AND t.tgname = 'update_places_timestamp'
    ) THEN
        CREATE TRIGGER update_places_timestamp
        BEFORE UPDATE ON places
        FOR EACH ROW
        EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;