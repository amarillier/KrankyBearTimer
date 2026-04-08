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
	"github.com/IamFaizanKhalid/lock"
	"github.com/go-vgo/robotgo"
	_ "github.com/go-vgo/robotgo/base"
	_ "github.com/go-vgo/robotgo/key"
	_ "github.com/go-vgo/robotgo/mouse"
	_ "github.com/go-vgo/robotgo/screen"
	_ "github.com/go-vgo/robotgo/window"
	"github.com/itchyny/volume-go"
)

// var clock fyne.Window
var clockTimeText *canvas.Text
var clockDateText *canvas.Text
var clockUtcText *canvas.Text
var clockBackground *canvas.Rectangle

// Additional timezone display elements
var timezone1Text *canvas.Text
var timezone2Text *canvas.Text
var timezone3Text *canvas.Text
var timezone4Text *canvas.Text
var timezone5Text *canvas.Text
var weatherText *canvas.Text

// clockVBox removed - not needed since we don't do dynamic rebuilding anymore

// Position preservation removed - we no longer recreate the window, just show/hide elements

func desktopclock(a fyne.App) { // , w fyne.Window, bg fyne.Canvas) {
	// Stop any existing update loop before creating a new clock window
	if clockUpdateLoopRunning && clockUpdateLoopStop != nil {
		close(clockUpdateLoopStop)
		clockUpdateLoopStop = nil
		clockUpdateLoopRunning = false
	}

	if clock != nil {
		// If clock already exists, just update colors and request focus
		// Don't try to check Content() as it can crash if window is invalid
		// updateClockColors will handle nil checks internally
		updateClockColors()
		// Use defer/recover to safely try to request focus
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Window is invalid, clear reference
					clock = nil
				}
			}()
			if clock != nil {
				clock.RequestFocus()
			}
		}()
		if clock != nil {
			return
		}
		// If clock was cleared due to error, fall through to create new one
	}

	// Create stop channel for graceful shutdown
	clockUpdateLoopStop = make(chan bool)
	clockUpdateLoopRunning = true

	// Create the clock window
	{
		var tre, tgr, tbl, ta uint8
		colors := strings.Split(timecolor, ",")
		col, _ := strconv.ParseUint(colors[0], 10, 8)
		tre = uint8(col)
		col, _ = strconv.ParseUint(colors[1], 10, 8)
		tgr = uint8(col)
		col, _ = strconv.ParseUint(colors[2], 10, 8)
		tbl = uint8(col)
		col, _ = strconv.ParseUint(colors[3], 10, 8)
		ta = uint8(col)

		var bre, bgr, bbl, ba uint8
		colors = strings.Split(bgcolor, ",")
		col, _ = strconv.ParseUint(colors[0], 10, 8)
		bre = uint8(col)
		col, _ = strconv.ParseUint(colors[1], 10, 8)
		bgr = uint8(col)
		col, _ = strconv.ParseUint(colors[2], 10, 8)
		bbl = uint8(col)
		col, _ = strconv.ParseUint(colors[3], 10, 8)
		ba = uint8(col)

		var dre, dgr, dbl, da uint8
		colors = strings.Split(datecolor, ",")
		col, _ = strconv.ParseUint(colors[0], 10, 8)
		dre = uint8(col)
		col, _ = strconv.ParseUint(colors[1], 10, 8)
		dgr = uint8(col)
		col, _ = strconv.ParseUint(colors[2], 10, 8)
		dbl = uint8(col)
		col, _ = strconv.ParseUint(colors[3], 10, 8)
		da = uint8(col)

		var ure, ugr, ubl, ua uint8
		colors = strings.Split(utccolor, ",")
		col, _ = strconv.ParseUint(colors[0], 10, 8)
		ure = uint8(col)
		col, _ = strconv.ParseUint(colors[1], 10, 8)
		ugr = uint8(col)
		col, _ = strconv.ParseUint(colors[2], 10, 8)
		ubl = uint8(col)
		col, _ = strconv.ParseUint(colors[3], 10, 8)
		ua = uint8(col)
		clockName := appName + ": Clock"
		// clock = a.NewWindow("Kranky Bear Clock")
		clock = a.NewWindow(clockName)
		clock.SetIcon(resourceKrankyBearFedoraRedPng)

		now := time.Now()
		// timeFormat := `15:04:05`
		// timeFormat := `3:04:05 PM (MST)`
		timeFormat := ``
		if showhr12 == 1 {
			timeFormat += `3:04`
		} else {
			timeFormat += `15:04`
		}
		if showseconds == 1 {
			timeFormat += `:05`
		}
		if showhr12 == 1 {
			timeFormat += ` PM` // this needs to be added AFTER seconds if 12 hour
		}
		if showtimezone == 1 {
			timeFormat += ` (MST)`
		}

		// Get the local time zone and offset
		_, offset := now.Zone()
		offsetHours := offset / 3600
		offsetMinutes := (offset % 3600) / 60
		offsetString := fmt.Sprintf(" (local is  %+02d:%02d)", offsetHours, offsetMinutes) // ZZZ
		// utcFormat := `(UTC 3:04 PM Z07)`
		utcFormat := `UTC 3:04 PM` //   (` + offsetString + `)`  // ZZZ
		nonutcFormat := `3:04 PM`
		dateFormat := ` Monday, January 2, 2006 `

		// nowtime := canvas.NewText(now.Format(timeFormat), color.RGBA{R: 255, G: 123, B: 31, A: 255})
		nowtime := canvas.NewText(now.Format(timeFormat), color.RGBA{R: tre, G: tgr, B: tbl, A: ta})
		nowtime.TextStyle = fyne.TextStyle{Bold: true}
		// nowtime.TextStyle = fyne.TextStyle{Monospace: true} // EXAMPLE FONT TYPE
		nowtime.Alignment = fyne.TextAlignCenter
		nowtime.TextSize = float32(timesize)
		clockTimeText = nowtime // store reference for dynamic updates

		// utctime := canvas.NewText(now.Format(utcFormat), color.RGBA{R: 255, G: 123, B: 31, A: 255})
		utctime := canvas.NewText(now.Format(utcFormat), color.RGBA{R: ure, G: ugr, B: ubl, A: ua})
		utctime.TextStyle = fyne.TextStyle{Bold: true}
		utctime.Alignment = fyne.TextAlignCenter
		utctime.TextSize = float32(utcsize)
		clockUtcText = utctime // store reference for dynamic updates

		// nowdate := canvas.NewText(now.Format(dateFormat), color.RGBA{R: 208, G: 145, B: 38, A: 255})
		nowdate := canvas.NewText(now.Format(dateFormat), color.RGBA{R: dre, G: dgr, B: dbl, A: da})
		nowdate.TextStyle = fyne.TextStyle{Bold: true}
		nowdate.Alignment = fyne.TextAlignCenter
		nowdate.TextSize = float32(datesize)
		clockDateText = nowdate // store reference for dynamic updates

		//background := canvas.NewRectangle(color.RGBA{R: 0, G: 0, B: 255, A: 255})
		// bgcolor := color.RGBA{R: 0, G: 143, B: 251, A: 255}
		bgcolor := color.RGBA{R: bre, G: bgr, B: bbl, A: ba}
		background := canvas.NewRectangle(bgcolor)
		clockBackground = background // store reference for dynamic updates

		vbox := container.NewVBox()

		// Populate vbox with clock elements
		vbox.Add(nowtime)
		if showdate == 1 {
			vbox.Add(nowdate)
		}
		if showutc == 1 {
			vbox.Add(utctime)
		}

		// Always create all timezone text elements and add them to vbox
		// We'll show/hide them dynamically based on enabled state
		timezone1Text = canvas.NewText("", color.RGBA{R: ure, G: ugr, B: ubl, A: ua})
		timezone1Text.TextStyle = fyne.TextStyle{Bold: true}
		timezone1Text.Alignment = fyne.TextAlignCenter
		timezone1Text.TextSize = float32(utcsize)
		vbox.Add(timezone1Text)
		if !(timezone1Enabled == 1 && (timezone1Name != "" || timezone1Offset != "")) {
			timezone1Text.Hide()
		}

		timezone2Text = canvas.NewText("", color.RGBA{R: ure, G: ugr, B: ubl, A: ua})
		timezone2Text.TextStyle = fyne.TextStyle{Bold: true}
		timezone2Text.Alignment = fyne.TextAlignCenter
		timezone2Text.TextSize = float32(utcsize)
		vbox.Add(timezone2Text)
		if !(timezone2Enabled == 1 && (timezone2Name != "" || timezone2Offset != "")) {
			timezone2Text.Hide()
		}

		timezone3Text = canvas.NewText("", color.RGBA{R: ure, G: ugr, B: ubl, A: ua})
		timezone3Text.TextStyle = fyne.TextStyle{Bold: true}
		timezone3Text.Alignment = fyne.TextAlignCenter
		timezone3Text.TextSize = float32(utcsize)
		vbox.Add(timezone3Text)
		if !(timezone3Enabled == 1 && (timezone3Name != "" || timezone3Offset != "")) {
			timezone3Text.Hide()
		}

		timezone4Text = canvas.NewText("", color.RGBA{R: ure, G: ugr, B: ubl, A: ua})
		timezone4Text.TextStyle = fyne.TextStyle{Bold: true}
		timezone4Text.Alignment = fyne.TextAlignCenter
		timezone4Text.TextSize = float32(utcsize)
		vbox.Add(timezone4Text)
		if !(timezone4Enabled == 1 && (timezone4Name != "" || timezone4Offset != "")) {
			timezone4Text.Hide()
		}

		timezone5Text = canvas.NewText("", color.RGBA{R: ure, G: ugr, B: ubl, A: ua})
		timezone5Text.TextStyle = fyne.TextStyle{Bold: true}
		timezone5Text.Alignment = fyne.TextAlignCenter
		timezone5Text.TextSize = float32(utcsize)
		vbox.Add(timezone5Text)
		if !(timezone5Enabled == 1 && (timezone5Name != "" || timezone5Offset != "")) {
			timezone5Text.Hide()
		}

		// Always create weather text element (we'll show/hide it dynamically)
		weatherText = canvas.NewText("", color.RGBA{R: ure, G: ugr, B: ubl, A: ua})
		weatherText.TextStyle = fyne.TextStyle{Bold: true}
		weatherText.Alignment = fyne.TextAlignCenter
		weatherText.TextSize = float32(utcsize)
		vbox.Add(weatherText)
		if !weatherEnabled {
			weatherText.Hide()
		}

		content := container.NewStack(background, vbox)

		updateClock := func() {
			now = time.Now()
			if automute == 1 {
				if now.Hour() == muteonhr && now.Minute() == muteonmin {
					k := wallClockMinuteKey(now)
					if k != lastAutomuteOnKey {
						lastAutomuteOnKey = k
						muted, _ := volume.GetMuted()
						jiggleconf = jiggle
						jiggle = 0 // disable jiggle while muted
						lastJiggleMinute = -1
						if !muted {
							currentvolume, _ = volume.GetVolume()
							volume.Mute()
						}
					}
				} else if now.Hour() == muteoffhr && now.Minute() == muteoffmin {
					k := wallClockMinuteKey(now)
					if k != lastAutomuteOffKey {
						lastAutomuteOffKey = k
						muted, _ := volume.GetMuted()
						jiggle = jiggleconf // restore jiggle value
						lastJiggleMinute = -1
						jiggleconf = 0
						if muted {
							volume.Unmute()
							// volume.SetVolume(20)
							volume.SetVolume(currentvolume)
						}
					}
				}
			}
			if now.Minute() == 0 {
				if hourchime == 1 {
					// Top of the hour: play once per local hour (no second==0; update may be late)
					currentHour := now.Hour()
					if lastChimeHour != currentHour {
						lastChimeHour = currentHour
						// Cache file existence check - only recheck if file changed
						chimeFilePath := sndDir + "/" + hourchimesound
						if !hourChimeFileChecked || hourChimeCachedFile != hourchimesound {
							hourChimeFileExists = checkFileExists(chimeFilePath)
							hourChimeFileChecked = true
							hourChimeCachedFile = hourchimesound
						}
						if !hourChimeFileExists {
							playBeep("updown")
						} else {
							playMp3(chimeFilePath)
						}
					}
				}
			}

			nowtime.Text = now.Format(timeFormat)
			fyne.Do(func() {
				nowtime.Refresh()
				nowdate.Refresh()
				// if screen is not locked and jiggle is on and minute modulo jiggle ...
				if jiggle == 0 {
					lastJiggleMinute = -1
				}
				if !lock.IsScreenLocked() && jiggle > 0 && now.Minute()%jiggle == 0 {
					if now.Minute() != lastJiggleMinute {
						robotgo.MoveRelative(1, 0)  // MoveSmoothRelative(200, 0)
						robotgo.MoveRelative(0, 1)  // MoveSmoothRelative(0, 200)
						robotgo.MoveRelative(-1, 0) // MoveSmoothRelative(-200, 0)
						robotgo.MoveRelative(0, -1) // MoveSmoothRelative(0, -200)
						lastJiggleMinute = now.Minute()
					}
				}
			})
			nowdate.Text = now.Format(dateFormat)
			if showutc == 1 {
				utc := now.UTC()
				utctime.Text = utc.Format(utcFormat) + offsetString
				fyne.Do(func() {
					utctime.Refresh()
				})
			}

			// Helper function to parse UTC offset string (e.g., "+5", "-3.5") to hours and minutes
			parseUTCOffset := func(offsetStr string) (hours int, minutes int, valid bool) {
				offsetStr = strings.TrimSpace(offsetStr)
				if offsetStr == "" {
					return 0, 0, false
				}
				// Check sign first
				negative := strings.HasPrefix(offsetStr, "-")
				// Remove leading + or - if present
				if strings.HasPrefix(offsetStr, "+") || strings.HasPrefix(offsetStr, "-") {
					offsetStr = offsetStr[1:]
				}
				// Parse the offset
				var offsetHours float64
				var err error
				if strings.Contains(offsetStr, ".") || strings.Contains(offsetStr, ",") {
					offsetStr = strings.ReplaceAll(offsetStr, ",", ".")
					offsetHours, err = strconv.ParseFloat(offsetStr, 64)
				} else {
					var h int
					h, err = strconv.Atoi(offsetStr)
					offsetHours = float64(h)
				}
				if err != nil {
					return 0, 0, false
				}
				// Apply sign
				if negative {
					offsetHours = -offsetHours
				}
				// Convert to hours and minutes
				hours = int(offsetHours)
				minutes = int((offsetHours - float64(hours)) * 60)
				if minutes < 0 {
					minutes = -minutes
				}
				return hours, minutes, true
			}

			// Helper function to format timezone time (without UTC prefix)
			formatTimezoneTime := func(tzTime time.Time, tzName string) string {
				_, tzOffset := tzTime.Zone()
				tzOffsetHours := tzOffset / 3600
				tzOffsetMinutes := (tzOffset % 3600) / 60
				tzOffsetString := fmt.Sprintf(" (%s %+02d:%02d)", tzName, tzOffsetHours, tzOffsetMinutes)

				// Use non-UTC format for timezone time
				return tzTime.Format(nonutcFormat) + tzOffsetString
			}

			// Helper function to format timezone time from UTC offset
			formatTimezoneTimeFromOffset := func(tzTime time.Time, offsetHours int, offsetMinutes int, offsetLabel string) string {
				tzOffsetString := fmt.Sprintf(" (UTC%+02d:%02d)", offsetHours, offsetMinutes)
				if offsetLabel != "" {
					tzOffsetString = fmt.Sprintf(" (%s UTC%+02d:%02d)", offsetLabel, offsetHours, offsetMinutes)
				}
				return tzTime.Format(nonutcFormat) + tzOffsetString
			}

			// Update additional timezones
			// Priority: UTC offset > timezone name
			if timezone1Enabled == 1 && timezone1Text != nil {
				if timezone1Offset != "" {
					// Use UTC offset
					hours, minutes, valid := parseUTCOffset(timezone1Offset)
					if valid {
						utcTime := now.UTC()
						offsetDuration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute
						tzTime := utcTime.Add(offsetDuration)
						timezone1Text.Text = formatTimezoneTimeFromOffset(tzTime, hours, minutes, "Custom")
						fyne.Do(func() {
							timezone1Text.Refresh()
						})
					}
				} else if timezone1Name != "" {
					if loc, err := time.LoadLocation(timezone1Name); err == nil {
						tzTime := now.In(loc)
						timezone1Text.Text = formatTimezoneTime(tzTime, timezone1Name)
						fyne.Do(func() {
							timezone1Text.Refresh()
						})
					}
				}
			}
			if timezone2Enabled == 1 && timezone2Text != nil {
				if timezone2Offset != "" {
					// Use UTC offset
					hours, minutes, valid := parseUTCOffset(timezone2Offset)
					if valid {
						utcTime := now.UTC()
						offsetDuration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute
						tzTime := utcTime.Add(offsetDuration)
						timezone2Text.Text = formatTimezoneTimeFromOffset(tzTime, hours, minutes, "Custom")
						fyne.Do(func() {
							timezone2Text.Refresh()
						})
					}
				} else if timezone2Name != "" {
					if loc, err := time.LoadLocation(timezone2Name); err == nil {
						tzTime := now.In(loc)
						timezone2Text.Text = formatTimezoneTime(tzTime, timezone2Name)
						fyne.Do(func() {
							timezone2Text.Refresh()
						})
					}
				}
			}
			if timezone3Enabled == 1 && timezone3Text != nil {
				if timezone3Offset != "" {
					// Use UTC offset
					hours, minutes, valid := parseUTCOffset(timezone3Offset)
					if valid {
						utcTime := now.UTC()
						offsetDuration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute
						tzTime := utcTime.Add(offsetDuration)
						timezone3Text.Text = formatTimezoneTimeFromOffset(tzTime, hours, minutes, "Custom")
						fyne.Do(func() {
							timezone3Text.Refresh()
						})
					}
				} else if timezone3Name != "" {
					if loc, err := time.LoadLocation(timezone3Name); err == nil {
						tzTime := now.In(loc)
						timezone3Text.Text = formatTimezoneTime(tzTime, timezone3Name)
						fyne.Do(func() {
							timezone3Text.Refresh()
						})
					}
				}
			}
			if timezone4Enabled == 1 && timezone4Text != nil {
				if timezone4Offset != "" {
					// Use UTC offset
					hours, minutes, valid := parseUTCOffset(timezone4Offset)
					if valid {
						utcTime := now.UTC()
						offsetDuration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute
						tzTime := utcTime.Add(offsetDuration)
						timezone4Text.Text = formatTimezoneTimeFromOffset(tzTime, hours, minutes, "Custom")
						fyne.Do(func() {
							timezone4Text.Refresh()
						})
					}
				} else if timezone4Name != "" {
					if loc, err := time.LoadLocation(timezone4Name); err == nil {
						tzTime := now.In(loc)
						timezone4Text.Text = formatTimezoneTime(tzTime, timezone4Name)
						fyne.Do(func() {
							timezone4Text.Refresh()
						})
					}
				}
			}
			if timezone5Enabled == 1 && timezone5Text != nil {
				if timezone5Offset != "" {
					// Use UTC offset
					hours, minutes, valid := parseUTCOffset(timezone5Offset)
					if valid {
						utcTime := now.UTC()
						offsetDuration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute
						tzTime := utcTime.Add(offsetDuration)
						timezone5Text.Text = formatTimezoneTimeFromOffset(tzTime, hours, minutes, "Custom")
						fyne.Do(func() {
							timezone5Text.Refresh()
						})
					}
				} else if timezone5Name != "" {
					if loc, err := time.LoadLocation(timezone5Name); err == nil {
						tzTime := now.In(loc)
						timezone5Text.Text = formatTimezoneTime(tzTime, timezone5Name)
						fyne.Do(func() {
							timezone5Text.Refresh()
						})
					}
				}
			}

			// Update weather temperature display if enabled
			if weatherEnabled && weatherText != nil {
				if currentWeatherData != nil {
					// Format temperature: Fahrenheit first, then Celsius
					// Helper function to convert Celsius to Fahrenheit
					celsiusToFahrenheit := func(c float64) float64 {
						return c*9/5 + 32
					}
					tempF := celsiusToFahrenheit(currentWeatherData.CurrentTemp)
					tempC := currentWeatherData.CurrentTemp
					weatherText.Text = fmt.Sprintf("Weather: %.1f°F / %.1f°C", tempF, tempC)
				} else {
					weatherText.Text = "Weather: Loading..."
				}
				fyne.Do(func() {
					weatherText.Refresh()
				})
			}
		}

		updateClock()
		go func() {
			// When seconds are hidden, run the full clock update once per wall-clock minute
			// (first 1s tick that sees the new minute). Relying only on local second==0 misses
			// updates if the goroutine wakes late (sleep, scheduling), which breaks jiggle/chimes.
			lastSlowClockMinuteEpoch := time.Now().Unix() / 60
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					now = time.Now()
					if showseconds == 1 {
						updateClock()
					} else {
						minuteEpoch := now.Unix() / 60
						if minuteEpoch != lastSlowClockMinuteEpoch {
							lastSlowClockMinuteEpoch = minuteEpoch
							updateClock()
						}
					}
					// lock screen / mute volume event handler, but only if enabled
					// and only unmute if we auto muted. If user had already muted, don't
					if slockmute == 1 {
						if lock.IsScreenLocked() {
							muted, _ := volume.GetMuted()
							if !muted {
								clockmutedvol = 1
								volume.Mute()
							}
						} else {
							lockmuted, _ := volume.GetMuted()
							if lockmuted && clockmutedvol == 1 {
								clockmutedvol = 0
								volume.Unmute()
							}
						}
					}
				case <-clockUpdateLoopStop:
					return
				}
			}
		}()

		clock.SetContent(content)
		clock.SetCloseIntercept(func() {
			// Stop the update loop gracefully
			if clockUpdateLoopRunning && clockUpdateLoopStop != nil {
				close(clockUpdateLoopStop)
				clockUpdateLoopStop = nil
				clockUpdateLoopRunning = false
			}
			clock.Close()
			clock = nil
			clockTimeText = nil
			clockDateText = nil
			clockUtcText = nil
			clockBackground = nil
			timezone1Text = nil
			timezone2Text = nil
			timezone3Text = nil
			timezone4Text = nil
			timezone5Text = nil
			weatherText = nil
		})
		clock.Resize(fyne.NewSize(content.MinSize().Width*1.2, content.MinSize().Height*1.1))
		// clock.Resize(fyne.NewSize(300, 200))
		// clock.ShowAndRun()  // for standalone clock app
		clock.CenterOnScreen()
		clock.Show() // for func inside KrankyBearTimer
		// Bring clock to front immediately
		clock.RequestFocus()
	}
}

