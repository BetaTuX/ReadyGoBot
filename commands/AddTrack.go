package commands

import (
	"ReadyGoBot/store"

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
		Name:        "addtrack",
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
		var fmtString LocalizedString
		options := i.ApplicationCommandData().Options
		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, opt := range options {
			optionMap[opt.Name] = opt
		}

		var id string
		if optionMap[idField] == nil {
			id = strcase.ToSnake(optionMap[labelField].StringValue())
		} else {
			id = optionMap[idField].StringValue()
		}

		trackUpdated := store.TrackStore.SetTrack(store.Track{
			Id:      id,
			Name:    optionMap[labelField].StringValue(),
			Picture: optionMap[pictureField].Value,
		})

		if trackUpdated {
			fmtString = LocalizedString{
				fallback: "You *updated* the track successfully! :tada:",
				localized: map[discordgo.Locale]string{
					discordgo.French: "Vous avez *mis à jour* le tracé avec succès ! :tada:",
				},
			}
		} else {
			fmtString = LocalizedString{
				fallback: "You *added* the track successfully! :tada:",
				localized: map[discordgo.Locale]string{
					discordgo.French: "Vous avez *ajouté* le tracé avec succès ! :tada:",
				},
			}
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmtString.getLocaleString(i.Locale),
			},
		})
	},
}
