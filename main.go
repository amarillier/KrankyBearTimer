package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
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

const (
	// appName    = "Kranky Bear Timer"
	appVersion = "0.9.1" // see FyneApp.toml
	appAuthor  = "Allan Marillier"
)

var appName = "Kranky Bear Timer"
var appNameCustom = ""
var appCopyright = "Copyright (c) Allan Marillier, 2024-" + strconv.Itoa(time.Now().Year())
var running = binding.NewBool()
var bg fyne.Canvas
var remain int
var notify int
var sound int
var traytimer int

var adhocTime int
var adhocbtn *widget.Button
var adhocmnu *fyne.MenuItem
var biobreakTime int
var lunchTime int
var endTime time.Time
var customTime time.Time
var endTimeSec int
var menu *fyne.Menu
var clock fyne.Window // clock window

var imgDir string
var timerbg string
var starttimer int

var sndDir string
var endsnd string
var oneminsnd string
var halfminsnd string

var showseconds int
var showtimezone int
var showdate int
var showutc int
var showhr12 int
var hourchime int
var slockmute int
var clockmutedvol int
var automute int
var currentvolume int
var muteonhr int
var muteonmin int
var muteoffhr int
var muteoffmin int
var bgcolor string
var timecolor string
var datecolor string
var utccolor string
var timefont string
var datefont string
var utcfont string
var timesize int
var datesize int
var utcsize int
var hourchimesound string
var startclock int
var processName string

/*
	minor difference from clock app which sets OS autostart,
	this in the timer app will influence opening the clock window
	when the timer app starts
*/

var debug int = 0
var abt fyne.Window
var hlp fyne.Window
var updt fyne.Window
var timerWidth float64
var timerHeight float64

// preferences stored via fyne preferences API land in
// ~/Library/Preferences/fyne/com.github.amarillier.KrankyBearTimer/preferences.json
// ~\AppData\Roaming\fyne\com.github.amarillier.KrankyBearTimer\preferences.json
// {"adhoc.default":300,"background.default":"blue","biobreak.default":600,"endsound.default":"baseball.mp3","halfminsound.default":"sosumi.mp3","lunch.default":3600,"notify.default":1,"oneminsound.default":"hero.mp3", "sound.default":1}

