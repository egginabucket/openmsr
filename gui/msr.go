package gui

import (
	"github.com/andlabs/ui"
	"github.com/egginabucket/openmsr/gui/tracks"
)

type MSR struct {
	Tracks [3]*tracks.Track
}

func (msr *MSR) setPreset(p tracks.Preset) {
	for _, t := range msr.Tracks {
		t.SetPreset(p)
	}
}

func (msr *MSR) selectPreset(cb *ui.Combobox) {
	msr.setPreset(tracks.Preset(cb.Selected()))
}

func MakeMainUI() ui.Control {
	var msr MSR
	hBox := ui.NewHorizontalBox()
	trackBox := ui.NewVerticalBox()
	for i := 0; i < 3; i++ {
		t := tracks.NewTrack(i + 1)
		trackBox.Append(t.Box, false)
		msr.Tracks[i] = t
	}
	hBox.Append(trackBox, false)
	sideForm := ui.NewForm()
	presetCB := ui.NewCombobox()
	presetCB.Append("ISO")
	presetCB.Append("AAMVA")
	presetCB.SetSelected(int(tracks.PresetISO))
	presetCB.OnSelected(msr.selectPreset)
	sideForm.SetPadded(true)
	sideForm.Append("Preset", presetCB, false)
	hBox.Append(sideForm, false)
	msr.setPreset(tracks.PresetISO)
	return hBox
}
