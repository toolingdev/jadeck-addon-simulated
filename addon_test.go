package main

import (
	"testing"

	"github.com/toolingdev/jadeck/sdk/addon"
	"github.com/toolingdev/jadeck/sdk/addon/addontest"
)

func TestOverThePipe(t *testing.T) {
	h := addontest.Start(t, spec())
	if caps := h.Open("sw", "sim-switcher", nil); len(caps) != 1 || caps[0] != addon.CapSwitcher {
		t.Fatalf("caps %v", caps)
	}
	var st addon.SwitcherState
	h.Call("sw", "switcher.state", nil, &st)
	if len(st.Inputs) == 0 {
		t.Fatal("no inputs")
	}
	h.Call("sw", "switcher.preview", map[string]int{"input": 3}, nil)
	h.Call("sw", "switcher.cut", nil, nil)
	h.Call("sw", "switcher.state", nil, &st)
	if st.Active != 3 {
		t.Fatalf("cut to 3, program is %d", st.Active)
	}
	h.Open("mx", "sim-mixer", map[string]any{"maxChannel": 32})
	h.Call("mx", "mixer.setMute", map[string]any{"channel": 5, "muted": true}, nil)
	var m struct{ Muted bool }
	h.Call("mx", "mixer.getMute", map[string]int{"channel": 5}, &m)
	if !m.Muted {
		t.Fatal("mute not kept")
	}
	if err := h.TryCall("sw", "mixer.probe", nil, nil); err == nil {
		t.Fatal("a switcher answered as a mixer")
	}
}
