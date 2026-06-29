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

## Setting up the database

`gator` does **not** create its tables automatically — the binary expects the
schema to already exist. Migrations are managed with
[Goose](https://github.com/pressly/goose) and live in the `sql/schema`
directory of this repository, so you'll need a clone of the repo to run them
(`go install` only installs the binary, not the migration files).

1. Install the Goose CLI:

   ```bash
   go install github.com/pressly/goose/v3/cmd/goose@latest
   ```

2. Clone this repository (if you haven't already) and move into the schema
   directory:

   ```bash
   git clone https://github.com/shadyendless/gator.git
   cd gator/sql/schema
   ```

3. Run all the "up" migrations against your database, using the same connection
   string you put in `.gatorconfig.json`:

   ```bash
   goose postgres "postgres://username:password@localhost:5432/gator" up
   ```

You only need to do this once (and again whenever new migrations are added). To
roll the most recent migration back, run `goose postgres "<db_url>" down`.

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
