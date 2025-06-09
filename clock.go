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
	"github.com/itchyny/volume-go"
)

// var clock fyne.Window

func desktopclock(a fyne.App) { // , w fyne.Window, bg fyne.Canvas) {
	if clock != nil { // &&  !clock.Content().Visible() {
		clock.RequestFocus()
	} else {
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
		clock.SetIcon(resourceKrankyBearPng)

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
		dateFormat := ` Monday, January 2, 2006 `

		// nowtime := canvas.NewText(now.Format(timeFormat), color.RGBA{R: 255, G: 123, B: 31, A: 255})
		nowtime := canvas.NewText(now.Format(timeFormat), color.RGBA{R: tre, G: tgr, B: tbl, A: ta})
		nowtime.TextStyle = fyne.TextStyle{Bold: true}
		// nowtime.TextStyle = fyne.TextStyle{Monospace: true} // EXAMPLE FONT TYPE
		nowtime.Alignment = fyne.TextAlignCenter
		nowtime.TextSize = float32(timesize)

		// utctime := canvas.NewText(now.Format(utcFormat), color.RGBA{R: 255, G: 123, B: 31, A: 255})
		utctime := canvas.NewText(now.Format(utcFormat), color.RGBA{R: ure, G: ugr, B: ubl, A: ua})
		utctime.TextStyle = fyne.TextStyle{Bold: true}
		utctime.Alignment = fyne.TextAlignCenter
		utctime.TextSize = float32(utcsize)

		// nowdate := canvas.NewText(now.Format(dateFormat), color.RGBA{R: 208, G: 145, B: 38, A: 255})
		nowdate := canvas.NewText(now.Format(dateFormat), color.RGBA{R: dre, G: dgr, B: dbl, A: da})
		nowdate.TextStyle = fyne.TextStyle{Bold: true}
		nowdate.Alignment = fyne.TextAlignCenter
		nowdate.TextSize = float32(datesize)

		//background := canvas.NewRectangle(color.RGBA{R: 0, G: 0, B: 255, A: 255})
		// bgcolor := color.RGBA{R: 0, G: 143, B: 251, A: 255}
		bgcolor := color.RGBA{R: bre, G: bgr, B: bbl, A: ba}
		background := canvas.NewRectangle(bgcolor)

		vbox := container.NewVBox()
		if showutc == 1 {
			if showdate == 1 {
				vbox = container.NewVBox(nowtime, nowdate, utctime)
			} else {
				vbox = container.NewVBox(nowtime, utctime)
			}
		} else {
			if showdate == 1 {
				vbox = container.NewVBox(nowtime, nowdate)
			} else {
				vbox = container.NewVBox(nowtime)
			}
		}
		content := container.NewStack(background, vbox)

		updateClock := func() {
			now = time.Now()
			if now.Hour() == muteonhr && now.Minute() == muteonmin && now.Second() == 0 {
				if automute == 1 {
					muted, _ := volume.GetMuted()
					if !muted {
						currentvolume, _ = volume.GetVolume()
						volume.Mute()
					}
				}
			} else if now.Hour() == muteoffhr && now.Minute() == muteoffmin && now.Second() == 0 {
				if automute == 1 {
					muted, _ := volume.GetMuted()
					if muted {
						volume.Unmute()
						// volume.SetVolume(20)
						volume.SetVolume(currentvolume)
					}
				}
			}
			if now.Minute() == 0 && now.Second() == 0 {
				if hourchime == 1 {
					if !checkFileExists(sndDir + "/" + hourchimesound) {
						playBeep("updown")
					} else {
						playMp3(sndDir + "/" + hourchimesound)
					}
				}
			}

			nowtime.Text = now.Format(timeFormat)
			fyne.Do(func() {
				nowtime.Refresh()
				nowdate.Refresh()
			})
			nowdate.Text = now.Format(dateFormat)
			if showutc == 1 {
				utc := now.UTC()
				utctime.Text = utc.Format(utcFormat) + offsetString
				fyne.Do(func() {
					utctime.Refresh()
				})
			}
		}

		updateClock()
		go func() {
			for range time.Tick(time.Second) {
				// updating frequently is something of a resource hog (CPU)
				// check here if seconds are displayed, update
				// if seconds are not displayed, check for seconds == 0
				// at the minute change, and only update the clock then
				now = time.Now()
				if showseconds == 1 || now.Second() == 0 {
					updateClock()
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
			}
		}()

		clock.SetContent(content)
		clock.SetCloseIntercept(func() {
			clock.Close()
			clock = nil
		})
		clock.Resize(fyne.NewSize(content.MinSize().Width*1.2, content.MinSize().Height*1.1))
		// clock.Resize(fyne.NewSize(300, 200))
		// clock.ShowAndRun()  // for standalone clock app
		clock.Show() // for func inside KrankyBearTimer
	}
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
