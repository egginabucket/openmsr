package gui

import (
	"fmt"

	"github.com/andlabs/ui"
	"github.com/egginabucket/openmsr/pkg/libmsr"
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
	enableCB *ui.Checkbox
	bpiCB    *ui.Combobox
	bpcCB    *ui.Combobox
	parityCB *ui.Combobox
	disabled bool
	//evenParity bool
	//data       []byte
	box *ui.Box
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

func (t *Track) setPreset(p Preset) {
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
	case PresetAAMVA:
		switch t.num {
		case 1:
			t.bpiCB.SetSelected(bpi210)
			t.bpcCB.SetSelected(7 - minBPC)
		case 2:
			t.bpiCB.SetSelected(bpi75)
			t.bpcCB.SetSelected(5 - minBPC)
		case 3:
			t.bpiCB.SetSelected(bpi210)
			t.bpcCB.SetSelected(7 - minBPC)
		}
		t.parityCB.SetSelected(parityOdd)
	}
	t.checkEdit()
}

func (t *Track) bpc() int {
	return minBPC + t.bpcCB.Selected()
}

func (t *Track) charOffset() byte {
	switch t.bpc() {
	case 4, 5:
		return '0'
	case 6, 7:
		return ' '
	}
	return ' '
}

func (t *Track) parityEven() bool {
	return t.parityCB.Selected() == parityEven
}

func (t *Track) decode(raw []byte) (chars []byte, parityOK []bool, lrcOK bool) {
	return libmsr.DecodeRaw(raw, t.charOffset(), t.bpc(), 8, t.parityEven())
}

func (t *Track) encode(chars []byte) []byte {
	return libmsr.EncodeRaw(chars, t.charOffset(), t.bpc(), 8, t.parityEven())
}

func (t *Track) checkEdit() {
	text := t.edit.Text()
	b := make([]byte, 0, len(text))
	min := rune(t.charOffset())
	max := min + (1 << rune(t.bpc()))
	changed := false
	for _, r := range text {
		if r < min || r > max {
			if max >= 'Z' && r >= 'a' && r <= 'z' {
				b = append(b, byte(r-('a'-'A')))
			}
			changed = true
		} else {
			b = append(b, byte(r))
		}
	}
	if changed {
		t.edit.SetText(string(b))
	}
}

func (t *Track) onEdit(e *ui.Entry) {
	t.checkEdit()
}

func (t *Track) onBPCChange(cb *ui.Combobox) {
	t.checkEdit()
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

func newTrack(num int) *Track {
	if num < 1 || num > 3 {
		panic("invalid track num")
	}
	var t Track
	t.num = num
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	t.enableCB = ui.NewCheckbox(fmt.Sprintf("Track %d", t.num))
	t.enableCB.SetChecked(true)
	t.enableCB.OnToggled(t.onToggle)
	hbox.Append(t.enableCB, true)

	t.bpiCB = ui.NewCombobox()
	t.bpiCB.Append("210 bits per inch")
	t.bpiCB.Append("75 bits per inch")
	t.bpcCB = ui.NewCombobox()
	for bpc := minBPC; bpc <= 7; bpc++ {
		t.bpcCB.Append(fmt.Sprintf("%d bits per char", bpc))
	}
	t.bpcCB.OnSelected(t.onBPCChange)
	t.parityCB = ui.NewCombobox()
	t.parityCB.Append("Odd parity")
	t.parityCB.Append("Even parity")

	hbox.Append(t.bpiCB, false)
	hbox.Append(t.bpcCB, false)
	hbox.Append(t.parityCB, false)
	//t.bpiCB.OnSelected(t.onBPISelected)

	t.edit = ui.NewEntry()
	t.edit.OnChanged(t.onEdit)
	vbox.Append(hbox, false)
	vbox.Append(t.edit, false)
	t.box = vbox
	return &t
}
