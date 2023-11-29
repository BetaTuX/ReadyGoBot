package commands

import (
	"fmt"

	"ReadyGoBot/store"

	"github.com/bwmarrin/discordgo"
)

var ListTracksCommand = RGCommand{
	Command: discordgo.ApplicationCommand{
		Name:        "listtrack",
		Description: "List all available tracks",
		DescriptionLocalizations: &map[discordgo.Locale]string{
			discordgo.French: "Lister tout les tracés disponibles",
		},
		Type: discordgo.ChatApplicationCommand,
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		fmtString := LocalizedString{
			fallback: "List of available tracks:",
			localized: map[discordgo.Locale]string{
				discordgo.French: "Les tracés disponibles sont les suivants :",
			},
		}
		trackListFmted := ""

		for _, track := range store.TrackStore.GetTracks() {
			trackListFmted = fmt.Sprintf("%s\n- %s", trackListFmted, track.Name)
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("%s\n%s", fmtString.getLocaleString(i.Locale), trackListFmted),
			},
		})
	},
}