// populateClockVBox removed - not used since we removed dynamic rebuilding

// rebuildClockContent removed - not used since we removed dynamic rebuilding

// updateClockColors dynamically updates the clock window colors
func updateClockColors() {
	if clock == nil {
		// Try to open the clock if it's not open
		if a := fyne.CurrentApp(); a != nil {
			desktopclock(a)
		}
		return
	}

	// Safely check if window elements exist - don't access window methods that might crash
	if clockTimeText == nil && clockDateText == nil && clockUtcText == nil && clockBackground == nil {
		return
	}

	// Parse time color
	var tre, tgr, tbl, ta uint8
	colors := strings.Split(timecolor, ",")
	if len(colors) == 4 {
		col, _ := strconv.ParseUint(colors[0], 10, 8)
		tre = uint8(col)
		col, _ = strconv.ParseUint(colors[1], 10, 8)
		tgr = uint8(col)
		col, _ = strconv.ParseUint(colors[2], 10, 8)
		tbl = uint8(col)
		col, _ = strconv.ParseUint(colors[3], 10, 8)
		ta = uint8(col)
	}

	// Parse background color
	var bre, bgr, bbl, ba uint8
	colors = strings.Split(bgcolor, ",")
	if len(colors) == 4 {
		col, _ := strconv.ParseUint(colors[0], 10, 8)
		bre = uint8(col)
		col, _ = strconv.ParseUint(colors[1], 10, 8)
		bgr = uint8(col)
		col, _ = strconv.ParseUint(colors[2], 10, 8)
		bbl = uint8(col)
		col, _ = strconv.ParseUint(colors[3], 10, 8)
		ba = uint8(col)
	}

	// Parse date color
	var dre, dgr, dbl, da uint8
	colors = strings.Split(datecolor, ",")
	if len(colors) == 4 {
		col, _ := strconv.ParseUint(colors[0], 10, 8)
		dre = uint8(col)
		col, _ = strconv.ParseUint(colors[1], 10, 8)
		dgr = uint8(col)
		col, _ = strconv.ParseUint(colors[2], 10, 8)
		dbl = uint8(col)
		col, _ = strconv.ParseUint(colors[3], 10, 8)
		da = uint8(col)
	}

	// Parse UTC color
	var ure, ugr, ubl, ua uint8
	colors = strings.Split(utccolor, ",")
	if len(colors) == 4 {
		col, _ := strconv.ParseUint(colors[0], 10, 8)
		ure = uint8(col)
		col, _ = strconv.ParseUint(colors[1], 10, 8)
		ugr = uint8(col)
		col, _ = strconv.ParseUint(colors[2], 10, 8)
		ubl = uint8(col)
		col, _ = strconv.ParseUint(colors[3], 10, 8)
		ua = uint8(col)
	}

	// Update colors - wrap in panic recovery to handle invalid windows
	fyne.Do(func() {
		defer func() {
			if r := recover(); r != nil {
				// Window was invalid, clear all references
				clock = nil
				clockTimeText = nil
				clockDateText = nil
				clockUtcText = nil
				clockBackground = nil
			}
		}()

		// Double check clock is still valid
		if clock == nil {
			return
		}

		updated := false
		if clockTimeText != nil {
			clockTimeText.Color = color.RGBA{R: tre, G: tgr, B: tbl, A: ta}
			clockTimeText.Refresh()
			updated = true
		}
		if clockDateText != nil {
			clockDateText.Color = color.RGBA{R: dre, G: dgr, B: dbl, A: da}
			clockDateText.Refresh()
			updated = true
		}
		if clockUtcText != nil {
			clockUtcText.Color = color.RGBA{R: ure, G: ugr, B: ubl, A: ua}
			clockUtcText.Refresh()
			updated = true
		}
		// Update timezone text colors
		if timezone1Text != nil {
			timezone1Text.Color = color.RGBA{R: ure, G: ugr, B: ubl, A: ua}
			timezone1Text.Refresh()
			updated = true
		}
		if timezone2Text != nil {
			timezone2Text.Color = color.RGBA{R: ure, G: ugr, B: ubl, A: ua}
			timezone2Text.Refresh()
			updated = true
		}
		if timezone3Text != nil {
			timezone3Text.Color = color.RGBA{R: ure, G: ugr, B: ubl, A: ua}
			timezone3Text.Refresh()
			updated = true
		}
		if timezone4Text != nil {
			timezone4Text.Color = color.RGBA{R: ure, G: ugr, B: ubl, A: ua}
			timezone4Text.Refresh()
			updated = true
		}
		if timezone5Text != nil {
			timezone5Text.Color = color.RGBA{R: ure, G: ugr, B: ubl, A: ua}
			timezone5Text.Refresh()
			updated = true
		}
		if weatherText != nil {
			weatherText.Color = color.RGBA{R: tre, G: tgr, B: tbl, A: ta}
			weatherText.Refresh()
			updated = true
		}
		if clockBackground != nil {
			clockBackground.FillColor = color.RGBA{R: bre, G: bgr, B: bbl, A: ba}
			clockBackground.Refresh()
			updated = true
		}
		// Refresh the entire clock window content to ensure changes are visible
		if updated && clock != nil {
			clock.Content().Refresh()
			// Also refresh the canvas to ensure visual update
			if clock.Canvas() != nil {
				clock.Canvas().Refresh(clock.Content())
			}
		}
	})

	// Also update countdown, alarm, and weather windows if they're open
	updateCountdownColors()
	updateAlarmColors()
	updateWeatherColors()
}

