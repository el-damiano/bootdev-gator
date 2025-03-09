package main

import (
	"context"
	"fmt"
	"time"

	"github.com/el-damiano/bootdev-gator/internal/database"
	"github.com/google/uuid"
)

type command struct {
	Name string
	Args []string
}

func handlerLogin(state *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("command %s expects [name] argument", cmd.Name)
	}

	username := cmd.Args[0]
	user, err := state.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("couldn't find user %s: %w", username, err)
	}

	err = state.config.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("%s set as current user\n", username)
	return nil
}

func handlerRegister(state *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("command %s expects [name] argument", cmd.Name)
	}

	username := cmd.Args[0]
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      username,
	}
	user, err := state.db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error creating user %w", err)
	}

	err = state.config.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	fmt.Printf("user %s was created\n", username)
	fmt.Printf("%+v\n", user)
	return nil
}

func handlerUsers(state *state, cmd command) error {
	_ = cmd
	users, err := state.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error retrieving users %w", err)
	}
	for _, user := range users {
		msg := fmt.Sprintf("* %s", user.Name)
		if user.Name == state.config.UserCurrent {
			msg += " (current)"
		}
		fmt.Println(msg)
	}
	return nil
}

func handlerReset(state *state, cmd command) error {
	_ = cmd
	err := state.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't reset the database: %w", err)
	}
	fmt.Printf("database reset, hope you're happy")
	return nil
}

func handlerAgg(state *state, cmd command) error {
	_ = state
	url := "https://www.wagslane.dev/index.xml"
	if len(cmd.Args) > 0 {
		url = cmd.Args[0]
	}

	feed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}

	fmt.Printf("Feed: %+v\n\n", feed)

	return nil
}

func handlerFeedAdd(state *state, cmd command) error {
	_ = state
	if len(cmd.Args) < 2 {
		return fmt.Errorf("command %s expects [name] and [url] arguments", cmd.Name)
	}

	name := cmd.Args[0]
	url := cmd.Args[1]
	user, err := state.db.GetUser(context.Background(), state.config.UserCurrent)
	if err != nil {
		return fmt.Errorf("error getting user: %w", err)
	}

	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}
	feed, err := state.db.CreateFeed(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error creating feed: %w", err)
	}

	fmt.Println("feed created successfully:")
	printFeed(feed)
	fmt.Println()
	fmt.Println("=====================================")
	return nil
}

func handlerFeedsList(state *state, cmd command) error {
	_ = state
	_ = cmd
	feeds, err := state.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds: %w", err)
	}

	fmt.Println("Feeds:")
	for _, feed := range feeds {
		fmt.Printf("User: %s\n", feed.UserName)
		fmt.Printf("Name: %s\n", feed.Name)
		fmt.Printf("URL: %s\n", feed.Url)
		fmt.Println()
	}

	return nil
}

func handlerFeedFollow(state *state, cmd command) error {
	_ = state
	if len(cmd.Args) < 1 {
		return fmt.Errorf("command %s expects [url] argument", cmd.Name)
	}

	user, err := state.db.GetUser(context.Background(), state.config.UserCurrent)
	if err != nil {
		return fmt.Errorf("error getting user: %w", err)
	}

	url := cmd.Args[0]
	feed, err := state.db.GetFeedsByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error getting feed: %w", err)
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := state.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error getting feeds: %w", err)
	}

	fmt.Printf("now following: %s\n", feedFollow.FeedName)

	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
}

type commandRegistry struct {
	reg map[string]func(*state, command) error
}

func (cmdReg *commandRegistry) register(name string, function func(*state, command) error) {
	cmdReg.reg[name] = function
}

func (cmdReg *commandRegistry) run(s *state, cmd command) error {
	command, ok := cmdReg.reg[cmd.Name]
	if !ok {
		return fmt.Errorf("command %s not found\n", cmd.Name)
	}
	return command(s, cmd)
}
