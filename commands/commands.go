package commands

import (
	"github.com/bwmarrin/discordgo"
)

type RGCommand struct {
	Command           discordgo.ApplicationCommand
	Handler           func(s *discordgo.Session, i *discordgo.InteractionCreate)
	ComponentHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func MapOptions(options []*discordgo.ApplicationCommandInteractionDataOption) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	return optionMap
}
