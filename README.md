# gator

`gator` is a command-line RSS feed aggregator. It lets you register users, add
RSS feeds, follow feeds added by others, continuously collect posts from those
feeds, and browse the latest posts from the feeds you follow.

## Prerequisites

To run `gator` you'll need the following installed on your machine:

- **[Go](https://go.dev/doc/install)** (1.26+) — used to build and install the CLI.
- **[PostgreSQL](https://www.postgresql.org/download/)** (15+) — the database where users, feeds, and posts are stored.

Make sure Postgres is running and that you've created a database for `gator` to
use, for example:

```bash
createdb gator
```

## Installation

Install the `gator` CLI with `go install`:

```bash
go install github.com/shadyendless/gator@latest
```

This builds the binary and places it in your `$GOBIN` (usually `~/go/bin`). Make
sure that directory is on your `PATH` so you can run `gator` from anywhere.

## Configuration

`gator` reads its configuration from a JSON file named `.gatorconfig.json` in
your home directory (`~/.gatorconfig.json`). Create it with your Postgres
connection string:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable"
}
```

- `db_url` — the connection string for your Postgres database.
- `current_user_name` — managed automatically by `gator`; it's set when you
  `register` or `login`, so you don't need to add it yourself.

## Usage

Run any command with:

```bash
gator <command> [args...]
```

### A few commands to get you started

| Command | Description |
| --- | --- |
| `gator register <name>` | Create a new user and log in as them. |
| `gator login <name>` | Switch to an existing user. |
| `gator addfeed <name> <url>` | Add a new RSS feed and follow it. |
| `gator feeds` | List all feeds that have been added. |
| `gator follow <url>` | Follow an existing feed by its URL. |
| `gator following` | List the feeds the current user follows. |
| `gator unfollow <url>` | Stop following a feed. |
| `gator agg <duration>` | Continuously fetch new posts (e.g. `gator agg 1m`). |
| `gator browse [limit]` | Show the latest posts from feeds you follow (defaults to 2). |
| `gator users` | List all registered users. |

### Example workflow

```bash
# Create a user (this also logs you in)
gator register alice

# Add and follow a feed
gator addfeed "Hacker News RSS" https://hnrss.org/newest

# Start the aggregator in one terminal to collect posts every minute
gator agg 1m

# In another terminal, browse the most recent posts
gator browse 5
```
