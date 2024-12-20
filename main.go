package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"fyne.io/systray"
)

// import "C" - was to be used for stay on top, processes, but not now

const (
	timerName      = "Tanium Timer"
	timerVersion   = "0.5" // see FyneApp.toml
	timerCopyright = "(c) Tanium, 2024"
	timerAuthor    = "Allan Marillier"
)

var running = binding.NewBool()
var bg fyne.Canvas
var remain int
var notify int
var sound int

var adhocTime int
var biobreakTime int
var lunchTime int

var imgDir string
var timerbg string
var starttimer int

var sndDir string
var endsnd string
var oneminsnd string
var halfminsnd string
var debug int = 0
var abt fyne.Window

// preferences stored via fyne preferences API land in
// ~/Library/Preferences/fyne/com.tanium.taniumtimer/preferences.json
// ~\AppData\Roaming\fyne\com.tanium.taniumtimer\preferences.json
// {"adhoc.default":300,"background.default":"blue","biobreak.default":600,"endsound.default":"baseball.mp3","halfminsound.default":"sosumi.mp3","lunch.default":3600,"notify.default":1,"oneminsound.default":"hero.mp3", "sound.default":1}

func main() {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	launchDir := filepath.Dir(exePath)

	if runtime.GOOS == "darwin" {
		if strings.HasPrefix(launchDir, "/Applications/TaniumTimer") {
			sndDir = launchDir + "/../Resources/Sounds"
			imgDir = launchDir + "/../Resources/Images"
		} else {
			sndDir = launchDir + "/Resources/Sounds"
			imgDir = launchDir + "/Resources/Images"
		}
	} else if runtime.GOOS == "windows" {
		sndDir = launchDir + "/Resources/Sounds"
		imgDir = launchDir + "/Resources/Images"
	}

	a := app.NewWithID("com.tanium.TaniumTimer")
	a.Settings().SetTheme(&appTheme{Theme: theme.DefaultTheme()})
	w := a.NewWindow(timerName)
	w.SetPadded(false)
	w.SetCloseIntercept(func() {
		a.Quit() // force quit, normal when somebody hits "x" to close
	})
	w.SetMaster()      // this sets this as master and closes all child windows
	w.CenterOnScreen() // run centered on primary (laptop) display

	prefs := strings.ReplaceAll((a.Storage().RootURI()).String(), "file://", "") + "/preferences.json"
	if !checkFileExists(prefs) {
		if debug == 1 {
			log.Println("prefs file does not exist")
		}
		// add some default prefs that can be modified via settings
		writeDefaultSettings(a)
	}
	// get default settings from preferences
	lunchTime = a.Preferences().IntWithFallback("lunch.default", 60*60)
	adhocTime = a.Preferences().IntWithFallback("adhoc.default", 5*60)
	biobreakTime = a.Preferences().IntWithFallback("biobreak.default", 10*60)
	notify = a.Preferences().IntWithFallback("notify.default", 1)
	sound = a.Preferences().IntWithFallback("sound.default", 1)
	timerbg = a.Preferences().StringWithFallback("background.default", "blue")
	endsnd = a.Preferences().StringWithFallback("endsound.default", "baseball.mp3")
	oneminsnd = a.Preferences().StringWithFallback("oneminsound.default", "hero.mp3")
	halfminsnd = a.Preferences().StringWithFallback("halfminsound.default", "sosumi.mp3")
	starttimer = a.Preferences().IntWithFallback("starttimer.default", 0)
	writeSettings(a)

	if len(os.Args) >= 2 {
		log.Println("arg count:", len(os.Args))
		if os.Args[1] == "debug" || os.Args[1] == "d" {
			debug = 1
			logInit()
			r, _ := os.Open("TaniumTimer0.txt")
			logLines, _ := lineCounter(r)
			r.Close()
			InfoLog.Println("logLines:", logLines)
			if logLines >= 100 {
				logRotate()
			}
			logInit()
			InfoLog.Println("Opening the application...")
			InfoLog.Println("Something has occurred...")
			WarningLog.Println("WARNING!!!..")
			ErrorLog.Println("Some error has occurred...")

			log.Println("debug mode:", debug)
			log.Println("exepath:", exePath)
			log.Println("launchdir:", launchDir)
			log.Println("Images:", imgDir)
			log.Println("Sounds:", sndDir)
			log.Println("endsnd:", endsnd)
			log.Println("oneminsnd:", oneminsnd)
			log.Println("halfminsnd:", halfminsnd)
			log.Println("starttimer:", starttimer)
			adhocTime = 65 // debug value - short for easy test
		}
	}

	timer := widget.NewRichText()
	updateTime(timer, adhocTime)

	if desk, ok := a.(desktop.App); ok {
		desk.SetSystemTrayIcon(resourceTaniumIconSvg)
		systray.SetTooltip(timerName)
		// systray.SetTitle(timerName)
		show := fyne.NewMenuItem("Show", func() {
			w.Show()
			w.Canvas().Focused()
		})
		hide := fyne.NewMenuItem("Hide", w.Hide)
		lunch := fyne.NewMenuItem("Lunch ("+strconv.Itoa(lunchTime/60)+")", func() {
			startTimer(lunchTime, "Lunch", w.Canvas(), w)
		})
		biobreak := fyne.NewMenuItem("Bio Break ("+strconv.Itoa(biobreakTime/60)+")", func() {
			startTimer(biobreakTime, "Bio Break", w.Canvas(), w)
		})
		adhoc := fyne.NewMenuItem("Ad Hoc", func() {
			startTimer(adhocTime, "Ad Hoc Timer", w.Canvas(), w)
		})
		stop := fyne.NewMenuItem("Stop", func() {
			remain = -1 // don't notify when the user stops it
		})
		about := fyne.NewMenuItem("About", func() {
			aboutText := timerName + " v " + timerVersion
			aboutText += "\n" + timerCopyright
			aboutText += "\n\n" + timerAuthor + ", using Go and fyne GUI"

			if abt == nil || !abt.Content().Visible() {
				abt = a.NewWindow(timerName + ": About")
				abt.Resize(fyne.NewSize(50, 100))
				abt.SetContent(widget.NewLabel(aboutText))
				abt.SetCloseIntercept(func() {
					abt.Close()
					abt = nil
				})
				abt.CenterOnScreen() // run centered on primary (laptop) display
				abt.Show()
			} else {
				// abt.Show()
				// abt = nil
				certs := []fyne.Resource{resourceTcnPng, resourceTccPng, resourceTcbePng}
				rand.Seed(time.Now().UnixNano())
				randomIndex := rand.Intn(len(certs))
				egg := a.NewWindow(timerName + ": easter egg")
				eggimage := canvas.NewImageFromResource(certs[randomIndex])
				// eggimage := canvas.NewImageFromResource(resourceTCNSvg)
				// eggimage := canvas.NewImageFromResource(resourceTcnPng)

				eggimage.FillMode = canvas.ImageFillOriginal
				text := "Whoo-hoo! You found the Easter egg!\n"
				text += "\n" + dadjoke()
				eggtext := widget.NewLabel(text)
				content := container.NewVBox(eggimage, eggtext)
				egg.SetContent(content)
				egg.CenterOnScreen() // run centered on primary (laptop) display
				for j := 0; j <= 2; j++ {
					playBeep("down")
					egg.Show()
					time.Sleep(time.Second / 3)
					egg.Hide()
					time.Sleep(time.Second / 3)
				}
				egg.Show()
			}
		})
		help := fyne.NewMenuItem("Help", func() {
			hlp := a.NewWindow(timerName + ": Help")
			hlpText := `More help will be added later
For now we're adding as we go:
- Ad hoc timer minimum is 5 minutes, with 5 minute increments
	- NOTE ad hoc default is updated in preferences to current value any time it is changed
- Bio break timer default is 10 minutes
- Lunch break timer default is 60 minutes
- Each of these break times can be modified using Settings, set in minutes
- Timer text color is green until 2 1/2 minutes remain,
	- color is orange from 2 1/2 minutes to 30 seconds
	- color is red from 30 seconds to completion
- optional setting to enable auto starting at boot

- System tray notifications and sound alerts are both optional, enabled by default
- Tone / beep alerts are at 60 seconds, at 30 seconds, and at completion
- Timer window flashes on/off at timer end (in addition to desktop notification & beep)
- A timer that has been hidden behind another window or minimized will be
	brought to the front / focused at 60 seconds, and at timer completion

- See Settings Info tab for more detail on settings / preferences

- Default settings will be created on first run if they don't exist
`
			hlpText += "\n" + timerName + " v " + timerVersion
			hlpText += "\n" + timerCopyright
			hlpText += "\n\n" + timerAuthor + ", using Go and fyne GUI"

			plnText := `- Open with timer window focused
	- this is currently MacOS LaunchPad behavior, but only allows one app
	- To run more than one simultaneously, in terminal: open -n -a TaniumTimer 
- Add a separate date / time clock tab so we can see current time
- Add custom timer button, allow user to type no. minutes
- Add lab timer button
- Add more selectable timer buttons - list? Readable from prefs
- Timer show progress bar? Cute but not really necessary, countdown is very clear
- Center + / - below ad hoc button in canvas
- Reset timer name in window title to 'Tanium Timer' after user stop or timer end
- Test if already running bring to front and exit, optional setting to allow multiple timers
- Add pause/resume buttons to pause and resume a running timer
- Allow selectable png /svg images as backgrounds
	- Change window size and perspective h vs w to match background image sizes
	- Possible: add png to svg conversion, or simply display png rather than svg
- Allow optional always on top, save in prefs
- Settings allow user selectable mid / wav
- Add stop sounds button to stop the mp3/wav playing when enabled
`

			bugText := `
- Activating tray menus causes running timer display to not show updates
	until Help, About, Settings etc are selected
	- But timer does continue to countdown, fix to run systray, settings etc in parallel
- Settings changes to background and timer default times are saved immediately.
	- Times are effective immediately, but timer button times and background
		do not currently refresh to new settings
	`
			link, err := url.Parse("https://www.tanium.com/end-user-license-agreement-policy")
			if err != nil {
				fyne.LogError("Could not parse URL", err)
			}
			hyperlink := widget.NewHyperlink("https://www.tanium.com/end-user-license-agreement-policy", link)
			hyperlink.Alignment = fyne.TextAlignLeading
			licText := `Tanium Timer is “Beta Software” as defined in the license agreement found at the link below. 
Please take a moment to read the license agreement:
 
In addition, please note that:
Tanium Timer is intended for internal Tanium use, however no proprietary
information or features are included, so pending Tanium legal and other 
approvals this application may be made available to others. Tanium Timer
provides no guarantees as to stability of operations or suitability for any
purpose, but every attempt has been made to make this application reliable.

Using this application (and reading this text) is considered acceptance of
the terms of the License Agreement, and acknowledgement that this is Beta
Software and the additional terms above
`

			settingsText := `Settings are a separate tray menu item
Settings contains defaults as below, which can be modified, and also reset to defaults:
{"adhoc.default":300,"background.default":"blue","biobreak.default":600,
	"endsound.default":"baseball.mp3","halfminsound.default":"sosumi.mp3",
	"lunch.default":3600,"notify.default":1,"oneminsound.default":"hero.mp3",
	"sound.default":1}

Tanium Timer looks for directories named Resources/Images and Resources/Sounds,
containing images and sounds.

IMAGES:
Background blue refers to a compiled in resource with Tanium blue background. 
Other supported compiled in backgrounds are: stone, almond, converge24 and converge24a
Future additions will allow selecting images of your choice, png, SVG,
	jpg maybe and specifying size - height / width. Manual window resizing
	is already possible

SOUNDS:
Built in tones include 'ding', 'down', 'up', and 'updown'. These are always available
	and will be listed first in sound selectors
The sounds directory as distributed also contains a number of other .mp3 files
including baseball.mp3, grandfatherclock.mp3, hero.mp3, pinball.mp3, sosumi.mp3
When selecting sounds, the sound will be played as a preview when possible.
When selected sounds are not present (removed from Sounds), Tanium Timer defaults
	to playing built in tones ding, down, up or updown
Future additions will allow also choosing from any .mid or .wav sound files of your
	choice if located in the Sounds directory
`
			lic := widget.NewLabel(licText)
			tabs := container.NewDocTabs(
				container.NewTabItem("Help", widget.NewLabel(hlpText)),
				container.NewTabItem("Known Issues", widget.NewLabel(bugText)),
				container.NewTabItem("Planned Updates", widget.NewLabel(plnText)),
				container.NewTabItem("Settings Info", widget.NewLabel(settingsText)),
				container.NewTabItem("License", container.NewVBox(lic, hyperlink)),
			)
			tabs.SetTabLocation(container.TabLocationTop)
			tabs.Show()
			hlp.Resize(fyne.NewSize(800, 300))
			hlp.SetContent(tabs)
			hlp.CenterOnScreen() // run centered on primary (laptop) display
			hlp.Show()
		})
		settings := fyne.NewMenuItem("Settings", func() {
			makeSettings(a, w, bg)
		})
		clock := fyne.NewMenuItem("Clock", func() {
			clock(a) // , w, bg)
		})
		menu := fyne.NewMenu(a.Metadata().Name, show, hide, fyne.NewMenuItemSeparator(), lunch, biobreak, adhoc, stop, fyne.NewMenuItemSeparator(), clock, about, help, settings)
		desk.SetSystemTrayMenu(menu)
		systray.SetTooltip(timerName)
		// systray.SetTitle(timerName)

		running.AddListener(binding.NewDataListener(func() {
			busy, _ := running.Get()
			lunch.Disabled = busy
			biobreak.Disabled = busy
			adhoc.Disabled = busy
			stop.Disabled = !busy
			menu.Refresh()
		}))

	}

	less := widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {
		if adhocTime <= 5*60 { // min bound
			playBeep("ding")
			return
		}
		adhocTime -= 60 * 5
		updateTime(timer, adhocTime)
		a.Preferences().SetInt("adhoc.default", adhocTime)
	})
	less.Importance = widget.WarningImportance // orange
	more := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		adhocTime += 60 * 5
		updateTime(timer, adhocTime)
		a.Preferences().SetInt("adhoc.default", adhocTime)
	})
	more.Importance = widget.WarningImportance // orange

	timeRow := container.NewCenter(padTime(timer))
	lessmoreRow := container.NewHBox(container.NewCenter(less), container.NewCenter(more))

	lunch := widget.NewButton("Lunch ("+strconv.Itoa(lunchTime/60)+")", func() {
		startTimer(lunchTime, "Lunch", w.Canvas(), w)
	})
	lunch.Importance = widget.SuccessImportance // green
	biobreak := widget.NewButton("Bio Break ("+strconv.Itoa(biobreakTime/60)+")", func() {
		startTimer(biobreakTime, "Bio Break", w.Canvas(), w)
	})
	biobreak.Importance = widget.MediumImportance // white
	adhoc := widget.NewButton("Ad Hoc", func() {
		startTimer(adhocTime, "Ad Hoc", w.Canvas(), w)
	})
	adhoc.Importance = widget.WarningImportance // orange
	quit := widget.NewButton("Quit", func() {
		a.Quit()
	})
	quit.Importance = widget.HighImportance // red
	content := container.NewCenter(container.NewVBox(timeRow,
		container.NewGridWithColumns(2, biobreak, lunch, adhoc, quit), lessmoreRow))

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
	w.Resize(fyne.NewSize(content.MinSize().Width*1.8, content.MinSize().Height*1.8))
	bg.FillMode = canvas.ImageFillContain
	bg.Translucency = 0.5 // 0.85
	w.SetContent(container.NewStack(
		bg,
		container.NewPadded(container.NewPadded(content))))
	w.ShowAndRun()
}

