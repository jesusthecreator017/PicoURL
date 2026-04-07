CREATE TABLE IF NOT EXISTS urls (
    short_url TEXT PRIMARY KEY,
    original_url TEXT NOT NULL,
    click_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);