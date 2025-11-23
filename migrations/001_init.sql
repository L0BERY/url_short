CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    click_count INTEGER DEFAULT 0
);

CREATE INDEX idx_short_code ON urls(short_code);
-- CREATE INDEX idx_created_at ON urls(created_at);