// updateTimezoneVisibility shows/hides timezone elements based on enabled state
func updateTimezoneVisibility() {
	if clock == nil {
		return
	}

	fyne.Do(func() {
		defer func() {
			if r := recover(); r != nil {
				// If window is invalid, just return
			}
		}()

		if timezone1Text != nil {
			if timezone1Enabled == 1 && (timezone1Name != "" || timezone1Offset != "") {
				timezone1Text.Show()
			} else {
				timezone1Text.Hide()
			}
		}
		if timezone2Text != nil {
			if timezone2Enabled == 1 && (timezone2Name != "" || timezone2Offset != "") {
				timezone2Text.Show()
			} else {
				timezone2Text.Hide()
			}
		}
		if timezone3Text != nil {
			if timezone3Enabled == 1 && (timezone3Name != "" || timezone3Offset != "") {
				timezone3Text.Show()
			} else {
				timezone3Text.Hide()
			}
		}
		if timezone4Text != nil {
			if timezone4Enabled == 1 && (timezone4Name != "" || timezone4Offset != "") {
				timezone4Text.Show()
			} else {
				timezone4Text.Hide()
			}
		}
		if timezone5Text != nil {
			if timezone5Enabled == 1 && (timezone5Name != "" || timezone5Offset != "") {
				timezone5Text.Show()
			} else {
				timezone5Text.Hide()
			}
		}
		if weatherText != nil {
			if weatherEnabled {
				weatherText.Show()
			} else {
				weatherText.Hide()
			}
		}

		// Refresh the clock content to update layout
		if clock != nil && clock.Content() != nil {
			clock.Content().Refresh()
		}
	})
}

