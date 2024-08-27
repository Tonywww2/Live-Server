package pages

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/clipboard"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"live_server_ui/config"
	"live_server_ui/settings"
)

var (
	currentPage      = 1
	grid             *fyne.Container
	backButton       fyne.CanvasObject
	forwardButton    fyne.CanvasObject
	recordListButton fyne.CanvasObject

	showingLives           []map[string]interface{}
	showingNameLabels      [10]*widget.Label
	showingTimeLabels      [10]*widget.Label
	showingStreamingLabels [10]*widget.Label
	infoButtons            [10]*widget.Button

	searchEntry         *widget.Entry
	pageLabel           *widget.Label
	infoName            *widget.Label
	infoID              *widget.Label
	infoTime            *widget.Label
	infoRtmp            *widget.Label
	infoStreamed        *widget.Label
	infoCopyID          *widget.Button
	infoCopyRtmp        *widget.Button
	infoCheckRtmpStream *widget.Button
	infoImg             *canvas.Image

	imgSize = fyne.NewSize(300, 300)

	idx = 0
)

func CreateClientContainer() *fyne.Container {
	pageLabel = widget.NewLabel("Page 1")

	createNewLiveButton := widget.NewButtonWithIcon("Create", theme.ContentAddIcon(), func() {
		settings.NewLiveWindow.Show()
	})

	searchEntry = widget.NewEntry()
	searchEntry.SetPlaceHolder("Search For Name")

	searchButton := widget.NewButtonWithIcon("Search", theme.SearchIcon(), Search)

	backButton = widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		currentPage = max(1, currentPage-1)
		UpdateShowingLives()
	})

	forwardButton = widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		currentPage = min((len(settings.CachedLivesOriginal)/10)+1, currentPage+1)
		UpdateShowingLives()
	})

	recordListButton = widget.NewButtonWithIcon("Records", theme.HistoryIcon(), func() {

		response, err := http.Get(config.Config.GetRecordsUrl)
		if err != nil {
			dialog.ShowError(err, settings.MainWindow)
			return

		}
		defer response.Body.Close()
		body, er := io.ReadAll(response.Body)
		if er != nil {
			dialog.ShowError(er, settings.MainWindow)
			return

		}
		dialog.ShowInformation("Records", string(body), settings.MainWindow)

	})

	for i := 0; i < 10; i++ {
		showingNameLabels[i] = widget.NewLabel("")
		showingTimeLabels[i] = widget.NewLabel("")
		showingStreamingLabels[i] = widget.NewLabel("")
		infoButtons[i] = widget.NewButtonWithIcon("Detail", theme.InfoIcon(), getFunc(i))
	}

	grid = container.NewGridWithRows(13,
		container.NewGridWithColumns(3,
			container.NewHBox(layout.NewSpacer(), createNewLiveButton),
			searchEntry,
			container.NewHBox(searchButton, recordListButton),
		),
		container.NewGridWithColumns(4, widget.NewLabel("Name"), widget.NewLabel("Created Time"), widget.NewLabel("Is Streaming(ed)")),
		createCustomHBox(gI(), showingNameLabels[idx-1], showingTimeLabels[idx-1], showingStreamingLabels[idx-1], infoButtons[idx-1]),
		createCustomHBox(gI(), showingNameLabels[idx-1], showingTimeLabels[idx-1], showingStreamingLabels[idx-1], infoButtons[idx-1]),
		createCustomHBox(gI(), showingNameLabels[idx-1], showingTimeLabels[idx-1], showingStreamingLabels[idx-1], infoButtons[idx-1]),
		createCustomHBox(gI(), showingNameLabels[idx-1], showingTimeLabels[idx-1], showingStreamingLabels[idx-1], infoButtons[idx-1]),
		createCustomHBox(gI(), showingNameLabels[idx-1], showingTimeLabels[idx-1], showingStreamingLabels[idx-1], infoButtons[idx-1]),
		createCustomHBox(gI(), showingNameLabels[idx-1], showingTimeLabels[idx-1], showingStreamingLabels[idx-1], infoButtons[idx-1]),
		createCustomHBox(gI(), showingNameLabels[idx-1], showingTimeLabels[idx-1], showingStreamingLabels[idx-1], infoButtons[idx-1]),
		createCustomHBox(gI(), showingNameLabels[idx-1], showingTimeLabels[idx-1], showingStreamingLabels[idx-1], infoButtons[idx-1]),
		createCustomHBox(gI(), showingNameLabels[idx-1], showingTimeLabels[idx-1], showingStreamingLabels[idx-1], infoButtons[idx-1]),
		createCustomHBox(gI(), showingNameLabels[idx-1], showingTimeLabels[idx-1], showingStreamingLabels[idx-1], infoButtons[idx-1]),
		container.NewGridWithColumns(5, layout.NewSpacer(), backButton, container.NewHBox(layout.NewSpacer(), pageLabel, layout.NewSpacer()), forwardButton, layout.NewSpacer()),
	)
	Search()
	return grid
}

