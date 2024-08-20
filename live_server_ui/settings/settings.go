package settings

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"net/http"
)

const (
	CreateLiveURL = "http://localhost:8082/createLive"
	GetAllLiveURL = "http://localhost:8082/getAllLive"
	ToStreamURL   = "http://localhost:8082/pushVideoToStream"
	ToRtmpURL     = "http://localhost:8082/pushStreamToRtmp"
	EndStreamUrl  = "http://localhost:8082/endStream"
)

var (
	CachedLives   []string
	StreamIdEntry *widget.SelectEntry
	MainWindow    fyne.Window
)

type LoadedMap interface {
}

func ToString(m LoadedMap) string {
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

type ErrorString struct {
	S string
}

func (e *ErrorString) Error() string {
	return e.S
}

func TreatError(err error, response *http.Response) {
	if err != nil {
		//fmt.Println("Error")
		dialog.ShowError(err, MainWindow)
	} else if response.StatusCode != 200 {
		dialog.ShowError(&ErrorString{S: response.Status}, MainWindow)
	} else {
		//fmt.Println("Success")
		dialog.ShowConfirm("Success", "Success", nil, MainWindow)

	}
}
