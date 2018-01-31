package openpomodoro

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"sort"
)

// Client holds the location of the directory and files.
type Client struct {
	Directory    string
	CurrentFile  string
	HistoryFile  string
	SettingsFile string
}

type State struct {
	Pomodoro *Pomodoro
	History  *History
	Settings *Settings
}

const (
	// FilePerm are the permissions set when creating files.
	FilePerm = 0644
)

// NewClient returns a new Client with the given directory. If the directory is
// an empty string, the default directory of ~/.pomodoro is used.
func NewClient(directory string) (*Client, error) {
	var d string
	var u *user.User
	var err error

	if directory == "" {
		u, err = user.Current()
		if err != nil {
			return nil, err
		}
		d = path.Join(u.HomeDir, ".pomodoro")
	} else {
		d, err = filepath.Abs(directory)
		if err != nil {
			return nil, err
		}
	}

	c := &Client{
		Directory:    d,
		CurrentFile:  path.Join(d, "current"),
		HistoryFile:  path.Join(d, "history"),
		SettingsFile: path.Join(d, "settings"),
	}

	return c, nil
}

// CurrentState returns a State with the current Pomodoro, history, and
// settings.
func (c *Client) CurrentState() (*State, error) {
	state := &State{}

	p, err := c.Pomodoro()
	if err != nil {
		return state, err
	}
	state.Pomodoro = p

	h, err := c.History()
	if err != nil {
		return state, err
	}
	state.History = h

	s, err := c.Settings()
	if err != nil {
		return state, err
	}
	state.Settings = s

	return state, nil
}

// History returns all Pomodoros from the `history` file.
func (c *Client) History() (*History, error) {
	ps := []*Pomodoro{}

	b, err := ioutil.ReadFile(c.HistoryFile)
	if err != nil {
		if os.IsNotExist(err) {
			return &History{Pomodoros: ps}, nil
		} else {
			return nil, err
		}
	}

	lines := bytes.Split(b, charNewline)

	for _, line := range lines {
		if bytesAllWhitespace(line) {
			continue
		}

		p := NewPomodoro()
		p.UnmarshalText(line)
		ps = append(ps, p)
	}

	return &History{Pomodoros: ps}, nil
}

// Pomodoro returns the current Pomodoro from the `current` file.
func (c *Client) Pomodoro() (*Pomodoro, error) {
	b, err := ioutil.ReadFile(c.CurrentFile)
	if err != nil {
		if os.IsNotExist(err) {
			return EmptyPomodoro(), nil
		}
		return nil, err
	}

	if len(b) == 0 {
		return EmptyPomodoro(), nil
	}

	p := NewPomodoro()
	p.UnmarshalText(b)

	return p, nil
}

// Settings returns the settings from the `settings` file.
func (c *Client) Settings() (*Settings, error) {
	s, err := c.readSettings()
	if err != nil {
		return nil, err
	}

	s.SetDefaults(&DefaultSettings)

	return s, nil
}

// Start starts a Pomodoro by writing the current timestamp along with
// configured defaults to the `current` file, and also records the Pomodoro in
// the `history` file.
func (c *Client) Start(p *Pomodoro) error {
	err := c.ensureDirectory()
	if err != nil {
		return err
	}

	current, err := c.Pomodoro()
	if err != nil {
		return err
	}

	if current.IsActive() {
		err = c.Cancel()
		if err != nil {
			return err
		}
	}

	p.StartTime = timeFunc()

	s, err := c.Settings()
	if err != nil {
		return err
	}

	p.ApplySettings(s)

	if err := c.writeCurrent(p); err != nil {
		return err
	}

	if err := c.appendHistory(p); err != nil {
		return err
	}

	return nil
}

// Finish ends the current Pomodoro by emptying the `current` file, and appending
// the `history` with the final duration.
func (c *Client) Finish() error {
	p, err := c.Pomodoro()
	if err != nil {
		return err
	}

	err = c.Clear()
	if err != nil {
		return err
	}

	p.Duration = timeFunc().Sub(p.StartTime)
	return c.updateHistory(p)
}

// Cancel cancels any current Pomodoro by emptying the `current` file, and
// removing the entry from the `history` file.
func (c *Client) Cancel() error {
	err := c.ensureDirectory()
	if err != nil {
		return err
	}

	p, err := c.Pomodoro()
	if err != nil {
		return err
	}

	if p.IsInactive() {
		return nil
	}

	err = c.writeCurrent(EmptyPomodoro())
	if err != nil {
		return err
	}

	return c.deleteHistory(p)
}

// Clear clears the current Pomodoro by emptying the `current` file.
func (c *Client) Clear() error {
	err := c.ensureDirectory()
	if err != nil {
		return err
	}

	return c.writeCurrent(EmptyPomodoro())
}

func (c *Client) ensureDirectory() error {
	return os.MkdirAll(c.Directory, 0755)
}

func (c *Client) writeCurrent(p *Pomodoro) error {
	var b []byte
	var err error

	if !p.IsInactive() {
		b, err = p.MarshalText()

		if err != nil {
			return err
		}
	}

	return ioutil.WriteFile(c.CurrentFile, b, FilePerm)
}

func (c *Client) appendHistory(p *Pomodoro) error {
	if p.IsInactive() {
		return nil
	}

	b, err := p.MarshalText()

	b = bytes.Replace(b, charNewline, charSpace, -1)

	f, err := os.OpenFile(c.HistoryFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, FilePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		return err
	}

	_, err = f.Write(charNewline)
	return err
}

func (c *Client) updateHistory(p *Pomodoro) error {
	history, err := c.History()
	if err != nil {
		return err
	}

	history.Update(p)

	return c.writeHistory(history)
}

func (c *Client) deleteHistory(p *Pomodoro) error {
	history, err := c.History()
	if err != nil {
		return err
	}

	history.Delete(p)

	return c.writeHistory(history)
}

func (c *Client) writeHistory(h *History) error {
	sort.Sort(h)

	b, err := h.MarshalText()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(c.HistoryFile, b, FilePerm)
}

func (c *Client) readSettings() (*Settings, error) {
	b, err := ioutil.ReadFile(c.SettingsFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	s := &Settings{}
	err = s.UnmarshalText(b)
	if err != nil {
		return nil, err
	}

	return s, nil
}