// updateTimezoneText immediately updates timezone text elements when timezone name or offset changes
func updateTimezoneText() {
	if clock == nil {
		return
	}

	now := time.Now()

	// Helper function to parse UTC offset string (e.g., "+5", "-3.5") to hours and minutes
	parseUTCOffset := func(offsetStr string) (hours int, minutes int, valid bool) {
		offsetStr = strings.TrimSpace(offsetStr)
		if offsetStr == "" {
			return 0, 0, false
		}
		// Check sign first
		negative := strings.HasPrefix(offsetStr, "-")
		// Remove leading + or - if present
		if strings.HasPrefix(offsetStr, "+") || strings.HasPrefix(offsetStr, "-") {
			offsetStr = offsetStr[1:]
		}
		// Parse the offset
		var offsetHours float64
		var err error
		if strings.Contains(offsetStr, ".") || strings.Contains(offsetStr, ",") {
			offsetStr = strings.ReplaceAll(offsetStr, ",", ".")
			offsetHours, err = strconv.ParseFloat(offsetStr, 64)
		} else {
			var h int
			h, err = strconv.Atoi(offsetStr)
			offsetHours = float64(h)
		}
		if err != nil {
			return 0, 0, false
		}
		// Apply sign
		if negative {
			offsetHours = -offsetHours
		}
		// Convert to hours and minutes
		hours = int(offsetHours)
		minutes = int((offsetHours - float64(hours)) * 60)
		if minutes < 0 {
			minutes = -minutes
		}
		return hours, minutes, true
	}

	// Helper function to format timezone time (without UTC prefix)
	formatTimezoneTime := func(tzTime time.Time, tzName string) string {
		_, tzOffset := tzTime.Zone()
		tzOffsetHours := tzOffset / 3600
		tzOffsetMinutes := (tzOffset % 3600) / 60
		tzOffsetString := fmt.Sprintf(" (%s %+02d:%02d)", tzName, tzOffsetHours, tzOffsetMinutes)

		nonutcFormat := `3:04 PM`
		return tzTime.Format(nonutcFormat) + tzOffsetString
	}

	// Helper function to format timezone time from UTC offset
	formatTimezoneTimeFromOffset := func(tzTime time.Time, offsetHours int, offsetMinutes int, offsetLabel string) string {
		tzOffsetString := fmt.Sprintf(" (UTC%+02d:%02d)", offsetHours, offsetMinutes)
		if offsetLabel != "" {
			tzOffsetString = fmt.Sprintf(" (%s UTC%+02d:%02d)", offsetLabel, offsetHours, offsetMinutes)
		}
		nonutcFormat := `3:04 PM`
		return tzTime.Format(nonutcFormat) + tzOffsetString
	}

	// Update timezone text elements if they exist and are enabled
	// Priority: UTC offset > timezone name
	fyne.Do(func() {
		if timezone1Enabled == 1 && timezone1Text != nil {
			if timezone1Offset != "" {
				// Use UTC offset
				hours, minutes, valid := parseUTCOffset(timezone1Offset)
				if valid {
					utcTime := now.UTC()
					offsetDuration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute
					tzTime := utcTime.Add(offsetDuration)
					timezone1Text.Text = formatTimezoneTimeFromOffset(tzTime, hours, minutes, "Custom")
					timezone1Text.Refresh()
				}
			} else if timezone1Name != "" {
				if loc, err := time.LoadLocation(timezone1Name); err == nil {
					tzTime := now.In(loc)
					timezone1Text.Text = formatTimezoneTime(tzTime, timezone1Name)
					timezone1Text.Refresh()
				}
			}
		}
		if timezone2Enabled == 1 && timezone2Text != nil {
			if timezone2Offset != "" {
				// Use UTC offset
				hours, minutes, valid := parseUTCOffset(timezone2Offset)
				if valid {
					utcTime := now.UTC()
					offsetDuration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute
					tzTime := utcTime.Add(offsetDuration)
					timezone2Text.Text = formatTimezoneTimeFromOffset(tzTime, hours, minutes, "Custom")
					timezone2Text.Refresh()
				}
			} else if timezone2Name != "" {
				if loc, err := time.LoadLocation(timezone2Name); err == nil {
					tzTime := now.In(loc)
					timezone2Text.Text = formatTimezoneTime(tzTime, timezone2Name)
					timezone2Text.Refresh()
				}
			}
		}
		if timezone3Enabled == 1 && timezone3Text != nil {
			if timezone3Offset != "" {
				// Use UTC offset
				hours, minutes, valid := parseUTCOffset(timezone3Offset)
				if valid {
					utcTime := now.UTC()
					offsetDuration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute
					tzTime := utcTime.Add(offsetDuration)
					timezone3Text.Text = formatTimezoneTimeFromOffset(tzTime, hours, minutes, "Custom")
					timezone3Text.Refresh()
				}
			} else if timezone3Name != "" {
				if loc, err := time.LoadLocation(timezone3Name); err == nil {
					tzTime := now.In(loc)
					timezone3Text.Text = formatTimezoneTime(tzTime, timezone3Name)
					timezone3Text.Refresh()
				}
			}
		}
		if timezone4Enabled == 1 && timezone4Text != nil {
			if timezone4Offset != "" {
				// Use UTC offset
				hours, minutes, valid := parseUTCOffset(timezone4Offset)
				if valid {
					utcTime := now.UTC()
					offsetDuration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute
					tzTime := utcTime.Add(offsetDuration)
					timezone4Text.Text = formatTimezoneTimeFromOffset(tzTime, hours, minutes, "Custom")
					timezone4Text.Refresh()
				}
			} else if timezone4Name != "" {
				if loc, err := time.LoadLocation(timezone4Name); err == nil {
					tzTime := now.In(loc)
					timezone4Text.Text = formatTimezoneTime(tzTime, timezone4Name)
					timezone4Text.Refresh()
				}
			}
		}
		if timezone5Enabled == 1 && timezone5Text != nil {
			if timezone5Offset != "" {
				// Use UTC offset
				hours, minutes, valid := parseUTCOffset(timezone5Offset)
				if valid {
					utcTime := now.UTC()
					offsetDuration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute
					tzTime := utcTime.Add(offsetDuration)
					timezone5Text.Text = formatTimezoneTimeFromOffset(tzTime, hours, minutes, "Custom")
					timezone5Text.Refresh()
				}
			} else if timezone5Name != "" {
				if loc, err := time.LoadLocation(timezone5Name); err == nil {
					tzTime := now.In(loc)
					timezone5Text.Text = formatTimezoneTime(tzTime, timezone5Name)
					timezone5Text.Refresh()
				}
			}
		}
	})
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942

// To-do:

// a few notes, format specific
// timeFormat := `3:04:05 PM (MST)`
// clock.SetText(now.Format("Mon Jan 2 15:04:05 2006"))
// clock.SetText(now.Format("15:04:05`nMonday, January 2, 2006"))

// show seconds
// clockFormat := `15:04:05
//Monday, January 2, 2006`

// no show seconds - not always valuable when we update every second
// anyway, but still - user preference ...
// clockFormat := `15:04`
//clockFormat := `15:04
//   Monday, January 2, 2006`
//clock.SetText(now.Format(clockFormat))
//clock.Alignment = fyne.TextAlignCenter
