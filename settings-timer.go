package main

import (
	"log"
	"regexp"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/spiretechnology/go-autostart/v2"

	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
)

var settingsti fyne.Window
var settingsth fyne.Window
var bgImage []string

func makeSettingsTimer(a fyne.App, w fyne.Window, bg fyne.Canvas) {
	// settings window
	if settingsti != nil { // &&  !settingst.Content().Visible() {
		settingsti.Show()
		teapot(a, settingsti)
	} else {
		settingsti = a.NewWindow(appName + ": Settings")
		settingsti.SetIcon(resourceKrankyBearPng)
		settingsText := `All updates are applied / saved immediately.
	Timer background now updates dynamically - no restart required!
	Time changes take immediate effect.`
		setText := widget.NewLabel(settingsText)
		setText.TextStyle = fyne.TextStyle{Bold: true}
		todoText := widget.NewLabel("Still to be added: allow .mid and .wav sounds as well as selectable background images in addition to built in images")
		todoText.TextStyle = fyne.TextStyle{Italic: true, Bold: true}

		pngfiles, err := listMatchingFiles(imgDir, "*.png")
		if err != nil {
			log.Fatal(err)
		}
		jpgfiles, err := listMatchingFiles(imgDir, "*.jpg")
		if err != nil {
			log.Fatal(err)
		}

		mp3files, err := listMatchingFiles(sndDir, "*.mp3")
		if err != nil {
			log.Fatal(err)
		}
		mp3 := []string{"ding", "down", "up", "updown"}
		//for _, file := range mp3files {
		//	mp3 = append(mp3, file)
		//}
		mp3 = append(mp3, mp3files...)

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
		systraytimer := widget.NewCheck("", func(value bool) {
			if debug == 1 {
				log.Println("system tray timer set to", value)
			}
			switch value {
			case true:
				traytimer = 1
			case false:
				traytimer = 0
			}
			a.Preferences().SetInt("traytimer.default", traytimer)
		})
		startatboot := widget.NewCheck("", func(value bool) {
			if debug == 1 {
				log.Println("startatboot set to", value)
			}
			autoTimer := autostart.New(autostart.Options{
				Label:       "com.github.amarillier.KrankyBearTimer",
				Name:        "KrankyBearTimer",
				Description: "Kranky Bear Timer",
				Mode:        autostart.ModeUser,
				Arguments:   []string{},
			})
			switch value {
			case true:
				starttimer = 1
				autoTimer.Enable()
			case false:
				starttimer = 0
				autoTimer.Disable()
			}
			a.Preferences().SetInt("starttimer.default", starttimer)
		})
		tanium := regexp.MustCompile(`^(?i)tanium`)
		if tanium.MatchString(appNameCustom) {
			// if Tanium branded, add selectable Tanium backgrounds
			bgImage = []string{"almond", "blue", "stone", "gray", "board1", "board2", "board3", "board4", "board5", "board6"}
		} else {
			bgImage = []string{"board1", "board2", "board3", "board4", "board5", "board6"}
		}
		bgImage = append(bgImage, pngfiles...)
		bgImage = append(bgImage, jpgfiles...)
		// background := widget.NewSelect([]string{"almond", "blue", "stone", "gray"}, func(value string) {
		background := widget.NewSelect(bgImage, func(value string) {
			if debug == 1 {
				log.Println("background set to", value)
			}
			timerbg = value
			a.Preferences().SetString("background.default", timerbg)
			// Dynamically update the timer background
			updateTimerBackground()
		})

		endsound := widget.NewSelect(mp3, func(value string) {
			if debug == 1 {
				log.Println("endsound set to", value)
			}
			endsnd = value
			switch endsnd {
			case "up", "down", "updown", "ding":
				playBeep(endsnd)
			default:
				playMp3(sndDir + "/" + endsnd)
			}
			a.Preferences().SetString("endsound.default", endsnd)
		})
		oneminsound := widget.NewSelect(mp3, func(value string) {
			if debug == 1 {
				log.Println("oneminsound set to", value)
			}
			oneminsnd = value
			switch oneminsnd {
			case "up", "down", "updown", "ding":
				playBeep(oneminsnd)
			default:
				playMp3(sndDir + "/" + oneminsnd)
			}
			a.Preferences().SetString("oneminsound.default", oneminsnd)
		})
		halfminsound := widget.NewSelect(mp3, func(value string) {
			if debug == 1 {
				log.Println("halfminsound set to", value)
			}
			halfminsnd = value
			switch halfminsnd {
			case "up", "down", "updown", "ding":
				playBeep(halfminsnd)
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
			writeDefaultSettingsTimer(a)
			notifications.SetChecked(true)
			soundalerts.SetChecked(true)
			systraytimer.SetChecked(false)
			startatboot.SetChecked(false)
			timerbg = "board1"
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
			systraytimer.Refresh()
			startatboot.Refresh()
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
			// Dynamically update the timer background after reset
			updateTimerBackground()
		})
		reset.Importance = widget.SuccessImportance // green
		close := widget.NewButton("Close settings", func() {
			settingsti.Close()
			settingsti = nil
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
		if traytimer == 1 {
			systraytimer.SetChecked(true)
		} else {
			systraytimer.SetChecked(false)
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
			widget.NewFormItem("System Tray Timer (N/A for Windows)", systraytimer),
			widget.NewFormItem("Auto Start at Boot", startatboot),
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

		settingsti.Resize(fyne.NewSize(500, 300))
		settingsti.CenterOnScreen() // run centered on primary (laptop) display
		settingsti.SetContent(container.NewVBox(setText, setform, todoText))
		settingsti.SetCloseIntercept(func() {
			settingsti.Close()
			settingsti = nil
		})
		settingsti.Show()
	}
}

func writeDefaultSettingsTimer(a fyne.App) {
	// write default prefs for timer settings only
	a.Preferences().SetInt("adhoc.default", 300)
	a.Preferences().SetInt("biobreak.default", 600)
	a.Preferences().SetInt("lunch.default", 3600)
	a.Preferences().SetInt("notify.default", 1)
	a.Preferences().SetInt("sound.default", 1)
	a.Preferences().SetInt("systraytimer.default", 0)
	a.Preferences().SetInt("starttimer.default", 0)
	a.Preferences().SetString("background.default", "board1")
	a.Preferences().SetString("endsound.default", "baseball.mp3")
	a.Preferences().SetString("oneminsound.default", "hero.mp3")
	a.Preferences().SetString("halfminsound.default", "sosumi.mp3")
	// a.Preferences().SetString("timername.default","")
	// example prefs:
	//{"adhoc.default":300,"background.default":"blue","biobreak.default":600,"endsound.default":"baseball.mp3","halfminsound.default":"sosumi.mp3","lunch.default":3600,"notify.default":1,"oneminsound.default":"hero.mp3"}
}

func writeSettingsTimer(a fyne.App) {
	// write current timer settings to global prefs
	a.Preferences().SetInt("adhoc.default", adhocTime)
	a.Preferences().SetInt("biobreak.default", biobreakTime)
	a.Preferences().SetInt("lunch.default", lunchTime)
	a.Preferences().SetInt("notify.default", notify)
	a.Preferences().SetInt("sound.default", sound)
	a.Preferences().SetInt("systraytimer.default", traytimer)
	a.Preferences().SetInt("starttimer.default", starttimer)
	a.Preferences().SetString("background.default", timerbg)
	a.Preferences().SetString("endsound.default", endsnd)
	a.Preferences().SetString("oneminsound.default", oneminsnd)
	a.Preferences().SetString("halfminsound.default", halfminsnd)
}

func makeSettingsTheme(a fyne.App, w fyne.Window, bg fyne.Canvas) {
	// allow modifying the fyne theme
	// this is dependent on fyne_settings in ~/go/pkg/mod/fyne.io/fyne/v2/cmd/fyne_settings/settings
	// but here I use a customized version to add a button 'Apply & Close'
	// modify as shown below
	if settingsth != nil { // &&  !settingsc.Content().Visible() {
		settingsth.Show()
		teapot(a, settingsth)
	} else {
		s := settings.NewSettings()
		settingsth = a.NewWindow(appName + ": Theme Settings")
		settingsth.SetIcon(resourceKrankyBearPng)

		appearance := s.LoadAppearanceScreen(w)
		tabs := container.NewAppTabs(
			&container.TabItem{Text: "Theme Appearance - affects all fyne based apps", Icon: s.AppearanceIcon(), Content: appearance})
		tabs.SetTabLocation(container.TabLocationLeading)
		settingsth.SetContent(tabs)

		settingsth.Resize(fyne.NewSize(520, 520))
		settingsth.CenterOnScreen() // run centered on primary (laptop) display
		settingsth.Show()
		settingsth.SetCloseIntercept(func() {
			settingsth.Close()
			settingsth = nil
		})
	}
}

// modify the latest ~/go/pkg/mod/fyne.io/fyne/v2/cmd/fyne_settings/settings/appearance.go

// add to function LoadAppearanceScreen last part with Apply & Close button:
/*
bottom := container.NewHBox(layout.NewSpacer(),
		&widget.Button{Text: "Apply", Importance: widget.HighImportance, OnTapped: func() {
			if s.fyneSettings.Scale == 0.0 {
				s.chooseScale(1.0)
			}
			err := s.save()
			if err != nil {
				fyne.LogError("Failed on saving", err)
			}

			s.appliedScale(s.fyneSettings.Scale)
		}},
		&widget.Button{Text: "Apply & Close", Importance: widget.WarningImportance, OnTapped: func() {
			if s.fyneSettings.Scale == 0.0 {
				s.chooseScale(1.0)
			}
			err := s.save()
			if err != nil {
				fyne.LogError("Failed on saving", err)
			}

			s.appliedScale(s.fyneSettings.Scale)
			w.Close()
		}},
	)
*/

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
