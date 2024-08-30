package pages

import (
	"bytes"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"live_server_ui/config"
	"live_server_ui/settings"
)

func PushVideoPage() *container.TabItem {
	path := widget.NewLabel("")

	fileDialog := dialog.NewFileOpen(func(f fyne.URIReadCloser, err error) {
		if f != nil {
			pathUri := strings.Split(f.URI().String(), "//")[1]
			defer f.Close()

			fileByte, err := os.ReadFile(pathUri)

			if errDealing(err, &settings.LiveInfoWindow) {
				return
			}

			if strings.Split(http.DetectContentType(fileByte), "/")[0] != "video" {
				dialog.ShowInformation("Wrong File Type", "Wrong File Type", settings.LiveInfoWindow)
				return
			}

			payload := &bytes.Buffer{}
			writer := multipart.NewWriter(payload)
			file, errFile1 := os.Open(pathUri)
			defer file.Close()
			part1, errFile1 := writer.CreateFormFile("file", filepath.Base(pathUri))
			_, errFile1 = io.Copy(part1, file)
			if errDealing(errFile1, &settings.LiveInfoWindow) {
				return
			}

			err2 := writer.Close()
			if errDealing(err2, &settings.LiveInfoWindow) {
				return
			}

			client := &http.Client{}
			req, err3 := http.NewRequest("POST", config.Config.UploadVideoUrl, payload)
			//fmt.Println(config.Config.UploadVideoUrl)
			if errDealing(err3, &settings.LiveInfoWindow) {
				return
			}

			req.Header.Set("Content-Type", writer.FormDataContentType())
			//fmt.Println(req)
			res, err4 := client.Do(req)
			if errDealing(err4, &settings.LiveInfoWindow) {
				return
			}
			defer res.Body.Close()

			body, err5 := io.ReadAll(res.Body)
			if errDealing(err5, &settings.LiveInfoWindow) {
				return
			}

			if res.StatusCode != http.StatusOK {
				dialog.ShowInformation(strconv.Itoa(res.StatusCode), string(body), settings.LiveInfoWindow)
				return
			}

			path.SetText(config.Config.VideoUrl + strings.Split(string(body), "live_videos")[1])
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
		container.NewHBox(widget.NewLabel("Path:"), widget.NewButtonWithIcon("", theme.FolderOpenIcon(), fileDialog.Show), path),
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
			"type":       {"fmp4"},
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
