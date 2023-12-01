package commands

import (
	"ReadyGoBot/commands/guards"
	"ReadyGoBot/store"
	"ReadyGoBot/utils"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type LeaderboardRequest struct {
	firstInteration *discordgo.Interaction
}

const (
	trackLeaderboardSelectComponentId = "listhotlap_track_select"
)

var (
	ongoingLeaderboardSelection = make(map[string]LeaderboardRequest)
)

var GetLeaderboardCommand = RGCommand{
	Command: discordgo.ApplicationCommand{
		Name:        "leaderboard",
		Description: "Prints the leaderboard for a track",
		DescriptionLocalizations: &map[discordgo.Locale]string{
			discordgo.French: "Afficher un tableau des temps pour un tracé",
		},
		Type: discordgo.ChatApplicationCommand,
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if trackOptions, trackListError := guards.TrackListIsPopulated(s, i); trackListError == nil {
			ongoingLeaderboardSelection[i.Member.User.ID] = LeaderboardRequest{
				firstInteration: i.Interaction,
			}
			localizedFmtString := utils.LocalizedString{
				Fallback: "Select which track leaderboard you want to display",
				Localized: map[discordgo.Locale]string{
					discordgo.French: "Sélectionnez un tracé pour afficher son tableau des temps",
				},
			}
			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: localizedFmtString.GetLocaleString(i.Locale),
					Flags:   discordgo.MessageFlagsEphemeral,
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.SelectMenu{
									CustomID:    trackLeaderboardSelectComponentId,
									Placeholder: trackSelectPlaceholder.GetLocaleString(i.Locale),
									Options:     trackOptions,
								},
							},
						},
					},
				},
			}); err != nil {
				log.Println("GetLeaderboard: trace:\n", err)
			}
		}
	},
	ComponentHandlers: map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		trackLeaderboardSelectComponentId: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if ongoingLeaderboard, selectionExists := ongoingLeaderboardSelection[i.Member.User.ID]; selectionExists {
				defer func() {
					delete(ongoingLeaderboardSelection, i.Member.User.ID)
					if err := s.InteractionResponseDelete(ongoingLeaderboard.firstInteration); err != nil {
						log.Println("AddHotlap: track_select: Couldn't delete base message\n", err)
					}
				}()
				introFmt := utils.LocalizedString{
					Fallback: "# **%s**\n## Leaderboard:",
					Localized: map[discordgo.Locale]string{
						discordgo.French: "# **%s**\n## Liste des temps :",
					},
				}
				selectedTrackId := i.MessageComponentData().Values[0]
				selectedTrack := store.TrackStore.GetTrack(selectedTrackId)
				if laps, err := guards.TrackHasHotlapRegistered(s, i, selectedTrackId); err == nil {
					lapStrings := make([]string, 0, len(laps))

					for position, v := range laps {
						var prefix string
						username := fmt.Sprintf("<@%s>", v.DriverUid)
						formattedLaptime := utils.FormatDuration(v.Time, i.Locale)

						if 0 <= position && position < 3 {
							prefixes := [...]string{
								"# :first_place:",
								"## :second_place",
								"### :third_place:",
							}

							prefix = prefixes[position]
						} else {
							prefix = fmt.Sprintf("%d", position+1)
						}

						lapStrings = append(lapStrings, fmt.Sprintf("%s %s\t%s", prefix, username, formattedLaptime))
					}

					contentString := fmt.Sprintf(
						"%s\n%s",
						fmt.Sprintf(introFmt.GetLocaleString(i.Locale), selectedTrack.Name),
						strings.Join(lapStrings, "\n"),
					)
					embeds := []*discordgo.MessageEmbed{
						{
							Type: discordgo.EmbedTypeImage,
							Image: &discordgo.MessageEmbedImage{
								URL: selectedTrack.Picture.URL,
							},
						},
					}
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content:         contentString,
							Embeds:          embeds,
							AllowedMentions: nil,
						},
					})
				}
			}
		},
	},
}
