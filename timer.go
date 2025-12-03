package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"fyne.io/systray"
)

// Timer-specific variables
var remain int
var notify int
var sound int
var traytimer int
var adhocTime int
var biobreakTime int
var lunchTime int
var endTime time.Time
var customTime time.Time
var endTimeSec int
var timerWindow fyne.Window        // main timer window
var timerBgImage *canvas.Image     // timer background image
var timerContent fyne.CanvasObject // timer window content (for dynamic background updates)
var timerbg string
var endsnd string
var oneminsnd string
var halfminsnd string
var starttimer int
var timerPaused bool    // track if timer is paused
var elapsedTime int     // track elapsed time for elapsed timer
var elapsedPaused bool  // track if elapsed timer is paused
var elapsedRunning bool // track if elapsed timer is running (for clean stop)

// formatTimer formats time in seconds to MM:SS format
func formatTimer(time int) string {
	secs := time % 60
	mins := (time - secs) / 60
	return fmt.Sprintf("%02d:%02d", mins, secs)
}

// centerTime centers the timer display
func centerTime(t *widget.RichText) fyne.CanvasObject {
	return container.New(layout.NewCenterLayout(), t)
}

// padTime pads the timer display
func padTime(t *widget.RichText) fyne.CanvasObject {
	pad := theme.Padding()
	return container.New(layout.NewCustomPaddedLayout(-3.5*pad, -2.5*pad, pad, pad), t)
}

