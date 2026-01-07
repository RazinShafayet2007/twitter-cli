# Twitter CLI

A command-line Twitter clone built to learn backend system design, SQL, and Go.

## Features

- ✅ User management (create, login, logout)
- ✅ Post creation and deletion
- ✅ Social graph (follow/unfollow)
- ✅ Personalized feed
- ✅ Likes and retweets
- ✅ User profiles
- ✅ Engagement statistics
- ✅ Direct messaging (send, inbox, conversation, unread, delete, search)
- ✅ User blocking (block, unblock, list blocked)
- ✅ Notifications (list, read, clear unread count)
- ✅ Hashtags (search, trending)
- ✅ User Mentions (parsing, notifications, list mentions)
- ✅ Image Support (upload, view, open)
- ✅ Replies and threads (create replies, view threads)

## Installation

### Quick Install (Recommended)

**Linux & macOS:**

```bash
curl -fsSL https://raw.githubusercontent.com/RazinShafayet2007/twitter-cli/main/scripts/install.sh | bash
```

**What this does:**

- Downloads the correct binary for your system
- Installs to `~/.twitter-cli/bin/`
- Adds to your PATH automatically

### Manual Installation

**Download pre-built binaries:**

Visit the [releases page](https://github.com/RazinShafayet2007/twitter-cli/releases) and download for your platform:

- `twt-linux-amd64` - Linux (64-bit)
- `twt-linux-arm64` - Linux (ARM64)
- `twt-darwin-amd64` - macOS (Intel)
- `twt-darwin-arm64` - macOS (Apple Silicon)
- `twt-windows-amd64.exe` - Windows (64-bit)

```bash
# Example for Linux:
wget https://github.com/RazinShafayet2007/twitter-cli/releases/latest/download/twt-linux-amd64
chmod +x twt-linux-amd64
sudo mv twt-linux-amd64 /usr/local/bin/twt
```

### Build from Source

Requires Go 1.21+:

```bash
git clone https://github.com/RazinShafayet2007/twitter-cli.git
cd twitter-cli
go install
```

### Verify Installation

```bash
twt --version
twt --help
```

## Uninstallation

```bash
curl -fsSL https://raw.githubusercontent.com/RazinShafayet2007/twitter-cli/main/scripts/uninstall.sh | bash
```

Or manually:

```bash
rm -rf ~/.twitter-cli
# Remove the PATH line from your shell RC file
```

## Configuration

Twitter CLI stores data in `~/.twitter-cli/`:

```
~/.twitter-cli/
├── bin/          # Binary location
├── data.db       # SQLite database
└── config.json   # Current session
```

To use a different database:

```bash
twt --db /path/to/custom.db <command>
```

## Quick Start
```bash
# Create a user
twt user create alice

# Login
twt login alice

# Create a post
twt post "Hello, world!"

# View your profile
twt profile alice

# Follow someone
twt user create bob
twt follow bob

# View your feed
twt feed
```

## Usage

### User Management
```bash
# Create a new user
twt user create <username>

# Login as a user
twt login <username>

# Show current user
twt whoami

# Logout
twt logout
```

### Posting
```bash
# Create a post
twt post "Your message here"

# View a user's posts
twt profile <username>

# View a specific post with stats
twt show <post_id>

# Delete your own post
twt delete <post_id>
```

### Replies and Threads
```bash
# Reply to a post
twt reply <post_id> "Great point!"
# View entire conversation thread
twt thread <post_id>
```

### Images
```bash
# Post with images (max 4)
twt post "Check this out!" --image photo.jpg

# Post with multiple images
twt post "My vacation" --image beach.png --image sunset.jpg

# Download images from a post
twt image download <post_id>

# Open image in default viewer
twt image open <post_id> <image_index>
```

### Social
```bash
# Follow a user
twt follow <username>

# Unfollow a user
twt unfollow <username>

# List who you're following
twt following

# List who follows you (or another user)
twt followers [username]

# View user statistics
twt stats [username]
```

### Direct Messaging
```bash
# Send a direct message
twt message send <username> "Your message here"

# View your inbox
twt message inbox

# View a conversation with a specific user
twt message conversation <username>

# Check your unread message count
twt message unread

# Delete a message (by ID)
twt message delete <message_id>

# Search messages
twt message search <query>
```

### User Blocking
```bash
# Block a user
twt block <username>

# Unblock a user
twt unblock <username>

# List users you have blocked
twt blocked
```

### Notifications
```bash
# View your notifications
twt notifications

# Mark all notifications as read
twt notifications read

# Clear all read notifications
twt notifications clear

# View only unread notifications
twt notifications --unread
```

### Feed
```bash
# View your personalized feed
twt feed

# Limit number of posts
twt feed --limit 10

# Pagination
twt feed --limit 10 --offset 20
```

### Engagement
```bash
# Like a post
twt like <post_id>

# Unlike a post
twt unlike <post_id>

# See who liked a post
twt likes <post_id>

# Retweet a post
twt retweet <post_id>
```

### Hashtags & Mentions
```bash
# Posts can include hashtags and mentions
twt post "Hello @alice check out #golang"

# View posts with a specific hashtag
twt hashtag golang

# View trending hashtags
twt trending

# View posts that mention you
twt mentions
```

## Architecture

### Data Model

- **Users**: Basic user accounts with unique usernames
- **Posts**: Text posts with timestamps, supports retweets
- **Follows**: Many-to-many relationship between users
- **Likes**: Many-to-many relationship between users and posts
- **Messages**: Direct messages between users
- **Blocks**: Records of one user blocking another
- **Notifications**: System notifications for user interactions

### Technology Stack

- **Language**: Go
- **Database**: SQLite with foreign key constraints
- **CLI Framework**: Cobra
- **ID Generation**: ULID (sortable, unique identifiers)

### Project Structure
```
├── CHANGELOG.md
├── CONTRIBUTING.md
├── README.md
├── cmd
│   ├── block.go
│   ├── feed.go
│   ├── hashtag.go
│   ├── image.go
│   ├── mentions.go
│   ├── message.go
│   ├── notifications.go
│   ├── post.go
│   ├── root.go
│   ├── social.go
│   └── user.go
├── go.mod
├── go.sum
├── internal
│   ├── config
│   │   └── config.go
│   ├── db
│   │   ├── db.go
│   │   └── schema.sql
│   ├── display
│   │   └── format.go
│   ├── errors
│   │   └── errors.go
│   ├── media
│   │   └── media.go
│   ├── models
│   │   ├── media.go
│   │   ├── message.go
│   │   ├── notification.go
│   │   ├── post.go
│   │   ├── social.go
│   │   └── user.go
│   ├── parser
│   │   └── parser.go
│   ├── store
│   │   ├── hashtag_store.go
│   │   ├── media_store.go
│   │   ├── mention_store.go
│   │   ├── message_store.go
│   │   ├── notification_store.go
│   │   ├── post_store.go
│   │   ├── social_store.go
│   │   ├── user_store.go
│   │   └── user_store_test.go
│   └── validation
│       └── validation.go
├── main.go                        # Entry point
├── package.json
├── scripts
│   ├── build-release.sh 
│   ├── install.sh
│   ├── migrate-blocks.sh
│   ├── migrate-hashtags-mentions.sh
│   ├── migrate-media.sh
│   ├── migrate-messages.sh
│   ├── migrate-notifications.sh
│   └── uninstall.sh
├── test_scenario.sh
├── version.txt
├── .changeset
│   ├── config.json
│   └── wet-mammals-win.md
├── .github
│   └── workflows
│       ├── changesets.yml
│       ├── release.yml
│       ├── require-changeset.yml
│       └── tag-release.yml
└── .gitignore                
```

## Database Schema

```sql
-- Users
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    created_at INTEGER NOT NULL
);

-- Posts
CREATE TABLE posts (
    id TEXT PRIMARY KEY,
    author_id TEXT NOT NULL,
    text TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    is_retweet INTEGER DEFAULT 0,
    original_post_id TEXT,
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Follows
CREATE TABLE follows (
    follower_id TEXT NOT NULL,
    followee_id TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    PRIMARY KEY (follower_id, followee_id)
);

-- Likes
CREATE TABLE likes (
    user_id TEXT NOT NULL,
    post_id TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    PRIMARY KEY (user_id, post_id)
);

-- Messages
CREATE TABLE messages (
    id TEXT PRIMARY KEY,
    sender_id TEXT NOT NULL,
    receiver_id TEXT NOT NULL,
    text TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    read INTEGER DEFAULT 0,
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Blocks
CREATE TABLE blocks (
    blocker_id TEXT NOT NULL,
    blocked_id TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    PRIMARY KEY (blocker_id, blocked_id),
    FOREIGN KEY (blocker_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (blocked_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Notifications
CREATE TABLE notifications (
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

-- Hashtags
CREATE TABLE hashtags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tag TEXT UNIQUE NOT NULL,
    created_at INTEGER NOT NULL
);

-- Post <-> Hashtags join table
CREATE TABLE post_hashtags (
    post_id TEXT NOT NULL,
    hashtag_id INTEGER NOT NULL,
    PRIMARY KEY (post_id, hashtag_id),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (hashtag_id) REFERENCES hashtags(id) ON DELETE CASCADE
);

-- Mentions
CREATE TABLE mentions (
    post_id TEXT NOT NULL,
    mentioned_user_id TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    PRIMARY KEY (post_id, mentioned_user_id),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (mentioned_user_id) REFERENCES users(id) ON DELETE CASCADE
);
```
## Development

### Running tests
```bash
go test ./...
```

### Building
```bash
go build -o twt
```

### Linting
```bash
go fmt ./...
go vet ./...
```

## Learning Outcomes

This project demonstrates:

- ✅ **Relational database design** with proper foreign keys and constraints
- ✅ **Complex SQL queries** with JOINs, subqueries, and aggregations
- ✅ **Many-to-many relationships** (follows, likes, blocks)
- ✅ **Self-referential relationships** (retweets, direct messages, notifications)
- ✅ **Feed generation algorithms** (combining multiple data sources)
- ✅ **CLI application architecture** with Cobra
- ✅ **State management** (session persistence)
- ✅ **Input validation and sanitization**
- ✅ **Error handling patterns** in Go
- ✅ **Test-driven development** basics
- ✅ **Direct messaging implementation**
- ✅ **User blocking functionality**
- ✅ **Message search capabilities**
- ✅ **Notification system** (real-time user feedback)
- ✅ **Hashtag support**
- ✅ **User mentions**
- ✅ **Image handling** (storage, metadata, CLI viewing)
- ✅ **Threads/replies implementation**
- ✅ **Automated Release Workflow** (Changesets, GitHub Actions)

## Limitations & Future Improvements

Current limitations:

- No comments/replies (threads)
- Single-user local system (no server)

Potential enhancements:

- [ ] Export data to JSON
- [ ] Import from real Twitter
- [ ] Web UI
- [ ] Multi-user server mode

## Contributing

This is a learning project, but suggestions are welcome! Open an issue or PR.

### Release Process

This project uses [Changesets](https://github.com/changesets/changesets) for automated versioning and releases.

1.  **Create a Branch**: `git checkout -b feat/my-feature`
2.  **Make Changes**: Write your code.
3.  **Add a Changeset**: Run `npx changeset` to create a changelog entry.
4.  **Push & PR**: Open a Pull Request.
5.  **Merge**: When merged to `main`, a "Version Packages" PR is created automatically.
6.  **Release**: Merging the "Version Packages" PR tags the release on GitHub.

## License

MIT License - feel free to use for learning purposes.

## Acknowledgments

Built as a practical exercise in system design and backend development.

Inspired by Twitter's core functionality.

