package pages

import (
	"bytes"
	"encoding/json"
	"fmt"
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

			if errDealing(err) {
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
			if errDealing(errFile1) {
				return
			}

			err2 := writer.Close()
			if errDealing(err2) {
				return
			}

			client := &http.Client{}
			req, err3 := http.NewRequest("POST", config.Config.UploadUrl, payload)
			if errDealing(err3) {
				return
			}

			req.Header.Set("Content-Type", writer.FormDataContentType())
			res, err4 := client.Do(req)
			if errDealing(err4) {
				return
			}
			defer res.Body.Close()

			body, err5 := io.ReadAll(res.Body)
			if errDealing(err5) {
				return
			}

			if res.StatusCode != http.StatusOK {
				dialog.ShowInformation(strconv.Itoa(res.StatusCode), string(body), settings.NewLiveWindow)
				return
			}
			posterLabel.SetText(config.Config.ImgUrl + strings.Split(string(body), "live_posters")[1])

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

func errDealing(err error) bool {
	if err != nil {
		dialog.ShowError(err, settings.NewLiveWindow)
		return true
	}
	return false
}

func CreatGetAllPage() *container.TabItem {
	resultLabel := widget.NewEntry()
	resultLabel.MultiLine = true

	getButton := widget.NewButton("Get", func() {
		response, err := http.Get(config.Config.GetAllLiveURL)
		defer response.Body.Close()
		if err != nil || response.StatusCode != 200 {
			dialog.ShowError(err, settings.MainWindow)
			//panic(err)
			return
		}

		var result []map[string]interface{}
		body, err := io.ReadAll(response.Body)
		if err == nil {
			err = json.Unmarshal(body, &result)
		}

		fmt.Println(result)

		resultText := "Result: \n"
		i := 0
		settings.CachedLives = make([]string, 0)
		for _, v := range result {
			text := "\"" + settings.ToString(v["Name"]) + "\", " +
				"\"" + settings.ToString(v["Poster"]) + "\", " +
				"\"" + settings.ToString(v["StartTime"]) + "\", " +
				"\"" + settings.ToString(v["RtmpAddr"]) + "\", " +
				settings.ToString(v["IsStreamed"])
			resultText += "\"" + settings.ToString(v["StreamID"]) + "\"" + "\n" + text + "\n\n"

			settings.CachedLives = append(settings.CachedLives, settings.ToString(v["Name"])+"//"+settings.ToString(v["StreamID"]))
			i++
		}
		resultLabel.SetText(resultText)

	})

	resultScroll := container.NewScroll(resultLabel)
	resultScroll.SetMinSize(fyne.NewSize(720, 300))

	getAllPage := container.NewVBox(
		widget.NewLabel("Get All Lives"),
		getButton,
		resultScroll,
	)

	return container.NewTabItem("Get All Lives", getAllPage)

}
