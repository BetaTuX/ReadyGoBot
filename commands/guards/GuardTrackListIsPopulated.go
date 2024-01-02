package guards

import (
	"ReadyGoBot/store"
	"ReadyGoBot/utils"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func getTrackOptions(defaultValue string) []discordgo.SelectMenuOption {
	tracks := store.TrackStore.GetTracks()
	trackOptions := make([]discordgo.SelectMenuOption, 0, len(tracks))

	for _, track := range tracks {
		trackOptions = append(trackOptions, discordgo.SelectMenuOption{
			Label:   track.Name,
			Value:   track.TrackId,
			Default: defaultValue == track.TrackId,
		})
	}

	return trackOptions
}

func TrackListIsPopulated(s *discordgo.Session, i *discordgo.InteractionCreate) ([]discordgo.SelectMenuOption, error) {

	trackOptions := getTrackOptions("")

	if len(trackOptions) <= 0 {
		noTrackStr := utils.LocalizedString{
			Fallback: "Sadly no tracks have been set up yet. Please add tracks with the `addtrack` slash command.",
			Localized: map[discordgo.Locale]string{
				discordgo.French: "Malheureusement aucun tracé n'est disponible. Ajoutez-en via la commande `addtrack` puis réessayez.",
			},
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: noTrackStr.GetLocaleString(i.Locale),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return nil, fmt.Errorf("no tracks found in storage")
	}
	return trackOptions, nil
}
