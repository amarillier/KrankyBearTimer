package main

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

// Alarm data structure
type Alarm struct {
	Enabled   bool
	Time      string // HH:MM format
	Repeat    string // "once" or "daily"
	Sunday    bool
	Monday    bool
	Tuesday   bool
	Wednesday bool
	Thursday  bool
	Friday    bool
	Saturday  bool
	SoundFile string
	LastFired string // Date when last fired (YYYY-MM-DD)
}

var alarms [5]*Alarm
var alarmWindow fyne.Window
var alarmTexts [5]*canvas.Text
var alarmBackground *canvas.Rectangle
var alarmEnabledChecks [5]*widget.Check // Store enabled checkboxes for refresh

// dayNames maps day of week to names
var dayNames = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

// loadAlarms loads alarm preferences
func loadAlarms(a fyne.App) {
	for i := 0; i < 5; i++ {
		alarms[i] = &Alarm{
			Enabled:   a.Preferences().BoolWithFallback(fmt.Sprintf("alarm%d.enabled", i+1), false),
			Time:      a.Preferences().StringWithFallback(fmt.Sprintf("alarm%d.time", i+1), "00:00"),
			Repeat:    a.Preferences().StringWithFallback(fmt.Sprintf("alarm%d.repeat", i+1), "once"),
			Sunday:    a.Preferences().BoolWithFallback(fmt.Sprintf("alarm%d.sunday", i+1), false),
			Monday:    a.Preferences().BoolWithFallback(fmt.Sprintf("alarm%d.monday", i+1), false),
			Tuesday:   a.Preferences().BoolWithFallback(fmt.Sprintf("alarm%d.tuesday", i+1), false),
			Wednesday: a.Preferences().BoolWithFallback(fmt.Sprintf("alarm%d.wednesday", i+1), false),
			Thursday:  a.Preferences().BoolWithFallback(fmt.Sprintf("alarm%d.thursday", i+1), false),
			Friday:    a.Preferences().BoolWithFallback(fmt.Sprintf("alarm%d.friday", i+1), false),
			Saturday:  a.Preferences().BoolWithFallback(fmt.Sprintf("alarm%d.saturday", i+1), false),
			SoundFile: a.Preferences().StringWithFallback(fmt.Sprintf("alarm%d.sound", i+1), ""),
			LastFired: a.Preferences().StringWithFallback(fmt.Sprintf("alarm%d.lastfired", i+1), ""),
		}
	}
}

// saveAlarm saves a single alarm to preferences
func saveAlarm(a fyne.App, index int) {
	if alarms[index] == nil {
		return
	}
	alarm := alarms[index]
	pref := fmt.Sprintf("alarm%d", index+1)
	a.Preferences().SetBool(pref+".enabled", alarm.Enabled)
	a.Preferences().SetString(pref+".time", alarm.Time)
	a.Preferences().SetString(pref+".repeat", alarm.Repeat)
	a.Preferences().SetBool(pref+".sunday", alarm.Sunday)
	a.Preferences().SetBool(pref+".monday", alarm.Monday)
	a.Preferences().SetBool(pref+".tuesday", alarm.Tuesday)
	a.Preferences().SetBool(pref+".wednesday", alarm.Wednesday)
	a.Preferences().SetBool(pref+".thursday", alarm.Thursday)
	a.Preferences().SetBool(pref+".friday", alarm.Friday)
	a.Preferences().SetBool(pref+".saturday", alarm.Saturday)
	a.Preferences().SetString(pref+".sound", alarm.SoundFile)
	a.Preferences().SetString(pref+".lastfired", alarm.LastFired)
}

