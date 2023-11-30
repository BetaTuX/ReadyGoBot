package utils

import (
	"fmt"
	"math"
	"time"

	"github.com/bwmarrin/discordgo"
)

func FormatDuration(duration time.Duration, locale discordgo.Locale) string {
	fmtStr := LocalizedString{
		Fallback: "%dm %.3fsec",
	}

	return fmt.Sprintf(
		fmtStr.GetLocaleString(locale),
		int(duration.Minutes()),
		duration.Seconds()-(math.Trunc(duration.Minutes())*60),
	)
}
