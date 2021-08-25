package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Name        string
	Description string
	Aliases     []string
	Command     func(s *discordgo.Session, m *discordgo.MessageCreate)
}

type CommandRouter struct {
	Prefix     string
	Name       string
	Commands   []Command
	IgnoreCase bool
	commandMap map[string]Command
}

func newCommandRouter(prefix string, commands []Command) (router CommandRouter) {
	router = CommandRouter{
		Prefix:     prefix,
		Commands:   commands,
		commandMap: make(map[string]Command),
		IgnoreCase: true,
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
		commandName := cmd.Name
		if cr.IgnoreCase {
			commandName = strings.ToLower(commandName)
		}
		cr.commandMap[commandName] = cmd
		if len(cmd.Aliases) > 0 {
			for _, alias := range cmd.Aliases {
				commandName = alias
				if cr.IgnoreCase {
					commandName = strings.ToLower(commandName)
				}
				cr.commandMap[commandName] = cmd
			}
		}
	}
	cr.commandMap["help"] = Command{
		Name:        "help",
		Description: "Show help for this bot",
		Command:     cr.helpCommand,
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
		if cr.IgnoreCase {
			commandName = strings.ToLower(commandName)
		}
		cmd, ok := cr.commandMap[commandName]
		if ok {
			cmd.Command(s, m)
		}
	}
}

func (cr CommandRouter) helpCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	var fields []*discordgo.MessageEmbedField
	for _, cmd := range cr.Commands {
		var nameArr []string
		nameArr = append(nameArr, cmd.Name)
		nameArr = append(nameArr, cmd.Aliases...)
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   strings.Join(nameArr, ", "),
			Value:  cmd.Description,
			Inline: false,
		})
	}

	message := &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       fmt.Sprintf("%s Help", cr.Name),
		Description: "Commands available",
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       0x34a1eb,
		Fields:      fields,
	}
	s.ChannelMessageSendEmbed(m.ChannelID, message)
}
