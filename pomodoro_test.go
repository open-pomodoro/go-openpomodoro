package openpomodoro

import (
	"encoding"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_PomodoroInterfaces(t *testing.T) {
	var _ encoding.TextMarshaler = Pomodoro{}
	var _ encoding.TextUnmarshaler = &Pomodoro{}
	var _ json.Marshaler = Pomodoro{}
}

func TestPomodoro_MarshalJSON(t *testing.T) {
	p := &Pomodoro{
		StartTime:   time.Date(2016, 06, 14, 12, 0, 0, 0, time.UTC),
		Duration:    25 * time.Minute,
		Tags:        []string{"a", "b"},
		Description: "A description",
	}
	b, err := p.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t,
		`{"start_time":"2016-06-14T12:00:00Z","description":"A description","duration":25,"tags":["a","b"]}`,
		string(b))
}

func TestPomodoro_MarshalText(t *testing.T) {
	p := &Pomodoro{
		StartTime:   time.Date(2016, 06, 14, 12, 0, 0, 0, time.UTC),
		Duration:    25 * time.Minute,
		Tags:        []string{"a", "b"},
		Description: "A description",
	}
	b, err := p.MarshalText()
	assert.Nil(t, err)
	assert.Equal(t,
		"2016-06-14T12:00:00Z description=\"A description\" duration=25 tags=a,b",
		string(b))
}

func Test_Matches(t *testing.T) {
	timestamp, err := time.Parse(TimeFormat, "2026-06-14T12:34:56-04:00")
	require.Nil(t, err)

	a := &Pomodoro{StartTime: timestamp}

	b := &Pomodoro{StartTime: timestamp}
	assert.True(t, a.Matches(b))

	b = &Pomodoro{StartTime: timestamp.Add(time.Minute)}
	assert.False(t, a.Matches(b))

	b = &Pomodoro{StartTime: timestamp.Add(500 * time.Millisecond)}
	assert.True(t, a.Matches(b))

	b = &Pomodoro{StartTime: timestamp.Add(-500 * time.Millisecond)}
	assert.True(t, a.Matches(b))
}

func Test_MarshalText(t *testing.T) {
	timestamp, err := time.Parse(TimeFormat, "2026-06-14T12:34:56-04:00")
	require.Nil(t, err)

	var p *Pomodoro
	var actual []byte
	var expected string

	p = &Pomodoro{
		StartTime: timestamp,
		Duration:  25 * time.Minute,
	}
	expected = `2026-06-14T12:34:56-04:00 duration=25`
	actual, err = p.MarshalText()
	require.Nil(t, err)
	assert.Equal(t, expected, string(actual))

	p = &Pomodoro{
		StartTime:   timestamp,
		Duration:    25 * time.Minute,
		Description: "working on stuff",
		Tags:        []string{"work", "stuff"},
	}
	expected = `2026-06-14T12:34:56-04:00 description="working on stuff" duration=25 tags=work,stuff`
	actual, err = p.MarshalText()
	require.Nil(t, err)
	assert.Equal(t, expected, string(actual))
}

func Test_UnmarshalText_timeOnly(t *testing.T) {
	p := &Pomodoro{}
	err := p.UnmarshalText([]byte(`2026-06-14T12:34:56-04:00`))
	require.Nil(t, err)

	startTime, err := time.Parse(TimeFormat, "2026-06-14T12:34:56-04:00")
	require.Nil(t, err)
	expected := &Pomodoro{StartTime: startTime}

	assert.Equal(t, expected, p)
}

func Test_UnmarshalText_timeOnlyWithNewline(t *testing.T) {
	p := &Pomodoro{}
	err := p.UnmarshalText([]byte(`2026-06-14T12:34:56-04:00
`))
	require.Nil(t, err)

	startTime, err := time.Parse(TimeFormat, "2026-06-14T12:34:56-04:00")
	require.Nil(t, err)
	expected := &Pomodoro{StartTime: startTime}

	assert.Equal(t, expected, p)
}

func Test_UnmarshalText_singleLineAttributes(t *testing.T) {
	p := &Pomodoro{}
	err := p.UnmarshalText([]byte(`2026-06-14T12:34:56-04:00 description="working on stuff" duration=25 tags=work,stuff`))
	require.Nil(t, err)

	startTime, err := time.Parse(TimeFormat, "2026-06-14T12:34:56-04:00")
	require.Nil(t, err)
	expected := &Pomodoro{
		StartTime:   startTime,
		Description: "working on stuff",
		Duration:    25 * time.Minute,
		Tags:        []string{"work", "stuff"},
	}

	assert.Equal(t, expected, p)
}

func Test_UnmarshalText_empty(t *testing.T) {
	p := &Pomodoro{}
	err := p.UnmarshalText([]byte(``))
	require.Nil(t, err)
	assert.True(t, p.IsInactive())
}

func Test_UnmarshalText_whitespace(t *testing.T) {
	p := &Pomodoro{}
	err := p.UnmarshalText([]byte(" \n "))
	require.Nil(t, err)
	assert.True(t, p.IsInactive())
}

