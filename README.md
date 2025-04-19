# Boot.dev RSS Aggregator

CLI RSS blog aggregator microservice using Go, codename Gator, built as a
guided project of the Boot.dev "Back-end Developer Career Path".

# Features? highlights?

- RSS feed aggregation
- continuous aggregation of posts from RSS feeds
- follow/unfollow feeds
- local user profiles to differ which feeds a user follows

# Prerequisites

- postgresql (17.4+)
- go (1.24.0+)

# Install

1. Ensure `$GOPATH` is setup correctly, and add it to your `$PATH` in your
`.bashrc` / `.zshrc` like so:

```bash
export GOPATH=$(go env GOPATH)
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN
```

And then install it with:

```bash
go install github.com/el-damiano/bootdev-gator@latest
```

Alternatively, you can forego `$GOPATH` and just do the following.

```bash
git clone https://github.com/el-damiano/bootdev-gator.git &&
cd bootdev-gator &&
go install .
```

2. Set up a `~/.gatorconfig.json` file, for example:

```json
{
  "db_url": "postgres://postgres:@localhost:5432/gator?sslmode=disable",
  "current_user_name": "John Gator"
}
```

# Uninstall

Either remove the file under `$GOBIN/bootdev-gator`, or go to where you have
cloned the repo and run:

```bash
go clean -i -cache -modcache
```

# Usage

Usage:
  bootdev-gator [command]

Available Commands:
  login [username]        Switch to another user
  register [username]     Register a new user
  users                   List all users
  agg [interval]          Start aggregation
  addfeed [name] [url]    Add a feed to follow
  feeds                   List out all feeds followed by all users
  follow [url]            Follow a feed
  following               List out all feeds followed by current user
  unfollow [url]          Unfollow a feed
  browse [number]         List out post titles of current user
