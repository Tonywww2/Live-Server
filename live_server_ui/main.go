package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"live_server_ui/pages"
	"live_server_ui/settings"
	"time"
)

func main() {
	a := app.NewWithID("live_server_ui")
	settings.MainWindow = a.NewWindow("Main")
	settings.MainWindow.Resize(fyne.NewSize(720, 480))
	settings.MainWindow.SetMaster()

	settings.StreamIdEntry = widget.NewSelectEntry(settings.CachedLives)

	tab := container.NewAppTabs(
		pages.CreateLivePage(),
		pages.CreatGetAllPage(),
		pages.PushVideoPage(),
		pages.PushRtmpPage(),
		pages.EndStreamPage(),
	)
	settings.MainWindow.SetContent(tab)

	go func() {
		for range time.Tick(time.Second) {
			//fmt.Println(entry1.Text)
			settings.StreamIdEntry.SetOptions(settings.CachedLives)
		}
	}()

	settings.MainWindow.Show()
	a.Run()
}