func Test_UnmarshalText_multipleEntries(t *testing.T) {
	p := &Pomodoro{}
	err := p.UnmarshalText([]byte(`2026-06-14T12:34:56-04:00 description="working on stuff" duration=25 tags=work,stuff
2026-06-14T12:34:56-04:00 description="working on stuff" duration=25 tags=work,stuff`))
	assert.Error(t, err)
}

func Test_ApplySettings_empty(t *testing.T) {
	p := &Pomodoro{}

	s := &Settings{
		DefaultPomodoroDuration: 25 * time.Minute,
		DefaultTags:             []string{"work"},
	}

	p.ApplySettings(s)

	assert.Equal(t, p.Duration, 25*time.Minute)
	assert.Equal(t, p.Tags, []string{"work"})
}

func Test_ApplySettings_existing(t *testing.T) {
	p := &Pomodoro{
		Duration: 30 * time.Minute,
		Tags:     []string{"play"},
	}

	s := &Settings{
		DefaultPomodoroDuration: 25 * time.Minute,
		DefaultTags:             []string{"work"},
	}

	p.ApplySettings(s)

	assert.Equal(t, p.Duration, 30*time.Minute)
	assert.Equal(t, p.Tags, []string{"play"})
}

func Test_DurationMinutes(t *testing.T) {
	p := Pomodoro{}

	p.Duration = 30 * time.Minute
	assert.Equal(t, 30, p.DurationMinutes())

	p.Duration = 29*time.Minute + 30*time.Second
	assert.Equal(t, 30, p.DurationMinutes())

	p.Duration = 29*time.Minute + 29*time.Second
	assert.Equal(t, 29, p.DurationMinutes())
}

func Test_EndTime(t *testing.T) {
	start, err := time.Parse(TimeFormat, "2026-06-14T12:34:56-04:00")
	require.Nil(t, err)
	expected, err := time.Parse(TimeFormat, "2026-06-14T12:59:56-04:00")
	require.Nil(t, err)

	p := Pomodoro{StartTime: start, Duration: 25 * time.Minute}
	assert.Equal(t, expected, p.EndTime())
}

func Test_IsActive(t *testing.T) {
	timeFunc = time.Now

	p := NewPomodoro()
	p.Duration = 25 * time.Minute

	cases := map[time.Duration]bool{
		24 * time.Minute: true,
		25 * time.Minute: false,
		26 * time.Minute: false,
		time.Hour:        false,
		-time.Hour:       true,
		0 * time.Second:  true,
	}

	for duration, expected := range cases {
		p.StartTime = timeFunc().Add(-duration)
		assert.Equal(t, expected, p.IsActive(), duration.String())
	}
}

func Test_IsDone(t *testing.T) {
	timeFunc = time.Now

	p := NewPomodoro()
	p.Duration = 25 * time.Minute

	cases := map[time.Duration]bool{
		24 * time.Minute: false,
		25 * time.Minute: true,
		26 * time.Minute: true,
		time.Hour:        true,
		-time.Hour:       false,
		0 * time.Second:  false,
	}

	for duration, expected := range cases {
		p.StartTime = timeFunc().Add(-duration)
		assert.Equal(t, expected, p.IsDone(), duration.String())
	}
}

func Test_IsInactive_true(t *testing.T) {
	assert.True(t, EmptyPomodoro().IsInactive())
}

func Test_IsInactive_false(t *testing.T) {
	timestamp, err := time.Parse(
		TimeFormat,
		"2026-06-14T12:34:56-04:00",
	)
	require.Nil(t, err)
	p := Pomodoro{StartTime: timestamp}

	assert.False(t, p.IsInactive())
}

func Test_Remaining(t *testing.T) {
	timeFunc = time.Now

	p := NewPomodoro()
	p.Duration = 25 * time.Minute

	assert.Equal(t, 0, p.Remaining().Seconds())

	cases := map[time.Duration]time.Duration{
		0 * time.Minute:  25 * time.Minute,
		1 * time.Minute:  24 * time.Minute,
		24 * time.Minute: 1 * time.Minute,
		25 * time.Minute: 0 * time.Minute,
		26 * time.Minute: -1 * time.Minute,
	}

	for duration, expected := range cases {
		p.StartTime = timeFunc().Add(-duration)
		assert.InDelta(t, expected.Seconds(), p.Remaining().Seconds(), 1)
	}
}

func Test_RemainingMinutes(t *testing.T) {
	p := NewPomodoro()
	p.Duration = 25 * time.Minute

	assert.Equal(t, 0, p.RemainingMinutes())

	cases := map[time.Duration]int{
		0 * time.Minute:  25,
		1 * time.Minute:  24,
		24 * time.Minute: 1,
		25 * time.Minute: 0,
		26 * time.Minute: -1,

		29 * time.Second:                25,
		30 * time.Second:                24,
		24*time.Minute + 29*time.Second: 1,
		24*time.Minute + 30*time.Second: 0,
	}

	for duration, expected := range cases {
		p.StartTime = timeFunc().Add(-duration)
		assert.Equal(t, expected, p.RemainingMinutes())
	}
}
