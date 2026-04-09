package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Faizan-KS/gator-rss/internal/config"
	"github.com/Faizan-KS/gator-rss/internal/database"
)

// *******PLEASE PRACTICE THESE BELOW FUNCTION AGAIN***********
// scrapeFeeds is a helper function

func handlerGetPostsForUser(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	for _, x := range feeds {
		posts, err := s.db.GetPostsForUser(context.Background(), x.FeedID)
		if err != nil {
			return err
		}
		fmt.Println("\nPosts for the feed\n", posts)
	}
	return nil
}

func scrapeFeeds(ctx context.Context, s *state) error {
	fmt.Println("\nAdding posts for link")

	feeds, err := s.db.GetAllFeeds(ctx)
	if err != nil {
		fmt.Println("GetNextFeedToFetch error:")
		return err
	}
	for _, feed := range feeds {
		err := s.db.MarkFeedFetched(ctx, feed.ID)
		if err != nil {
			return err
		}
		rss, err := config.FetchFeed(ctx, feed.Url)
		if err != nil {
			return err
		}
		fmt.Printf("\nFetching from %s\n", feed.Url)

		for i := 0; i < len(rss.Channel.Item) && i < 3; i++ {
			item := rss.Channel.Item[i]
			t, err := time.Parse(time.RFC1123Z, item.PubDate)
			if err != nil {
				return err
			}
			err = s.db.CreatePost(ctx, database.CreatePostParams{
				Title:       item.Title,
				Url:         sql.NullString{String: item.Link, Valid: true},
				Description: sql.NullString{String: item.Description, Valid: true},
				PublishedAt: sql.NullTime{Time: t, Valid: true},
				FeedID:      feed.ID,
			})
			if err != nil {
				return err
			}
			fmt.Println("Post Added")
		}
		time.Sleep(5 * time.Second)
	}
	fmt.Println("\nFeed aggregation cycle complete")
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: agg <time_between_reqs>")
	}

	duration, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %s\n\n", duration)

	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)
		<-sigChan
		fmt.Println("\nShutting down aggregator...")
		cancel()
	}()

	//Immediately run it once
	if err := scrapeFeeds(ctx, s); err != nil {
		fmt.Println("error scraping feeds:", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := scrapeFeeds(ctx, s); err != nil {
				fmt.Println("error scraping feeds:", err)
			}
		}
	}
}

//------------------------------------------------------------

func handlerGetAllFeeds(s *state, cmd command) error {
	allFeeds, err := s.db.GetAllFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feeds: %w", err)
	}
	if len(allFeeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}
	fmt.Printf("Found %d feeds\n", len(allFeeds))
	fmt.Println("-------------------------------")
	for _, feed := range allFeeds {
		user, err := s.db.GetUsersById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("couldn't get user: %w", err)
		}
		fmt.Printf("Name: %s\n", feed.Name)
		fmt.Printf("URL: %s\n", feed.Url)
		fmt.Printf("Name_of_User: %s\n", user.Name)
		fmt.Println("----------------------------")
	}
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	// user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	// if err != nil {
	// 	return err
	// }

	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		Name:   name,
		Url:    url,
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't follow feed: %w", err)
	}

	fmt.Println("Feed created successfully:")
	fmt.Println(feed)
	fmt.Println("=====================================")
	return nil
}
