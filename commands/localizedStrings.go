package commands

import "github.com/bwmarrin/discordgo"

type LocalizedString struct {
	fallback  string
	localized map[discordgo.Locale]string
}

func (localizedStr *LocalizedString) getLocaleString(locale discordgo.Locale) string {
	fmtString, ok := localizedStr.localized[locale]

	if !ok {
		fmtString = localizedStr.fallback
	}
	return fmtString
}
