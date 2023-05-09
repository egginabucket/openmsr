package gui

import (
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
)

var mainwin *ui.Window

func SetupUI() {
	mainwin = ui.NewWindow("OpenMSR", 640, 480, true)
	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	tab := ui.NewTab()
	mainwin.SetChild(tab)
	mainwin.SetMargined(true)
	tab.Append("MSR", MakeMainUI(mainwin))
	tab.SetMargined(0, true)

	mainwin.Show()
}