// startTimer starts a countdown timer with the specified duration and name
func startTimer(timer int, name string, c fyne.Canvas, w fyne.Window) {
	remain = timer
	busy, _ := running.Get()
	if busy {
		return
	}
	w.SetTitle(appName + ": " + name)
	running.Set(true)

	if desk, ok := fyne.CurrentApp().(desktop.App); ok {
		_, month, _ := time.Now().Date()
		if month == time.December {
			desk.SetSystemTrayIcon(resourceKrankyBearChristmasGrinchPng)
		} else {
			desk.SetSystemTrayIcon(resourceKrankyBearFedoraRedPng)
		}
		systray.SetTooltip(appName)
		// systray.SetTitle(timerName)
	}

	ticker := widget.NewRichText()
	fyne.Do(func() {
		updateTime(ticker, remain)
	})

	// Reset paused state when starting new timer
	timerPaused = false

	stop := widget.NewButton("Stop", nil)
	stop.Importance = widget.WarningImportance // orange
	pause := widget.NewButton("Pause", nil)
	pause.Importance = widget.SuccessImportance // green
	resume := widget.NewButton("Resume", nil)
	resume.Importance = widget.SuccessImportance // green
	resume.Hide()                                // Initially hidden, shown when paused

	buttonRow := container.NewHBox(
		layout.NewSpacer(),
		pause,
		resume,
		stop,
		layout.NewSpacer(),
	)
	overlay := container.NewPadded(container.NewVBox(
		// padTime(ticker),
		centerTime(ticker),
		buttonRow))
	p := widget.NewModalPopUp(overlay, c)

	pause.OnTapped = func() {
		timerPaused = true
		pause.Hide()
		resume.Show()
		pause.Refresh()
		resume.Refresh()
	}

	resume.OnTapped = func() {
		timerPaused = false
		pause.Show()
		resume.Hide()
		pause.Refresh()
		resume.Refresh()
	}

	stop.OnTapped = func() {
		remain = -1 // don't notify
		timerPaused = false
		w.SetTitle(appName)
		if desk, ok := fyne.CurrentApp().(desktop.App); ok {
			_, month, _ := time.Now().Date()
			if month == time.December {
				desk.SetSystemTrayIcon(resourceKrankyBearChristmasGrinchPng)
			} else {
				desk.SetSystemTrayIcon(resourceKrankyBearFedoraRedPng)
			}
			systray.SetTooltip(appName)
			systray.SetTitle("")
			stop.Disable()
		}
		p.Hide()
	}
	// Resize popup to fit content with minimal padding (1.2x required size)
	overlayMinSize := overlay.MinSize()
	p.Resize(fyne.NewSize(overlayMinSize.Width*1.2, overlayMinSize.Height*1.2))

	go func() {
		for remain > 0 {
			// Check if paused - if so, just wait without decrementing
			if timerPaused {
				time.Sleep(time.Second)
				continue
			}

			fyne.Do(func() {
				updateTime(ticker, remain)
			})
			// system tray tooltip is not supported on Windows!
			if traytimer == 1 && runtime.GOOS != "windows" {
				if _, ok := fyne.CurrentApp().(desktop.App); ok {
					systray.SetTitle(formatTimer(remain))
				}
			}
			if remain == 60 {
				fyne.Do(func() {
					w.Show() // in case it has been hidden
				})
				if sound == 1 {
					switch oneminsnd {
					case "up", "down", "updown", "ding":
						playBeep(oneminsnd) // built in sounds
					default:
						if !checkFileExists(sndDir + "/" + oneminsnd) {
							playBeep("up")
						} else {
							playMp3(sndDir + "/" + oneminsnd) // Basso, Blow, Hero, Funk, Glass, Ping, Purr, Sosumi, Submarine,
						}
					}
				}
			} else if remain == 30 {
				fyne.Do(func() {
					w.Show() // in case it has been hidden
				})
				if sound == 1 {
					switch halfminsnd {
					case "up", "down", "updown", "ding":
						playBeep(halfminsnd) // built in sounds
					default:
						if !checkFileExists(sndDir + "/" + halfminsnd) {
							for j := 0; j <= 2; j++ {
								playBeep("down")
							}
						} else {
							playMp3(sndDir + "/" + halfminsnd) // Basso, Blow, Hero, Funk, Glass, Ping, Purr, Sosumi, Submarine,
						}
					}
				}
			}

			remain--
			time.Sleep(time.Second)
		}
		fyne.Do(func() {
			w.SetTitle(appName)
		})

		running.Set(false)
		timerPaused = false // Reset paused state when timer finishes
		if remain == 0 {
			updateTime(ticker, remain)
			fyne.Do(func() {
				stop.Disable()
				pause.Disable()
				resume.Hide()
				w.Show() // in case it has been hidden
			})
			if notify == 1 {
				fyne.CurrentApp().SendNotification(fyne.NewNotification(name+" done", "Your "+strings.ToLower(name)+" timer finished"))
				if sound == 1 {
					switch endsnd {
					case "up", "down", "updown", "ding":
						playBeep(endsnd) // built in sounds
					default:
						if !checkFileExists(sndDir + "/" + endsnd) {
							playBeep("updown")
							for i := 0; i < 3; i++ {
							}
						} else {
							playMp3(sndDir + "/" + endsnd) // grandfatherClock, baseball, pinball
						}
					}
				}
				for i := 0; i < 3; i++ {
					fyne.Do(func() {
						w.Hide()
					})
					time.Sleep(time.Second / 2)
					fyne.Do(func() {
						w.Show()
					})
					time.Sleep(time.Second / 2)
				}
			}
		}
		if desk, ok := fyne.CurrentApp().(desktop.App); ok {
			_, month, _ := time.Now().Date()
			if month == time.December {
				desk.SetSystemTrayIcon(resourceKrankyBearChristmasGrinchPng)
			} else {
				desk.SetSystemTrayIcon(resourceKrankyBearFedoraRedPng)
			}
			systray.SetTooltip(appName)
			systray.SetTitle("")
		}
		fyne.Do(func() {
			p.Hide()
		})
	}()
	fyne.Do(func() {
		p.Show()
	})
}

// updateTime updates the timer display with the current remaining time
func updateTime(out *widget.RichText, time int) {
	fyne.Do(func() {
		out.ParseMarkdown("# " + formatTimer(time))
		themeTimer(out, time)
	})
}

