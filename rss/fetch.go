package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
    client := &http.Client{
        Timeout: 5 *  time.Second,
    }

    req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
    if err != nil {
        return nil, fmt.Errorf("Failed to create request: %w", err)
    }

    req.Header.Set("User-Agent", "gator")

    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("Failed to execute request: %w", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("Failed to read response: %w", err)
    }

    var result RSSFeed
    err = xml.Unmarshal(body, &result)
    if err != nil {
        return nil, fmt.Errorf("Failed to parse response: %w", err)
    }

    result.Channel.Title = html.UnescapeString(result.Channel.Title)
    result.Channel.Description = html.UnescapeString(result.Channel.Description)

    for i := range result.Channel.Item {
        result.Channel.Item[i].Title = html.UnescapeString(result.Channel.Item[i].Title)
        result.Channel.Item[i].Description = html.UnescapeString(result.Channel.Item[i].Description)
    }

    return &result, nil
}
