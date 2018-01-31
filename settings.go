package openpomodoro

import (
	"bytes"
	"time"

	"github.com/justincampbell/go-logfmt"
)

// Settings is a collection of user settings, which can come from a file, env
// var, or set from the client program.
type Settings struct {
	DailyGoal               int           `logfmt:"daily_goal"`
	DefaultBreakDuration    time.Duration `logfmt:"default_break_duration,m"`
	DefaultPomodoroDuration time.Duration `logfmt:"default_pomodoro_duration,m"`
	DefaultTags             []string      `logfmt:"default_tags"`
}

// DefaultSettings are used as a starting point before settings are overridden
// by the user.
var DefaultSettings = Settings{
	DailyGoal:               0,
	DefaultBreakDuration:    5 * time.Minute,
	DefaultPomodoroDuration: 25 * time.Minute,
	DefaultTags:             []string{},
}

// SetDefaults fills in settings values from another setting struct if the
// existing values are considered to not be set yet.
func (s *Settings) SetDefaults(d *Settings) {
	if s.DailyGoal == 0 {
		s.DailyGoal = d.DailyGoal
	}

	if s.DefaultBreakDuration == 0 {
		s.DefaultBreakDuration = d.DefaultBreakDuration
	}

	if s.DefaultPomodoroDuration == 0 {
		s.DefaultPomodoroDuration = d.DefaultPomodoroDuration
	}

	if len(s.DefaultTags) == 0 {
		s.DefaultTags = d.DefaultTags
	}
}

// UnmarshalText updates settings by parsing each key/value pair in logfmt.
func (s *Settings) UnmarshalText(b []byte) error {
	b = bytes.Replace(b, charNewline, charSpace, -1)
	return logfmt.Unmarshal(b, s)
}
