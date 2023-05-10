package gui

import (
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
	deviceCB       *ui.Combobox
	refreshButton  *ui.Button
	availableDevs  []*usb.DeviceInfo
	trackBox       *ui.Box
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
	progBar *ui.ProgressBar
	mu      sync.Mutex
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

func (a *App) lookforDevices() {
	//a.availableDevs = make([]*usb.DeviceInfo, 0)
	hids, err := usb.EnumerateHid(libmsr.VendorID, libmsr.ProductID)
	if err != nil {
		a.throwErr(err)
		return
	}
	raws, err := usb.EnumerateRaw(libmsr.VendorID, libmsr.ProductID)
	if err != nil {
		a.throwErr(err)
		return
	}
	newDevs := make([]*usb.DeviceInfo, 0, len(hids)+len(raws))
	for _, di := range hids {
		newDevs = append(newDevs, &di)
	}
	for _, di := range raws {
		newDevs = append(newDevs, &di)
	}
	a.availableDevs = append(a.availableDevs, newDevs...)
	for _, d := range newDevs {
		a.deviceCB.Append(d.Product)
	}
}

func (a *App) selectDevice(cb *ui.Combobox) {
	if a.device != nil {
		a.disconnect()
	}
	if cb.Selected() == 0 {
		return // already disconnected
	}
	d, err := a.availableDevs[cb.Selected()-1].Open()
	if err != nil {
		cb.SetSelected(0)
		a.throwErr(err)
		return
	}
	a.connect(d)
}

func (a *App) connect(d usb.Device) {
	a.mu.Lock()
	defer a.reset(nil)
	defer a.mu.Unlock()
	a.setDeviceAble(true)
	a.device = libmsr.NewDevice(d)
}

func (a *App) disconnect() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.deviceCB.SetSelected(connNone)
	a.setDeviceAble(false)
	if err := a.device.Close(); err != nil {
		a.throwErr(err)
	}
	a.device = nil
}

func (a *App) controls() []ui.Control {
	return append(a.deviceControls(),
		a.deviceCB,
		a.refreshButton,
		a.openButton,
		a.saveButton,
		a.trackBox,
	)
}

func (a *App) setFrozen(f bool) {
	fn := enableDisable(!f)
	for _, c := range a.controls() {
		fn(c)
	}
}

func (a *App) freeze() {
	a.mu.Lock()
	a.progBar.Show()
	a.setFrozen(true)
}

func (a *App) unfreeze() {
	a.setFrozen(false)
	a.progBar.Hide()
	a.mu.Unlock()
}

func (a *App) selectPreset(r *ui.RadioButtons) {
	a.mu.Lock()
	defer a.mu.Unlock()
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

func (a *App) setDefaults() {
	a.presetRadio.SetSelected(int(PresetISO))
	a.selectPreset(a.presetRadio)
	a.infoTypeCB.SetSelected(infoIEC7813)
}

func (a *App) reset(*ui.Button) {
	a.mu.Lock()
	defer a.setDefaults()
	defer a.mu.Unlock()
	if err := a.device.Reset(); err != nil {
		a.throwErr(err)
	}
	isHiCo, err := a.device.IsHiCo()
	if err != nil {
		a.throwErr(err)
	} else {
		if isHiCo {
			a.coRadio.SetSelected(hiCo)
		} else {
			a.coRadio.SetSelected(loCo)
		}
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
	var tracks [3][]byte
	for i, track := range a.tracks {
		if !track.disabled {
			tracks[i] = []byte(track.edit.Text())
		}
	}
	err := a.device.WriteISOTracks(tracks[0], tracks[1], tracks[2])
	if err != nil {
		a.throwErr(err)
	}
}

func (a *App) writeRaw() {
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
}

func (a *App) setCoercivity(r *ui.RadioButtons) {
	a.mu.Lock()
	defer a.mu.Unlock()
	var err error
	switch r.Selected() {
	case hiCo:
		err = a.device.SetHiCo()
	case loCo:
		err = a.device.SetLoCo()
	}
	if err != nil {
		a.throwErr(err)
	}
}

func (a *App) showInfo(*ui.Button) {
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

func (a *App) openFile(*ui.Button) {
	path := ui.OpenFile(a.win)
	if path == "" {
		return
	}
}

func (a *App) saveFile(b *ui.Button) {
	path := ui.SaveFile(a.win)
	if path == "" {
		return
	}
}

func (a *App) onClosing(*ui.Window) bool {
	if a.device != nil {
		a.device.Close()
	}
	return true
}

func MakeMainUI(win *ui.Window) ui.Control {
	var a App
	a.win = win
	a.win.OnClosing(a.onClosing)
	vBox := ui.NewVerticalBox()
	vBox.SetPadded(true)
	hBox := ui.NewHorizontalBox()
	hBox.SetPadded(true)
	a.trackBox = ui.NewVerticalBox()
	a.trackBox.SetPadded(true)
	for i := 0; i < 3; i++ {
		t := newTrack(i + 1)
		a.trackBox.Append(t.box, false)
		a.tracks[i] = t
	}

	sideMenu := ui.NewVerticalBox()
	sideMenu.SetPadded(true)
	//sideForm := ui.NewForm()
	a.deviceCB = ui.NewCombobox()
	a.deviceCB.Append("Disconnected")
	a.deviceCB.OnSelected(a.selectDevice)
	a.refreshButton = ui.NewButton("Refresh")
	a.refreshButton.OnClicked(func(*ui.Button) { a.lookforDevices() })
	a.coRadio = ui.NewRadioButtons()
	a.coRadio.Append("Hi-Co")
	a.coRadio.Append("Lo-Co")

	a.coRadio.OnSelected(a.setCoercivity)
	a.presetRadio = ui.NewRadioButtons()
	a.presetRadio.Append("ISO")
	a.presetRadio.Append("AAMVA")
	a.presetRadio.OnSelected(a.selectPreset)
	a.infoTypeCB = ui.NewCombobox()
	a.infoTypeCB.Append("ISO/IEC 7813")
	a.infoTypeCB.Append("AAMVA DL")
	a.showInfoButton = ui.NewButton("Show info")
	a.showInfoButton.OnClicked(a.showInfo)
	sideMenu.Append(a.deviceCB, false)
	sideMenu.Append(a.refreshButton, false)
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
	a.progBar = ui.NewProgressBar()
	a.progBar.Hide()
	a.progBar.SetValue(-1)
	vBox.Append(a.progBar, false)
	vBox.Append(buttonBox, false)

	hBox.Append(sideMenu, false)
	hBox.Append(a.trackBox, true)

	a.setDeviceAble(false)
	a.setDefaults()
	a.deviceCB.SetSelected(connNone)
	a.lookforDevices()
	if len(a.availableDevs) > 0 {
		d, err := a.availableDevs[0].Open()
		if err == nil {
			a.device = libmsr.NewDevice(d)
			a.deviceCB.SetSelected(1)
			a.connect(d)
		}
	}
	return vBox
}
