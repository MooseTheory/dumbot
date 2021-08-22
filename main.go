package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/lus/dgc"
	"github.com/pelletier/go-toml"
)

type Config struct {
	Discord DiscordInfo
}

type DiscordInfo struct {
	Token string
}

var (
	config Config
	dg     *discordgo.Session
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
	dg, err = discordgo.New("Bot " + config.Discord.Token)
	if err != nil {
		fatalLog(err)
	}

	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fatalLog(err)
	}

	sc := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sc)
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				// TIMER!
				setStatus(dg)
			case <-sc:
				dg.Close()
				done <- true
			}
		}
	}()
	setStatus(dg)
	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	<-done
	fmt.Println("Goodbye!")
	dg.Close()
}

func setStatus(s *discordgo.Session) {
	// s.Update
	// s.UpdateStatus(0, "#help for a list of commands")
}

func pingHandler(ctx *dgc.Ctx) {
	ctx.RespondText("Ping")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
