# Gator - RSS Feed Aggregator

A command-line RSS feed aggregator that allows you to follow multiple RSS feeds, aggregate posts, and browse them in your terminal.

## Prerequisites

Before installing Gator, ensure you have the following installed:

- **Go 1.25+** - [Installation Guide](https://go.dev/doc/install)
- **PostgreSQL 15+** - [Installation Guide](https://www.postgresql.org/download/)

## Installation

### 1. Install Gator

```bash
go install github.com/Professor-Goo/gator@latest
```

This compiles the binary and places it in your `$GOPATH/bin` directory (typically `~/go/bin`).

### 2. Database Setup

Start PostgreSQL and create the `gator` database:

```bash
# Start PostgreSQL (varies by system)
# macOS: brew services start postgresql@15
# Linux: sudo service postgresql start

# Create database
psql postgres
CREATE DATABASE gator;
\q
```

Set your PostgreSQL user password if needed:

```bash
# Linux
sudo -u postgres psql
ALTER USER postgres PASSWORD 'your_password';
\q
```

### 3. Run Migrations

Clone the repository to run migrations:

```bash
git clone https://github.com/Professor-Goo/gator.git
cd gator

# Install Goose migration tool
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations
cd sql/schema
goose postgres "postgres://postgres:your_password@localhost:5432/gator?sslmode=disable" up
```

### 4. Configuration

Create a configuration file at `~/.gatorconfig.json`:

```json
{
  "db_url": "postgres://postgres:your_password@localhost:5432/gator?sslmode=disable"
}
```

Replace `your_password` with your PostgreSQL password.

## Usage

### User Management

**Register a new user:**
```bash
gator register <username>
```

**Login as existing user:**
```bash
gator login <username>
```

**List all users:**
```bash
gator users
```

### Feed Management

**Add a feed:**
```bash
gator addfeed "Feed Name" "https://example.com/feed.xml"
```

**List all feeds:**
```bash
gator feeds
```

**Follow a feed:**
```bash
gator follow "https://example.com/feed.xml"
```

**List feeds you're following:**
```bash
gator following
```

**Unfollow a feed:**
```bash
gator unfollow "https://example.com/feed.xml"
```

### Post Aggregation

**Start the aggregator (runs continuously):**
```bash
gator agg <duration>
```

Examples:
- `gator agg 1m` - Fetch feeds every minute
- `gator agg 30s` - Fetch feeds every 30 seconds
- `gator agg 5m` - Fetch feeds every 5 minutes

The aggregator runs continuously and fetches new posts from all feeds in the database. Press `Ctrl+C` to stop.

**Browse posts:**
```bash
# Show 2 most recent posts (default)
gator browse

# Show specific number of posts
gator browse 10
```

### Utility Commands

**Reset database (delete all users and data):**
```bash
gator reset
```

## Recommended RSS Feeds

- **TechCrunch**: `https://techcrunch.com/feed/`
- **Hacker News**: `https://news.ycombinator.com/rss`
- **Boot.dev Blog**: `https://blog.boot.dev/index.xml`
- **The Changelog**: `https://changelog.com/feed`

## Typical Workflow

1. **Setup:**
   ```bash
   gator register myusername
   ```

2. **Add feeds:**
   ```bash
   gator addfeed "TechCrunch" "https://techcrunch.com/feed/"
   gator addfeed "Hacker News" "https://news.ycombinator.com/rss"
   ```

3. **Start aggregator (Terminal 1):**
   ```bash
   gator agg 1m
   ```

4. **Browse posts (Terminal 2):**
   ```bash
   gator browse 5
   ```

## Architecture

- **Language**: Go 1.25
- **Database**: PostgreSQL 16
- **Migrations**: Goose
- **Query Generation**: SQLC (type-safe SQL)
- **RSS Parsing**: encoding/xml

## Project Structure

```
gator/
├── main.go              # CLI application entry point
├── internal/
│   ├── config/          # Configuration management
│   └── database/        # Generated SQLC code
├── sql/
│   ├── schema/          # Goose migrations
│   └── queries/         # SQLC SQL queries
└── sqlc.yaml            # SQLC configuration
```

## Development

**Run from source:**
```bash
git clone https://github.com/Professor-Goo/gator.git
cd gator
go run . <command>
```

**Build binary:**
```bash
go build -o gator
./gator <command>
```

**Run tests:**
```bash
go test ./...
```

## Database Schema

- **users** - User accounts
- **feeds** - RSS feed sources
- **feed_follows** - User-to-feed relationships (many-to-many)
- **posts** - Aggregated posts from feeds

## Troubleshooting

**"unknown command" error:**
- Ensure you've run `go install` and `~/go/bin` is in your PATH

**Database connection errors:**
- Verify PostgreSQL is running
- Check your `~/.gatorconfig.json` has the correct connection string
- Ensure migrations have been run

**No posts showing:**
- Run the `agg` command first to fetch posts
- Ensure you're following feeds with `gator following`
- Wait a few minutes for aggregation to complete

## License

This project is part of the Boot.dev backend development course.

## Contributing

This is a learning project. Feel free to fork and experiment!

## Author

Built by [Professor-Goo](https://github.com/Professor-Goo) as part of the Boot.dev curriculum.