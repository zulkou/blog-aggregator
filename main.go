package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/zulkou/blog-aggregator/internal/config"
)

func main() {
    cfg, err := config.Read()
    if err != nil {
        fmt.Println(err)
    }

    s := state {
        cfg: &cfg,
    }
    cmds := commands{
        handlers: make(map[string]func(*state, command) error),
    }

    cmds.register("login", handlerLogin)

    args := os.Args

    if len(args) < 2 {
        fmt.Println("require command name")
        os.Exit(1)
    }
    
    cmd := command{
        name: args[1],
        args: args[2:],
    }

    err = cmds.run(&s, cmd)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

type state struct {
    cfg     *config.Config
}

type command struct {
    name    string
    args    []string
}

type commands struct {
    handlers    map[string]func(s *state, cmd command) error
}

func handlerLogin(s *state, cmd command) error {
    if len(cmd.args) == 0 || len(cmd.args) > 1 {
        return errors.New("only accepts one argument")
    }

    err := s.cfg.SetUser(cmd.args[0])
    if err != nil {
        return fmt.Errorf("user not found: %w", err)
    }

    fmt.Printf("State assigned to %s, Welcome!\n", s.cfg.CurrentUserName)
    return nil
}

func (c *commands) register(name string, f func(s *state, cmd command) error) {
    c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
    handler, exists := c.handlers[cmd.name]
    if !exists {
        return fmt.Errorf("command not found: %s", cmd.name)
    }

    return handler(s, cmd)
}

