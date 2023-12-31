package commands

import (
	"ReadyGoBot/store"
	"ReadyGoBot/utils"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/iancoleman/strcase"
)

const (
	idField      string = "id"
	labelField   string = "label"
	pictureField string = "picture"
)

var AddTrackCommand = RGCommand{
	Command: discordgo.ApplicationCommand{
		Name:        "trackadd",
		Description: "Add a track to the leaderboard",
		DescriptionLocalizations: &map[discordgo.Locale]string{
			discordgo.French: "Ajouter une piste à tableau des scores",
		},
		Type: discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name: labelField,
				NameLocalizations: map[discordgo.Locale]string{
					discordgo.French: "label",
				},
				Description: "The name of the track",
				DescriptionLocalizations: map[discordgo.Locale]string{
					discordgo.French: "Le nom du circuit",
				},
				Type:     discordgo.ApplicationCommandOptionString,
				Required: true,
			},
			{
				Name: pictureField,
				NameLocalizations: map[discordgo.Locale]string{
					discordgo.French: "image",
				},
				Description: "A picture of the track (either a visual or the track plan)",
				DescriptionLocalizations: map[discordgo.Locale]string{
					discordgo.French: "Un visuel pour le circuit ou le tracé",
				},
				Type:     discordgo.ApplicationCommandOptionAttachment,
				Required: true,
			},
			{
				Name: idField,
				NameLocalizations: map[discordgo.Locale]string{
					discordgo.French: "id",
				},
				Description: "The internal track ID",
				DescriptionLocalizations: map[discordgo.Locale]string{
					discordgo.French: "L'identifiant interne du circuit",
				},
				Type: discordgo.ApplicationCommandOptionString,
			},
		},
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		var fmtString utils.LocalizedString
		data := i.ApplicationCommandData()
		optionMap := MapOptions(data.Options)
		pictureId := optionMap[pictureField].Value
		trackName := optionMap[labelField].StringValue()
		embeds := make([]*discordgo.MessageEmbed, 0)
		var picture *discordgo.MessageAttachment

		if pictureId != nil {
			picture = data.Resolved.Attachments[pictureId.(string)]
			embeds = append(embeds, &discordgo.MessageEmbed{
				Type: discordgo.EmbedTypeImage,
				Image: &discordgo.MessageEmbedImage{
					URL: picture.URL,
				},
			})
		}

		var id string
		if optionMap[idField] == nil {
			id = strcase.ToSnake(trackName)
		} else {
			id = optionMap[idField].StringValue()
		}

		trackUpdated := store.TrackStore.SetTrack(store.Track{
			TrackId: id,
			Name:    trackName,
			Picture: store.PictureAttachment(*picture),
		})

		if trackUpdated {
			fmtString = utils.LocalizedString{
				Fallback: "%s has been updated! :new:",
				Localized: map[discordgo.Locale]string{
					discordgo.French: "%s a été mis à jour ! :new:",
				},
			}
		} else {
			fmtString = utils.LocalizedString{
				Fallback: "New track available! :tada:\n# %s",
				Localized: map[discordgo.Locale]string{
					discordgo.French: "Nouveau circuit disponible ! :tada:\n# %s",
				},
			}
		}

		if respErr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf(fmtString.GetLocaleString(i.Locale), trackName),
				Embeds:  embeds,
			},
		}); respErr != nil {
			log.Println("AddTrack: trace:\n", respErr)
		}
	},
}
