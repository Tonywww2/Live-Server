package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	CachedLives   []string
	StreamIdEntry *widget.SelectEntry
)

type LoadedMap interface {
}

func name(m LoadedMap) string {
	switch (m).(type) {
	case string:
		return (m).(string)

	case bool:
		if (m).(bool) {
			return "True"
		}
		return "False"
	}

	return "ERROR"
}

func main() {
	a := app.NewWithID("live_server_ui")
	window := a.NewWindow("Main")
	window.Resize(fyne.NewSize(720, 480))
	window.SetMaster()

	windowSuccess := a.NewWindow("Success!")
	windowSuccess.SetContent(container.NewHBox(
		widget.NewLabel("Success!"),
		widget.NewButton("Confirm", func() {
			windowSuccess.Hide()
		}),
	))
	windowSuccess.SetCloseIntercept(func() { windowSuccess.Hide() })
	windowSuccess.RequestFocus()

	windowError := a.NewWindow("ERROR!")
	windowError.SetContent(container.NewHBox(
		widget.NewLabel("ERROR!"),
		widget.NewButton("Confirm", func() {
			windowError.Hide()
		}),
	))
	windowError.SetCloseIntercept(func() { windowError.Hide() })
	windowError.RequestFocus()

	StreamIdEntry = widget.NewSelectEntry(CachedLives)

	tab := container.NewAppTabs(
		createLivePage(&window, &windowSuccess, &windowError),
		creatGetAllPage(&windowError),
		pushVideoPage(&window, &windowSuccess, &windowError),
		pushRtmpPage(&windowSuccess, &windowError),
		endStreamPage(&windowSuccess, &windowError),
	)
	window.SetContent(tab)

	go func() {
		for range time.Tick(time.Second) {
			//fmt.Println(entry1.Text)
			StreamIdEntry.SetOptions(CachedLives)
		}
	}()

	window.Show()

	a.Run()
}

func createLivePage(window *fyne.Window, windowSuccess *fyne.Window, windowError *fyne.Window) *container.TabItem {
	nameEntry := widget.NewEntry()
	posterLabel := widget.NewLabel("")

	file := dialog.NewFileOpen(func(f fyne.URIReadCloser, err error) {
		if f != nil {
			posterUri := strings.Split(f.URI().String(), "//")[1]
			posterLabel.SetText(posterUri)
			defer f.Close()

		}
	}, *window)

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
			result, err := http.PostForm("http://localhost:8082/createLive", payload)
			if err != nil || result.StatusCode != 200 {
				fmt.Println("Error")
				(*windowError).Show()
			} else {
				fmt.Println("Success")
				(*windowSuccess).Show()

			}

		}),
	)

	return container.NewTabItem("Create New Live", livePage)
}

func creatGetAllPage(windowError *fyne.Window) *container.TabItem {
	resultLabel := widget.NewEntry()
	resultLabel.MultiLine = true

	getButton := widget.NewButton("Get", func() {
		response, err := http.Get("http://localhost:8082/getAllLive")
		if err != nil || response.StatusCode != 200 {
			fmt.Println("Error")
			(*windowError).Show()
		}

		var result map[string]map[string]interface{}
		body, err := ioutil.ReadAll(response.Body)
		if err == nil {
			err = json.Unmarshal(body, &result)
		}

		resultText := "Result: \n"
		i := 0
		CachedLives = make([]string, 0)
		for k, v := range result {
			text := name(v["Name"]) + ", " +
				name(v["Poster"]) + ", " +
				name(v["StartTime"]) + ", " +
				name(v["RtmpAddr"]) + ", " +
				name(v["IsStreamed"])
			resultText += k + "\n" + text + "\n\n"

			CachedLives = append(CachedLives, name(v["Name"])+"//"+k)
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

func pushVideoPage(window *fyne.Window, windowSuccess *fyne.Window, windowError *fyne.Window) *container.TabItem {
	path := widget.NewLabel("")

	file := dialog.NewFileOpen(func(f fyne.URIReadCloser, err error) {
		if f != nil {
			pathUri := strings.Split(f.URI().String(), "//")[1]
			path.SetText(pathUri)
			defer f.Close()
		}
	}, *window)

	pushButton := widget.NewButton("Push", func() {

		temp := strings.Split(StreamIdEntry.Text, "//")
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
		response, err := http.PostForm("http://localhost:8082/pushVideoToStream", payload)
		if err != nil || response.StatusCode != 200 {
			fmt.Println("Error")
			(*windowError).Show()
		} else {
			fmt.Println("Success")
			(*windowSuccess).Show()

		}
		fmt.Println(response)
	})

	getAllPage := container.NewVBox(
		widget.NewLabel("Push Video"),
		widget.NewLabel("Stream ID: "),
		StreamIdEntry,
		container.NewHBox(widget.NewLabel("Path:"), widget.NewButtonWithIcon("", theme.FolderOpenIcon(), file.Show), path),
		pushButton,
	)

	return container.NewTabItem("Push video to Live", getAllPage)
}

func pushRtmpPage(windowSuccess *fyne.Window, windowError *fyne.Window) *container.TabItem {
	rtmp := widget.NewEntry()

	pushButton := widget.NewButton("Push", func() {
		temp := strings.Split(StreamIdEntry.Text, "//")
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
		response, err := http.PostForm("http://localhost:8082/pushStreamToRtmp", payload)
		if err != nil || response.StatusCode != 200 {
			fmt.Println("Error")
			(*windowError).Show()
		} else {
			fmt.Println("Success")
			(*windowSuccess).Show()

		}
		fmt.Println(response)
	})

	pushRtmp := container.NewVBox(
		widget.NewLabel("Stream to Rtmp"),
		widget.NewLabel("Stream ID: "),
		StreamIdEntry,
		widget.NewLabel("Path:"),
		rtmp,
		pushButton,
	)

	return container.NewTabItem("Push video to Live", pushRtmp)
}

func endStreamPage(windowSuccess *fyne.Window, windowError *fyne.Window) *container.TabItem {
	endButton := widget.NewButton("End", func() {
		temp := strings.Split(StreamIdEntry.Text, "//")
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
		response, err := http.PostForm("http://localhost:8082/endStream", payload)
		if err != nil || response.StatusCode != 200 {
			fmt.Println("Error")
			(*windowError).Show()
		} else {
			fmt.Println("Success")
			(*windowSuccess).Show()

		}
		fmt.Println(response)
	})

	endStream := container.NewVBox(
		widget.NewLabel("End Stream"),
		widget.NewLabel("Stream ID: "),
		StreamIdEntry,
		endButton,
	)

	return container.NewTabItem("End Stream", endStream)
}
