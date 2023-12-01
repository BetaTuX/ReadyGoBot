package commands

import (
	"ReadyGoBot/commands/guards"
	"ReadyGoBot/store"
	"ReadyGoBot/utils"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type HotlapRegistration struct {
	firstInteration *discordgo.Interaction
	laptime         time.Duration
}

const (
	addTrackSelectComponentId = "addhotlap_track_select"

	timeOptionId = "time"
)

var (
	trackSelectPlaceholder = utils.LocalizedString{
		Fallback: "Select a track",
		Localized: map[discordgo.Locale]string{
			discordgo.French: "Sélectionnez un tracé",
		},
	}
	ongoingRegisteration = make(map[string]HotlapRegistration)
)

func parseTimeParam(userInput string) (time.Duration, error) {
	timeParts := make([]int, 3)
	var err error

	regex := regexp.MustCompile(`(?m)(?P<m>[0-9]+):(?P<s>[0-9]{0,2})\.?(?P<ms>[0-9]{0,3})?`)
	regexResults := regex.FindStringSubmatch(userInput)
	for i := 0; i < 3; i++ {
		if timeParts[i], err = strconv.Atoi(regexResults[i+1]); err != nil {
			return time.Duration(0), err
		}
	}
	msString := strings.Replace(fmt.Sprintf("%-3d", timeParts[2]), " ", "0", -1)
	timeStr := fmt.Sprintf("%dm%ds%sms", timeParts[0], timeParts[1], msString)
	laptime, parseError := time.ParseDuration(timeStr)

	return laptime, parseError
}

var AddHotlapCommand = RGCommand{
	Command: discordgo.ApplicationCommand{
		Name:        "hotlapadd",
		Description: "Share your lastest hotlap",
		DescriptionLocalizations: &map[discordgo.Locale]string{
			discordgo.French: "Partagez votre dernier meilleur temps",
		},
		Type: discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name: timeOptionId,
				NameLocalizations: map[discordgo.Locale]string{
					discordgo.French: "temps",
				},
				Description: "Your best laptime (MM:SS.00)",
				DescriptionLocalizations: map[discordgo.Locale]string{
					discordgo.French: "Votre meilleur temps (MM:SS.00)",
				},
				Required: true,
				Type:     discordgo.ApplicationCommandOptionString,
			},
		},
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		optionMap := MapOptions(i.ApplicationCommandData().Options)
		laptime, parseError := parseTimeParam(optionMap[timeOptionId].StringValue())

		if parseError == nil {
			ongoingRegisteration[i.Member.User.ID] = HotlapRegistration{
				firstInteration: i.Interaction,
				laptime:         laptime,
			}
			localizedFmtString := utils.LocalizedString{
				Fallback: "Hotlap time: `%s`\n\nTrack:",
				Localized: map[discordgo.Locale]string{
					discordgo.French: "Meilleur temps : `%s`\n\nTracé :",
				},
			}

			if trackOptions, trackListError := guards.TrackListIsPopulated(s, i); trackListError == nil {
				if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf(localizedFmtString.GetLocaleString(i.Locale), laptime.String()),
						Flags:   discordgo.MessageFlagsEphemeral,
						Components: []discordgo.MessageComponent{
							discordgo.ActionsRow{
								Components: []discordgo.MessageComponent{
									discordgo.SelectMenu{
										CustomID:    addTrackSelectComponentId,
										Placeholder: trackSelectPlaceholder.GetLocaleString(i.Locale),
										Options:     trackOptions,
									},
								},
							},
						},
					},
				}); err != nil {
					log.Println("AddHotlap: trace:\n", err)
				}
			}
		}
	},

	ComponentHandlers: map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		addTrackSelectComponentId: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			driverUid := i.Member.User.ID
			matchedRegistration, isRegistrated := ongoingRegisteration[driverUid]
			if isRegistrated {
				selectedTrackId := strings.Join(i.MessageComponentData().Values, ";")
				selectedTrack := store.TrackStore.GetTrack(selectedTrackId)
				timeImproved := store.HotlapStore.SetHotlap(selectedTrackId, driverUid, matchedRegistration.laptime)
				defer func() {
					// Delete original message once the announcement has been made
					delete(ongoingRegisteration, driverUid)
					if err := s.InteractionResponseDelete(matchedRegistration.firstInteration); err != nil {
						log.Println("AddHotlap: track_select: Couldn't delete base message\n", err)
					}
				}()
				if timeImproved {
					congratsFmtStr := utils.LocalizedString{
						Fallback: "Congrats <@%s> for his/her hottest lap yet!\n%s on %s :checkered_flag::tada:",
						Localized: map[discordgo.Locale]string{
							discordgo.French: "Félicitation <@%s> pour son meilleurs temps !\n%s sur %s :checkered_flag::tada:",
						},
					}
					contentStr := fmt.Sprintf(
						congratsFmtStr.GetLocaleString(i.Locale),
						driverUid,
						utils.FormatDuration(matchedRegistration.laptime, i.Locale),
						selectedTrack.Name,
					)
					if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: contentStr,
							Embeds: []*discordgo.MessageEmbed{
								{
									Type: discordgo.EmbedTypeImage,
									Image: &discordgo.MessageEmbedImage{
										URL: selectedTrack.Picture.URL,
									},
								},
							},
						},
					}); err != nil {
						log.Println("AddHotlap: track_select: Couldn't post new success message\n", err)
					}
				} else {
					hottestLap, _ := store.HotlapStore.GetDriverHottestLap(driverUid, selectedTrackId)
					outFmtStr := utils.LocalizedString{
						Fallback: "No improvments since last session\nYour hottest lap stands at: %s",
						Localized: map[discordgo.Locale]string{
							discordgo.French: "Pas d'améliorations depuis la dernière session\nVotre meilleur temps est : %s",
						},
					}
					contentStr := fmt.Sprintf(
						outFmtStr.GetLocaleString(i.Locale),
						utils.FormatDuration(hottestLap.Time, i.Locale),
					)
					if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: contentStr,
							Flags:   discordgo.MessageFlagsEphemeral,
						},
					}); err != nil {
						log.Println("AddHotlap: track_select: Couldn't post new success message\n", err)
					}
				}
			}
		},
	},
}
