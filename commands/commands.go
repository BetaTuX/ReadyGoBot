package commands

import (
	"github.com/bwmarrin/discordgo"
)

type RGCommand struct {
	Command discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
}
