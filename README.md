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

## Installation

### Prerequisites

- Go 1.21 or higher
- SQLite3

### Build from source
```bash
git clone https://github.com/YOUR_USERNAME/twitter-cli.git
cd twitter-cli
go build -o twt
```

### Install globally (optional)
```bash
go install
```

Or copy the binary:
```bash
sudo cp twt /usr/local/bin/
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

## Architecture

### Data Model

- **Users**: Basic user accounts with unique usernames
- **Posts**: Text posts with timestamps, supports retweets
- **Follows**: Many-to-many relationship between users
- **Likes**: Many-to-many relationship between users and posts

### Technology Stack

- **Language**: Go
- **Database**: SQLite with foreign key constraints
- **CLI Framework**: Cobra
- **ID Generation**: ULID (sortable, unique identifiers)

### Project Structure
```
twitter-cli/
├── cmd/                    # CLI commands
│   ├── root.go            # Main CLI setup
│   ├── user.go            # User commands
│   ├── post.go            # Post commands
│   ├── feed.go            # Feed command
│   └── social.go          # Social commands
├── internal/
│   ├── db/                # Database setup
│   ├── models/            # Data structures
│   ├── store/             # Database operations
│   ├── config/            # Configuration management
│   ├── display/           # Output formatting
│   └── validation/        # Input validation
└── main.go                # Entry point
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
```

## Configuration

Database and config files are stored in `~/.twitter-cli/`:
- `data.db` - SQLite database
- `config.json` - Current user session

To use a different database location:
```bash
twt --db /path/to/database.db <command>
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
- ✅ **Many-to-many relationships** (follows, likes)
- ✅ **Self-referential relationships** (retweets)
- ✅ **Feed generation algorithms** (combining multiple data sources)
- ✅ **CLI application architecture** with Cobra
- ✅ **State management** (session persistence)
- ✅ **Input validation and sanitization**
- ✅ **Error handling patterns** in Go
- ✅ **Test-driven development** basics

## Limitations & Future Improvements

Current limitations:
- No comments/replies (threads)
- No direct messages
- No media uploads
- No hashtags or mentions
- No notifications
- Single-user local system (no server)

Potential enhancements:
- [ ] Search functionality
- [ ] Hashtag support
- [ ] User mentions (@username)
- [ ] Threads/replies
- [ ] Export data to JSON
- [ ] Import from real Twitter
- [ ] Web UI
- [ ] Multi-user server mode

## Contributing

This is a learning project, but suggestions are welcome! Open an issue or PR.

## License

MIT License - feel free to use for learning purposes.

## Acknowledgments

Built as a practical exercise in system design and backend development.

Inspired by Twitter's core functionality.