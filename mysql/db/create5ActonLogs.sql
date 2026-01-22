CREATE TABLE action_logs (
    id SERIAL PRIMARY KEY,
    shift_id INTEGER,
    user_id INTEGER,
    date_id INTEGER,
    action_type VARCHAR(50) NOT NULL, -- "CREATE", "UPDATE", "DELETE"
    diff_payload JSONB,
    is_sent BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (shift_id) REFERENCES shifts(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (date_id) REFERENCES dates(id)
);