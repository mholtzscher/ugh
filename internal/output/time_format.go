package output

import (
	"time"

	"github.com/mholtzscher/ugh/internal/config"
)

type TimeFormatter struct {
	location *time.Location
	layout   string
}

func NewTimeFormatter(cfg config.Display) *TimeFormatter {
	layout := cfg.DatetimeFormat
	if layout == "" {
		layout = "2006-01-02 15:04"
	}
	location := loadLocation(cfg.Timezone)
	return &TimeFormatter{
		location: location,
		layout:   layout,
	}
}

func (f *TimeFormatter) Format(t time.Time) string {
	return t.In(f.location).Format(f.layout)
}

func loadLocation(tz string) *time.Location {
	if tz == "local" || tz == "" {
		return time.Local
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.UTC
	}
	return loc
}
