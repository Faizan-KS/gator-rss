package main

import (
	"context"
	"fmt"

	"github.com/Faizan-KS/gator-rss/internal/database"
)

func handlerDeleteFeedByURL(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	x := database.DeleteFeedByURLParams{
		UserID: user.ID,
		Url:    cmd.Args[0],
	}
	err := s.db.DeleteFeedByURL(context.Background(), x)
	if err != nil {
		return err
	}
	fmt.Println("unfollowed the feed")
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	followList, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't fetch follow list for user: %w", err)
	}
	if len(followList) == 0 {
		fmt.Println("Follow list empty!! Start following by entering the url")
		return nil
	}
	for _, list := range followList {
		fmt.Printf("Feed of : %s\n", list.FeedOf)
		fmt.Printf("Feed Name : %s\n", list.FeedName)
		fmt.Printf("Feed URL : %s\n", list.FeedUrl)
		fmt.Printf("Feed ID : %s\n", list.FeedID)
		fmt.Println("-------------------------")
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	url := cmd.Args[0]
	feeds, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}

	newFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feeds.ID,
	})

	if err != nil {
		return fmt.Errorf("couldn't create follow for the user: %w", err)
	}

	fmt.Println(newFollow)
	return nil
}
