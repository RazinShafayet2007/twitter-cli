#!/bin/bash

DB_PATH="$HOME/.twitter-cli/data.db"

echo "Adding notifications table..."

sqlite3 "$DB_PATH" << 'EOF'
CREATE TABLE IF NOT EXISTS notifications (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    actor_id TEXT NOT NULL,
    type TEXT NOT NULL,
    target_id TEXT,
    created_at INTEGER NOT NULL,
    read INTEGER DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (actor_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_read ON notifications(user_id, read);

SELECT 'Notifications table created!';
EOF

echo "✓ Migration complete"