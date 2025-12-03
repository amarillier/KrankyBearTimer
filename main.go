package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
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

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
)

const (
	// appName    = "Kranky Bear Timer"
	appVersion = "0.9.5" // see FyneApp.toml
	appAuthor  = "Allan Marillier"
)

var appName = "Kranky Bear Timer"
var appNameCustom = ""
var appCopyright = "Copyright (c) Allan Marillier, 2024-" + strconv.Itoa(time.Now().Year())
var running = binding.NewBool()
var bg fyne.Canvas

// Timer-specific variables moved to timer.go
var adhocbtn *widget.Button
var adhocmnu *fyne.MenuItem
var menu *fyne.Menu
var clock fyne.Window           // clock window
var selectedMenu *fyne.MenuItem // Selected End Time menu item
var elapsedMenu *fyne.MenuItem  // Elapsed Time menu item
var stopMenu *fyne.MenuItem     // Stop menu item
var lunchMenu *fyne.MenuItem    // Lunch timer menu item
var biobreakMenu *fyne.MenuItem // Bio Break timer menu item
var lunchBtn *widget.Button     // Lunch timer button
var biobreakBtn *widget.Button  // Bio Break timer button
var elapsedBtn *widget.Button   // Elapsed timer button

// Shared directories (used by timer, audio, etc.)
var imgDir string
var sndDir string

var showseconds int
var showtimezone int
var showdate int
var showutc int
var showhr12 int
var hourchime int
var slockmute int
var clockmutedvol int
var automute int
var jiggle int
var jiggleconf int
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
var prefs string
var lastChimeHour int = -1              // Track last hour when chime played to prevent double playback
var clockUpdateLoopRunning bool = false // Prevent multiple update loops from running
var clockUpdateLoopStop chan bool       // Channel to stop the update loop
var hourChimeFileExists bool = false    // Cache file existence check
var hourChimeFileChecked bool = false   // Track if we've checked the file
var hourChimeCachedFile string = ""     // Track which file was cached

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

// Countdown window variables
var countdown fyne.Window
var countdownDate1 string
var countdownDate2 string
var countdownDate3 string
var countdownDesc1 string
var countdownDesc2 string
var countdownDesc3 string
var countdownTitleText *canvas.Text
var countdownHelpText *canvas.Text
var countdownBackground *canvas.Rectangle
var countdownDaysText1 *canvas.Text
var countdownDaysText2 *canvas.Text
var countdownDaysText3 *canvas.Text

// Additional timezones (up to 5)
var timezone1Enabled int
var timezone1Name string
var timezone1Offset string // UTC offset (e.g., "+5", "-3.5")
var timezone2Enabled int
var timezone2Name string
var timezone2Offset string
var timezone3Enabled int
var timezone3Name string
var timezone3Offset string
var timezone4Enabled int
var timezone4Name string
var timezone4Offset string
var timezone5Enabled int
var timezone5Name string
var timezone5Offset string

