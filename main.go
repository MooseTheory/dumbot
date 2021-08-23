package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pelletier/go-toml"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type Config struct {
	Discord DiscordInfo
}

type DiscordInfo struct {
	Token string
}

var (
	config    Config
	session   *discordgo.Session
	client    *reddit.Client
	timeZones = map[string]string{
		"Eastern": "America/New_York",
		"Pacific": "America/Los_Angeles",
	}
	previousFashionCheck = time.Unix(0, 0)
	currentFashionReport *discordgo.MessageEmbed
)

func init() {
	contents, err := ioutil.ReadFile("config.toml")
	if err != nil {
		panic(err)
	}
	err = toml.Unmarshal(contents, &config)
	if err != nil {
		panic(err)
	}
	doLog("Starting")
}

func doLog(message string) {
	fmt.Println(message)
}

func fatalLog(err error) {
	doLog(err.Error())
	panic(err)
}

func main() {
	var err error
	// Set up discord stuff
	session, err = discordgo.New("Bot " + config.Discord.Token)
	if err != nil {
		fatalLog(err)
	}

	session.Identify.Intents = discordgo.IntentsGuildMessages

	commands := []Command{
		{
			Name:        "maintenance",
			Aliases:     []string{"m", "maint"},
			Description: "Fetch maintenance information",
			Command:     maintenanceCommand,
		},
		{
			Name:        "fashionreport",
			Aliases:     []string{"f", "FashionReport"},
			Description: "Fetch the latest fashion report",
			Command:     fashionReport,
		},
	}

	commandRouter := newCommandRouter("!", commands)
	commandRouter.IgnoreCase = true
	commandRouter.initialize(session)
	// End Discord stuff

	// Set up reddit stuff
	client, err = reddit.NewReadonlyClient()
	// End reddit stuff

	err = session.Open()
	if err != nil {
		fatalLog(err)
	}

	sc := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGTERM)
	go func() {
		for {
			select {
			case <-sc:
				doLog("Attempting graceful shutdown")
				session.Close()
				done <- true
			}
		}
	}()
	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	<-done
	fmt.Println("Goodbye!")
	session.Close()
}

func fashionReport(s *discordgo.Session, m *discordgo.MessageCreate) {
	if previousFashionCheck.Add(60 * time.Second).Before(time.Now()) {
		ctx := context.Background()
		previousFashionCheck = time.Now()
		searchOpts := *&reddit.ListPostSearchOptions{
			Sort: "new",
			ListPostOptions: reddit.ListPostOptions{
				ListOptions: reddit.ListOptions{
					Limit: 1,
				},
			},
		}
		// We aren't using the response part of this for anything currently, so it gets ignored
		posts, _, err := client.Subreddit.SearchPosts(ctx, "author:kaiyoko Fashion Report", "ffxiv", &searchOpts)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}
		// For debugging, I'm leaving this here, resp has the rate information.
		// fmt.Printf("posts: %+v\nresp: %+v\n", posts[0], resp)
		currentFashionReport = createFashionEmbed(posts[0])
	}
	s.ChannelMessageSendEmbed(m.ChannelID, currentFashionReport)
}

func maintenanceCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	allMaintenance, err := getMaintenance()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Could not fetch maintenance information")
	}
	curTime := time.Now()
	nextMaintenance := MaintenanceInfo{
		Title: "ERROR",
	}
	for _, maint := range allMaintenance {
		// I think these come in chronological order, so the "last" one we
		// encounter should be the next one.
		if maint.Start.After(curTime) {
			nextMaintenance = maint
		}
	}
	if nextMaintenance.Title == "ERROR" {
		s.ChannelMessageSend(m.ChannelID, "There isn't currently any scheduled maintenance! Enjoy the game!")
	} else {
		embed, err := createMaintenanceEmbed(nextMaintenance)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}
}

func createFashionEmbed(post *reddit.Post) *discordgo.MessageEmbed {
	// var fields []*discordgo.MessageEmbedField
	return &discordgo.MessageEmbed{
		Type:  discordgo.EmbedTypeImage,
		Title: post.Title,
		Image: &discordgo.MessageEmbedImage{
			URL: post.URL,
		},
	}
}

func createMaintenanceEmbed(maint MaintenanceInfo) (*discordgo.MessageEmbed, error) {
	// Do we need to do this every time? I don't wanna be dumb when we switch
	// from daylight saving to "standard" time
	easternLoc, err := time.LoadLocation(timeZones["Eastern"])
	if err != nil {
		return nil, err
	}
	pacificLoc, err := time.LoadLocation(timeZones["Pacific"])
	if err != nil {
		return nil, err
	}

	var fields []*discordgo.MessageEmbedField
	easternFieldText := fmt.Sprintf("From %s until %s", maint.Start.In(easternLoc).Format("02 Jan 2006, 3:04PM"), maint.End.In(easternLoc).Format("02 Jan 2006, 3:04PM"))
	pacificFieldText := fmt.Sprintf("From %s until %s", maint.Start.In(pacificLoc).Format("02 Jan 2006, 3:04PM"), maint.End.In(pacificLoc).Format("02 Jan 2006, 3:04PM"))
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Eastern",
		Value:  easternFieldText,
		Inline: false,
	})
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Pacific",
		Value:  pacificFieldText,
		Inline: false,
	})

	return &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       "Upcoming Maintenance",
		Description: maint.Title,
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       0xffff00,
		Fields:      fields,
	}, nil
}
