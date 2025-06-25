//go:generate fyne bundle -o bundled.go -a Resources/Images/KrankyBear.png
//go:generate fyne bundle -o bundled.go -a Resources/Images/KrankyBearBeret.png
//go:generate fyne bundle -o bundled.go -a Resources/Images/KrankyBearBeanieMultiColor.png
//go:generate fyne bundle -o bundled.go -a Resources/Images/KrankyBearChristmas.png
//go:generate fyne bundle -o bundled.go -a Resources/Images/KrankyBearChristmasGrinch.png
//go:generate fyne bundle -o bundled.go -a Resources/Images/KrankyBearCowboyBrown.png
//go:generate fyne bundle -o bundled.go -a Resources/Images/KrankyBearFedoraRed.png
//go:generate fyne bundle -o bundled.go -a Resources/Images/KrankyBearHardHat.png
//go:generate fyne bundle -o bundled.go -a Resources/Images/KrankyBearHogwartsSorting.png
//go:generate fyne bundle -o bundled.go -a Resources/Images/KrankyBearTrapperRedPlaid.png
//go:generate fyne bundle -o bundled.go -a Resources/Images/KrankyBearVikingHelmet.png
//go:generate fyne bundle -o bundled.go -a Resources/Images/http418.png

//go:generate fyne bundle -o bundled.go -a Resources/Images/SchoolBoard1.png
//go:generate fyne bundle -o bundled.go -a Resources/Images/SchoolBoard2.png
//go:generate fyne bundle -o bundled.go -a Resources/Images/SchoolBoard3.png
//go:generate fyne bundle -o bundled.go -a Resources/Images/SchoolBoard4.png
//go:generate fyne bundle -o bundled.go -a Resources/Images/SchoolBoard5.png
//go:generate fyne bundle -o bundled.go -a Resources/Images/SchoolBoard6.png

///go:generate fyne bundle -o bundled.go -a Resources/Images/TBlue.png
///go:generate fyne bundle -o bundled.go -a Resources/Images/TStone.png
///go:generate fyne bundle -o bundled.go -a Resources/Images/TAlmond.png
///go:generate fyne bundle -o bundled.go -a Resources/Images/TGrayTeach.png
///go:generate fyne bundle -o bundled.go -a Resources/Images/Tcn256.png
///go:generate fyne bundle -o bundled.go -a Resources/Images/Tcc256.png
///go:generate fyne bundle -o bundled.go -a Resources/Images/Tcbe256.png

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
