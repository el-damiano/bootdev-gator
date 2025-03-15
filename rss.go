package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/el-damiano/bootdev-gator/internal/database"
	"github.com/google/uuid"
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for _, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}

	return &feed, nil
}

func scrapeFeeds(state *state) error {
	feedToFetch, err := state.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	state.db.MarkFeedFetched(context.Background(), feedToFetch.ID)

	feed, err := fetchFeed(context.Background(), feedToFetch.Url)
	if err != nil {
		return err
	}

	for _, item := range feed.Channel.Item {
		var nullTime sql.NullTime
		feedDate, err := time.Parse(time.RFC822, item.PubDate)
		if err == nil {
			nullTime.Time = feedDate
			nullTime.Valid = true
		}

		params := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UdpatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: nullTime,
			FeedID:      feedToFetch.ID,
		}

		_, err = state.db.CreatePost(context.Background(), params)
		if err != nil {
			return err
		}
	}

	return nil
}
