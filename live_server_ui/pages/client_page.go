package pages

import (
	"bufio"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/clipboard"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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
		body, err := io.ReadAll(response.Body)
		if err != nil {
			dialog.ShowError(err, settings.MainWindow)
			return

		}
		var res map[int]string
		err = json.Unmarshal(body, &res)
		if err != nil {
			return
		}
		var recordLists []fyne.CanvasObject
		for _, v := range res {
			str := strings.Split(v, "/")
			recordUrl, err := url.Parse(config.Config.RecordsUrl + strings.Split(v, "record")[1])
			if err != nil {
				dialog.ShowError(err, settings.MainWindow)
				return
			}
			recordLists = append(recordLists, &widget.Hyperlink{Text: str[len(str)-1], URL: recordUrl})
		}
		vList := container.NewVScroll(
			container.NewVBox(recordLists...),
		)

		vList.SetMinSize(fyne.NewSize(200, 300))

		dialog.ShowCustom("Records", "OK", vList, settings.MainWindow)

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
		container.NewGridWithColumns(4, widget.NewLabel("Name"), widget.NewLabel("Created Time"), widget.NewLabel("Is Streaming(ed)"), widget.NewLabel("Details")),
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
			//fmt.Println(img)
			if img == "" {
				infoImg.File = "icon.png"
			} else {
				im := saveImage(img)
				fmt.Println(im)
				if im == nil {
					infoImg.File = "icon.png"
				} else {
					infoImg.File = ""
					infoImg.Image = im
				}
			}

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

	infoCheckRtmpStream = widget.NewButtonWithIcon("Check", theme.QuestionIcon(), func() {
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

func saveImage(url string) image.Image {
	response, err := http.Get(url)
	if err != nil {
		return nil
	}

	defer response.Body.Close()
	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(response.Body, 32*1024)

	//file, err := os.Create(strings.Split(url, "./cache_img/"+"live_posters")[1])
	if err != nil {
		panic(err)
	}
	img, err := png.Decode(reader)
	if err != nil {
		return nil
	}
	//fmt.Println(img)
	return img
}
