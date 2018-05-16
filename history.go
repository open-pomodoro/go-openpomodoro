package openpomodoro

import (
	"bytes"
	"encoding/json"
	"sort"
	"time"
)

// History is a collection of Pomodoros.
type History struct {
	Pomodoros []*Pomodoro `json:"pomodoros"`
}

// sort.Interface
func (h History) Len() int           { return len(h.Pomodoros) }
func (h History) Swap(i, j int)      { h.Pomodoros[i], h.Pomodoros[j] = h.Pomodoros[j], h.Pomodoros[i] }
func (h History) Less(i, j int) bool { return h.Pomodoros[i].StartTime.Before(h.Pomodoros[j].StartTime) }

// MarshalJSON implements json.Marshaler.
func (h History) MarshalJSON() ([]byte, error) {
	// This is required so that json.Marshal ignores that we also implement
	// encoding.TextMarshaler via MarshalText.
	type alias History
	return json.Marshal((alias)(h))
}

// MarshalText implements encoding.TextMarshaler. It returns a byte slice of
// each Pomodoro in the History also marshaled, separated by a newline.
func (h History) MarshalText() ([]byte, error) {
	var bs [][]byte

	for _, p := range h.Pomodoros {
		b, err := p.MarshalText()
		if err != nil {
			return nil, err
		}

		bs = append(bs, b)
	}

	bs = append(bs, nil)

	return bytes.Join(bs, charNewline), nil
}

// Latest sorts the collection and then returns the latest Pomodoro.
func (h *History) Latest() *Pomodoro {
	sort.Sort(h)

	n := len(h.Pomodoros)
	if n == 0 {
		return nil
	}

	return h.Pomodoros[n-1]
}

// Count returns the total Pomodoro count.
func (h *History) Count() int {
	return len(h.Pomodoros)
}

// Date returns a new History collection for the given date.
func (h *History) Date(date time.Time) *History {
	y, m, d := date.Date()

	today := time.Date(y, m, d, 0, 0, 0, 0, date.Location())
	tomorrow := today.AddDate(0, 0, 1)

	return h.Range(today, tomorrow)
}

// Range returns a new History collection between the start and end times.
func (h *History) Range(start time.Time, end time.Time) *History {
	result := &History{}
	for _, pomodoro := range h.Pomodoros {
		if t := pomodoro.StartTime; t.Before(start) || t.After(end) {
			continue
		}
		result.Pomodoros = append(result.Pomodoros, pomodoro)
	}

	return result
}

// Update replaces a Pomodoro within a History collection in place. If the
// Pomodoro does not exist in the collection, it is appended and then the
// collection is sorted.
func (h *History) Update(p *Pomodoro) {
	for i, needle := range h.Pomodoros {
		if needle.Matches(p) {
			h.Pomodoros[i] = p
			return
		}
	}

	h.Pomodoros = append(h.Pomodoros, p)
	sort.Sort(h)
}

// Delete removes a Pomodoro from a History collection in place.
func (h *History) Delete(p *Pomodoro) {
	new := &History{}

	for _, needle := range h.Pomodoros {
		if !needle.Matches(p) {
			new.Pomodoros = append(new.Pomodoros, needle)
		}
	}

	h.Pomodoros = new.Pomodoros
}
