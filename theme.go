//go:generate fyne bundle -o bundled.go -a Resources/Images/TaniumTimer.svg
//go:generate fyne bundle -o bundled.go -a Resources/Images/TaniumTimer2.svg
//go:generate fyne bundle -o bundled.go -a Resources/Images/TaniumIcon.svg
//go:generate fyne bundle -o bundled.go -a Resources/Images/TaniumTrayIcon.svg
//go:generate fyne bundle -o bundled.go -a Resources/Images/TaniumTrainBlue.svg
//go:generate fyne bundle -o bundled.go -a Resources/Images/TaniumTrainStone.svg
//go:generate fyne bundle -o bundled.go -a Resources/Images/TaniumTrainAlmond.svg
//go:generate fyne bundle -o bundled.go -a Resources/Images/TaniumConverge2024.svg
//go:generate fyne bundle -o bundled.go -a Resources/Images/TaniumConverge2024a.svg
///go:generate fyne bundle -o bundled.go -a Resources/Images/redSettingsGear.svg

///go:generate fyne bundle -o bundled.go -a Resources/Sounds/boing.mp3
///go:generate fyne bundle -o bundled.go -a Resources/Sounds/Basso.mp3
///go:generate fyne bundle -o bundled.go -a Resources/Sounds/Sosumi.mp3
///go:generate fyne bundle -o bundled.go -a Resources/Sounds/Submarine.mp3
///go:generate fyne bundle -o bundled.go -a Resources/Sounds/baseball.mp3
///go:generate fyne bundle -o bundled.go -a Resources/Sounds/grandfatherClock.mp3
///go:generate fyne bundle -o bundled.go -a Resources/Sounds/pinball.mp3

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type appTheme struct {
	fyne.Theme
}

func (a *appTheme) Size(n fyne.ThemeSizeName) float32 {
	if n == theme.SizeNameHeadingText {
		return a.Theme.Size(n) * 1.5
	}

	return a.Theme.Size(n)
}

func themeTimer(text *widget.RichText, time int) {
	seg := text.Segments[0].(*widget.TextSegment)
	if time <= 30 {
		seg.Style.ColorName = theme.ColorNameError
	} else if time < 150 {
		seg.Style.ColorName = theme.ColorNameWarning
	} else {
		seg.Style.ColorName = theme.ColorNameSuccess
	}
	text.Refresh()
}
