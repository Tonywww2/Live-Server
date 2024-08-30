package pages

import (
	"bytes"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"io"
	"live_server_ui/config"
	"live_server_ui/settings"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func CreateLivePage() *container.TabItem {
	nameEntry := widget.NewEntry()
	posterLabel := widget.NewLabel("")

	file := dialog.NewFileOpen(func(f fyne.URIReadCloser, err error) {
		if f != nil {
			defer f.Close()
			posterUri := strings.Split(f.URI().String(), "//")[1]
			fileByte, err := os.ReadFile(posterUri)

			if errDealing(err, &settings.NewLiveWindow) {
				return
			}

			if strings.Split(http.DetectContentType(fileByte), "/")[0] != "image" {
				dialog.ShowInformation("Wrong File Type", "Wrong File Type", settings.NewLiveWindow)
				return
			}

			payload := &bytes.Buffer{}
			writer := multipart.NewWriter(payload)
			file, errFile1 := os.Open(posterUri)
			defer file.Close()
			part1, errFile1 := writer.CreateFormFile("file", filepath.Base(posterUri))
			_, errFile1 = io.Copy(part1, file)
			if errDealing(errFile1, &settings.NewLiveWindow) {
				return
			}

			err2 := writer.Close()
			if errDealing(err2, &settings.NewLiveWindow) {
				return
			}

			client := &http.Client{}
			req, err3 := http.NewRequest("POST", config.Config.UploadUrl, payload)
			if errDealing(err3, &settings.NewLiveWindow) {
				return
			}

			req.Header.Set("Content-Type", writer.FormDataContentType())
			res, err4 := client.Do(req)
			if errDealing(err4, &settings.NewLiveWindow) {
				return
			}
			defer res.Body.Close()

			body, err5 := io.ReadAll(res.Body)
			if errDealing(err5, &settings.NewLiveWindow) {
				return
			}

			if res.StatusCode != http.StatusOK {
				dialog.ShowInformation(strconv.Itoa(res.StatusCode), string(body), settings.NewLiveWindow)
				return
			}
			posterLabel.SetText(config.Config.PosterUrl + strings.Split(string(body), "live_posters")[1])

		}
	}, settings.NewLiveWindow)

	posterButton := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), file.Show)
	posterSelectBox := container.NewHBox(widget.NewLabel("Poster: "), posterButton, posterLabel)

	livePage := container.NewVBox(
		widget.NewLabel("Name: "),
		nameEntry,
		posterSelectBox,
		widget.NewButton("Create", func() {
			payload := url.Values{
				"name":   {nameEntry.Text},
				"poster": {posterLabel.Text},
			}
			response, err := http.PostForm(config.Config.CreateLiveURL, payload)
			settings.TreatError(err, response)
			Search()
		}),
	)

	return container.NewTabItem("Create New Live", livePage)
}

func errDealing(err error, window *fyne.Window) bool {
	if err != nil {
		dialog.ShowError(err, *window)
		return true
	}
	return false
}
