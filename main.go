package main

import (
	"fmt"

	"github.com/zulkou/blog-aggregator/internal/config"
)

func main() {
    cfg, err := config.Read()
    if err != nil {
        fmt.Println(err)
    }

    err = cfg.SetUser("zulkou")
    if err != nil {
        fmt.Println(err)
    }

    cfg, err = config.Read()
    if err != nil {
        fmt.Println(err)
    }

    fmt.Printf("db_url: \"%v\"\ncurrent_user_name: \"%v\"", cfg.DBURL, cfg.CurrentUserName)
}