func formatTimer(time int) string {
	secs := time % 60
	mins := (time - secs) / 60
	return fmt.Sprintf("%02d:%02d", mins, secs)
}

func padTime(t *widget.RichText) fyne.CanvasObject {
	pad := theme.Padding()
	return container.New(layout.NewCustomPaddedLayout(-3.5*pad, -2.5*pad, pad, pad), t)
}

func startTimer(timer int, name string, c fyne.Canvas, w fyne.Window) {
	remain = timer
	busy, _ := running.Get()
	if busy {
		return
	}
	w.SetTitle(timerName + ": " + name)
	running.Set(true)
	if desk, ok := fyne.CurrentApp().(desktop.App); ok {
		desk.SetSystemTrayIcon(resourceTaniumIconSvg)
		systray.SetTooltip(timerName)
		// systray.SetTitle(timerName)
	}
	ticker := widget.NewRichText()
	updateTime(ticker, remain)

	stop := widget.NewButton("Stop", nil)
	overlay := container.NewPadded(container.NewVBox(
		padTime(ticker),
		stop))
	p := widget.NewModalPopUp(overlay, c)
	stop.OnTapped = func() {
		remain = -1 // don't notify
		w.SetTitle(timerName)
		if desk, ok := fyne.CurrentApp().(desktop.App); ok {
			desk.SetSystemTrayIcon(resourceTaniumIconSvg)
			systray.SetTooltip(timerName)
			systray.SetTitle("")
			stop.Disable()
		}
		p.Hide()
	}
	go func() {
		for remain > 0 {
			updateTime(ticker, remain)
			if _, ok := fyne.CurrentApp().(desktop.App); ok {
				systray.SetTitle(formatTimer(remain))
			}
			if remain == 60 {
				w.Show() // in case it has been hidden
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
				w.Show() // in case it has been hidden
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
		w.SetTitle(timerName)
		running.Set(false)
		if remain == 0 {
			updateTime(ticker, remain)
			stop.Disable()
			w.Show() // in case it has been hidden
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
					w.Hide()
					time.Sleep(time.Second / 2)
					w.Show()
					time.Sleep(time.Second / 2)
				}
			}
		}
		if desk, ok := fyne.CurrentApp().(desktop.App); ok {
			desk.SetSystemTrayIcon(resourceTaniumIconSvg)
			systray.SetTooltip(timerName)
			systray.SetTitle("")
		}
		p.Hide()
	}()
	p.Show()
}

func updateTime(out *widget.RichText, time int) {
	out.ParseMarkdown("# " + formatTimer(time))
	themeTimer(out, time)
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942

// This timer is based on an original project named Fomato by Andy Williams, and heavily redeveloped - Allan Marillier, 2024

/*
To-do:
Allow optional always on top, save in prefs
https://www.google.com/search?q=fyne+golang+always+on+top&oq=fyne+golang+always+on+top&gs_lcrp=EgZjaHJvbWUyBggAEEUYOTIKCAEQABiABBiiBDIKCAIQABiABBiiBDIKCAMQABiABBiiBDIKCAQQABiABBiiBNIBCDg5MTBqMGoxqAIAsAIA&sourceid=chrome&ie=UTF-8

*/
