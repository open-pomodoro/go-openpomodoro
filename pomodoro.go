package openpomodoro

import (
	"bytes"
	"math"
	"time"

	"github.com/justincampbell/go-logfmt"
)

const (
	// TimeFormat is the format we generate and expect to parse timestamps in.
	TimeFormat = time.RFC3339
)

var (
	charNewline = []byte("\n")
	charSpace   = []byte(" ")
	timeFunc    = time.Now
)

// Pomodoro holds a single Pomodoro and related information.
type Pomodoro struct {
	StartTime   time.Time
	Description string        `logfmt:"description"`
	Duration    time.Duration `logfmt:"duration,m"`
	Tags        []string      `logfmt:"tags"`
}

// NewPomodoro returns a Pomodoro with defaults set.
func NewPomodoro() *Pomodoro {
	return &Pomodoro{
		Duration: DefaultSettings.DefaultPomodoroDuration,
	}
}

// EmptyPomodoro returns an empty Pomodoro.
func EmptyPomodoro() *Pomodoro {
	return &Pomodoro{}
}

// String return a string representation of the Pomodoro.
func (p Pomodoro) String() string {
	b, _ := p.MarshalText()
	return string(b)
}

// Matches returns whether or not another Pomodoro has the same StartTime.
func (p Pomodoro) Matches(o *Pomodoro) bool {
	delta := p.StartTime.Sub(o.StartTime)
	return delta >= -time.Second && delta <= time.Second
}

// MarshallText marshals the Pomodoro's start time and attributes into a text
// string.
func (p Pomodoro) MarshalText() ([]byte, error) {
	timestamp := []byte(p.StartTime.Format(TimeFormat))
	attributes, err := logfmt.Encode(p)
	if err != nil {
		return nil, err
	}

	return bytes.Join([][]byte{timestamp, attributes}, charSpace), nil
}

// UnmarshalText updates a Pomodoro's timestamp and attributes from a byte
// string.
func (p *Pomodoro) UnmarshalText(b []byte) error {
	b = bytes.TrimSpace(b)
	parts := bytes.SplitN(b, charSpace, 2)

	var timestamp []byte
	var attributes []byte

	switch len(parts) {
	case 0:
		return nil
	case 1:
		timestamp = parts[0]
	case 2:
		if parts[0] == nil {
			return nil
		}
		timestamp = parts[0]
		attributes = parts[1]
	default:
		return nil
	}

	if bytesAllWhitespace(timestamp) {
		return nil
	}

	startTime, err := time.Parse(TimeFormat, string(timestamp))
	if err != nil {
		return err
	}

	p.StartTime = startTime

	err = logfmt.Unmarshal(attributes, p)
	if err != nil {
		return err
	}

	return nil
}

// ApplySettings sets the Pomodoro's defaults from settings if they are
// considered to be missing.
func (p *Pomodoro) ApplySettings(s *Settings) {
	if p.Duration == 0 {
		p.Duration = s.DefaultPomodoroDuration
	}

	if len(p.Tags) == 0 {
		p.Tags = s.DefaultTags
	}
}

// DurationMinutes returns the Pomodoro's duration in minutes.
func (p *Pomodoro) DurationMinutes() int {
	return round(p.Duration.Minutes())
}

// EndTime returns the time the Pomodoro would end.
func (p *Pomodoro) EndTime() time.Time {
	return p.StartTime.Add(p.Duration)
}

// IsActive returns whether or not a Pomodoro is active.
func (p *Pomodoro) IsActive() bool {
	return !p.IsInactive() && !p.IsDone()
}

// IsDone returns whether or not a Pomodoro was active and is now done.
func (p *Pomodoro) IsDone() bool {
	if p.IsInactive() {
		return false
	}
	return timeFunc().After(p.EndTime())
}

// IsInactive returns whether or not a Pomodoro is empty/not set/etc.
func (p *Pomodoro) IsInactive() bool {
	return p.StartTime.IsZero()
}

// Remaining returns the remaining duration of the Pomodoro.
func (p *Pomodoro) Remaining() time.Duration {
	if p.IsInactive() {
		return time.Duration(0)
	}

	return p.EndTime().Sub(timeFunc())
}

// RemainingMinutes returns the remaining duration of the Pomodoro in minutes.
// Partial minutes are rounded up and down normally, so that there are 25
// minutes remaining for 30 seconds after the Pomodoro starts, and 0 for 30
// seconds before it completes.
func (p *Pomodoro) RemainingMinutes() int {
	return round(p.Remaining().Minutes())
}

func bytesAllWhitespace(b []byte) bool {
	return len(bytes.TrimSpace(b)) == 0
}

func round(f float64) int {
	return int(math.Floor(f + 0.5))
}
