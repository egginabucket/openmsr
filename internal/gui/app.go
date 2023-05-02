package gui

import (
	"io"
	"sync"

	"github.com/andlabs/ui"
	"github.com/egginabucket/openmsr/pkg/libmsr"
	"github.com/egginabucket/openmsr/pkg/libtracks"
	"github.com/karalabe/usb"
)

const (
	hiCo = iota
	loCo
)

const (
	infoIEC7813 = iota
	infoAAMVA
)

const (
	connNone = iota
	connHID
	connRaw
)

type App struct {
	device         *libmsr.Device
	win            *ui.Window
	connCB         *ui.RadioButtons
	presetRadio    *ui.RadioButtons
	coRadio        *ui.RadioButtons
	infoTypeCB     *ui.Combobox
	showInfoButton *ui.Button
	tracks         [3]*Track
	resetButton,
	readButton,
	writeButton,
	eraseButton,
	openButton,
	saveButton *ui.Button
	mu sync.Mutex
}

func enableDisable(m bool) func(ui.Control) {
	if m {
		return func(c ui.Control) { c.Enable() }
	} else {
		return func(c ui.Control) { c.Disable() }
	}
}

func (a *App) throwErr(err error) {
	if err == usb.ErrDeviceClosed {

	} else {
		ui.MsgBoxError(a.win, "Error", err.Error())
	}
}

func (a *App) selectPreset(r *ui.RadioButtons) {
	preset := Preset(r.Selected())
	switch preset {
	case PresetAAMVA:
		a.infoTypeCB.SetSelected(infoAAMVA)
	case PresetISO:
		a.infoTypeCB.SetSelected(infoIEC7813)
	}
	for _, t := range a.tracks {
		t.setPreset(preset)
	}
}

func (a *App) deviceControls() []ui.Control {
	return []ui.Control{
		a.presetRadio,
		a.coRadio,
		a.resetButton,
		a.readButton,
		a.writeButton,
		a.eraseButton,
	}
}

func (a *App) setDeviceAble(m bool) {
	fn := enableDisable(m)
	for _, c := range a.deviceControls() {
		fn(c)
	}
}

func (a *App) controls() []ui.Control {
	return append(a.deviceControls(), a.openButton, a.saveButton)
}

func (a *App) setFrozen(f bool) {
	fn := enableDisable(!f)
	for _, c := range a.controls() {
		fn(c)
	}
	for _, t := range a.tracks {
		fn(t.enableCB)
		if !t.disabled {
			for _, c := range t.controls() {
				fn(c)
			}
		}
	}
}

func (a *App) freeze() {
	a.mu.Lock()
	a.setFrozen(true)
}

func (a *App) unfreeze() {
	a.setFrozen(false)
	a.mu.Unlock()
}

func (a *App) reset(b *ui.Button) {
	a.mu.Lock()
	defer a.mu.Unlock()
	err := a.device.Reset()
	if err != nil {
		a.throwErr(err)
	}
}

func (a *App) read() {
	a.freeze()
	defer a.unfreeze()
	rawTracks, err := a.device.ReadRawTracks()
	if err != nil {
		a.throwErr(err)
		return
	}
	for i, t := range a.tracks {
		chars, _, _ := t.decode(rawTracks[i])
		t.edit.SetText(string(chars))
	}
}

func (a *App) write() {
	a.freeze()
	defer a.unfreeze()
	var rawTracks [3][]byte
	for i, track := range a.tracks {
		track.edit.SetReadOnly(true)
		if !track.disabled {
			chars := []byte(track.edit.Text())
			rawTracks[i] = track.encode(chars)
		}
	}
	err := a.device.WriteRawTracks(rawTracks[0], rawTracks[1], rawTracks[2])
	if err != nil {
		a.throwErr(err)
	}
}

func (a *App) erase() {
	a.freeze()
	defer a.unfreeze()
	var tracks [3]bool
	for i, t := range a.tracks {
		if !t.disabled {
			tracks[i] = true
			t.edit.SetText("")
		}
	}
	err := a.device.Erase(tracks[0], tracks[1], tracks[2])
	if err != nil {
		a.throwErr(err)
	}
	for i, t := range tracks {
		if t {
			a.tracks[i].edit.SetText("")
		}
	}
}

func (a *App) setCoercivity(r *ui.RadioButtons) {
	a.mu.Lock()
	defer a.mu.Unlock()
	switch r.Selected() {
	case hiCo:
		a.device.SetHiCo()
	case loCo:
		a.device.SetLoCo()
	}
}

