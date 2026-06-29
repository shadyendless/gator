package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/shadyendless/gator/internal/commands"
	"github.com/shadyendless/gator/internal/database"
	"github.com/shadyendless/gator/internal/middleware"
	"github.com/shadyendless/gator/internal/state"
	"github.com/shadyendless/gator/internal/xml"
)

func main() {
	s, err := state.New()
	if err != nil {
		fmt.Printf("[ERROR]: %v\n", err)
		os.Exit(1)
	}

	comms := commands.New()
	comms.Register("login", handlerLogin)
	comms.Register("register", handlerRegister)
	comms.Register("reset", handlerReset)
	comms.Register("users", handlerUsers)
	comms.Register("agg", handlerAgg)
	comms.Register("addfeed", middleware.LoggedIn(handlerAddFeed, &s))
	comms.Register("feeds", handlerFeeds)
	comms.Register("follow", middleware.LoggedIn(handlerFollow, &s))
	comms.Register("following", middleware.LoggedIn(handlerFollowing, &s))
	comms.Register("unfollow", middleware.LoggedIn(handlerUnfollow, &s))
	comms.Register("browse", middleware.LoggedIn(handlerBrowse, &s))

	if len(os.Args) < 2 {
		fmt.Println("[ERROR]: Not enough arguments were passed")
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	if err = comms.Run(&s, commands.Command{
		Name: command,
		Args: args,
	}); err != nil {
		fmt.Printf("[ERROR]: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func handlerLogin(s *state.State, cmd commands.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("username is required")
	}

	username := cmd.Args[0]
	user, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		return err
	}

	if err = s.Config.SetUser(user.Name); err != nil {
		return err
	}

	fmt.Printf("The user has been set to: %s\n", user.Name)
	return nil
}

func handlerRegister(s *state.State, cmd commands.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("name is required")
	}

	name := cmd.Args[0]
	user, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	})

	if err != nil {
		return err
	}

	fmt.Printf("Registered the following user: %v\n", user)
	if err = handlerLogin(s, commands.Command{
		Name: "login",
		Args: []string{user.Name},
	}); err != nil {
		return err
	}

	return nil
}

func handlerReset(s *state.State, cmd commands.Command) error {
	if err := s.Db.Reset(context.Background()); err != nil {
		return err
	}

	return nil
}

func handlerUsers(s *state.State, cmd commands.Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		fmt.Printf("* %s", user.Name)

		if user.Name == s.Config.CurrentUserName {
			fmt.Print(" (current)")
		}

		fmt.Print("\n")
	}

	return nil
}

func handlerAgg(s *state.State, cmd commands.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("duration is required")
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func handlerAddFeed(s *state.State, cmd commands.Command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("name and url are required")
	}

	name := cmd.Args[0]
	feedURL := cmd.Args[1]

	feed, err := s.Db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       feedURL,
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}

	fmt.Print(feed)

	return nil
}

func handlerFeeds(s *state.State, cmd commands.Command) error {
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("- %s (%s)\n", feed.Name, feed.Url)
		fmt.Printf("    Added by %s\n", feed.CreatedBy)
	}

	return nil
}

func handlerFollow(s *state.State, cmd commands.Command, user database.User) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("you must provide a url")
	}

	feedUrl := cmd.Args[0]
	feed, err := s.Db.GetFeedByUrl(context.Background(), feedUrl)
	if err != nil {
		return err
	}

	feedFollow, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID:    user.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	fmt.Printf("%s followed %s\n", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func handlerFollowing(s *state.State, cmd commands.Command, user database.User) error {
	feeds, err := s.Db.GetFeedFollows(context.Background(), user.ID)
	if err != nil {
		return err
	}

	fmt.Printf("%s is following:\n", user.Name)
	for _, feed := range feeds {
		fmt.Printf(" - %s\n", feed)
	}

	return nil
}

func handlerUnfollow(s *state.State, cmd commands.Command, user database.User) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("you must provide a url")
	}

	feedUrl := cmd.Args[0]
	feed, err := s.Db.GetFeedByUrl(context.Background(), feedUrl)
	if err != nil {
		return err
	}

	return s.Db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
}

func handlerBrowse(s *state.State, cmd commands.Command, user database.User) error {
	limit := int32(2)

	if len(cmd.Args) > 0 {
		parseResult, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return err
		}

		limit = int32(parseResult)
	}

	posts, err := s.Db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf(`%s
Published on %s
%s
Read at: %s
---------
`, post.Title, post.PublishedAt.Format(time.RFC1123), post.Description.String, post.Url)
	}

	return nil
}

func scrapeFeeds(s *state.State) error {
	dbFeed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("Fetching feed \"%s\"...\n", dbFeed.Name)
	xmlFeed, err := xml.FetchFeed(context.Background(), dbFeed.Url)
	if err != nil {
		return err
	}
	fmt.Printf("Fetched %d item(s) from \"%s\"\n", len(xmlFeed.Channel.Item), dbFeed.Name)

	err = s.Db.MarkFeedFetched(context.Background(), dbFeed.ID)
	if err != nil {
		return err
	}

	for _, item := range xmlFeed.Channel.Item {
		pubTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			fmt.Printf("Could not parse the PubDate for \"%s\"\n  Received: %s\n", item.Title, item.PubDate)
			continue
		}

		_, err = s.Db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
			},
			PublishedAt: pubTime,
			FeedID:      dbFeed.ID,
		})
	}

	return nil
}
