package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/HarryCoburn/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type command struct {
	name string
	args []string
}

type commands struct {
	commands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	elem, ok := c.commands[cmd.name]
	if !ok {
		return fmt.Errorf("Command %s is not registered", cmd.name)
	}
	err := elem(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if cmd.args == nil {
		return fmt.Errorf("Login command requires a username.")
	}

	// Check if the user exists

	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("Could not login with user %s. Error: %v", cmd.args[0], err)
	}

	s.cfg.SetUser(user.Name)

	fmt.Printf("Username set to %s\n", user.Name)
	return nil

}

func handlerRegister(s *state, cmd command) error {
	if cmd.args == nil {
		return fmt.Errorf("Register command requires a username.")
	}
	fmt.Printf("Attempting to register user: %s\n", cmd.args[0])
	var params database.CreateUserParams
	params.ID = uuid.New()
	params.CreatedAt = time.Now()
	params.UpdatedAt = params.CreatedAt
	params.Name = cmd.args[0]
	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Could not register user: %v", err)
	}
	s.cfg.SetUser(cmd.args[0])
	fmt.Printf("The returned user is: %v\n Setting current user to this.", user)

	return nil
}

func handlerReset(s *state, cmd command) error {
	fmt.Println("Resetting users...")
	err := s.db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("Could not reset user table. Reason: %v", err)
	}
	fmt.Println("Reset successful.")
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("User table likely empty. Error: %v", err)
	}
	currentUser := s.cfg.Current_user_name
	for _, user := range users {
		if user == currentUser {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Include a time duration string when using agg (e.g. 1m).")
	}
	fmt.Printf("Collecting feeds every %s\n", cmd.args[0])
	duration, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("Given time duration is in an incorrect format: %v", err)
	}
	ticker := time.NewTicker(duration)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("AddFeed requires a name and a URL in quotes.")
	}
	var params database.CreateFeedParams
	params.ID = uuid.New()
	params.CreatedAt = time.Now()
	params.UpdatedAt = params.CreatedAt
	params.Name = cmd.args[0]
	params.Url = cmd.args[1]
	params.UserID = user.ID

	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Could not add feed: %v", err)
	}
	fmt.Println(feed)

	newCmd := command{}
	newCmd.name = "follow"
	newCmd.args = []string{cmd.args[1]}
	handlerFollow(s, newCmd, user)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Feeds table likely empty. Error: %v", err)
	}

	for _, feed := range feeds {
		fmt.Printf("Name: %s\n", feed.Name)
		fmt.Printf("URL: %s\n", feed.Url)
		userName, err := s.db.GetUserFromId(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("Could not find userID of feed creator. Error: %v", err)
		}
		fmt.Printf("Created by: %s\n\n", userName)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	url := cmd.args[0]
	var params database.CreateFeedFollowParams
	params.ID = uuid.New()
	params.CreatedAt = time.Now()
	params.UpdatedAt = params.CreatedAt
	params.UserID = user.ID
	feed, err := s.db.GetFeedFromURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Given URL not found in feeds table. %v", err)
	}
	params.FeedID = feed.ID

	feed_follow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error creating a feed_follow row. %v", err)
	}
	fmt.Printf("Feed Follow created\n")
	fmt.Printf("Created by: %s\n", feed_follow.UserName)
	fmt.Printf("Feed Name: %s\n", feed_follow.FeedName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	followedFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("Error in handlerFollowing. %v", err)
	}
	if len(followedFeeds) == 0 {
		fmt.Println("Current user is following no feeds")
		return nil
	}
	for _, feed := range followedFeeds {
		fmt.Println(feed)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	feed, err := s.db.GetFeedFromURL(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("Given URL not found in feeds table. %v", err)
	}
	var params database.DeleteFeedFollowParams
	params.FeedID = feed.ID
	params.UserID = user.ID
	err = s.db.DeleteFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Could not unfollow feed: %v", err)
	}
	return nil
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	var params database.MarkFeedFetchedParams
	params.ID = nextFeed.ID
	params.LastFetchedAt = sql.NullTime{Time: time.Now(), Valid: true}
	s.db.MarkFeedFetched(context.Background(), params)
	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}
	for _, item := range feed.Channel.Item {
		fmt.Println(item.Title)
	}
	return nil
}
