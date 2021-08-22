package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/pelletier/go-toml"
)

type Config struct {
	Discord DiscordInfo
}

type DiscordInfo struct {
	Token string
}

var (
	config  Config
	session *discordgo.Session
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
	session, err = discordgo.New("Bot " + config.Discord.Token)
	if err != nil {
		fatalLog(err)
	}

	session.Identify.Intents = discordgo.IntentsGuildMessages

	commands := []Command{
		{
			"ping",
			"Responds with pong",
			pongFunc,
		},
	}

	commandRouter := newCommandRouter("!", commands)
	commandRouter.initialize(session)

	err = session.Open()
	if err != nil {
		fatalLog(err)
	}

	sc := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sc)
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

func pongFunc(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "pong")
}

// func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
// 	// Ignore all messages created by the bot itself
// 	// This isn't required in this specific example but it's a good practice.
// 	if m.Author.ID == s.State.User.ID {
// 		return
// 	}
// 	// If the message is "ping" reply with "Pong!"
// 	if m.Content == "ping" {
// 		s.ChannelMessageSend(m.ChannelID, "Pong!")
// 	}

// 	// If the message is "pong" reply with "Ping!"
// 	if m.Content == "pong" {
// 		s.ChannelMessageSend(m.ChannelID, "Ping!")
// 	}
// }
