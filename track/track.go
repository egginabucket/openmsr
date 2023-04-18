package track

import (
	"fmt"

	"github.com/andlabs/ui"
)

const (
	bpi210 = iota
	bpi75
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
	num         int
	bitsPerInch int
	//bitsPerChar int
	bpiCB      *ui.Combobox
	bpcCB      *ui.Combobox
	parityCB   *ui.Combobox
	disabled   bool
	evenParity bool
	//data       []byte
	Box *ui.Box
	//cb       *ui.Checkbox
	edit *ui.MultilineEntry
}

func (t *Track) onToggle(cb *ui.Checkbox) {
	t.disabled = !cb.Checked()
	if t.disabled {
		t.edit.Disable()
	} else {
		t.edit.Enable()
	}
}

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

func NewTrack(num int) *Track {
	var t Track
	t.num = num
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox := ui.NewHorizontalBox()
	//hbox.SetPadded(true)
	hbox.Append(ui.NewLabel(fmt.Sprintf("Track %d", t.num)), false)
	cb := ui.NewCheckbox("Enabled")
	cb.SetChecked(true)
	cb.OnToggled(t.onToggle)
	hbox.Append(cb, false)

	t.bpiCB = ui.NewCombobox()
	t.bpiCB.Append("210")
	t.bpiCB.Append("75")
	t.bpiCB.OnSelected(t.onBPISelected)
	if t.num == 2 {
		t.bpiCB.SetSelected(bpi75)
	} else {
		t.bpiCB.SetSelected(bpi210)
	}

	t.edit = ui.NewMultilineEntry()
	vbox.Append(t.edit, true)
	t.Box = vbox
	return &t
}
