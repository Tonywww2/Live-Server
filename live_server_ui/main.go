package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"time"

	"live_server_ui/config"
	"live_server_ui/pages"
	"live_server_ui/settings"
)

func main() {
	config.LoadConfig()

	a := app.NewWithID("live_server_ui")
	settings.MainWindow = a.NewWindow("Main")
	settings.MainWindow.Resize(fyne.NewSize(720, 540))
	settings.MainWindow.SetMaster()
	settings.MainWindow.CenterOnScreen()

	settings.NewLiveWindow = a.NewWindow("Create")
	settings.NewLiveWindow.Resize(fyne.NewSize(320, 320))
	settings.NewLiveWindow.SetCloseIntercept(func() {
		settings.NewLiveWindow.Hide()
	})

	settings.LiveInfoWindow = a.NewWindow("Info")
	settings.LiveInfoWindow.Resize(fyne.NewSize(520, 320))
	settings.LiveInfoWindow.SetCloseIntercept(func() {
		settings.LiveInfoWindow.Hide()
	})

	settings.StreamIdEntry = widget.NewSelectEntry(settings.CachedLives)

	//tab := container.NewAppTabs(
	//	pages.CreateLivePage(),
	//	pages.CreatGetAllPage(),
	//	pages.PushVideoPage(),
	//	pages.PushRtmpPage(),
	//	pages.EndStreamPage(),
	//)
	//settings.MainWindow.SetContent(tab)

	settings.MainWindow.SetContent(pages.CreateClientContainer())

	settings.NewLiveWindow.SetContent(container.NewAppTabs(pages.CreateLivePage()))

	settings.LiveInfoWindow.SetContent(container.NewAppTabs(
		pages.CreateLiveInfoContainer(),
		pages.PushVideoPage(),
		pages.PushRtmpPage(),
		pages.EndStreamPage(),
	))

	go func() {
		for range time.Tick(time.Second) {
			//fmt.Println(entry1.Text)
			settings.StreamIdEntry.SetOptions(settings.CachedLives)
		}
	}()

	settings.MainWindow.Show()
	a.Run()
}