// startElapsedTimer starts an elapsed time timer that counts up from 0
func startElapsedTimer(c fyne.Canvas, w fyne.Window) {
	busy, _ := running.Get()
	if busy {
		return
	}
	w.SetTitle(appName + ": Elapsed Time")
	running.Set(true)

	if desk, ok := fyne.CurrentApp().(desktop.App); ok {
		_, month, _ := time.Now().Date()
		if month == time.December {
			desk.SetSystemTrayIcon(resourceKrankyBearChristmasGrinchPng)
		} else {
			desk.SetSystemTrayIcon(resourceKrankyBearFedoraRedPng)
		}
		systray.SetTooltip(appName)
	}

	elapsedTime = 0
	elapsedPaused = false
	elapsedRunning = true

	ticker := widget.NewRichText()
	fyne.Do(func() {
		updateElapsedTime(ticker, elapsedTime)
	})

	stop := widget.NewButton("Stop", nil)
	stop.Importance = widget.WarningImportance // orange
	pause := widget.NewButton("Pause", nil)
	pause.Importance = widget.SuccessImportance // green
	resume := widget.NewButton("Resume", nil)
	resume.Importance = widget.SuccessImportance // green
	resume.Hide()                                // Initially hidden, shown when paused

	buttonRow := container.NewHBox(
		layout.NewSpacer(),
		pause,
		resume,
		stop,
		layout.NewSpacer(),
	)
	overlay := container.NewPadded(container.NewVBox(
		centerTime(ticker),
		buttonRow))
	p := widget.NewModalPopUp(overlay, c)

	pause.OnTapped = func() {
		elapsedPaused = true
		pause.Hide()
		resume.Show()
		pause.Refresh()
		resume.Refresh()
	}

	resume.OnTapped = func() {
		elapsedPaused = false
		pause.Show()
		resume.Hide()
		pause.Refresh()
		resume.Refresh()
	}

	stop.OnTapped = func() {
		elapsedRunning = false // Signal to stop the timer
		elapsedPaused = false
		w.SetTitle(appName)
		if desk, ok := fyne.CurrentApp().(desktop.App); ok {
			_, month, _ := time.Now().Date()
			if month == time.December {
				desk.SetSystemTrayIcon(resourceKrankyBearChristmasGrinchPng)
			} else {
				desk.SetSystemTrayIcon(resourceKrankyBearFedoraRedPng)
			}
			systray.SetTooltip(appName)
			systray.SetTitle("")
			stop.Disable()
		}
		p.Hide()
	}

	// Resize popup to fit content with minimal padding (1.2x required size)
	overlayMinSize := overlay.MinSize()
	p.Resize(fyne.NewSize(overlayMinSize.Width*1.2, overlayMinSize.Height*1.2))

	go func() {
		for elapsedRunning {
			// Check if paused - if so, just wait without incrementing
			if elapsedPaused {
				time.Sleep(time.Second)
				continue
			}

			fyne.Do(func() {
				updateElapsedTime(ticker, elapsedTime)
			})

			// system tray tooltip is not supported on Windows!
			if traytimer == 1 && runtime.GOOS != "windows" {
				if _, ok := fyne.CurrentApp().(desktop.App); ok {
					systray.SetTitle(formatTimer(elapsedTime))
				}
			}

			elapsedTime++
			time.Sleep(time.Second)
		}

		fyne.Do(func() {
			w.SetTitle(appName)
		})

		running.Set(false)
		elapsedPaused = false
		if desk, ok := fyne.CurrentApp().(desktop.App); ok {
			_, month, _ := time.Now().Date()
			if month == time.December {
				desk.SetSystemTrayIcon(resourceKrankyBearChristmasGrinchPng)
			} else {
				desk.SetSystemTrayIcon(resourceKrankyBearFedoraRedPng)
			}
			systray.SetTooltip(appName)
			systray.SetTitle("")
		}
		fyne.Do(func() {
			p.Hide()
		})
	}()
	fyne.Do(func() {
		p.Show()
	})
}

