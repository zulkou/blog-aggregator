package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/zulkou/blog-aggregator/internal/database"
	"github.com/zulkou/blog-aggregator/rss"
)

type command struct {
    name    string
    args    []string
}

type commands struct {
    handlers    map[string]func(s *state, cmd command) error
}

func (c *commands) register(name string, f func(s *state, cmd command) error) {
    c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
    handler, exists := c.handlers[cmd.name]
    if !exists {
        return fmt.Errorf("command not found: %s\n", cmd.name)
    }

    return handler(s, cmd)
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
    return func(s *state, cmd command) error {
        if s.cfg.CurrentUserName == "" {
            return errors.New("You need to logged in to use this function")
        }

        user, err := s.db.GetUserByName(context.Background(), s.cfg.CurrentUserName)
        if err != nil {
            return fmt.Errorf("Authentication error: %w", err)
        }

        return handler(s, cmd, user)
    }
}

func scrapeFeeds(s *state) error {
    feed, err := s.db.GetNextFeedToFetch(context.Background())
    if err != nil {
        fmt.Printf("Reading DB fault: %v\n", err)
        return fmt.Errorf("Failed to fetch next feed: %w", err)
    }

    err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
        ID: feed.ID,
        LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
        UpdatedAt: time.Now(),
    })
    if err != nil {
        fmt.Printf("Marking fault: %v\n", err)
        return fmt.Errorf("Failed to mark fetched feed: %w", err)
    }

    rssfeed, err := rss.FetchFeed(context.Background(), feed.Url)
    if err != nil {
        fmt.Printf("Fetching fault: %v\n", err)
        return fmt.Errorf("Failed to fetch feed content: %w", err)
    }

    for _, rssitem := range(rssfeed.Channel.Item) {
        fmt.Printf("- %s\n", rssitem.Title)
    }

    return nil
}

func handlerLogin(s *state, cmd command) error {
    if len(cmd.args) != 1 {
        return errors.New("The login command expects ONE argument")
    }

    name := cmd.args[0]

    _, err := s.db.GetUserByName(context.Background(), name)
    if err != nil {
        return fmt.Errorf("Failed to retrieve %v: %w\n", name, err)
    }

    err = s.cfg.SetUser(name)
    if err != nil {
        return fmt.Errorf("User not found: %w\n", err)
    }

    fmt.Printf("State assigned to %s, Welcome!\n", s.cfg.CurrentUserName)
    return nil
}

func handlerRegister(s *state, cmd command) error {
    if len(cmd.args) != 1 {
        return errors.New("The register command expects ONE argument")
    }

    name := cmd.args[0]
    
    _, err := s.db.GetUserByName(context.Background(), name)
    if err == nil {
        return fmt.Errorf("User %s already exists\n", name)
    } else if !errors.Is(err, sql.ErrNoRows) {
        return fmt.Errorf("Error checking if user exists: %v\n", err)
    }

    _, err = s.db.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: name})
    if err != nil {
        return fmt.Errorf("Failed to create user: %w\n", err)
    }
    
    err = s.cfg.SetUser(name)
    if err != nil {
        return errors.New(fmt.Sprintf(("User failed to login: %s"), name))
    }

    fmt.Printf("User %s successfully created!\nWelcome %s!\n", name, name)

    return nil
}

func handlerReset(s *state, cmd command) error {
    if len(cmd.args) != 0 {
        return errors.New("The reset command expects ZERO arguments")
    }

    err := s.db.DeleteUsers(context.Background())
    if err != nil {
        return fmt.Errorf("Failed resetting database: %v\n", err)
    }

    fmt.Println("Database successfully resetted")

    return nil
}

func handlerUsers(s *state, cmd command) error {
    if len(cmd.args) != 0 {
        return errors.New("The users command expects ZERO arguments")
    }

    users, err := s.db.GetUsers(context.Background())
    if err != nil {
        return fmt.Errorf("Failed fetching users: %w", err)
    }

    for _, user := range(users) {
        if s.cfg.CurrentUserName == user.Name {
            fmt.Printf("* %s (current)\n", user.Name)
        } else {
            fmt.Printf("* %s\n", user.Name) 
        }
    }

    return nil
}

func handlerAgg(s *state, cmd command) error {
    if len(cmd.args) != 1 {
        return errors.New("The agg command expects ONE arguments")
    }

    time_between_reqs, err := time.ParseDuration(cmd.args[0])
    if err != nil {
        return fmt.Errorf("Failed to parse input args: %w", err)
    }

    fmt.Printf("Collecting feeds every %v\n", time_between_reqs)

    ticker := time.NewTicker(time_between_reqs)
    for ; ; <-ticker.C {
        scrapeFeeds(s)
    }

    return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
    if len(cmd.args) != 2 {
        return errors.New("The addfeed command expects TWO arguments")
    }

    url := cmd.args[0]
    name := cmd.args[1]

    feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
        ID: uuid.New(),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Name: name,
        Url: url,
        UserID: user.ID,
    })

    if err != nil {
        return fmt.Errorf("Failed to store feed to db: %w", err)
    }

    _, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
        ID: uuid.New(),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        UserID: user.ID,
        FeedID: feed.ID,
    })

    if err != nil {
        return fmt.Errorf("Failed to auto-follow after feed creation: %w", err)
    }

    fmt.Printf("Name: %v\nURL: %v\n-- Current user automatically follow created Feed --\n", feed.Name, feed.Url)

    return nil
}

func handlerFeeds(s *state, cmd command) error {
    if len(cmd.args) != 0 {
        return errors.New("The feeds command expects ZERO arguments")
    }

    feeds, err := s.db.GetFeeds(context.Background())
    if err != nil {
        return fmt.Errorf("Failed to fetch feeds: %w", err)
    }

    for _, feed := range(feeds) {
        user, err := s.db.GetUserByID(context.Background(), feed.UserID)
        if err != nil {
            return fmt.Errorf("Failed to fetch feed's user: %w", err)
        }

        fmt.Printf("---\nName: %v\nURL: %v\nUser: %v\n", feed.Name, feed.Url, user.Name)
    }

    return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
    if len(cmd.args) != 1 {
        return errors.New("The follow command expects ONE argument")
    }

    feedName := cmd.args[0]

    feed, err := s.db.GetFeedByURL(context.Background(), feedName)
    if err != nil {
        return fmt.Errorf("Failed to retrieve feed with provided URL: %w", err)
    }

    follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
        ID: uuid.New(),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        UserID: user.ID,
        FeedID: feed.ID,
    })

    if err != nil {
        return fmt.Errorf("Failed to create follow data: %w", err)
    }

    fmt.Printf("Current user: %s\nSuccess on following %s\n", follow.UserName, follow.FeedName)

    return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
    if len(cmd.args) != 0 {
        return errors.New("The following command expects ZERO arguments")
    }

    feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
    if err != nil {
        return fmt.Errorf("Failed to fetch feed data for current user: %w", err)
    }

    fmt.Printf("Feeds followed by %s:\n", user.Name)
    for _, feed := range(feeds) {
        fmt.Printf("- %s\n", feed.FeedName)
    }

    return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
    if len(cmd.args) != 1 {
        return errors.New("The unfollow command expects ONE argument")
    }

    feedName := cmd.args[0]

    feed, err := s.db.GetFeedByURL(context.Background(), feedName)
    if err != nil {
        return fmt.Errorf("Failed to fetch feed: %w", err)
    }

    err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
        UserID: user.ID,
        FeedID: feed.ID,
    })

    fmt.Printf("Successfully unfollowed %s", feed.Name)

    return nil
}
