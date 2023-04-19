package tracks

import (
	"fmt"

	"github.com/andlabs/ui"
)

const minBPC = 4

const (
	bpi210 = iota
	bpi75
)

const (
	parityOdd = iota
	parityEven
)

type Preset int

const (
	PresetISO Preset = iota
	PresetAAMVA
)

/*
const (
	bpc4 = iota
	bpc5
	bpc6
	bpc7
)
*/

type Track struct {
	num int
	//bitsPerInch int
	//bitsPerChar int
	bpiCB    *ui.Combobox
	bpcCB    *ui.Combobox
	parityCB *ui.Combobox
	disabled bool
	//evenParity bool
	//data       []byte
	Box *ui.Box
	//cb       *ui.Checkbox
	edit *ui.Entry
}

func (t *Track) controls() []ui.Control {
	return []ui.Control{
		t.edit,
		t.bpiCB,
		t.bpcCB,
		t.parityCB,
	}
}

func (t *Track) IsDisabled() bool {
	return t.disabled
}

func (t *Track) onToggle(cb *ui.Checkbox) {
	t.disabled = !cb.Checked()
	for _, c := range t.controls() {
		if t.disabled {
			c.Disable()
		} else {
			c.Enable()
		}
	}
}

func (t *Track) SetPreset(p Preset) {
	switch p {
	case PresetISO:
		switch t.num {
		case 1:
			t.bpiCB.SetSelected(bpi210)
			t.bpcCB.SetSelected(7 - minBPC)
		case 2:
			t.bpiCB.SetSelected(bpi75)
			t.bpcCB.SetSelected(5 - minBPC)
		case 3:
			t.bpiCB.SetSelected(bpi210)
			t.bpcCB.SetSelected(5 - minBPC)
		}
		t.parityCB.SetSelected(parityOdd)
	}
}

func (t *Track) Decode(d []byte) string {
	return ""
}

func (t *Track) Encode(s string) []byte {
	return nil
}

/*
func (t *Track) onBPISelected(cb *ui.Combobox) {
	switch cb.Selected() {
	case -1:
		cb.SetSelected(bpi210)
	case bpi210:
		t.bitsPerInch = 210
	case bpi75:
		t.bitsPerInch = 75
	}
}

func (t *Track) onBPCSelected(cb *ui.Combobox) {
	switch cb.Selected() {
	case 0:
	}
}
*/

func NewTrack(num int) *Track {
	if num < 1 || num > 3 {
		panic("invalid track num")
	}
	var t Track
	t.num = num
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox := ui.NewHorizontalBox()
	//hbox.SetPadded(true)

	cb := ui.NewCheckbox(fmt.Sprintf("Track %d", t.num))
	cb.SetChecked(true)
	cb.OnToggled(t.onToggle)
	hbox.Append(cb, true)

	t.bpiCB = ui.NewCombobox()
	t.bpiCB.Append("210 bits per inch")
	t.bpiCB.Append("75 bits per inch")
	t.bpcCB = ui.NewCombobox()
	for bpc := minBPC; bpc <= 7; bpc++ {
		t.bpcCB.Append(fmt.Sprintf("%d bits per char", bpc))
	}
	t.parityCB = ui.NewCombobox()
	t.parityCB.Append("Odd parity")
	t.parityCB.Append("Even parity")

	hbox.Append(t.bpiCB, false)
	hbox.Append(t.bpcCB, false)
	hbox.Append(t.parityCB, false)
	//t.bpiCB.OnSelected(t.onBPISelected)

	t.edit = ui.NewEntry()
	vbox.Append(t.edit, false)
	vbox.Append(hbox, false)
	t.Box = vbox
	return &t
}