// preferences stored via fyne preferences API land in
// ~/Library/Preferences/fyne/com.github.amarillier.KrankyBearTimer/preferences.json
// ~\AppData\Roaming\fyne\com.github.amarillier.KrankyBearTimer\preferences.json

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

	// Initialize speaker at application launch for faster audio playback
	speakerSampleRate = beep.SampleRate(48000)
	speaker.Init(speakerSampleRate, speakerSampleRate.N(time.Second/10))

	a := app.NewWithID("com.github.amarillier.KrankyBearTimer")

	// Load alarms and start alarm checker
	loadAlarms(a)
	startAlarmChecker(a)
	// Load weather settings and start weather refresh if enabled
	loadWeatherSettings(a)
	startWeatherRefresh(a)
	a.Settings().SetTheme(&appTheme{Theme: theme.DefaultTheme()})
	w := a.NewWindow(appName)
	timerWindow = w // store reference for dynamic updates
	_, month, _ := time.Now().Date()
	if month == time.December {
		w.SetIcon(resourceKrankyBearChristmasGrinchPng)
	} else {
		w.SetIcon(resourceKrankyBearFedoraRedPng)
	}
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

	prefs = strings.ReplaceAll((a.Storage().RootURI()).String(), "file://", "") + "/preferences.json"
	if !checkFileExists(prefs) {
		if debug == 1 {
			log.Println("prefs file does not exist")
		}
		// add some default prefs that can be modified via settings
		// Initialize both timer and clock defaults
		writeDefaultSettingsTimer(a)
		writeDefaultSettings(a) // clock defaults from settings-clock.go
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
	jiggle = a.Preferences().IntWithFallback("jiggle.default", 0)
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
	// Load countdown dates
	countdownDate1 = a.Preferences().StringWithFallback("countdown.date1", "")
	countdownDesc1 = a.Preferences().StringWithFallback("countdown.desc1", "")
	countdownDate2 = a.Preferences().StringWithFallback("countdown.date2", "")
	countdownDesc2 = a.Preferences().StringWithFallback("countdown.desc2", "")
	countdownDate3 = a.Preferences().StringWithFallback("countdown.date3", "")
	countdownDesc3 = a.Preferences().StringWithFallback("countdown.desc3", "")
	// Load additional timezones
	timezone1Enabled = a.Preferences().IntWithFallback("timezone1.enabled", 0)
	timezone1Name = a.Preferences().StringWithFallback("timezone1.name", "")
	timezone1Offset = a.Preferences().StringWithFallback("timezone1.offset", "")
	timezone2Enabled = a.Preferences().IntWithFallback("timezone2.enabled", 0)
	timezone2Name = a.Preferences().StringWithFallback("timezone2.name", "")
	timezone2Offset = a.Preferences().StringWithFallback("timezone2.offset", "")
	timezone3Enabled = a.Preferences().IntWithFallback("timezone3.enabled", 0)
	timezone3Name = a.Preferences().StringWithFallback("timezone3.name", "")
	timezone3Offset = a.Preferences().StringWithFallback("timezone3.offset", "")
	timezone4Enabled = a.Preferences().IntWithFallback("timezone4.enabled", 0)
	timezone4Name = a.Preferences().StringWithFallback("timezone4.name", "")
	timezone4Offset = a.Preferences().StringWithFallback("timezone4.offset", "")
	timezone5Enabled = a.Preferences().IntWithFallback("timezone5.enabled", 0)
	timezone5Name = a.Preferences().StringWithFallback("timezone5.name", "")
	timezone5Offset = a.Preferences().StringWithFallback("timezone5.offset", "")
	// Load saved end time from preferences
	endTimeStr := a.Preferences().StringWithFallback("endtime.default", "")
	if endTimeStr != "" {
		// Parse the saved time (HH:MM format)
		if parsedTime, err := time.Parse("15:04", endTimeStr); err == nil {
			now := time.Now()
			// Set customTime to today at the saved time
			customTime = time.Date(now.Year(), now.Month(), now.Day(), parsedTime.Hour(), parsedTime.Minute(), 0, 0, now.Location())
			endTime = customTime
		}
	} else {
		// "Mon Jan 2 15:04:05 MST 2006"
		endTime, _ = time.Parse("15:04", "00:00") // set default midnight
	}

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
		updateAlert(a, updtmsg)
	}

	if desk, ok := a.(desktop.App); ok {
		_, month, _ := time.Now().Date()
		if month == time.December {
			desk.SetSystemTrayIcon(resourceKrankyBearChristmasGrinchPng)
		} else {
			desk.SetSystemTrayIcon(resourceKrankyBearFedoraRedPng)
		}
		if startclock == 1 {
			desktopclock(a)
			// Bring clock to front after a short delay to ensure it appears above timer window
			go func() {
				time.Sleep(200 * time.Millisecond)
				fyne.Do(func() {
					if clock != nil {
						clock.RequestFocus()
					}
				})
			}()
		}
		systray.SetTooltip(appName)
		//systray.SetTitle(timerName)
		show := fyne.NewMenuItem("Show", func() {
			w.Show()
			w.Canvas().Focused()
		})
		hide := fyne.NewMenuItem("Hide", w.Hide)
		lunchMenu = fyne.NewMenuItem("Lunch ("+strconv.Itoa(lunchTime/60)+")", func() {
			startTimer(lunchTime, "Lunch", w.Canvas(), w)
		})
		biobreakMenu = fyne.NewMenuItem("Bio Break ("+strconv.Itoa(biobreakTime/60)+")", func() {
			startTimer(biobreakTime, "Bio Break", w.Canvas(), w)
		})
		adhocmnu = fyne.NewMenuItem("Ad Hoc ("+strconv.Itoa(adhocTime/60)+")", func() {
			startTimer(adhocTime, "Ad Hoc Timer", w.Canvas(), w)
		})
		selectedMenu = fyne.NewMenuItem("Selected End Time", func() {
			// Open dialog to set time and start timer
			setEndTime(a, w, bg, true)
		})
		elapsedMenu = fyne.NewMenuItem("Elapsed Time", func() {
			startElapsedTimer(w.Canvas(), w)
		})
		stopMenu = fyne.NewMenuItem("Stop", func() {
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

			kb := canvas.NewImageFromResource(resourceKrankyBearFedoraRedPng)
			_, month, _ := time.Now().Date()
			if month == time.December {
				kb = canvas.NewImageFromResource(resourceKrankyBearChristmasGrinchPng)
			}
			text := widget.NewLabel(aboutText)
			kb.FillMode = canvas.ImageFillOriginal
			content := container.NewHBox(kb, text)

			if abt == nil {
				abt = a.NewWindow(appName + ": About")
				_, month, _ := time.Now().Date()
				if month == time.December {
					abt.SetIcon(resourceKrankyBearChristmasGrinchPng)
				} else {
					abt.SetIcon(resourceKrankyBearFedoraRedPng)
				}
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
				_, month, _ := time.Now().Date()
				if month == time.December {
					hlp.SetIcon(resourceKrankyBearChristmasGrinchPng)
				} else {
					hlp.SetIcon(resourceKrankyBearFedoraRedPng)
				}
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
		settingsClock := fyne.NewMenuItem("Settings (Clock)", func() {
			makeSettingsClock(a, w, bg)
		})
		settingsTimer := fyne.NewMenuItem("Settings (Timer)", func() {
			makeSettingsTimer(a, w, bg)
		})
		settingsTheme := fyne.NewMenuItem("Settings (Theme)", func() {
			makeSettingsTheme(a, w, bg)
		})
		prefsEdit := fyne.NewMenuItem("Preferences manual edit", func() {
			var cmd *exec.Cmd

			switch runtime.GOOS {
			case "windows":
				cmd = exec.Command("cmd", "/d", "/c", "start", prefs)
			case "darwin": // macOS
				cmd = exec.Command("open", prefs)
			case "linux":
				cmd = exec.Command("xdg-open", prefs)
			default:
				fmt.Printf("Unsupported operating system: %s\n", runtime.GOOS)
				return
			}
			err := cmd.Run()
			if err != nil {
				playBeep("down")
			}
		})
		clock := fyne.NewMenuItem("Clock", func() {
			if clock == nil {
				desktopclock(a)
			} else {
				clock.RequestFocus()
			}
		})
		countdownDates := fyne.NewMenuItem("Countdown Dates", func() {
			makeCountdownDates(a, w, bg)
		})
		alarmsMenu := fyne.NewMenuItem("Alarms", func() {
			makeAlarmsWindow(a, w, bg)
		})
		weatherMenu := fyne.NewMenuItem("Weather", func() {
			makeWeatherWindow(a, w, bg)
		})
		updtchk := fyne.NewMenuItem("Check for update", func() {
			// throw away updateAvail here, use _, unneeded for manual check
			updtmsg, _ := updateChecker("amarillier", "KrankyBearTimer", "Kranky Bear Timer", "https://github.com/amarillier/KrankyBearTimer/releases/latest")
			if updt == nil {
				updateAlert(a, updtmsg)
			} else {
				updt.RequestFocus()
			}
		})
		menu = fyne.NewMenu(a.Metadata().Name, show, hide, clock, alarmsMenu, countdownDates, weatherMenu, fyne.NewMenuItemSeparator(), lunchMenu, biobreakMenu, adhocmnu, selectedMenu, elapsedMenu, stopMenu, fyne.NewMenuItemSeparator(), about, updtchk, help, settingsClock, settingsTimer, settingsTheme, prefsEdit)
		desk.SetSystemTrayMenu(menu)
		_, month, _ = time.Now().Date()
		if month == time.December {
			desk.SetSystemTrayIcon(resourceKrankyBearChristmasGrinchPng)
		} else {
			desk.SetSystemTrayIcon(resourceKrankyBearFedoraRedPng)
		}
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
		newMenuOps := fyne.NewMenu("Operations", show, hide, clock, alarmsMenu, countdownDates, weatherMenu, fyne.NewMenuItemSeparator(), quit)
		newMenuTimers := fyne.NewMenu("Timers", lunchMenu, biobreakMenu, adhocmnu, selectedMenu, elapsedMenu, stopMenu)
		// NB Mac intercepts about item below and puts it where they want to put it!
		// Under 'KrankyBear Timer / About' main section, not under Help
		newMenuHelp := fyne.NewMenu("Help", about, updtchk, help)
		newMenuSettings := fyne.NewMenu("Settings", settingsClock, settingsTimer, settingsTheme, prefsEdit)
		barmenu := fyne.NewMainMenu(newMenuOps, newMenuTimers, newMenuHelp, newMenuSettings)
		w.SetMainMenu(barmenu)
		// barmenu.Refresh()
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
	less.Importance = widget.DangerImportance // red
	more := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		adhocTime += 60 * 5
		adhocbtn.SetText("Ad Hoc (" + strconv.Itoa(adhocTime/60) + ")")
		adhocmnu.Label = "Ad Hoc (" + strconv.Itoa(adhocTime/60) + ")"
		menu.Refresh()
		a.Preferences().SetInt("adhoc.default", adhocTime)
	})
	more.Importance = widget.DangerImportance // red

	lunchBtn = widget.NewButton("Lunch ("+strconv.Itoa(lunchTime/60)+")", func() {
		startTimer(lunchTime, "Lunch", w.Canvas(), w)
	})
	lunchBtn.Importance = widget.SuccessImportance // green
	biobreakBtn = widget.NewButton("Bio Break ("+strconv.Itoa(biobreakTime/60)+")", func() {
		startTimer(biobreakTime, "Bio Break", w.Canvas(), w)
	})
	biobreakBtn.Importance = widget.MediumImportance // white
	adhocbtn = widget.NewButton("Ad Hoc ("+strconv.Itoa(adhocTime/60)+")", func() {
		startTimer(adhocTime, "Ad Hoc", w.Canvas(), w)
	})
	adhocbtn.Importance = widget.DangerImportance // red to match +/- buttons

	// Ad hoc +/- buttons row
	adhocControls := container.NewHBox(container.NewCenter(less), container.NewCenter(more), layout.NewSpacer())

	endtime := widget.NewButton("Selected End Time", func() {
		// Always open the dialog - it will start the timer automatically when time is set
		setEndTime(a, w, bg, true)
	})
	endtime.Importance = widget.WarningImportance // orange

	elapsedBtn = widget.NewButton("Elapsed Time", func() {
		startElapsedTimer(w.Canvas(), w)
	})
	elapsedBtn.Importance = widget.WarningImportance // orange

	// Set up listener for enabling/disabling timer menu items when a timer is running
	running.AddListener(binding.NewDataListener(func() {
		busy, _ := running.Get()
		lunchMenu.Disabled = busy
		biobreakMenu.Disabled = busy
		adhocmnu.Disabled = busy
		selectedMenu.Disabled = busy
		elapsedMenu.Disabled = busy
		stopMenu.Disabled = !busy
		menu.Refresh()
	}))

	// Layout: Grid with biobreak/lunch in row 1, adhoc/endtime in row 2
	// Then adhocControls under adhoc, and elapsed under endtime
	// Create containers to align elapsed under endtime and controls under adhoc
	leftColumn := container.NewVBox(adhocbtn, adhocControls)
	rightColumn := container.NewVBox(endtime, elapsedBtn)

	// Recreate grid with the column containers
	topRow := container.NewGridWithColumns(2, biobreakBtn, lunchBtn)
	bottomRow := container.NewGridWithColumns(2, leftColumn, rightColumn)

	content := container.NewCenter(container.NewVBox(topRow, bottomRow))

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
	bg.Translucency = 0.5                                            // 0.85
	timerBgImage = bg                                                // store reference for dynamic updates
	timerContent = container.NewPadded(container.NewPadded(content)) // store content reference for dynamic updates
	w.SetContent(container.NewStack(
		bg,
		timerContent))
	w.ShowAndRun()
	if updt != nil {
		updt.RequestFocus()
	}
}

// Timer functions moved to timer.go: formatTimer, centerTime, padTime, startTimer, updateTime, setEndTime, isValidCustomTime, updateTimerBackground

func updateAlert(a fyne.App, updtmsg string) {
	// open a window to show the update message
	// no need to test for updt window open at first start
	var kbimg *canvas.Image
	releaselink, rerr := url.Parse("https://github.com/amarillier/KrankyBearTimer/releases/latest")
	if rerr != nil {
		fyne.LogError("Could not parse URL", rerr)
	}
	myreleaselink := widget.NewHyperlink("https://github.com/amarillier/KrankyBearTimer/releases/latest", releaselink)
	myreleaselink.Alignment = fyne.TextAlignLeading

	releasenoteslink, rnerr := url.Parse("https://github.com/amarillier/KrankyBearTimer/blob/main/ReleaseNotes.txt")
	if rnerr != nil {
		fyne.LogError("Could not parse URL", rnerr)
	}
	myreleasenoteslink := widget.NewHyperlink("https://github.com/amarillier/KrankyBearTimer/blob/main/ReleaseNotes.txt", releasenoteslink)
	myreleasenoteslink.Alignment = fyne.TextAlignLeading

	if strings.Contains(updtmsg, "newer version") {
		kbimg = canvas.NewImageFromResource(resourceKrankyBearHardHatPng)
		kbimg.FillMode = canvas.ImageFillOriginal
	} else if strings.Contains(updtmsg, "running the latest") {
		kbimg = canvas.NewImageFromResource(resourceKrankyBearFedoraRedPng)
		kbimg.FillMode = canvas.ImageFillOriginal
	} else {
		alert := sndDir + "/KrankyBearGrowl.mp3"
		alert = sndDir + "/uhOh.mp3"
		if !checkFileExists(alert) {
			playBeep("up")
		} else {
			playMp3(alert) // Basso, Blow, Hero, Funk, Glass, Ping, Purr, Sosumi, Submarine,
		}
		kbimg = canvas.NewImageFromResource(resourceKrankyBearVikingHelmetPng)
		kbimg.FillMode = canvas.ImageFillOriginal
	}

	text := widget.NewLabel(updtmsg)
	content := container.NewVBox(kbimg, text, myreleaselink, myreleasenoteslink)
	updt = a.NewWindow(appName + ": Update Check")
	_, month, _ := time.Now().Date()
	if month == time.December {
		updt.SetIcon(resourceKrankyBearChristmasGrinchPng)
	} else {
		updt.SetIcon(resourceKrankyBearBeretPng)
	}
	updt.Resize(fyne.NewSize(50, 100))
	updt.SetContent(content)
	updt.SetCloseIntercept(func() {
		updt.Close()
		updt = nil
	})
	// updt.CenterOnScreen() // run centered on primary (laptop) display
	updt.Show()
}

// updateTimerBackground moved to timer.go

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942

// This timer is based on an original project named Fomato by Andy Williams, and heavily redeveloped - Allan Marillier, 2024
