package openpomodoro

import (
	"encoding"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_SettingsInterfaces(t *testing.T) {
	var _ encoding.TextUnmarshaler = &Settings{}
}

func Test_SetDefaults_empty(t *testing.T) {
	s := &Settings{}
	d := &DefaultSettings
	s.SetDefaults(&DefaultSettings)

	assert.Equal(t, d, s)
}

func Test_SetDefaults_filled(t *testing.T) {
	s := &Settings{
		DailyGoal:               10,
		DefaultBreakDuration:    10 * time.Minute,
		DefaultPomodoroDuration: 20 * time.Minute,
		DefaultTags:             []string{"work"},
	}

	expected := &Settings{}
	expected.SetDefaults(s)

	d := &DefaultSettings
	s.SetDefaults(&DefaultSettings)

	assert.NotEqual(t, d, s)
	assert.Equal(t, expected, s)
}

func Test_Settings_UnmarshalText(t *testing.T) {
	s := &Settings{}

	err := s.UnmarshalText([]byte(`
	  daily_goal=8
	  default_break_duration=10
	  default_pomodoro_duration=20
	  default_tags=billable,work
	`))
	require.Nil(t, err)

	assert.Equal(t, 8, s.DailyGoal)
	assert.Equal(t, 10*time.Minute, s.DefaultBreakDuration)
	assert.Equal(t, 20*time.Minute, s.DefaultPomodoroDuration)
	assert.Equal(t, []string{"billable", "work"}, s.DefaultTags)
}
