package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
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
	client := &http.Client{}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Could not build request from URL: %s", feedURL)
	}

	request.Header.Set("User-Agent", "gator")

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Could not complete request: %s", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not read the response: %s", err)
	}
	var rssResponse RSSFeed

	if err := xml.Unmarshal(body, &rssResponse); err != nil {
		return nil, fmt.Errorf("Could not unmarshal XML return: %v", err)
	}

	rssResponse.Channel.Title = html.UnescapeString(rssResponse.Channel.Title)
	rssResponse.Channel.Description = html.UnescapeString(rssResponse.Channel.Description)
	for _, item := range rssResponse.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}
	return &rssResponse, nil
}
