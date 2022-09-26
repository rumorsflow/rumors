CREATE TABLE IF NOT EXISTS {{main.feeds}} (
    [[id]] INTEGER PRIMARY KEY,
    [[by]] INTEGER NOT NULL,
    [[host]] TEXT NOT NULL,
    [[title]] TEXT DEFAULT NULL,
    [[link]] TEXT UNIQUE NOT NULL,
    [[enabled]] BOOLEAN DEFAULT FALSE NOT NULL,
    [[created]] TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS main.feeds_by_idx ON {{feeds}} ([[by]]);
CREATE INDEX IF NOT EXISTS main.feeds_host_idx ON {{feeds}} ([[host]]);
CREATE INDEX IF NOT EXISTS main.feeds_enabled_idx ON {{feeds}} ([[enabled]]);
CREATE INDEX IF NOT EXISTS main.feeds_created_idx ON {{feeds}} ([[created]]);
