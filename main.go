package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/zulkou/blog-aggregator/internal/config"
	"github.com/zulkou/blog-aggregator/internal/database"
)

type state struct {
    cfg     *config.Config
    db      *database.Queries
}

func run() int {
    cfg, err := config.Read()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error in reading config file: %v\n", err)
        return 1
    }

    db, err := sql.Open("postgres", cfg.DBURL)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
        return 1
    }
    defer db.Close()
    dbQueries := database.New(db)
    
    s := state {
        cfg: &cfg,
        db: dbQueries,
    }
 
    cmds := commands{
        handlers: make(map[string]func(*state, command) error),
    }

    cmds.register("login", handlerLogin)
    cmds.register("register", handlerRegister)
    cmds.register("reset", handlerReset)
    cmds.register("users", handlerUsers)
    cmds.register("agg", handlerAgg)

    args := os.Args

    if len(args) < 2 {
        fmt.Fprintf(os.Stderr, "Commands not specified\n")
        return 1
    }
    
    cmd := command{
        name: args[1],
        args: args[2:],
    }

    err = cmds.run(&s, cmd)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error running command: %v\n", err)
        return 1
    }

    return 0
}

func main() {
    exitCode := run()
    os.Exit(exitCode)
}

