package commands

import (
	"fmt"
	"strings"

	"ReadyGoBot/store"
	"ReadyGoBot/utils"

	"github.com/bwmarrin/discordgo"
)

var ListTracksCommand = RGCommand{
	Command: discordgo.ApplicationCommand{
		Name:        "tracklist",
		Description: "List all available tracks",
		DescriptionLocalizations: &map[discordgo.Locale]string{
			discordgo.French: "Lister tout les tracés disponibles",
		},
		Type: discordgo.ChatApplicationCommand,
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		fmtString := utils.LocalizedString{
			Fallback: "List of available tracks:",
			Localized: map[discordgo.Locale]string{
				discordgo.French: "Les tracés disponibles sont les suivants :",
			},
		}
		tracks := store.TrackStore.GetTracks()
		trackNameList := make([]string, 0, len(tracks))

		for _, track := range tracks {
			trackNameList = append(trackNameList, track.Name)
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("%s\n- %s", fmtString.GetLocaleString(i.Locale), strings.Join(trackNameList, "\n")),
			},
		})
	},
}
