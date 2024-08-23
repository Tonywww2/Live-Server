package pages

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
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
	currentPage   = 1
	grid          *fyne.Container
	backButton    fyne.CanvasObject
	forwardButton fyne.CanvasObject

	showingLives           []map[string]interface{}
	showingNameLabels      [10]*widget.Label
	showingTimeLabels      [10]*widget.Label
	showingStreamingLabels [10]*widget.Label
	infoButtons            [10]*widget.Button

	pageLabel           *widget.Label
	infoName            *widget.Label
	infoID              *widget.Label
	infoTime            *widget.Label
	infoRtmp            *widget.Label
	infoStreamed        *widget.Label
	infoCopyID          *widget.Button
	infoCopyRtmp        *widget.Button
	infoCheckRtmpStream *widget.Button

	idx = 0
)

func CreateClientContainer() *fyne.Container {
	pageLabel = widget.NewLabel("Page 1")

	createNewLiveButton := widget.NewButtonWithIcon("Create", theme.ContentAddIcon(), func() {
		settings.NewLiveWindow.Show()
	})

	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search For Name")

	searchButton := widget.NewButtonWithIcon("Search", theme.SearchIcon(), func() {
		params := url.Values{}
		params.Set("name", searchEntry.Text)
		parseURL, err := url.Parse(config.Config.FuzzySearchLiveURL)
		if err != nil {
			log.Println("err")
		}
		parseURL.RawQuery = params.Encode()
		response, err := http.Get(parseURL.String())
		if err != nil || response.StatusCode != 200 {
			dialog.ShowError(err, settings.MainWindow)
			panic(err)
		}

		var result []map[string]interface{}
		body, err := io.ReadAll(response.Body)
		if err == nil {
			err = json.Unmarshal(body, &result)
		}

		//fmt.Println(result)
		settings.CachedLivesOriginal = result

		UpdateShowingLives()

	})

	backButton = widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		currentPage = max(1, currentPage-1)
		UpdateShowingLives()
	})

	forwardButton = widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		currentPage = min((len(settings.CachedLivesOriginal)/10)+1, currentPage+1)
		UpdateShowingLives()
	})

	for i := 0; i < 10; i++ {
		showingNameLabels[i] = widget.NewLabel("")
		showingTimeLabels[i] = widget.NewLabel("")
		showingStreamingLabels[i] = widget.NewLabel("")
		infoButtons[i] = widget.NewButtonWithIcon("Detail", theme.InfoIcon(), getFunc(i))
	}

	grid = container.NewGridWithRows(13,
		container.NewGridWithColumns(3,
			createNewLiveButton,
			searchEntry,
			container.NewHBox(searchButton, layout.NewSpacer()),
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

	return grid
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

			settings.StreamIdEntry.SetText(settings.ToString(showingLives[i]["StreamID"]))

			infoCheckRtmpStream.SetIcon(theme.QuestionIcon())

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

	infoCopyID = widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		clipboard.Write(clipboard.FmtText, []byte(infoID.Text))
	})

	infoCopyRtmp = widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		clipboard.Write(clipboard.FmtText, []byte(infoRtmp.Text))
	})

	infoCheckRtmpStream = widget.NewButtonWithIcon("Streaming", theme.QuestionIcon(), func() {
		response, err := http.Get("http://localhost:8080/rtmp/api/list")
		if err != nil {
			fmt.Println(err)
		}

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

	return container.NewTabItem("Info", container.NewGridWithRows(5,
		container.NewGridWithColumns(2, widget.NewLabel("Name: "),
			infoName,
		),
		container.NewGridWithColumns(2, widget.NewLabel("Stream ID: "),
			container.NewHBox(infoID, layout.NewSpacer(), infoCopyID),
		),
		container.NewGridWithColumns(2, widget.NewLabel("Started Time: "),
			infoTime,
		),
		container.NewGridWithColumns(2, widget.NewLabel("Rtmp Address: "),
			container.NewHBox(infoRtmp, layout.NewSpacer(), infoCheckRtmpStream, infoCopyRtmp),
		),
		container.NewGridWithColumns(2, widget.NewLabel("Is Streamed: "),
			infoStreamed,
		),
		//container.NewGridWithColumns(3, infoStartStream, infoPushToRtmp, infoStopStream),
	))
}
