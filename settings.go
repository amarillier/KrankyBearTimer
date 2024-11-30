package main

import (
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var fileButton *widget.Button
var selectedFile *widget.Label
var fileURI fyne.URI

func makeSettings(a fyne.App, w fyne.Window, bg fyne.Canvas) {
	// settings window
	settings := a.NewWindow(timerName + ": Settings")
	settingsText := `All updates are applied / saved immediately.
	Note: timer background and timer buttons on the timer window do not currently auto refresh,
	restart is required. Time changes do take immediate effect, refresh of background and buttons is planned`
	setText := widget.NewLabel(settingsText)
	setText.TextStyle = fyne.TextStyle{Bold: true}
	todoText := widget.NewLabel("Still to be added: allow .mid and .wav sounds as well as selectable background images in addition to built in images")
	todoText.TextStyle = fyne.TextStyle{Italic: true, Bold: true}

	mp3files, err := listMatchingFiles(sndDir, "*.mp3")
	if err != nil {
		log.Fatal(err)
	}
	mp3 := []string{"ding", "down", "up", "updown"}
	for _, file := range mp3files {
		mp3 = append(mp3, file)
	}

	notifications := widget.NewCheck("", func(value bool) {
		if debug == 1 {
			log.Println("notifications set to", value)
		}
		switch value {
		case true:
			notify = 1
		case false:
			notify = 0
		}
		a.Preferences().SetInt("notify.default", notify)
	})
	soundalerts := widget.NewCheck("", func(value bool) {
		if debug == 1 {
			log.Println("sound alerts set to", value)
		}
		switch value {
		case true:
			sound = 1
		case false:
			sound = 0
		}
		a.Preferences().SetInt("sound.default", sound)
	})
	background := widget.NewSelect([]string{"almond", "blue", "stone", "converge24", "converge24a", "taniumtimer2"}, func(value string) {
		if debug == 1 {
			log.Println("background set to", value)
		}
		timerbg = value
		bg := canvas.NewImageFromResource(resourceTaniumTrainBlueSvg)
		switch timerbg {
		case "taniumtimer2":
			bg = canvas.NewImageFromResource(resourceTaniumTimer2Svg)
		case "blue":
			bg = canvas.NewImageFromResource(resourceTaniumTrainBlueSvg)
		case "stone":
			bg = canvas.NewImageFromResource(resourceTaniumTrainStoneSvg)
		case "almond":
			bg = canvas.NewImageFromResource(resourceTaniumTrainAlmondSvg)
		case "taniumtimer":
			bg = canvas.NewImageFromResource(resourceTaniumTimerSvg)
		case "converge24":
			bg = canvas.NewImageFromResource(resourceTaniumConverge2024Svg)
		case "converge24a":
			bg = canvas.NewImageFromResource(resourceTaniumConverge2024aSvg)
		default:
			bg = canvas.NewImageFromResource(resourceTaniumTrainBlueSvg)
		}
		bg.FillMode = canvas.ImageFillContain
		bg.Translucency = 0.5 // 0.85
		// bg.Refresh()          // WHY does this not refresh????
		w.Canvas().Refresh(bg)
		a.Preferences().SetString("background.default", timerbg)
	})

	//endsound := widget.NewSelect([]string{"updown", "up", "down", "ding", "baseball.mp3", "pinball.mp3", "grandfatherclock.mp3"}, func(value string) {
	endsound := widget.NewSelect(mp3, func(value string) {
		if debug == 1 {
			log.Println("endsound set to", value)
		}
		endsnd = value // strings.Replace(value, "builtin ", "", 1)
		switch endsnd {
		case "up", "down", "updown", "ding":
			playBeep(endsnd) // built in sounds
		default:
			playMp3(sndDir + "/" + endsnd)
		}
		a.Preferences().SetString("endsound.default", endsnd)
	})
	// oneminsound := widget.NewSelect([]string{"updown", "up", "down", "ding", "hero.mp3", "sosumi.mp3", "baseball.mp3", "pinball.mp3", "grandfatherclock.mp3"}, func(value string) {
	oneminsound := widget.NewSelect(mp3, func(value string) {
		if debug == 1 {
			log.Println("oneminsound set to", value)
		}
		oneminsnd = value // strings.Replace(value, "builtin ", "", 1)
		switch oneminsnd {
		case "up", "down", "updown", "ding":
			playBeep(oneminsnd) // built in sounds
		default:
			playMp3(sndDir + "/" + oneminsnd)
		}
		a.Preferences().SetString("oneminsound.default", oneminsnd)
	})
	// halfminsound := widget.NewSelect([]string{"updown", "up", "down", "ding", "sosumi.mp3", "hero.mp3", "baseball.mp3", "pinball.mp3", "grandfatherclock.mp3"}, func(value string) {
	halfminsound := widget.NewSelect(mp3, func(value string) {
		if debug == 1 {
			log.Println("halfminsound set to", value)
		}
		halfminsnd = value // strings.Replace(value, "builtin ", "", 1)
		switch halfminsnd {
		case "up", "down", "updown", "ding":
			playBeep(halfminsnd) // built in sounds
		default:
			playMp3(sndDir + "/" + halfminsnd)
		}
		a.Preferences().SetString("halfminsound.default", halfminsnd)
	})
	adhoc := widget.NewRadioGroup([]string{"5", "10", "15"}, func(value string) {
		if debug == 1 {
			log.Println("adhoc time set to", value)
		}
		adhocTime, _ = strconv.Atoi(value)
		adhocTime *= 60
		a.Preferences().SetInt("adhoc.default", adhocTime)
	})
	adhoc.Horizontal = true
	biobreak := widget.NewRadioGroup([]string{"5", "10", "15", "20"}, func(value string) {
		if debug == 1 {
			log.Println("bio break set to", value)
		}
		biobreakTime, _ = strconv.Atoi(value)
		biobreakTime *= 60
		a.Preferences().SetInt("biobreak.default", biobreakTime)
	})
	biobreak.Horizontal = true
	lunch := widget.NewRadioGroup([]string{"30", "45", "60", "90", "120"}, func(value string) {
		if debug == 1 {
			log.Println("lunch break set to", value)
		}
		lunchTime, _ = strconv.Atoi(value)
		lunchTime *= 60
		a.Preferences().SetInt("lunch.default", lunchTime)
	})
	lunch.Horizontal = true
	reset := widget.NewButton("Reset default settings", func() {
		if debug == 1 {
			log.Println("preferences reset to defaults")
		}
		writeDefaultSettings(a)
		notifications.SetChecked(true)
		soundalerts.SetChecked(true)
		timerbg = "blue"
		endsnd = "baseball.mp3"
		oneminsnd = "hero.mp3"
		halfminsnd = "sosumi.mp3"
		lunchTime = 3600
		biobreakTime = 600
		adhocTime = 300
		background.Selected = timerbg
		endsound.Selected = endsnd
		oneminsound.Selected = oneminsnd
		halfminsound.Selected = halfminsnd

		background.Refresh()
		endsound.Refresh()
		oneminsound.Refresh()
		halfminsound.Refresh()
		adhoc.Refresh()
		biobreak.Refresh()
		lunch.Refresh()
		switch adhocTime {
		case 300:
			adhoc.SetSelected("5")
		case 600:
			adhoc.SetSelected("10")
		case 900:
			adhoc.SetSelected("15")
		}
		switch biobreakTime {
		case 300:
			biobreak.SetSelected("5")
		case 600:
			biobreak.SetSelected("10")
		case 900:
			biobreak.SetSelected("15")
		case 1200:
			biobreak.SetSelected("20")
		}
		switch lunchTime {
		case 1800:
			lunch.SetSelected("30")
		case 2700:
			lunch.SetSelected("45")
		case 3600:
			lunch.SetSelected("60")
		case 5400:
			lunch.SetSelected("90")
		case 7200:
			lunch.SetSelected("120")
		}
	})
	reset.Importance = widget.SuccessImportance // green
	close := widget.NewButton("Close settings", func() {
		settings.Close()
	})
	close.Importance = widget.WarningImportance // orange

	// fileButton := widget.NewButton("File", func() { showFilePicker(settings) })
	// allow for file selectors

	if notify == 1 {
		notifications.SetChecked(true)
	} else {
		notifications.SetChecked(false)
	}
	if sound == 1 {
		soundalerts.SetChecked(true)
	} else {
		soundalerts.SetChecked(false)
	}
	background.Selected = timerbg
	endsound.Selected = endsnd
	oneminsound.Selected = oneminsnd
	halfminsound.Selected = halfminsnd
	switch adhocTime {
	case 300:
		adhoc.SetSelected("5")
	case 600:
		adhoc.SetSelected("10")
	case 900:
		adhoc.SetSelected("15")
	}
	switch biobreakTime {
	case 300:
		biobreak.SetSelected("5")
	case 600:
		biobreak.SetSelected("10")
	case 900:
		biobreak.SetSelected("15")
	case 1200:
		biobreak.SetSelected("20")
	}
	switch lunchTime {
	case 1800:
		lunch.SetSelected("30")
	case 2700:
		lunch.SetSelected("45")
	case 3600:
		lunch.SetSelected("60")
	case 5400:
		lunch.SetSelected("90")
	case 7200:
		lunch.SetSelected("120")
	}

	setform := widget.NewForm(
		widget.NewFormItem("Notifications", notifications),
		widget.NewFormItem("Sound alerts", soundalerts),
		widget.NewFormItem("Background", background),
		widget.NewFormItem("Timer end sound", endsound),
		widget.NewFormItem("One minute sound", oneminsound),
		widget.NewFormItem("Half minute sound", halfminsound),
		widget.NewFormItem("Ad hoc break", adhoc),
		widget.NewFormItem("Bio break", biobreak),
		widget.NewFormItem("Lunch break", lunch),
		// widget.NewFormItem("File picker", fileButton),
		widget.NewFormItem("", reset),
		widget.NewFormItem("", close),
	)

	settings.Resize(fyne.NewSize(500, 300))
	// settings.SetIcon(resourceRedSettingsGearSvg)
	settings.CenterOnScreen() // run centered on primary (laptop) display
	settings.SetContent(container.NewVBox(setText, setform, todoText))
	settings.Show()

}

func writeDefaultSettings(a fyne.App) {
	// write default prefs that can be modified via settings
	a.Preferences().SetInt("adhoc.default", 300)
	a.Preferences().SetInt("biobreak.default", 600)
	a.Preferences().SetInt("lunch.default", 3600)
	a.Preferences().SetInt("notify.default", 1)
	a.Preferences().SetInt("sound.default", 1)
	a.Preferences().SetString("background.default", "blue")
	a.Preferences().SetString("endsound.default", "baseball.mp3")
	a.Preferences().SetString("oneminsound.default", "hero.mp3")
	a.Preferences().SetString("halfminsound.default", "sosumi.mp3")
	// example prefs:
	//{"adhoc.default":300,"background.default":"blue","biobreak.default":600,"endsound.default":"baseball.mp3","halfminsound.default":"sosumi.mp3","lunch.default":3600,"notify.default":1,"oneminsound.default":"hero.mp3"}
}

func writeSettings(a fyne.App) {
	// write current settings to global prefs
	a.Preferences().SetInt("adhoc.default", adhocTime)
	a.Preferences().SetInt("biobreak.default", biobreakTime)
	a.Preferences().SetInt("lunch.default", lunchTime)
	a.Preferences().SetInt("notify.default", notify)
	a.Preferences().SetInt("sound.default", sound)
	a.Preferences().SetString("background.default", timerbg)
	a.Preferences().SetString("endsound.default", endsnd)
	a.Preferences().SetString("oneminsound.default", oneminsnd)
	a.Preferences().SetString("halfminsound.default", halfminsnd)
}

func showFilePicker(w fyne.Window) {
	// Show file picker and return selected file
	// https://dev.to/cjr29/learning-go-building-a-file-picker-using-fyneio-33le
	dialog.ShowFileOpen(func(f fyne.URIReadCloser, err error) {
		saveFile := "NoFileYet"
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if f == nil {
			return
		}
		saveFile = f.URI().Path()
		fileURI = f.URI()
		selectedFile.SetText(saveFile)
	}, w)
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