func (a *App) showInfo(b *ui.Button) {
	var in libtracks.Informer
	var err error
	switch a.infoTypeCB.Selected() {
	case infoIEC7813:
		in, err = libtracks.NewIEC7813Tracks(a.tracks[0].edit.Text(), a.tracks[1].edit.Text())
	case infoAAMVA:
		in, err = libtracks.NewAAMVATracks(a.tracks[0].edit.Text(), a.tracks[1].edit.Text(), a.tracks[2].edit.Text())
	}
	if err != nil {
		a.throwErr(err)
		return
	}
	win := ui.NewWindow("Card Info", 640, 480, false)
	win.SetMargined(true)
	hBox := ui.NewHorizontalBox()
	vBox := ui.NewVerticalBox()
	entry := ui.NewMultilineEntry()
	entry.SetReadOnly(true)
	entry.SetText(in.Info().String())
	hBox.Append(entry, true)
	vBox.Append(hBox, true)
	win.SetChild(vBox)
	win.OnClosing(func(*ui.Window) bool { return true })
	win.Show()
}

func (a *App) openFile(b *ui.Button) {
	if path := ui.OpenFile(a.win); path != "" {
		a.throwErr(io.EOF)
	}
}

func (a *App) saveFile(b *ui.Button) {
	if path := ui.SaveFile(a.win); path != "" {
		a.throwErr(io.EOF)
	}
}

func MakeMainUI(win *ui.Window) ui.Control {
	var a App
	a.win = win
	hids, err := usb.EnumerateHid(libmsr.VendorID, libmsr.ProductID)
	if err != nil {
		panic(err)
	}
	hid, err := hids[0].Open()
	if err != nil {
		panic(err)
	}
	a.device = libmsr.NewDevice(hid)
	vBox := ui.NewVerticalBox()
	hBox := ui.NewHorizontalBox()
	hBox.SetPadded(true)
	trackBox := ui.NewVerticalBox()
	trackBox.SetPadded(true)
	for i := 0; i < 3; i++ {
		t := newTrack(i + 1)
		trackBox.Append(t.box, false)
		a.tracks[i] = t
	}

	sideMenu := ui.NewVerticalBox()
	sideMenu.SetPadded(true)
	//sideForm := ui.NewForm()
	connectionCB := ui.NewCombobox()
	connectionCB.Append("MSR605X (HID)")
	connectionCB.Append("Disconnected")
	connectionCB.Append("MSR605 (Raw)")
	connectionCB.SetSelected(connHID)
	a.coRadio = ui.NewRadioButtons()
	a.coRadio.Append("Hi-Co")
	a.coRadio.Append("Lo-Co")
	isHiCo, err := a.device.IsHiCo()
	if err != nil {
		panic(err)
	}
	if isHiCo {
		a.coRadio.SetSelected(hiCo)
	} else {
		a.coRadio.SetSelected(loCo)
	}

	a.coRadio.OnSelected(a.setCoercivity)
	a.presetRadio = ui.NewRadioButtons()
	a.presetRadio.Append("ISO")
	a.presetRadio.Append("AAMVA")
	a.presetRadio.SetSelected(int(PresetISO))
	a.presetRadio.OnSelected(a.selectPreset)
	a.infoTypeCB = ui.NewCombobox()
	a.infoTypeCB.Append("ISO/IEC 7813")
	a.infoTypeCB.Append("AAMVA DL")
	a.infoTypeCB.SetSelected(infoIEC7813)
	a.showInfoButton = ui.NewButton("Show info")
	a.showInfoButton.OnClicked(a.showInfo)

	sideMenu.Append(a.coRadio, false)
	sideMenu.Append(a.presetRadio, false)
	sideMenu.Append(a.infoTypeCB, false)
	sideMenu.Append(a.showInfoButton, false)

	vBox.Append(hBox, true)

	buttonBox := ui.NewHorizontalBox()
	buttonBox.SetPadded(true)

	a.readButton = ui.NewButton("Read")
	a.readButton.OnClicked(func(*ui.Button) { go a.read() })
	a.writeButton = ui.NewButton("Write")
	a.writeButton.OnClicked(func(*ui.Button) { go a.write() })
	a.eraseButton = ui.NewButton("Erase")
	a.eraseButton.OnClicked(func(*ui.Button) { go a.erase() })
	a.openButton = ui.NewButton("Open file")
	a.openButton.OnClicked(a.openFile)
	a.saveButton = ui.NewButton("Save file")
	a.saveButton.OnClicked(a.saveFile)
	a.resetButton = ui.NewButton("Reset")
	a.resetButton.OnClicked(a.reset)

	buttonBox.Append(a.readButton, true)
	buttonBox.Append(a.writeButton, true)
	buttonBox.Append(a.eraseButton, true)
	buttonBox.Append(a.resetButton, true)
	buttonBox.Append(a.openButton, false)
	buttonBox.Append(a.saveButton, false)
	vBox.Append(buttonBox, false)

	hBox.Append(sideMenu, false)
	hBox.Append(trackBox, true)

	a.selectPreset(a.presetRadio)
	//a.setFrozen(true)
	return vBox
}
