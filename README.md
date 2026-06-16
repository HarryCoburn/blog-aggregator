# Gator RSS Reader (Boot.dev Project)

This is a CLI RSS reader, though at this point it only outputs descriptions, not the full content.

## Requirements

To install the program, you will need:

- Go
- Postgresql

### Configuration File

Gator expects a .gatorconfig.json file in your home directory. "db_url" is the URL for your Postgresql database. You will need to disable sslmode. "current_user_name" can be null. This will get set when you start using commands.

Example:

```{"db_url":"postgres://postgres:@localhost:5432/gator?sslmode=disable", "current_user_name":"harry"}```

## Installation

You can download the repo and use ```go run . <command>```. If you want to install it, use ```go install``` and you can run the program using the command ```gator <command>```.

## Commands

Use these commands by typing ```gator command```

reset: Resets the database to a clean start.

#### User Commands
register <user>: Registers a username. It will automatically set the current user to what you register. Users are tied to a set of feeds.
login <user>: Changes the current username to <user>.
users: Lists the users in the database.

#### Feed Commands
feeds: Lists the feeds in the database across all users.
follow <url>: Follow an RSS feed at the given URL.
following: Lists the feeds the current user is following.
unfollow <url>: Unfollows an RSS feed at the given URL.
agg <time>: Queries the feeds your following for a <time> duration. e.g. 1m for every minute. At this time, this must be run, then then canceled with Ctrl+C before you can browse the posts. You'll get output in your console every time the feed is queried. TODO: make this a service in the background.
browse <num>: Browse the latest posts you have downloaded. Default is two.


## TODO
- Proper tests that work with web querying.
- Agg service in the background.
- Getting the content of the posts, not just the meta descriptions.