// checkAlarms checks if any alarms should fire
func checkAlarms(a fyne.App) {
	now := time.Now()
	currentTime := now.Format("15:04")
	currentDate := now.Format("2006-01-02")
	currentDay := int(now.Weekday()) // 0 = Sunday, 6 = Saturday

	for i := 0; i < 5; i++ {
		if alarms[i] == nil || !alarms[i].Enabled {
			continue
		}
		alarm := alarms[i]

		// Check if time matches
		if alarm.Time != currentTime {
			continue
		}

		// Check if already fired today
		if alarm.LastFired == currentDate {
			continue
		}

		// Check repeat type and day
		if alarm.Repeat == "once" {
			// Fire once alarm
			alarmName := fmt.Sprintf("Alarm %d", i+1)
			showAlarmNotification(a, alarmName, alarm.Time, alarm.SoundFile)
			playAlarmSound(alarm.SoundFile)
			alarm.LastFired = currentDate
			// Disable after firing
			alarm.Enabled = false
			saveAlarm(a, i)
			// Refresh alarm window UI if it's open
			refreshAlarmWindowUI(i)
		} else if alarm.Repeat == "daily" {
			// Check if today is selected
			var daySelected bool
			switch currentDay {
			case 0:
				daySelected = alarm.Sunday
			case 1:
				daySelected = alarm.Monday
			case 2:
				daySelected = alarm.Tuesday
			case 3:
				daySelected = alarm.Wednesday
			case 4:
				daySelected = alarm.Thursday
			case 5:
				daySelected = alarm.Friday
			case 6:
				daySelected = alarm.Saturday
			}

			if daySelected {
				alarmName := fmt.Sprintf("Alarm %d", i+1)
				showAlarmNotification(a, alarmName, alarm.Time, alarm.SoundFile)
				playAlarmSound(alarm.SoundFile)
				alarm.LastFired = currentDate
				saveAlarm(a, i)
			}
		}
	}
}

// playAlarmSound plays the alarm sound file
func playAlarmSound(soundFile string) {
	if soundFile == "" {
		playBeep("updown")
		return
	}
	// Check if it's an MP3 file
	if strings.HasSuffix(strings.ToLower(soundFile), ".mp3") {
		if checkFileExists(soundFile) {
			playMp3(soundFile)
		} else {
			playBeep("updown")
		}
	} else if strings.HasSuffix(strings.ToLower(soundFile), ".wav") {
		if checkFileExists(soundFile) {
			playWav(soundFile)
		} else {
			playBeep("updown")
		}
	} else {
		playBeep("updown")
	}
}

// showAlarmNotification shows a notification when an alarm triggers
func showAlarmNotification(a fyne.App, alarmName string, alarmTime string, soundFile string) {
	title := fmt.Sprintf("%s - %s", alarmName, alarmTime)
	message := "Alarm triggered"
	if soundFile == "" {
		message += " (no sound file selected)"
	} else {
		// Extract just the filename for the message
		parts := strings.Split(soundFile, "/")
		filename := parts[len(parts)-1]
		message += fmt.Sprintf(" - Playing: %s", filename)
	}

	// Always send notification for alarms (regardless of notify setting)
	// This ensures alarms are visible even if no sound is selected
	fyne.CurrentApp().SendNotification(fyne.NewNotification(title, message))
}

// startAlarmChecker starts the alarm checking goroutine
func startAlarmChecker(a fyne.App) {
	go func() {
		for {
			checkAlarms(a)
			time.Sleep(time.Second) // Check every second
		}
	}()
}