func main() {
	exePath, err := os.Executable()
	processName = filepath.Base(os.Args[0])
	if err != nil {
		panic(err)
	}

	launchDir := filepath.Dir(exePath)

	if runtime.GOOS == "darwin" {
		if strings.HasPrefix(launchDir, "/Applications/KrankyBearTimer") {
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

	a := app.NewWithID("com.github.amarillier.KrankyBearTimer")
	a.Settings().SetTheme(&appTheme{Theme: theme.DefaultTheme()})
	w := a.NewWindow(appName)
	w.SetIcon(resourceKrankyBearPng)
	w.SetPadded(false)

	w.SetCloseIntercept(func() {
		width := w.Content().Size().Width
		height := w.Content().Size().Height
		timerWidth = float64(width)
		timerHeight = float64(height)
		a.Preferences().SetFloat("width.default", float64(width))
		a.Preferences().SetFloat("height.default", float64(height))
		w.Close()
		// NEVER use a.Quit(), this hangs!
		// a.Quit() // force quit, normal when somebody hits "x" to close
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
		a.Preferences().SetString("timername.default", "")
	}
	// get default timer settings from preferences
	lunchTime = a.Preferences().IntWithFallback("lunch.default", 60*60)
	adhocTime = a.Preferences().IntWithFallback("adhoc.default", 5*60)
	biobreakTime = a.Preferences().IntWithFallback("biobreak.default", 10*60)
	notify = a.Preferences().IntWithFallback("notify.default", 1)
	sound = a.Preferences().IntWithFallback("sound.default", 1)
	traytimer = a.Preferences().IntWithFallback("traytimer.default", 0)
	timerbg = a.Preferences().StringWithFallback("background.default", "board1")
	endsnd = a.Preferences().StringWithFallback("endsound.default", "baseball.mp3")
	oneminsnd = a.Preferences().StringWithFallback("oneminsound.default", "hero.mp3")
	halfminsnd = a.Preferences().StringWithFallback("halfminsound.default", "sosumi.mp3")
	starttimer = a.Preferences().IntWithFallback("starttimer.default", 0)
	// get default clock settings from preferences
	showseconds = a.Preferences().IntWithFallback("showseconds.default", 1)
	showtimezone = a.Preferences().IntWithFallback("showtimezone.default", 1)
	showdate = a.Preferences().IntWithFallback("showdate.default", 1)
	showutc = a.Preferences().IntWithFallback("showutc.default", 1)
	showhr12 = a.Preferences().IntWithFallback("showhr12.default", 1)
	slockmute = a.Preferences().IntWithFallback("slockmute.default", 0)
	automute = a.Preferences().IntWithFallback("automute.default", 0)
	muteonhr = a.Preferences().IntWithFallback("muteonhr.default", 20)
	muteonmin = a.Preferences().IntWithFallback("muteonmin.default", 0)
	muteoffhr = a.Preferences().IntWithFallback("muteoffhr.default", 8)
	muteoffmin = a.Preferences().IntWithFallback("muteoffmin.default", 0)
	hourchime = a.Preferences().IntWithFallback("hourchime.default", 1)
	bgcolor = a.Preferences().StringWithFallback("bgcolor.default", "0,143,251,255")      // blue
	timecolor = a.Preferences().StringWithFallback("timecolor.default", "255,123,31,255") // orange
	datecolor = a.Preferences().StringWithFallback("datecolor.default", "131,222,74,255") // red
	utccolor = a.Preferences().StringWithFallback("utccolor.default", "238,229,58.255")   // yellow
	timefont = a.Preferences().StringWithFallback("timefont.default", "arial")            // not yet!
	datefont = a.Preferences().StringWithFallback("datefont.default", "arial")            // not yet!
	utcfont = a.Preferences().StringWithFallback("utcfont.default", "arial")              // not yet!
	timesize = a.Preferences().IntWithFallback("timesize.default", 36)
	datesize = a.Preferences().IntWithFallback("datesize.default", 24)
	utcsize = a.Preferences().IntWithFallback("utcsize.default", 18)
	hourchimesound = a.Preferences().StringWithFallback("hourchimesound.default", "hero.mp3")
	startclock = a.Preferences().IntWithFallback("startclock.default", 0)
	// "Mon Jan 2 15:04:05 MST 2006"
	endTime, _ = time.Parse("15:04", "00:00") // set default midnight

	// Allow for user defined custom timer name to brand e.g. Tanium Timer
	appNameCustom = a.Preferences().StringWithFallback("timername.default", appName)
	if appNameCustom != "" && appNameCustom != "default" {
		// allow for Tanium branding, test also for Tanium backgrounds
		// not allowed any other times
		tanium := regexp.MustCompile(`^(?i)tanium`)
		if tanium.MatchString(appNameCustom) {
			// if it does not end with [Tt]imer, add it
			// if !strings.HasSuffix(timerName, "Timer") && !strings.HasSuffix(timerName, "timer"){
			timer := regexp.MustCompile(` (?i)timer`)
			if !timer.MatchString(appNameCustom) {
				appNameCustom += " Timer"
			}
		}
		appName = appNameCustom
		w.SetTitle(appName)
	} else {
		// if timer is not customized to Tanium, don't allow use of
		// built in Tanium backgrounds, but user added are ok
		if timerbg == "blue" || timerbg == "stone" || timerbg == "almond" || timerbg == "gray" {
			timerbg = "board1" // reset default non Tanium
		}
	}

	if len(os.Args) >= 2 {
		log.Println("arg count:", len(os.Args))
		if os.Args[1] == "debug" || os.Args[1] == "d" {
			debug = 1
			logInit()
			r, _ := os.Open("KrankyBearTimer0.txt")
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

	// check update first
	updtmsg, updateAvail := updateChecker("amarillier", "KrankyBearTimer", "Kranky Bear Timer", "https://github.com/amarillier/KrankyBearTimer/releases/latest")
	if updateAvail {
		// open a window to show the update message
		// no need to test for updt window open at first start
		// appName += " (update available)"
		// appNameCustom += " (update available)"
		kb := canvas.NewImageFromResource(resourceKrankyBearPng)
		kb.FillMode = canvas.ImageFillOriginal
		text := widget.NewLabel(updtmsg)
		content := container.NewHBox(kb, text)
		updt = a.NewWindow(appName + ": Update Check")
		updt.SetIcon(resourceKrankyBearPng)
		updt.Resize(fyne.NewSize(50, 100))
		updt.SetContent(content)
		updt.SetCloseIntercept(func() {
			updt.Close()
			updt = nil
		})
		// updt.CenterOnScreen() // run centered on primary (laptop) display
		updt.Show()

	}

	if desk, ok := a.(desktop.App); ok {
		desk.SetSystemTrayIcon(resourceKrankyBearPng)
		if startclock == 1 {
			desktopclock(a)
		}
		systray.SetTooltip(appName)
		//systray.SetTitle(timerName)
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
		adhocmnu = fyne.NewMenuItem("Ad Hoc ("+strconv.Itoa(adhocTime/60)+")", func() {
			startTimer(adhocTime, "Ad Hoc Timer", w.Canvas(), w)
		})
		selected := fyne.NewMenuItem("Selected End Time", func() {
			startTimer(endTimeSec, "Selected End Time", w.Canvas(), w)
		})
		stop := fyne.NewMenuItem("Stop", func() {
			remain = -1 // don't notify when the user stops it
		})
		about := fyne.NewMenuItem("About", func() {
			aboutText := appName + " v " + appVersion
			aboutText += "\n" + appCopyright + ", written using Go and fyne GUI"
			if appNameCustom != "" && appNameCustom != "default" {
				aboutText += "\n\n(Currently rebranded as " + appNameCustom + ")"
			}
			aboutText += "\n\nCreated by " + appAuthor + ", using Go and fyne GUI"
			aboutText += "\n\nNo obligation, it's rewarding to hear if you use this app."
			aboutText += "\n\nLooking about about and help or settings too too much might expose an easter egg!"

			kb := canvas.NewImageFromResource(resourceKrankyBearPng)
			text := widget.NewLabel(aboutText)
			kb.FillMode = canvas.ImageFillOriginal
			content := container.NewHBox(kb, text)

			if abt == nil {
				abt = a.NewWindow(appName + ": About")
				abt.SetIcon(resourceKrankyBearPng)
				abt.Resize(fyne.NewSize(50, 100))
				// abt.SetContent(widget.NewLabel(aboutText))
				abt.SetContent(content)
				abt.SetCloseIntercept(func() {
					abt.Close()
					abt = nil
				})
				abt.CenterOnScreen() // run centered on pr1imary (laptop) display
				abt.Show()
			} else {
				abt.RequestFocus()
				easterEgg(a, w)
			}
		})
		help := fyne.NewMenuItem("Help", func() {
			if hlp == nil {
				hlp = a.NewWindow(appName + ": Help")
				hlp.SetIcon(resourceKrankyBearPng)
				hlpText := `This application is primarily a timer to manage ad hoc, bio-break and lunch break times during training or other events. 
It also includes an optional desktop clock that can be set to auto start when the timer starts, or run on demand as needed.

NOTE: The timer main window can be rebranded from default Kranky Bear Timer to any name of your choice by setting
the timername.default preference in the settings menu. This is a manual configuration only, not available 
via settings, and will not be reset if the settings reset option is used.
Trainer colleagues, this will also enable some additional built in custom specific backgrounds.

- Ad hoc timer minimum is 5 minutes, with 5 minute increments
	- NOTE ad hoc default is updated in preferences to current value any time it is changed
- Bio break timer default is 10 minutes
- Lunch break timer default is 60 minutes
- Each of these break times can be modified using Settings, set in minutes
- A custom time can also be set using the 'Select End Time' button.
	- This time will be calculated in minutes from the current time when set, and is reset when the timer ends.
- Timer text color is green until 2 1/2 minutes remain,
	- color is orange from 2 1/2 minutes to 30 seconds
	- color is red from 30 seconds to completion
- optional setting to enable auto starting at boot

- System tray notifications and sound alerts are both optional, enabled by default
- System tray can display the countdown timer when enabled and a timer is running.
	This is disabled by default to save CPU cycles updating it.
	Minor, but you may see increased CPU usage when this is enabled.
- Tone / beep alerts are at 60 seconds, at 30 seconds, and at completion
- Timer window flashes on/off at timer end (in addition to desktop notification & beep)
- A timer that has been hidden behind another window or minimized will be
	brought to the front / focused at 60 seconds, and at timer completion

- The separate clock window allows optional display of seconds, date, UTC time, with customizable background and 
  text colors available, configured through a separate settings menu item
- autostart clock when starting the timer is also available, in clock settings
- Note: Displaying seconds can be quite resource intensive with clock display updates every second. 
  The app can be substantially less CPU intensive when seconds are not displayed, allowing the app to
  refresh the display every minute rather than every second

- See Settings Info tab for more detail on settings / preferences

- Default settings will be created on first run if they don't exist
`
				hlpText += "\n" + appName + " v " + appVersion
				hlpText += "\n" + appCopyright
				hlpText += ", written using Go and fyne GUI"

				plnText := `- Allow multiple time zones for clock, hh:mm only + offset
- Allow multiple alarm times with user selectable tones for each, one time, recurring etc.
- Allow settings set/save window locations to open timer/clock,
	unfortunately not implemented in the fyne library yet
- Open with timer window focused
	- this is currently MacOS LaunchPad behavior, but only allows one app
	- To run more than one simultaneously, in terminal: open -n -a KrankyBearTimer 
- Possibly add lab timer button, or more selectable timer buttons, as a list, readable from prefs
- Timer show progress bar? Cute but not really necessary, countdown is very clear
- Test if already running bring to front and exit, optional setting to allow multiple timers
- Add pause/resume buttons to pause and resume a running timer
- Allow optional always on top, save in prefs
- Possibly add stop sound button to stop audio tones, mp3/wav playing when enabled and already started. 
`

				bugText := `
- Activating tray menus causes running timer display to not show updates
	until Help, About, Settings etc are selected
	- But timer does continue to countdown, fix to run systray, settings etc in parallel
- Settings changes to background and timer default times are saved immediately.
	- Times are effective immediately, but timer button times and background
		do not currently refresh to new settings
- Font type settings in preferences are currently ignored, the app uses system theme defaults. (Future planned update)
- OpenGL drivers are required for some Windows systems, not a bug but a specific library requirement that might not allow some to use this app
	`
				link, err := url.Parse("https://github.com/amarillier/KrankyBearTimer/blob/main/license.txt")
				if err != nil {
					fyne.LogError("Could not parse URL", err)
				}
				hyperlink := widget.NewHyperlink("https://github.com/amarillier/KrankyBearTimer/blob/main/license.txt", link)
				hyperlink.Alignment = fyne.TextAlignLeading
				licText := `KrankyBearTimer is FREE Software” as defined in the license agreement below. 
 
This application is "FREE Software". 

This application is intended for any use by any individual, in any organization.

This application provides no guarantees as to stability of operations or suitability 
for any purpose, but every attempt has been made to make this application reliable.

This application may not be sold, no money may be asked by anyone for provision of, or any services related to this application.

Using this application (and reading this text) is considered acceptance of
the terms of the License Agreement, and acknowledgement that this is FREE
Software and the additional terms above.

See https://github.com/amarillier/KrankyBearTimer/
`

				settingsText := `Settings are a separate tray menu item
Settings contains defaults, which can be modified as well as reset to defaults in Settings menus. 

One exception is the default.timername which can be used to rebrand the timer main window with a custom name.
The timer name is set manually in the preferences file, and will not be reset if the settings are reset.
See below for the preferences.json file location.

KrankyBear Timer looks for directories named Resources/Images and Resources/Sounds,
containing optional user provided images and sounds.

IMAGES:
Some background images are included, compiled into the app, user selectable
.png and .jpg images can also be placed in the app's Resources/Images directory

App window size is detected and automatically saved in preferences when exiting the app to preserve preferences if the window is resized.

SOUNDS:
Built in tones include 'ding', 'down', 'up', and 'updown'. These are always available
	and will be listed first in sound selectors
The Resources/Sounds directory as distributed also contains a number of other .mp3 files
including baseball.mp3, grandfatherclock.mp3, hero.mp3, pinball.mp3, sosumi.mp3
When selecting sounds, the sound will be played as a preview when possible.
When selected sounds are not present (removed from Sounds), KrankyBear Timer defaults
	to playing built in tones ding, down, up or updown
Future additions may allow also choosing from other sound file types of your choice if located in the Sounds directory

Resources directory locations:
MacOS: /Applications/KrankyBearTimer.app/Contents/Resources
Windows: \Program Files/KrankyBearTimer\Contents\Resources
preferences.json file location:
MacOS: ~/Library/Preferences/fyne/com.github.amarillier.KrankyBearTimer/preferences.json
Windows: ~\AppData\Roaming\fyne\com.github.amarillier.KrankyBearTimer/preferences.json
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
				hlp.SetCloseIntercept(func() {
					hlp.Close()
					hlp = nil
				})
				hlp.CenterOnScreen() // run centered on primary (laptop) display
				hlp.Show()
			} else {
				hlp.RequestFocus()
				easterEgg(a, w)
			}
		})
		settingsTimer := fyne.NewMenuItem("Settings (Timer)", func() {
			makeSettingsTimer(a, w, bg)
		})
		settingsClock := fyne.NewMenuItem("Settings (Clock)", func() {
			makeSettingsClock(a, w, bg)
		})
		settingsTheme := fyne.NewMenuItem("Settings (Theme)", func() {
			makeSettingsTheme(a, w, bg)
		})
		clock := fyne.NewMenuItem("Clock", func() {
			if clock == nil {
				desktopclock(a)
			} else {
				clock.RequestFocus()
			}
		})
		updtchk := fyne.NewMenuItem("Check for update", func() {
			updtmsg, updateAvail := updateChecker("amarillier", "KrankyBearTimer", "Kranky Bear Timer", "https://github.com/amarillier/KrankyBearTimer/releases/latest")
			if updt == nil {
				kb := canvas.NewImageFromResource(resourceKrankyBearPng)
				kb.FillMode = canvas.ImageFillOriginal
				text := widget.NewLabel(updtmsg)
				content := container.NewHBox(kb, text)
				updt = a.NewWindow(appName + ": Update Check")
				updt.SetIcon(resourceKrankyBearPng)
				updt.Resize(fyne.NewSize(50, 100))
				// updt.SetContent(widget.NewLabel(updtmsg))
				updt.SetContent(content)
				updt.SetCloseIntercept(func() {
					updt.Close()
					updt = nil
				})
				updt.CenterOnScreen() // run centered on pr1imary (laptop) display
				updt.Show()
				// if !strings.Contains(updtmsg, "You are running the latest") {
				if updateAvail {
					if !checkFileExists(sndDir + "/KrankyBearGrowl.mp3") {
						playBeep("up")
					} else {
						playMp3(sndDir + "//KrankyBearGrowl.mp3") // Basso, Blow, Hero, Funk, Glass, Ping, Purr, Sosumi, Submarine,
					}
				}
			} else {
				updt.RequestFocus()
			}
		})
		menu = fyne.NewMenu(a.Metadata().Name, show, hide, fyne.NewMenuItemSeparator(), lunch, biobreak, adhocmnu, selected, stop, fyne.NewMenuItemSeparator(), clock, about, updtchk, help, settingsTimer, settingsClock, settingsTheme)
		desk.SetSystemTrayMenu(menu)
		desk.SetSystemTrayIcon(resourceKrankyBearPng)
		systray.SetTooltip(appName)

		// Menu items
		// compile / run with syntax below to force Mac to do menus like Windows
		// otherwise menus will be at the top of the display
		// https://github.com/fyne-io/fyne/issues/3988
		// go build -tags no_native_menus .
		// go run -tags no_native_menus .
		quit := fyne.NewMenuItem("Quit", func() {
			width := w.Content().Size().Width
			height := w.Content().Size().Height
			timerWidth = float64(width)
			timerHeight = float64(height)
			a.Preferences().SetFloat("width.default", float64(width))
			a.Preferences().SetFloat("height.default", float64(height))
			a.Quit()
		})
		newMenuOps := fyne.NewMenu("Operations", show, hide, clock, fyne.NewMenuItemSeparator(), quit)
		newMenuTimers := fyne.NewMenu("Timers", lunch, biobreak, adhocmnu, selected, stop)
		// NB Mac intercepts about item below and puts it where they want to put it!
		// Under 'KrankyBear Timer / About' main section, not under Help
		newMenuHelp := fyne.NewMenu("Help", about, help)
		newMenuSettings := fyne.NewMenu("Settings", settingsTimer, settingsClock, settingsTheme)
		barmenu := fyne.NewMainMenu(newMenuOps, newMenuTimers, newMenuHelp, newMenuSettings)
		w.SetMainMenu(barmenu)
		// barmenu.Refresh()

		running.AddListener(binding.NewDataListener(func() {
			busy, _ := running.Get()
			lunch.Disabled = busy
			biobreak.Disabled = busy
			adhocmnu.Disabled = busy
			selected.Disabled = busy
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
		adhocbtn.SetText("Ad Hoc (" + strconv.Itoa(adhocTime/60) + ")")
		adhocmnu.Label = "Ad Hoc (" + strconv.Itoa(adhocTime/60) + ")"
		menu.Refresh()
		a.Preferences().SetInt("adhoc.default", adhocTime)
	})
	less.Importance = widget.WarningImportance // orange
	more := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		adhocTime += 60 * 5
		adhocbtn.SetText("Ad Hoc (" + strconv.Itoa(adhocTime/60) + ")")
		adhocmnu.Label = "Ad Hoc (" + strconv.Itoa(adhocTime/60) + ")"
		menu.Refresh()
		a.Preferences().SetInt("adhoc.default", adhocTime)
	})
	more.Importance = widget.WarningImportance // orange
	endset := widget.NewButtonWithIcon("", theme.RadioButtonIcon(), func() {
		setEndTime(a, w, bg, "set")
		now := time.Now()
		endTimeSec = (customTime.Hour()*60*60 + customTime.Minute()*60) - (now.Hour()*60*60 + now.Minute()*60 + now.Second())
		// a.Preferences().SetInt("endTime.default", endTime)
	})
	endset.Importance = widget.WarningImportance // orange

	lessmoreRow := container.NewHBox(container.NewCenter(less), container.NewCenter(more), layout.NewSpacer(), endset)

	lunch := widget.NewButton("Lunch ("+strconv.Itoa(lunchTime/60)+")", func() {
		startTimer(lunchTime, "Lunch", w.Canvas(), w)
	})
	lunch.Importance = widget.SuccessImportance // green
	biobreak := widget.NewButton("Bio Break ("+strconv.Itoa(biobreakTime/60)+")", func() {
		startTimer(biobreakTime, "Bio Break", w.Canvas(), w)
	})
	biobreak.Importance = widget.MediumImportance // white
	adhocbtn = widget.NewButton("Ad Hoc ("+strconv.Itoa(adhocTime/60)+")", func() {
		startTimer(adhocTime, "Ad Hoc", w.Canvas(), w)
	})
	adhocbtn.Importance = widget.WarningImportance // orange
	endtime := widget.NewButton("Selected End Time", func() {
		now := time.Now()
		endTimeSec = (customTime.Hour()*60*60 + customTime.Minute()*60) - (now.Hour()*60*60 + now.Minute()*60 + now.Second())
		if endTimeSec <= 60 {
			playBeep("ding")
			setEndTime(a, w, bg, "run")
			now := time.Now()
			endTimeSec = (customTime.Hour()*60*60 + customTime.Minute()*60) - (now.Hour()*60*60 + now.Minute()*60 + now.Second())
		} else {
			now := time.Now()
			endTime = time.Date(now.Year(), now.Month(), now.Day(), customTime.Hour(), customTime.Minute(), 0, 0, now.Location())
			startTimer(endTimeSec, "Selected End Time", w.Canvas(), w)
		}
	})
	endtime.Importance = widget.WarningImportance // orange

	content := container.NewCenter(container.NewVBox(container.NewGridWithColumns(2, biobreak, lunch, adhocbtn, endtime), lessmoreRow))

	bg := canvas.NewImageFromResource(resourceSchoolBoard1Png)
	if strings.HasSuffix(timerbg, ".png") || strings.HasSuffix(timerbg, ".jpg") {
		// if it's a png or jpg file specified, test if it exists and use it
		// otherwise use resource based image
		if checkFileExists(imgDir + "/" + timerbg) {
			bg = canvas.NewImageFromFile(imgDir + "/" + timerbg)
		} else {
			bg = canvas.NewImageFromResource(resourceSchoolBoard1Png)
		}
	} else {
		switch timerbg {
		case "board1":
			bg = canvas.NewImageFromResource(resourceSchoolBoard1Png)
		case "board2":
			bg = canvas.NewImageFromResource(resourceSchoolBoard2Png)
		case "board3":
			bg = canvas.NewImageFromResource(resourceSchoolBoard3Png)
		case "board4":
			bg = canvas.NewImageFromResource(resourceSchoolBoard4Png)
		case "board5":
			bg = canvas.NewImageFromResource(resourceSchoolBoard5Png)
		case "board6":
			bg = canvas.NewImageFromResource(resourceSchoolBoard6Png)
		case "blue":
			bg = canvas.NewImageFromResource(resourceTBluePng)
		case "stone":
			bg = canvas.NewImageFromResource(resourceTStonePng)
		case "almond":
			bg = canvas.NewImageFromResource(resourceTAlmondPng)
		case "gray":
			bg = canvas.NewImageFromResource(resourceTGrayTeachPng)
		default:
			bg = canvas.NewImageFromResource(resourceSchoolBoard1Png)
		}
	}

	width := a.Preferences().FloatWithFallback("width.default", float64(content.MinSize().Width*1.8))
	height := a.Preferences().FloatWithFallback("height.default", float64(content.MinSize().Height*1.8))
	w.Resize(fyne.NewSize(float32(width), float32(height)))
	// w.Resize(fyne.NewSize(content.MinSize().Width*1.8, content.MinSize().Height*1.8))
	// w.Resize(fyne.NewSize(content.MinSize().Width*2.2, content.MinSize().Height*2.2))
	bg.FillMode = canvas.ImageFillContain
	// bg.FillMode = canvas.ImageFillOriginal
	bg.Translucency = 0.5 // 0.85
	w.SetContent(container.NewStack(
		bg,
		container.NewPadded(container.NewPadded(content))))
	w.ShowAndRun()
	if updt != nil {
		updt.RequestFocus()
	}
}

func formatTimer(time int) string {
	secs := time % 60
	mins := (time - secs) / 60
	return fmt.Sprintf("%02d:%02d", mins, secs)
}

func centerTime(t *widget.RichText) fyne.CanvasObject {
	return container.New(layout.NewCenterLayout(), t)
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
	w.SetTitle(appName + ": " + name)
	running.Set(true)

	if desk, ok := fyne.CurrentApp().(desktop.App); ok {
		desk.SetSystemTrayIcon(resourceKrankyBearPng)
		systray.SetTooltip(appName)
		// systray.SetTitle(timerName)
	}

	ticker := widget.NewRichText()
	fyne.Do(func() {
		updateTime(ticker, remain)
	})

	stop := widget.NewButton("Stop", nil)
	overlay := container.NewPadded(container.NewVBox(
		// padTime(ticker),
		centerTime(ticker),
		stop))
	p := widget.NewModalPopUp(overlay, c)
	//p.Resize(fyne.NewSize(300, 100))
	overlay.Resize(fyne.NewSize(100, 100))
	p.Resize(fyne.NewSize(w.Canvas().Size().Width*0.5, w.Canvas().Size().Height*0.5))
	stop.OnTapped = func() {
		remain = -1 // don't notify
		w.SetTitle(appName)
		if desk, ok := fyne.CurrentApp().(desktop.App); ok {
			desk.SetSystemTrayIcon(resourceKrankyBearPng)
			systray.SetTooltip(appName)
			systray.SetTitle("")
			stop.Disable()
		}
		p.Hide()
	}
	go func() {
		for remain > 0 {
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
		fyne.Do(func() {
			w.SetTitle(appName)
		})

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
			desk.SetSystemTrayIcon(resourceKrankyBearPng)
			systray.SetTooltip(appName)
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

func setEndTime(a fyne.App, w fyne.Window, bg fyne.Canvas, caller string) {
	var selectedTime time.Time
	var current string

	e := a.NewWindow("Select End Time")
	// Set window size to fit the input prompt
	e.Resize(fyne.NewSize(300, 150))
	now := time.Now()

	// check to see if predefined end / custom time is still
	// in the future, if not, set to the current time. If it is future,
	// default to that future time
	if customTime.Before(now) {
		current = fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute()+5)
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
	submitButton := widget.NewButton("Set", func() {
		enteredTime := timeEntry.Text
		if isValidCustomTime(enteredTime, "custom") {
			selectedTime, _ = time.Parse("15:04", enteredTime)
			customTime = time.Date(now.Year(), now.Month(), now.Day(), selectedTime.Hour(), selectedTime.Minute(), 0, 0, now.Location())
			if caller == "set" {
				messageLabel.SetText("Custom time: " + customTime.Format("Mon Jan 2 15:04:05 MST 2006"))
				time.Sleep(1 * time.Second)
			} else {
				messageLabel.SetText("Custom time: " + customTime.Format("Mon Jan 2 15:04:05 MST 2006"+"\n\nTime has been set\nPress the Selected End Time button again\nwhen ready to run the timer"))
				time.Sleep(4 * time.Second)
			}
			e.Close()
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
	endTime = customTime
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

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942

// This timer is based on an original project named Fomato by Andy Williams, and heavily redeveloped - Allan Marillier, 2024

/*
To-do:
- Allow optional always on top, save in prefs - may not be possible on Mac
https://www.google.com/search?q=fyne+golang+always+on+top&oq=fyne+golang+always+on+top&gs_lcrp=EgZjaHJvbWUyBggAEEUYOTIKCAEQABiABBiiBDIKCAIQABiABBiiBDIKCAMQABiABBiiBDIKCAQQABiABBiiBNIBCDg5MTBqMGoxqAIAsAIA&sourceid=chrome&ie=UTF-8
- Known problems - needs OpenGL drivers on some Windows
-
*/
