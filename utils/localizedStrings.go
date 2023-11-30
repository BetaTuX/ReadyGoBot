package utils

import "github.com/bwmarrin/discordgo"

type LocalizedString struct {
	Fallback  string
	Localized map[discordgo.Locale]string
}

func (localizedStr *LocalizedString) GetLocaleString(locale discordgo.Locale) string {
	fmtString, ok := localizedStr.Localized[locale]

	if !ok {
		fmtString = localizedStr.Fallback
	}
	return fmtString
}
