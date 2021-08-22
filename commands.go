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

func newCommandRouter(prefix string, commands []Command) (router CommandRouter) {
	router = CommandRouter{
		Prefix:     prefix,
		Commands:   commands,
		commandMap: make(map[string]Command),
	}
	return router
}

func (cr CommandRouter) initialize(session *discordgo.Session) {
	cr.buildCommandMap()
	session.AddHandler(cr.runCommand)
}

func (cr CommandRouter) buildCommandMap() {
	for k := range cr.commandMap {
		delete(cr.commandMap, k)
	}

	for _, cmd := range cr.Commands {
		cr.commandMap[cmd.Name] = cmd
	}
}

func (cr CommandRouter) addRouter(newCommand Command) {
	cr.Commands = append(cr.Commands, newCommand)
	cr.buildCommandMap()
}

func (cr CommandRouter) runCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		// We don't reply to ourself.
		return
	}
	if strings.HasPrefix(m.Content, cr.Prefix) {
		// We're handling this command!
		commandName := strings.TrimPrefix(m.Content, cr.Prefix)
		cmd, ok := cr.commandMap[commandName]
		if ok {
			cmd.Command(s, m)
		}
	}
}
