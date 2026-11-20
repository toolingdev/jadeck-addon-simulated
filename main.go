// Simulated devices: a generic jadeck add-on. A practice
// switcher, slides and mixer that touch no real gear.
package main

import (
	"encoding/json"

	"github.com/toolingdev/jadeck/sdk/addon"
)

// Version is set by the release build.
var Version = "1.0.0"

func main() { addon.Serve(spec()) }

func spec() addon.Spec {
	sw := map[string]any{"/api/switcher/status": map[string]any{"label": "Simulated", "connected": true,
		"inputs": []map[string]any{{"number": 1, "title": "PTZ Wide — Sanctuary"}, {"number": 2, "title": "PTZ Pulpit — Close"}}, "preview": 2, "active": 1}}
	return addon.Spec{
		Manifest: addon.Manifest{
			ID: "simulated", Name: "Simulated devices", Version: Version, Publisher: "jadeck", APIVersion: addon.APIVersion, Repo: "toolingdev/jadeck-addon-simulated", Generic: true,
			Description: "A practice switcher, slides and mixer that touch no real gear.",
			Reaches:     "Nothing — everything stays inside the controller",
			ConnectionTypes: []addon.ConnectionType{
				{ID: "sim-switcher", Label: "Simulated switcher", Capabilities: []string{addon.CapSwitcher}, Fields: []addon.Field{},
					Description: "Behaves like vMix with four inputs.", SampleTest: "Connected — 4 inputs, program 1, preview 2"},
				{ID: "sim-slides", Label: "Simulated slides", Capabilities: []string{addon.CapSlides}, Fields: []addon.Field{},
					Description: "Records the slide actions it gets.", SampleTest: "Ready — simulated"},
				{ID: "sim-mixer", Label: "Simulated mixer", Capabilities: []string{addon.CapMixer}, Fields: []addon.Field{},
					Description: "A 48-channel console that remembers mutes.", SampleTest: "Console answered — simulated"},
			},
			Widgets: []addon.WidgetDef{{Kind: "core", Type: "switcher.cutfade", Label: "Cut / Fade", Sample: sw}},
		},
		Open: map[string]addon.OpenFunc{
			"sim-switcher": func(json.RawMessage) (any, error) { return NewSwitcher(), nil },
			"sim-slides": func(raw json.RawMessage) (any, error) {
				var c struct {
					Actions []string `json:"actions"`
				}
				json.Unmarshal(raw, &c)
				return NewPresentation(c.Actions), nil
			},
			"sim-mixer": func(raw json.RawMessage) (any, error) {
				var c struct {
					MaxChannel int `json:"maxChannel"`
				}
				json.Unmarshal(raw, &c)
				return NewMixer(c.MaxChannel), nil
			},
		},
	}
}