// updateElapsedTime updates the elapsed time display
func updateElapsedTime(out *widget.RichText, time int) {
	fyne.Do(func() {
		out.ParseMarkdown("# " + formatTimer(time))
		// Use a neutral/success color for elapsed time (counting up)
		if len(out.Segments) > 0 {
			if seg, ok := out.Segments[0].(*widget.TextSegment); ok {
				seg.Style.ColorName = theme.ColorNameSuccess
			}
		}
		out.Refresh()
	})
}

// setEndTime allows the user to set a custom end time for the timer and optionally start it immediately
func setEndTime(a fyne.App, w fyne.Window, bg fyne.Canvas, autoStart bool) {
	var selectedTime time.Time
	var current string

	e := a.NewWindow("Select End Time")
	// Set window size to fit the input prompt
	e.Resize(fyne.NewSize(300, 150))
	now := time.Now()

	// check to see if predefined end / custom time is still
	// in the future, if not, try to load from preferences or set to current time + 5 minutes
	if customTime.Before(now) {
		// Try to load from preferences
		if a := fyne.CurrentApp(); a != nil {
			endTimeStr := a.Preferences().StringWithFallback("endtime.default", "")
			if endTimeStr != "" {
				// Parse the saved time (HH:MM format)
				if parsedTime, err := time.Parse("15:04", endTimeStr); err == nil {
					// Check if saved time today is still in the future
					savedTimeToday := time.Date(now.Year(), now.Month(), now.Day(), parsedTime.Hour(), parsedTime.Minute(), 0, 0, now.Location())
					if savedTimeToday.After(now) {
						customTime = savedTimeToday
						current = endTimeStr
					} else {
						// Saved time is in the past today, use it as default but user can change it
						current = endTimeStr
					}
				} else {
					current = fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute()+5)
				}
			} else {
				current = fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute()+5)
			}
		} else {
			current = fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute()+5)
		}
	} else {
		current = fmt.Sprintf("%02d:%02d", customTime.Hour(), customTime.Minute())
	}

	// Create a time entry widget
	timeEntry := widget.NewEntry()
	timeEntry.SetPlaceHolder(current)
	timeEntry.SetText(current)

	// Create a label to display messages
	messageLabel := widget.NewLabel("")

	// Create a button to submit the time
	submitButton := widget.NewButton("Set & Start", func() {
		enteredTime := timeEntry.Text
		if isValidCustomTime(enteredTime, "custom") {
			selectedTime, _ = time.Parse("15:04", enteredTime)
			customTime = time.Date(now.Year(), now.Month(), now.Day(), selectedTime.Hour(), selectedTime.Minute(), 0, 0, now.Location())
			endTime = customTime

			// Save to preferences (store as HH:MM string)
			if a := fyne.CurrentApp(); a != nil {
				timeStr := fmt.Sprintf("%02d:%02d", customTime.Hour(), customTime.Minute())
				a.Preferences().SetString("endtime.default", timeStr)
			}

			// Calculate the duration in seconds
			nowAfter := time.Now()
			endTimeSec = int(customTime.Sub(nowAfter).Seconds())

			// If duration is negative (time already passed), don't start
			if endTimeSec <= 0 {
				messageLabel.SetText("Selected time is in the past. Please choose a future time.")
				return
			}

			if autoStart {
				messageLabel.SetText("Starting timer...")
				e.Close()
				// Start the timer immediately
				go func() {
					time.Sleep(100 * time.Millisecond) // Small delay to ensure dialog is closed
					fyne.Do(func() {
						startTimer(endTimeSec, "Selected End Time", w.Canvas(), w)
					})
				}()
			} else {
				messageLabel.SetText("Custom time: " + customTime.Format("Mon Jan 2 15:04:05 MST 2006"))
				time.Sleep(1 * time.Second)
				e.Close()
			}
		} else {
			messageLabel.SetText("Enter a valid future time (HH:MM) at least 5 minutes from now")
		}
	})

	// Arrange the widgets in a vertical box
	content := container.NewVBox(
		timeEntry,
		submitButton,
		messageLabel,
	)

	e.SetContent(content)
	e.CenterOnScreen() // run centered on primary (laptop) display
	e.Show()
}

