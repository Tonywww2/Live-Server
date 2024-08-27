package pages

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"net/http"
	"net/url"
	"os"
	"strings"

	"live_server_ui/config"
	"live_server_ui/settings"
)

func PushVideoPage() *container.TabItem {
	path := widget.NewLabel("")

	file := dialog.NewFileOpen(func(f fyne.URIReadCloser, err error) {
		if f != nil {
			pathUri := strings.Split(f.URI().String(), "//")[1]
			defer f.Close()

			file, er := os.ReadFile(pathUri)

			if er != nil {
				dialog.ShowError(er, settings.LiveInfoWindow)
				return
			}

			if strings.Split(http.DetectContentType(file), "/")[0] != "video" {
				dialog.ShowInformation("Wrong File Type", "Wrong File Type", settings.LiveInfoWindow)
				return
			}
			path.SetText(pathUri)
		}
	}, settings.LiveInfoWindow)

	pushButton := widget.NewButton("Push", func() {

		temp := strings.Split(settings.StreamIdEntry.Text, "//")
		text := ""
		if len(temp) == 2 {
			text = temp[1]
		} else {
			text = temp[0]
		}

		payload := url.Values{
			"streamID": {text},
			"path":     {path.Text},
		}
		response, err := http.PostForm(config.Config.ToStreamURL, payload)
		settings.TreatError(err, response)
		//fmt.Println(response)
		Search()
	})

	getAllPage := container.NewVBox(
		widget.NewLabel("Push Video"),
		widget.NewLabel("Stream ID: "),
		settings.StreamIdEntry,
		container.NewHBox(widget.NewLabel("Path:"), widget.NewButtonWithIcon("", theme.FolderOpenIcon(), file.Show), path),
		pushButton,
	)

	return container.NewTabItem("Start Streaming by Video", getAllPage)
}

func PushRtmpPage() *container.TabItem {
	rtmp := widget.NewEntry()

	pushButton := widget.NewButton("Push", func() {
		temp := strings.Split(settings.StreamIdEntry.Text, "//")
		text := ""
		if len(temp) == 2 {
			text = temp[1]
		} else {
			text = temp[0]
		}

		payload := url.Values{
			"stream_id": {text},
			"rtmp_addr": {rtmp.Text},
		}
		response, err := http.PostForm(config.Config.ToRtmpURL, payload)
		settings.TreatError(err, response)
		//fmt.Println(response)
	})

	pushRtmp := container.NewVBox(
		widget.NewLabel("Stream to Rtmp"),
		widget.NewLabel("Stream ID: "),
		settings.StreamIdEntry,
		widget.NewLabel("Path:"),
		rtmp,
		pushButton,
	)

	return container.NewTabItem("Push to Rtmp", pushRtmp)
}

func EndStreamPage() *container.TabItem {
	endButton := widget.NewButton("End", func() {
		temp := strings.Split(settings.StreamIdEntry.Text, "//")
		text := ""
		if len(temp) == 2 {
			text = temp[1]
		} else {
			text = temp[0]
		}

		payload := url.Values{
			"streamPath": {text},
			"type":       {"flv"},
		}
		response, err := http.PostForm(config.Config.EndStreamUrl, payload)
		settings.TreatError(err, response)
		//fmt.Println(response)
	})

	endStream := container.NewVBox(
		widget.NewLabel("End Stream"),
		widget.NewLabel("Stream ID: "),
		settings.StreamIdEntry,
		endButton,
	)

	return container.NewTabItem("End Streaming", endStream)
}
