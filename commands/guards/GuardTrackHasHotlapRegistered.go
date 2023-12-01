package guards

import (
	"ReadyGoBot/store"
	"ReadyGoBot/utils"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func TrackHasHotlapRegistered(s *discordgo.Session, i *discordgo.InteractionCreate, trackId string) ([]store.Hotlap, error) {
	lapDisplayLimit := 10
	laps := store.HotlapStore.GetTrackHotlapList(trackId, lapDisplayLimit)

	if len(laps) <= 0 {
		noTrackStr := utils.LocalizedString{
			Fallback: "Sadly no hotlap have been recorded yet. Set your own hotlap with the `hotlapadd` command.",
			Localized: map[discordgo.Locale]string{
				discordgo.French: "Malheureusement aucun record n'a été ajouté pour le moment. Ajoutez le vôtre grâce à la commande `hotlapadd`.",
			},
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: noTrackStr.GetLocaleString(i.Locale),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return nil, fmt.Errorf("no hotlap record found for that track")
	}
	return laps, nil
}
