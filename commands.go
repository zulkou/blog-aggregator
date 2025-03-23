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

func handlerLogin(s *state, cmd command) error {
    if len(cmd.args) == 0 || len(cmd.args) > 1 {
        return errors.New("The login command expects ONE argument")
    }

    name := cmd.args[0]

    _, err := s.db.GetUser(context.Background(), name)
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
    if len(cmd.args) == 0 || len(cmd.args) > 1 {
        return errors.New("The register command expects ONE argument")
    }

    name := cmd.args[0]
    
    _, err := s.db.GetUser(context.Background(), name)
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
    if len(cmd.args) != 0 {
        return errors.New("The agg command expects ZERO arguments")
    }

    url := "https://www.wagslane.dev/index.xml"

    res, err := rss.FetchFeed(context.Background(), url)
    if err != nil {
        return fmt.Errorf("Failed fetching feeds: %w", err)
    }

    fmt.Println(res)

    return nil
}