// makeAlarmsWindow creates the alarm management window
func makeAlarmsWindow(a fyne.App, w fyne.Window, bg fyne.Canvas) {
	if alarmWindow != nil {
		alarmWindow.Show()
		alarmWindow.RequestFocus()
		return
	}

	alarmWindow = a.NewWindow(appName + ": Alarms")
	_, month, _ := time.Now().Date()
	if month == time.December {
		alarmWindow.SetIcon(resourceKrankyBearChristmasGrinchPng)
	} else {
		alarmWindow.SetIcon(resourceKrankyBearFedoraRedPng)
	}

	// Parse colors
	bre, bgr, bbl, ba := parseColor(bgcolor)
	tre, tgr, tbl, ta := parseColor(timecolor)

	// Background
	background := canvas.NewRectangle(color.RGBA{R: bre, G: bgr, B: bbl, A: ba})
	background.Resize(fyne.NewSize(550, 500))
	alarmBackground = background // Store for dynamic updates

	content := container.NewVBox()

	// Create UI for each alarm
	for i := 0; i < 5; i++ {
		alarmIndex := i
		if alarms[alarmIndex] == nil {
			alarms[alarmIndex] = &Alarm{}
		}
		alarm := alarms[alarmIndex]

		// Alarm label
		alarmLabelText := canvas.NewText(fmt.Sprintf("Alarm %d", i+1), color.RGBA{R: tre, G: tgr, B: tbl, A: ta})
		alarmTexts[i] = alarmLabelText

		// Enabled checkbox
		enabledCheck := widget.NewCheck("Enabled", func(checked bool) {
			alarm.Enabled = checked
			saveAlarm(a, alarmIndex)
		})
		enabledCheck.SetChecked(alarm.Enabled)
		alarmEnabledChecks[i] = enabledCheck // Store reference for refresh

		// Time entry
		timeEntry := widget.NewEntry()
		timeEntry.SetText(alarm.Time)
		timeEntry.SetPlaceHolder("HH:MM")
		timeEntry.OnChanged = func(text string) {
			alarm.Time = text
			saveAlarm(a, alarmIndex)
		}
		// Create a container to give the time entry a fixed width
		timeEntryContainer := container.NewWithoutLayout(timeEntry)
		timeEntryContainer.Resize(fyne.NewSize(100, timeEntry.MinSize().Height))
		timeEntry.Resize(fyne.NewSize(100, timeEntry.MinSize().Height))

		// Repeat selection
		repeatSelect := widget.NewSelect([]string{"once", "daily"}, func(selected string) {
			alarm.Repeat = selected
			saveAlarm(a, alarmIndex)
		})
		repeatSelect.SetSelected(alarm.Repeat)
		if alarm.Repeat == "" {
			repeatSelect.SetSelected("once")
		}

		// Day checkboxes (only show if daily)
		var dayChecks []fyne.CanvasObject
		sundayCheck := widget.NewCheck("Sun", func(checked bool) {
			alarm.Sunday = checked
			saveAlarm(a, alarmIndex)
		})
		mondayCheck := widget.NewCheck("Mon", func(checked bool) {
			alarm.Monday = checked
			saveAlarm(a, alarmIndex)
		})
		tuesdayCheck := widget.NewCheck("Tue", func(checked bool) {
			alarm.Tuesday = checked
			saveAlarm(a, alarmIndex)
		})
		wednesdayCheck := widget.NewCheck("Wed", func(checked bool) {
			alarm.Wednesday = checked
			saveAlarm(a, alarmIndex)
		})
		thursdayCheck := widget.NewCheck("Thu", func(checked bool) {
			alarm.Thursday = checked
			saveAlarm(a, alarmIndex)
		})
		fridayCheck := widget.NewCheck("Fri", func(checked bool) {
			alarm.Friday = checked
			saveAlarm(a, alarmIndex)
		})
		saturdayCheck := widget.NewCheck("Sat", func(checked bool) {
			alarm.Saturday = checked
			saveAlarm(a, alarmIndex)
		})

		sundayCheck.SetChecked(alarm.Sunday)
		mondayCheck.SetChecked(alarm.Monday)
		tuesdayCheck.SetChecked(alarm.Tuesday)
		wednesdayCheck.SetChecked(alarm.Wednesday)
		thursdayCheck.SetChecked(alarm.Thursday)
		fridayCheck.SetChecked(alarm.Friday)
		saturdayCheck.SetChecked(alarm.Saturday)

		dayChecks = []fyne.CanvasObject{
			sundayCheck, mondayCheck, tuesdayCheck, wednesdayCheck,
			thursdayCheck, fridayCheck, saturdayCheck,
		}

		// Update visibility when repeat changes
		repeatSelect.OnChanged = func(selected string) {
			alarm.Repeat = selected
			showDays := selected == "daily"
			for _, check := range dayChecks {
				if cb, ok := check.(*widget.Check); ok {
					if showDays {
						cb.Show()
					} else {
						cb.Hide()
					}
				}
			}
			saveAlarm(a, alarmIndex)
		}

		// Initially hide/show based on current repeat setting
		showDays := alarm.Repeat == "daily"
		for _, check := range dayChecks {
			if cb, ok := check.(*widget.Check); ok {
				if showDays {
					cb.Show()
				} else {
					cb.Hide()
				}
			}
		}

		// Sound file selection
		soundLabel := widget.NewLabel("Sound:")
		if alarm.SoundFile != "" {
			soundLabel.SetText("Sound: " + alarm.SoundFile)
		}
		soundButton := widget.NewButton("Select Sound File", func() {
			// Default to sounds directory if available
			var initialURI fyne.ListableURI
			if sndDir != "" && checkFileExists(sndDir) {
				var err error
				initialURI, err = storage.ListerForURI(storage.NewFileURI(sndDir))
				if err != nil {
					initialURI = nil
				}
			}

			fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err != nil {
					return
				}
				if reader == nil {
					return
				}
				defer reader.Close()
				alarm.SoundFile = reader.URI().Path()
				// Extract just the filename for display
				uriPath := reader.URI().Path()
				parts := strings.Split(uriPath, "/")
				filename := parts[len(parts)-1]
				soundLabel.SetText("Sound: " + filename)
				saveAlarm(a, alarmIndex)
			}, alarmWindow)

			// Set initial location to sounds directory if available
			if initialURI != nil {
				fileDialog.SetLocation(initialURI)
			}

			fileDialog.Show()
		})
		clearSoundButton := widget.NewButton("Clear Sound", func() {
			alarm.SoundFile = ""
			soundLabel.SetText("Sound:")
			saveAlarm(a, alarmIndex)
		})

		// Clear/Delete alarm button
		clearAlarmButton := widget.NewButton("Clear Alarm", func() {
			// Reset alarm to defaults
			alarm.Enabled = false
			alarm.Time = "00:00"
			alarm.Repeat = "once"
			alarm.Sunday = false
			alarm.Monday = false
			alarm.Tuesday = false
			alarm.Wednesday = false
			alarm.Thursday = false
			alarm.Friday = false
			alarm.Saturday = false
			alarm.SoundFile = ""
			alarm.LastFired = ""

			// Update UI
			enabledCheck.SetChecked(false)
			timeEntry.SetText("00:00")
			repeatSelect.SetSelected("once")
			sundayCheck.SetChecked(false)
			mondayCheck.SetChecked(false)
			tuesdayCheck.SetChecked(false)
			wednesdayCheck.SetChecked(false)
			thursdayCheck.SetChecked(false)
			fridayCheck.SetChecked(false)
			saturdayCheck.SetChecked(false)
			soundLabel.SetText("Sound:")

			// Hide day checkboxes when reset to "once"
			for _, check := range dayChecks {
				if cb, ok := check.(*widget.Check); ok {
					cb.Hide()
				}
			}

			saveAlarm(a, alarmIndex)
		})
		clearAlarmButton.Importance = widget.WarningImportance // orange

		// Alarm row container
		alarmRow := container.NewVBox(
			alarmLabelText,
			container.NewHBox(enabledCheck, widget.NewLabel("Time:"), timeEntryContainer, widget.NewLabel("               "), repeatSelect, widget.NewLabel("  "), clearAlarmButton),
			container.NewHBox(dayChecks[0], dayChecks[1], dayChecks[2], dayChecks[3], dayChecks[4], dayChecks[5], dayChecks[6]),
			container.NewHBox(soundLabel, soundButton, clearSoundButton),
			widget.NewSeparator(),
		)

		content.Add(alarmRow)
	}

	// Calculate minimum size needed for content
	// Width: Enabled (~80) + "Time:" (~40) + time entry (100) + spaces (~90) + repeat (~80) + spacing (~20) + Clear Alarm (~120) + padding (~40) = ~570px
	minWidth := float32(570)
	// Height: 3 alarms visible. Each alarm: label (~25) + main row (~40) + day row (~30) + sound row (~40) + separator (~5) = ~140px * 3 = ~420px + padding (~40) = ~460px
	minHeight := float32(460)

	scroll := container.NewScroll(content)
	// Don't set fixed min size on scroll, let it calculate naturally
	scroll.SetMinSize(fyne.NewSize(minWidth, minHeight))

	padded := container.NewPadded(scroll)

	// Create a container that will resize the background with the window
	mainContent := container.NewStack(background, padded)
	alarmWindow.SetContent(mainContent)

	// Update background size when content changes (will be called in updateAlarmColors)
	// The background rectangle will need to be resized to match the window content size

	alarmWindow.SetCloseIntercept(func() {
		alarmWindow.Hide()
		// Don't clear references on hide, just hide the window
	})

	// Initial color update
	updateAlarmColors()

	// Set initial window size to show 3 alarms minimum, allow resizing larger
	// Background initial size matches window
	background.Resize(fyne.NewSize(minWidth, minHeight))
	alarmWindow.Resize(fyne.NewSize(minWidth, minHeight))
	alarmWindow.CenterOnScreen()
	alarmWindow.Show()
}

