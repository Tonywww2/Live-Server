package pages

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"io/ioutil"
	"live_server_ui/settings"
	"net/http"
	"net/url"
	"strings"
)

func CreateLivePage() *container.TabItem {
	nameEntry := widget.NewEntry()
	posterLabel := widget.NewLabel("")

	file := dialog.NewFileOpen(func(f fyne.URIReadCloser, err error) {
		if f != nil {
			posterUri := strings.Split(f.URI().String(), "//")[1]
			posterLabel.SetText(posterUri)
			defer f.Close()

		}
	}, settings.MainWindow)

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
			response, err := http.PostForm(settings.CreateLiveURL, payload)
			settings.TreatError(err, response)

		}),
	)

	return container.NewTabItem("Create New Live", livePage)
}

func CreatGetAllPage() *container.TabItem {
	resultLabel := widget.NewEntry()
	resultLabel.MultiLine = true

	getButton := widget.NewButton("Get", func() {
		response, err := http.Get(settings.GetAllLiveURL)
		if err != nil || response.StatusCode != 200 {
			//fmt.Println("Error")
			dialog.ShowError(err, settings.MainWindow)
		}

		var result map[string]map[string]interface{}
		body, err := ioutil.ReadAll(response.Body)
		if err == nil {
			err = json.Unmarshal(body, &result)
		}

		resultText := "Result: \n"
		i := 0
		settings.CachedLives = make([]string, 0)
		for k, v := range result {
			text := "\"" + settings.ToString(v["Name"]) + "\", " +
				"\"" + settings.ToString(v["Poster"]) + "\", " +
				"\"" + settings.ToString(v["StartTime"]) + "\", " +
				"\"" + settings.ToString(v["RtmpAddr"]) + "\", " +
				settings.ToString(v["IsStreamed"])
			resultText += "\"" + k + "\"" + "\n" + text + "\n\n"

			settings.CachedLives = append(settings.CachedLives, settings.ToString(v["Name"])+"//"+k)
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
