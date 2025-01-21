package main

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

var c fyne.Window

func clock(a fyne.App) { // , w fyne.Window, bg fyne.Canvas) {
	if c != nil { // &&  !c.Content().Visible() {
		c.RequestFocus()
	} else {
		c = a.NewWindow("Tanium Clock")

		now := time.Now()
		// timeFormat := `15:04:05`
		timeFormat := `3:04:05 PM (MST)`
		utcFormat := `(UTC 3:04 PM Z07)`
		dateFormat := ` Monday, January 2, 2006 `
		// Create the first label
		nowtime := canvas.NewText(now.Format(timeFormat), color.RGBA{R: 255, G: 0, B: 0, A: 255})
		nowtime.TextStyle = fyne.TextStyle{Bold: true}
		nowtime.Alignment = fyne.TextAlignCenter
		// nowtime = canvas.TextStyle{Alignment: canvas.AlignmentCenter}
		nowtime.TextSize = 48
		utctime := canvas.NewText(now.Format(utcFormat), color.RGBA{R: 243, G: 119, B: 53, A: 255})
		utctime.TextStyle = fyne.TextStyle{Bold: true}
		utctime.Alignment = fyne.TextAlignCenter

		// Create the second label
		nowdate := canvas.NewText(now.Format(dateFormat), color.RGBA{R: 243, G: 119, B: 53, A: 255})
		nowdate.TextStyle = fyne.TextStyle{Bold: true}
		nowdate.Alignment = fyne.TextAlignCenter
		nowdate.TextSize = 24

		// Create a rectangle with blue background
		background := canvas.NewRectangle(color.RGBA{R: 0, G: 0, B: 255, A: 255})

		// Create a container to hold the labels and background
		content := container.NewStack(background, container.NewVBox(nowtime, nowdate, utctime))

		/**/
		updateClock := func() {
			now = time.Now()
			nowtime.Text = now.Format(timeFormat)
			utctime.Text = now.Format(utcFormat)
			nowdate.Text = now.Format(dateFormat)
			nowtime.Refresh()
			utctime.Refresh()
			nowdate.Refresh()
		}

		updateClock()
		go func() {
			for range time.Tick(time.Second) {
				updateClock()
			}
		}()
		/**/

		c.SetContent(content)
		c.SetCloseIntercept(func() {
			c.Close()
			c = nil
		})
		c.Resize(fyne.NewSize(content.MinSize().Width*1.2, content.MinSize().Height*1.1))
		// c.Resize(fyne.NewSize(300, 200))
		// c.ShowAndRun()  // for standalone clock app
		c.Show() // for func inside TaniumTimer
	}
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942

// To-do:
// allow background color choice
// allow clock color choice and text size
// allow date color choice and text size
// allow display / hide UTC time
// allow 12 / 24 hour clock display format
// allow display seconds or not
// allow displaying timezone or not

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
