#!/bin/bash

DB_PATH="$HOME/.twitter-cli/data.db"

echo "Adding media tables..."

sqlite3 "$DB_PATH" << 'EOF'
-- Media table
CREATE TABLE IF NOT EXISTS media (
    id TEXT PRIMARY KEY,
    post_id TEXT NOT NULL,
    file_path TEXT NOT NULL,
    file_name TEXT NOT NULL,
    file_type TEXT NOT NULL,
    file_size INTEGER NOT NULL,
    width INTEGER,
    height INTEGER,
    position INTEGER DEFAULT 0,
    created_at INTEGER NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_media_post ON media(post_id, position);

SELECT 'Media table created!';
EOF

echo "✓ Migration complete"