// Package sim provides simulated devices for mock mode and rehearsal:
// a switcher, a presentation and a mixer with honest little state
// machines and call logs. driver name "sim" for all three kinds.
package main

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/toolingdev/ptz-controller/sdk/addon"
)

// ---------------------------------------------------------------- switcher

type Switcher struct {
	mu    sync.Mutex
	st    addon.SwitcherState
	Calls []string
}

func NewSwitcher() *Switcher {
	return &Switcher{st: addon.SwitcherState{
		Version: "sim", Edition: "Simulated",
		Inputs: []addon.Input{
			{Number: 1, Key: "sim-1", Type: "Camera", Title: "PTZ Wide — Sanctuary", State: "Running"},
			{Number: 2, Key: "sim-2", Type: "Camera", Title: "PTZ Pulpit — Close", State: "Running"},
			{Number: 3, Key: "sim-3", Type: "NDI", Title: "Slides", State: "Running"},
			{Number: 4, Key: "sim-4", Type: "Video", Title: "Media / Announcements", State: "Paused"},
		},
		Preview: 2, Active: 1, Streaming: true,
	}}
}

func (s *Switcher) State(context.Context) (addon.SwitcherState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c := s.st
	c.Inputs = append([]addon.Input(nil), s.st.Inputs...)
	return c, nil
}

func (s *Switcher) Preview(_ context.Context, n int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, in := range s.st.Inputs {
		if in.Number == n {
			s.st.Preview = n
			s.Calls = append(s.Calls, fmt.Sprintf("preview %d", n))
			return nil
		}
	}
	return fmt.Errorf("switcher has no input %d", n)
}

func (s *Switcher) take(what string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.st.Active, s.st.Preview = s.st.Preview, s.st.Active
	s.Calls = append(s.Calls, what)
	return nil
}

func (s *Switcher) Cut(context.Context) error { return s.take("cut") }
func (s *Switcher) Fade(_ context.Context, ms int) error {
	return s.take(fmt.Sprintf("fade %d", ms))
}

// ------------------------------------------------------------ presentation

type Presentation struct {
	mu      sync.Mutex
	actions []string
	Sent    []string
}

func NewPresentation(actions []string) *Presentation {
	if len(actions) == 0 {
		actions = []string{"black", "clear", "logo", "next", "prev"}
	}
	sort.Strings(actions)
	return &Presentation{actions: actions}
}

func (p *Presentation) Info() addon.PresentationInfo {
	return addon.PresentationInfo{Available: true, Target: "simulated — actions are recorded, not sent"}
}

func (p *Presentation) Actions() []string { return append([]string(nil), p.actions...) }

func (p *Presentation) Do(_ context.Context, a string) (map[string]any, error) {
	if i := sort.SearchStrings(p.actions, a); i >= len(p.actions) || p.actions[i] != a {
		return nil, fmt.Errorf("%w %q", addon.ErrUnknownAction, a)
	}
	p.mu.Lock()
	p.Sent = append(p.Sent, a)
	p.mu.Unlock()
	return map[string]any{"action": a, "port": "simulated"}, nil
}

// ------------------------------------------------------------------- mixer

type Mixer struct {
	mu    sync.Mutex
	max   int
	muted map[int]bool
	// StuckChannels ignore sets (tests the read-back safety net).
	StuckChannels map[int]bool
	Calls         []string
}

func NewMixer(max int) *Mixer {
	if max <= 0 {
		max = 48
	}
	return &Mixer{max: max, muted: map[int]bool{}, StuckChannels: map[int]bool{}}
}

func (m *Mixer) Probe(context.Context) (string, error) { return "Simulated console", nil }
func (m *Mixer) MaxChannel() int                       { return m.max }

func (m *Mixer) check(ch int) error {
	if ch < 1 || ch > m.max {
		return fmt.Errorf("channel %d is outside 1..%d", ch, m.max)
	}
	return nil
}

func (m *Mixer) GetMute(_ context.Context, ch int) (bool, error) {
	if err := m.check(ch); err != nil {
		return false, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.muted[ch], nil
}

func (m *Mixer) SetMute(_ context.Context, ch int, muted bool) error {
	if err := m.check(ch); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Calls = append(m.Calls, fmt.Sprintf("mute %d %v", ch, muted))
	if !m.StuckChannels[ch] {
		m.muted[ch] = muted
	}
	return nil
}
