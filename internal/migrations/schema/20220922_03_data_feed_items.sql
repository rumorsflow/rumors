CREATE TABLE IF NOT EXISTS {{data.feed_items}} (
    [[id]] INTEGER PRIMARY KEY,
    [[feedId]] INTEGER NOT NULL,
    [[title]] TEXT NOT NULL,
    [[desc]] TEXT DEFAULT NULL,
    [[link]] TEXT NOT NULL,
    [[guid]] TEXT NOT NULL,
    [[pubDate]] TEXT NOT NULL,
    [[created]] TEXT NOT NULL,
    [[authors]] JSON DEFAULT "[]" NOT NULL,
    [[categories]] JSON DEFAULT "[]" NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS data.feed_items_unq_idx ON {{feed_items}} ([[link]], [[guid]], [[pubDate]]);
CREATE INDEX IF NOT EXISTS data.feed_items_feedId_idx ON {{feed_items}} ([[feedId]]);
CREATE INDEX IF NOT EXISTS data.feed_items_created_idx ON {{feed_items}} ([[created]]);
CREATE INDEX IF NOT EXISTS data.feed_items_authors_idx ON {{feed_items}} ([[authors]]);
CREATE INDEX IF NOT EXISTS data.feed_items_categories_idx ON {{feed_items}} ([[categories]]);
