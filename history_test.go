package openpomodoro

import (
	"encoding"
	"encoding/json"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	empty = History{}

	a = &Pomodoro{StartTime: time.Date(2016, 06, 13, 12, 0, 0, 0, time.UTC)}
	b = &Pomodoro{StartTime: time.Date(2016, 06, 14, 12, 0, 0, 0, time.UTC)}
	c = &Pomodoro{StartTime: time.Date(2016, 06, 15, 12, 0, 0, 0, time.UTC)}

	one = History{Pomodoros: []*Pomodoro{b}}

	many = History{Pomodoros: []*Pomodoro{a, b, c}}
)

func Test_HistoryInterfaces(t *testing.T) {
	var _ encoding.TextMarshaler = History{}
	var _ json.Marshaler = History{}
	var _ sort.Interface = History{}
}

func TestHistory_MarshalJSON(t *testing.T) {
	p := &Pomodoro{
		StartTime:   time.Date(2016, 06, 14, 12, 0, 0, 0, time.UTC),
		Duration:    25 * time.Minute,
		Tags:        []string{"a", "b"},
		Description: "A description",
	}
	h := &History{Pomodoros: []*Pomodoro{p}}
	b, err := h.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t,
		`{"pomodoros":[{"start_time":"2016-06-14T12:00:00Z","description":"A description","duration":25,"tags":["a","b"]}]}`,
		string(b))
}

func TestHistory_MarshalText(t *testing.T) {
	p := &Pomodoro{
		StartTime:   time.Date(2016, 06, 14, 12, 0, 0, 0, time.UTC),
		Duration:    25 * time.Minute,
		Tags:        []string{"a", "b"},
		Description: "A description",
	}
	h := &History{Pomodoros: []*Pomodoro{p}}
	b, err := h.MarshalText()
	assert.Nil(t, err)
	assert.Equal(t,
		"2016-06-14T12:00:00Z description=\"A description\" duration=25 tags=a,b\n",
		string(b))
}

func Test_Latest(t *testing.T) {
	assert.Nil(t, empty.Latest())
	assert.Equal(t, b, one.Latest())
	assert.Equal(t, c, many.Latest())
}

func Test_Count(t *testing.T) {
	assert.Equal(t, 0, empty.Count())
	assert.Equal(t, 1, one.Count())
	assert.Equal(t, 3, many.Count())
}

func Test_Date(t *testing.T) {
	assert.Equal(t, &one, many.Date(b.StartTime))
}

func Test_Range(t *testing.T) {
	start := time.Date(2016, 06, 14, 0, 0, 0, 0, time.UTC)
	end := time.Date(2016, 06, 15, 0, 0, 0, 0, time.UTC)

	assert.Equal(t, 0, empty.Range(start, end).Count())
	assert.Equal(t, 1, one.Range(start, end).Count())
	assert.Equal(t, 1, many.Range(start, end).Count())
}

func Test_Update(t *testing.T) {
	history := &History{}

	assert.Equal(t, &empty, history)

	history.Update(b)

	assert.Equal(t, &one, history)

	history.Update(b)

	assert.Equal(t, &one, history)

	history.Update(a)
	history.Update(b)
	history.Update(c)

	assert.Equal(t, &many, history)

	bNew := &Pomodoro{StartTime: b.StartTime, Description: "updated"}
	history.Update(bNew)

	assert.Equal(t,
		&History{Pomodoros: []*Pomodoro{a, bNew, c}},
		history,
	)
}

func Test_Delete(t *testing.T) {
	history := &History{Pomodoros: []*Pomodoro{a, b, c}}

	history.Delete(c)
	history.Delete(a)

	expected := &History{Pomodoros: []*Pomodoro{b}}
	assert.Equal(t, expected, history)
}