// updateAlarmColors updates the alarm window with current color settings
func updateAlarmColors() {
	if alarmWindow == nil {
		return
	}

	// Parse colors
	var tre, tgr, tbl, ta uint8
	var bre, bgr, bbl, ba uint8

	colors := strings.Split(timecolor, ",")
	if len(colors) == 4 {
		if col, err := strconv.ParseUint(colors[0], 10, 8); err == nil {
			tre = uint8(col)
		}
		if col, err := strconv.ParseUint(colors[1], 10, 8); err == nil {
			tgr = uint8(col)
		}
		if col, err := strconv.ParseUint(colors[2], 10, 8); err == nil {
			tbl = uint8(col)
		}
		if col, err := strconv.ParseUint(colors[3], 10, 8); err == nil {
			ta = uint8(col)
		}
	}

	colors = strings.Split(bgcolor, ",")
	if len(colors) == 4 {
		if col, err := strconv.ParseUint(colors[0], 10, 8); err == nil {
			bre = uint8(col)
		}
		if col, err := strconv.ParseUint(colors[1], 10, 8); err == nil {
			bgr = uint8(col)
		}
		if col, err := strconv.ParseUint(colors[2], 10, 8); err == nil {
			bbl = uint8(col)
		}
		if col, err := strconv.ParseUint(colors[3], 10, 8); err == nil {
			ba = uint8(col)
		}
	}

	fyne.Do(func() {
		defer func() {
			if r := recover(); r != nil {
				// If window is invalid, clear references
				alarmWindow = nil
				alarmBackground = nil
				for i := range alarmTexts {
					alarmTexts[i] = nil
				}
			}
		}()

		if alarmBackground != nil {
			alarmBackground.FillColor = color.RGBA{R: bre, G: bgr, B: bbl, A: ba}
			// Update background size to match window content size
			if alarmWindow != nil && alarmWindow.Content() != nil {
				contentSize := alarmWindow.Content().Size()
				alarmBackground.Resize(contentSize)
			}
			alarmBackground.Refresh()
		}

		for i := 0; i < 5; i++ {
			if alarmTexts[i] != nil {
				alarmTexts[i].Color = color.RGBA{R: tre, G: tgr, B: tbl, A: ta}
				alarmTexts[i].Refresh()
			}
		}

		if alarmWindow != nil && alarmWindow.Content() != nil {
			alarmWindow.Content().Refresh()
			if alarmWindow.Canvas() != nil {
				alarmWindow.Canvas().Refresh(alarmWindow.Content())
			}
		}
	})
}

