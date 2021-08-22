package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Name        string
	Description string
	Command     func(s *discordgo.Session, m *discordgo.MessageCreate)
}

type CommandRouter struct {
	Prefix     string
	Commands   []Command
	commandMap map[string]Command
}

func (cr CommandRouter) initialize(session *discordgo.Session) {
	cr.commandMap = make(map[string]Command)
	for _, cmd := range cr.Commands {
		cr.commandMap[cmd.Name] = cmd
	}
	session.AddHandler(cr.runCommand)
}

func (cr CommandRouter) runCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		// We don't reply to ourself.
		return
	}
	if strings.HasPrefix(m.Content, cr.Prefix) {
		// We're handling this command!
		commandName := strings.TrimPrefix(m.Content, cr.Prefix)
		cr.commandMap[commandName].Command(s, m)
	}
}