func Search() {
	params := url.Values{}
	params.Set("name", searchEntry.Text)
	parseURL, err := url.Parse(config.Config.FuzzySearchLiveURL)
	if err != nil {
		log.Println("err")
	}
	parseURL.RawQuery = params.Encode()
	response, err := http.Get(parseURL.String())
	if err != nil {
		dialog.ShowError(err, settings.MainWindow)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		dialog.ShowError(err, settings.MainWindow)
		return
	}

	var result []map[string]interface{}
	body, err := io.ReadAll(response.Body)
	if err == nil {
		err = json.Unmarshal(body, &result)
	}

	//fmt.Println(result)
	settings.CachedLivesOriginal = result

	UpdateShowingLives()

}

func UpdateShowingLives() {
	start := (currentPage - 1) * 10
	end := min(len(settings.CachedLivesOriginal), currentPage*10)
	showingLives = []map[string]interface{}{}
	showingLives = append(showingLives, settings.CachedLivesOriginal[start:end]...)

	for k, v := range showingLives {
		showingNameLabels[k].SetText(settings.ToString(v["Name"]))
		showingTimeLabels[k].SetText(settings.ToString(v["StartTime"]))
		showingStreamingLabels[k].SetText(settings.ToString(v["IsStreamed"]))

	}

	for i := len(showingLives); i < 10; i++ {
		showingNameLabels[i].SetText("")
		showingTimeLabels[i].SetText("")
	}
	pageLabel.SetText("Page " + strconv.Itoa(currentPage))

}

func createCustomHBox(i int, name *widget.Label, time *widget.Label, streaming *widget.Label, b *widget.Button) *fyne.Container {
	return container.NewGridWithColumns(4,
		container.NewHBox(widget.NewLabel(strconv.Itoa(i)), name),
		time,
		streaming,
		b)

}

func getFunc(i int) func() {
	return func() {
		if i < len(showingLives) {
			infoName.SetText(settings.ToString(showingLives[i]["Name"]))
			infoID.SetText(settings.ToString(showingLives[i]["StreamID"]))
			infoTime.SetText(settings.ToString(showingLives[i]["StartTime"]))
			infoRtmp.SetText(settings.ToString(showingLives[i]["RtmpAddr"]))
			infoStreamed.SetText(settings.ToString(showingLives[i]["IsStreamed"]))

			img := settings.ToString(showingLives[i]["Poster"])
			if img == "" {
				img = "icon.png"
			}
			infoImg.File = img

			settings.StreamIdEntry.SetText(settings.ToString(showingLives[i]["StreamID"]))

			infoCheckRtmpStream.SetIcon(theme.QuestionIcon())

		} else {
			infoName.SetText("")
			infoID.SetText("")
			infoTime.SetText("")
			infoRtmp.SetText("")
			infoStreamed.SetText("")

			infoImg.File = "icon.png"
			settings.StreamIdEntry.SetText("")
		}

		settings.LiveInfoWindow.Show()
	}

}

func gI() int {
	idx++
	return idx
}

func CreateLiveInfoContainer() *container.TabItem {
	infoName = widget.NewLabel("")
	infoID = widget.NewLabel("")
	infoTime = widget.NewLabel("")
	infoRtmp = widget.NewLabel("")
	infoStreamed = widget.NewLabel("")

	infoImg = canvas.NewImageFromFile("icon.png")
	infoImg.SetMinSize(imgSize)
	infoImg.FillMode = canvas.ImageFillContain

	infoCopyID = widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		clipboard.Write(clipboard.FmtText, []byte(infoID.Text))
	})

	infoCopyRtmp = widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		clipboard.Write(clipboard.FmtText, []byte(infoRtmp.Text))
	})

	infoCheckRtmpStream = widget.NewButtonWithIcon("Streaming", theme.QuestionIcon(), func() {
		response, err := http.Get(config.Config.RtmpListUrl)
		if err != nil {
			settings.TreatError(err, response)
			return
		}
		defer response.Body.Close()

		var result []map[string]interface{}
		body, err := io.ReadAll(response.Body)
		if err == nil && body != nil {
			err = json.Unmarshal(body, &result)
		}

		for _, v := range result {
			id, er := v["Path"]
			if er && id == infoID.Text {
				infoCheckRtmpStream.SetIcon(theme.ConfirmIcon())
				return
			}
		}
		infoCheckRtmpStream.SetIcon(theme.CancelIcon())

	})

	return container.NewTabItem("Info", container.NewVBox(infoImg, container.New(layout.NewFormLayout(),
		widget.NewLabel("Name: "),
		infoName,

		widget.NewLabel("Stream ID: "),
		container.NewHBox(infoID, layout.NewSpacer(), infoCopyID),

		widget.NewLabel("Started Time: "),
		infoTime,

		widget.NewLabel("Rtmp Address: "),
		container.NewHBox(infoRtmp, layout.NewSpacer(), infoCheckRtmpStream, infoCopyRtmp),

		widget.NewLabel("Is Streamed: "),
		infoStreamed,

		//container.NewGridWithColumns(3, infoStartStream, infoPushToRtmp, infoStopStream),
	)))

}