// parseColor parses a color string "R,G,B,A"
func parseColor(colorStr string) (r, g, b, a uint8) {
	parts := strings.Split(colorStr, ",")
	if len(parts) == 4 {
		if ri, err := strconv.Atoi(strings.TrimSpace(parts[0])); err == nil {
			r = uint8(ri)
		}
		if gi, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
			g = uint8(gi)
		}
		if bi, err := strconv.Atoi(strings.TrimSpace(parts[2])); err == nil {
			b = uint8(bi)
		}
		if ai, err := strconv.Atoi(strings.TrimSpace(parts[3])); err == nil {
			a = uint8(ai)
		}
	}
	return
}

// refreshAlarmWindowUI refreshes the UI for a specific alarm in the alarm window
func refreshAlarmWindowUI(alarmIndex int) {
	if alarmWindow == nil || alarmIndex < 0 || alarmIndex >= 5 {
		return
	}

	// Check if window is visible
	if !alarmWindow.Content().Visible() {
		return
	}

	// Update the enabled checkbox if it exists
	if alarmEnabledChecks[alarmIndex] != nil && alarms[alarmIndex] != nil {
		fyne.Do(func() {
			alarmEnabledChecks[alarmIndex].SetChecked(alarms[alarmIndex].Enabled)
			alarmEnabledChecks[alarmIndex].Refresh()
		})
	}
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