// isValidCustomTime checks if the entered time is valid in 24-hour format
// and / or is in the future compared to the current time.
func isValidCustomTime(t string, test string) bool {
	parts := strings.Split(t, ":")
	if len(parts) != 2 {
		return false
	}

	hours, err1 := strconv.Atoi(parts[0])
	minutes, err2 := strconv.Atoi(parts[1])

	if err1 != nil || err2 != nil {
		return false
	}

	if test == "custom" {
		now := time.Now()
		// allow 5 minute buffer, force selected time at least 5 minutes after current time
		customTime = time.Date(now.Year(), now.Month(), now.Day(), hours, minutes-5, 0, 0, now.Location())
		if customTime.After(now) {
			return true
		} else {
			return false
		}
	} else {
		if hours < 0 || hours > 23 || minutes < 0 || minutes > 59 {
			return false
		} else {
			return true
		}
	}
}

// updateTimerBackground dynamically updates the timer window background image
func updateTimerBackground() {
	if timerWindow == nil {
		return
	}
	if timerContent == nil {
		return
	}

	var newBg *canvas.Image
	if strings.HasSuffix(timerbg, ".png") || strings.HasSuffix(timerbg, ".jpg") {
		// if it's a png or jpg file specified, test if it exists and use it
		// otherwise use resource based image
		if checkFileExists(imgDir + "/" + timerbg) {
			newBg = canvas.NewImageFromFile(imgDir + "/" + timerbg)
		} else {
			newBg = canvas.NewImageFromResource(resourceSchoolBoard1Png)
		}
	} else {
		switch timerbg {
		case "board1":
			newBg = canvas.NewImageFromResource(resourceSchoolBoard1Png)
		case "board2":
			newBg = canvas.NewImageFromResource(resourceSchoolBoard2Png)
		case "board3":
			newBg = canvas.NewImageFromResource(resourceSchoolBoard3Png)
		case "board4":
			newBg = canvas.NewImageFromResource(resourceSchoolBoard4Png)
		case "board5":
			newBg = canvas.NewImageFromResource(resourceSchoolBoard5Png)
		case "board6":
			newBg = canvas.NewImageFromResource(resourceSchoolBoard6Png)
		case "blue":
			newBg = canvas.NewImageFromResource(resourceTBluePng)
		case "stone":
			newBg = canvas.NewImageFromResource(resourceTStonePng)
		case "almond":
			newBg = canvas.NewImageFromResource(resourceTAlmondPng)
		case "gray":
			newBg = canvas.NewImageFromResource(resourceTGrayTeachPng)
		default:
			newBg = canvas.NewImageFromResource(resourceSchoolBoard1Png)
		}
	}

	newBg.FillMode = canvas.ImageFillContain
	newBg.Translucency = 0.5

	// Update the stored reference
	timerBgImage = newBg

	// Recreate the stack with new background and existing content
	// Use fyne.Do to ensure we're on the UI thread
	fyne.Do(func() {
		newStack := container.NewStack(newBg, timerContent)
		timerWindow.SetContent(newStack)
		// Force a refresh of both the content and the window itself
		timerWindow.Content().Refresh()
		// Also refresh the canvas to ensure visual update
		if timerWindow.Canvas() != nil {
			timerWindow.Canvas().Refresh(timerWindow.Content())
		}
	})
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
