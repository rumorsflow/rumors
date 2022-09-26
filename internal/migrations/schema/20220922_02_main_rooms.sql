CREATE TABLE IF NOT EXISTS {{main.rooms}} (
    [[id]] BIGINT PRIMARY KEY,
    [[type]] TEXT NOT NULL,
    [[title]] TEXT NOT NULL,
    [[broadcast]] BOOLEAN DEFAULT FALSE NOT NULL,
    [[deleted]] BOOLEAN DEFAULT FALSE NOT NULL,
    [[created]] TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS main.rooms_title_idx ON {{rooms}} ([[title]]);
CREATE INDEX IF NOT EXISTS main.rooms_broadcast_idx ON {{rooms}} ([[broadcast]]);
CREATE INDEX IF NOT EXISTS main.rooms_deleted_idx ON {{rooms}} ([[deleted]]);
