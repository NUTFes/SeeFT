CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    mail VARCHAR(255) NOT NULL UNIQUE,
    grade_id INTEGER NOT NULL,
    department_id INTEGER NOT NULL,
    bureau_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    student_number INTEGER NOT NULL,
    tel VARCHAR(15) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (grade_id) REFERENCES grades (id) ON UPDATE CASCADE,
    FOREIGN KEY (department_id) REFERENCES departments (id) ON UPDATE CASCADE,
    FOREIGN KEY (bureau_id) REFERENCES bureaus (id) ON UPDATE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles (id) ON UPDATE CASCADE
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_dates_timestamp') THEN
        CREATE TRIGGER update_dates_timestamp
        BEFORE UPDATE ON users
        FOR EACH ROW
        EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;

