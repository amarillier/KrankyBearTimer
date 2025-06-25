package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	updatechecker "github.com/amarillier/go-update-checker"
	"github.com/itchyny/volume-go"
)

var (
	WarningLog *log.Logger
	InfoLog    *log.Logger
	ErrorLog   *log.Logger
)

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(error, os.ErrNotExist)
}

func daysUntil(targetDate string) (int, error) {
	// Parse the target date
	layout := "2006-01-02"
	target, err := time.Parse(layout, targetDate)
	if err != nil {
		return 0, err
	}
	// Get the current date
	current := time.Now()
	// Calculate the duration between the current date and the target date
	duration := target.Sub(current)
	// Convert the duration to days
	days := int(duration.Hours() / 24)
	return days, nil
}

// see error logging:
// https://rollbar.com/blog/golang-error-logging-guide/
func logInit() {
	// typically written to Resources/...
	file, err := os.OpenFile(processName+"0.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	InfoLog = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLog = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLog = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func lineCounter(r io.Reader) (int, error) {
	// count lines in a file, used for log rotation and possible other uses
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)
		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func logRotate() {
	if checkFileExists(processName + "2.txt") {
		f := os.Remove(processName + "2.txt")
		if f != nil {
			ErrorLog.Printf("Error attempting to remove %s", processName+"2.txt")
		}
	}
	if checkFileExists(processName + "1.txt") {
		os.Rename(processName+"1.txt", processName+"2.txt")
	}
	if checkFileExists(processName + "0.txt") {
		os.Rename(processName+"0.txt", processName+"1.txt")
	}
}

func easterEgg(a fyne.App, w fyne.Window) {
	muted, _ := volume.GetMuted()
	vol, _ := volume.GetVolume()
	var eggvol = 15
	var certs []fyne.Resource

	// allow for Tanium branding, show certs
	// improve this, and add other Kranky Bear images
	processName := "^tanium.*"
	processRegex := regexp.MustCompile(`(?i)^tanium`)
	if processRegex.MatchString(appNameCustom) || isProcessRunning(processName) {
		certs = []fyne.Resource{resourceKrankyBearPng, resourceKrankyBearBeretPng, resourceKrankyBearBeanieMultiColorPng, resourceKrankyBearChristmasPng, resourceKrankyBearChristmasGrinchPng, resourceKrankyBearCowboyBrownPng, resourceKrankyBearFedoraRedPng, resourceKrankyBearHardHatPng, resourceKrankyBearHogwartsSortingPng, resourceKrankyBearTrapperRedPlaidPng, resourceKrankyBearVikingHelmetPng, resourceHttp418Png, resourceTcn256Png, resourceTcc256Png, resourceTcbe256Png}
	} else {
		certs = []fyne.Resource{resourceKrankyBearPng, resourceKrankyBearBeretPng, resourceKrankyBearBeanieMultiColorPng, resourceKrankyBearChristmasPng, resourceKrankyBearChristmasGrinchPng, resourceKrankyBearCowboyBrownPng, resourceKrankyBearFedoraRedPng, resourceKrankyBearHardHatPng, resourceKrankyBearHogwartsSortingPng, resourceKrankyBearTrapperRedPlaidPng, resourceKrankyBearVikingHelmetPng, resourceHttp418Png}
	}

	randomIndex := rand.Intn(len(certs))
	egg := a.NewWindow(appName + ": easter egg")
	egg.SetIcon(resourceKrankyBearTrapperRedPlaidPng)
	eggimage := canvas.NewImageFromResource(certs[randomIndex])
	eggimage.FillMode = canvas.ImageFillOriginal
	text := "Whoo-hoo! You found the Easter egg!\n"
	text += "\n" + dadjoke()

	eggtext := widget.NewLabel(text)
	content := container.NewVBox(eggimage, eggtext)
	egg.SetContent(content)
	// egg.CenterOnScreen() // run centered on primary (laptop) display
	if muted {
		volume.Unmute()
		if vol <= 10 {
			volume.SetVolume(eggvol)
		}
	}
	playBeep("down")
	/*
		for j := 0; j <= 2; j++ {
			playBeep("down")
			egg.Show()
			time.Sleep(time.Second / 3)
			egg.Hide()
			time.Sleep(time.Second / 3)
		}
	*/
	if muted {
		if eggvol > vol {
			volume.SetVolume(vol)
		}
		volume.Mute()
	}
	w.RequestFocus()
	egg.Show()
}

func teapot(a fyne.App, w fyne.Window) {
	muted, _ := volume.GetMuted()
	vol, _ := volume.GetVolume()
	var teapotvol = 10

	link, err := url.Parse("https://www.rfc-editor.org/rfc/rfc2324.html")
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}
	hyperlink := widget.NewHyperlink("What is http 418? https://www.rfc-editor.org/rfc/rfc2324.html", link)
	hyperlink.Alignment = fyne.TextAlignLeading
	tpwin := a.NewWindow(appName + ": http: 418")
	tpwin.SetIcon(resourceKrankyBearTrapperRedPlaidPng)
	tpwinimage := canvas.NewImageFromResource(resourceHttp418Png)
	tpwinimage.FillMode = canvas.ImageFillOriginal
	text := "Whoo-hoo! You found another Easter egg!\n"

	tpwintext := widget.NewLabel(text)
	content := container.NewVBox(tpwinimage, tpwintext, hyperlink)
	tpwin.SetContent(content)
	// tpwin.CenterOnScreen() // run centered on primary (laptop) display
	// tpwin.Show()
	if muted {
		volume.Unmute()
		if vol <= 10 {
			volume.SetVolume(teapotvol)
		}
	}
	playBeep("down")
	/*
		for j := 0; j <= 2; j++ {
			// playBeep("down")
			fmt.Println("egg loop")
			tpwin.Show()
			time.Sleep(time.Second / 3)
			tpwin.Hide()
			time.Sleep(time.Second / 3)
		}
	*/
	if muted {
		if teapotvol > vol {
			volume.SetVolume(vol)
		}
		volume.Mute()
	}
	w.RequestFocus()
	tpwin.Show()
}

func listMatchingFiles(directory, pattern string) ([]string, error) {
	var matchingFiles []string

	// Read the directory
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	// Loop through the files and match the pattern
	for _, file := range files {
		if matched, err := filepath.Match(pattern, file.Name()); err != nil {
			return nil, err
		} else if matched {
			matchingFiles = append(matchingFiles, file.Name())
		}
	}
	return matchingFiles, nil
}

func dadjoke() string {
	// Define an array of jokes
	jokes := []string{
		"We're having Himalayan rabbit stew for dinner.\nI found Him a-layin in the middle of the road",
		"I went to the local zoo, but all they had was one dog.\nIt was a Shi-Tzu",
		"Wildlife biologists have proved that Pronghorn Antelope can jump higher than the average house.\nThis is due to the fact that the average house can't jump",
		"Where do rainbows go when they've been bad?\nTo prism, so they have time to reflect on what they've done",
		"Dogs can't operate MRI machines.\nBut catscan",
		"What do you call a dog who meditates?\nAware wolf",
		"Why did the old man fall down the well?\nHe couldn’t see that well",
		"The other day I bought a thesaurus, but when I got home and opened it, all the pages were blank.\nI have no words to describe how angry I am",
		"Why can't humans hear a dog whistle?\nBecause a dog can't whistle",
		"What is a dog's favorite form of transport?\nA waggin",
		"Lemoncello? Over in the clearance corner because nobody could get any good notes from it",
		"What's a forklift?\nUsually, food",
		"Did you know you can wear a canoe as a hat?\nIf you turn it over, it is capsized",
		"Eucaplyptus is the only plant named for what it would say after you prune it",
		"I was going to tell a time traveling joke, but you didn't like it",
		"Why did the chicken join a band?\nBecause it had the drumsticks",
		"Why did the scarecrow win an award?\nBecause he was outstanding in his field",
		"Why don't skeletons fight each other?\nThey don't have the guts",
		"I was going to tell a chemistry joke, but I knew I wouldn't get a reaction",
		"Why don't scientists trust atoms?\nBecause they make up everything",
		"What do you call fake spaghetti?\nAn impasta",
		"Why did the math book look sad?\nBecause it had too many problems",
		"I was going to cook alligator tonight, but I only have a crocpot",
		"A Japanese gardener asked me what I know about bonsai trees.\nI said, 'Very little'",
		"A horse walked into a bar and ordered a beer. The bartender said 'You come in here often, do you think you might be an alcoholic?'\nThe horse said 'I don'''t think I am, then vanished from existence.\nYou see, this joke is about Descartes, 'I think, therefore I am'. But to have explained that first would'''ve put Descartes before the horse.",
		"Davey Crocket was the only man ever to have three ears.\nA left ear, a right ear, and a wild front ear",
	}
	randomIndex := rand.Intn(len(jokes))
	joke := jokes[randomIndex]
	return (joke)
}

func isProcessRunning(processRegex string) bool {
	var cmd *exec.Cmd
	// if os.PathSeparator == '\\' { // Windows
	if runtime.GOOS == "windows" {
		cmd = exec.Command("tasklist")
	} else { // Unix-based systems
		cmd = exec.Command("pgrep", "-fl", processRegex) // "-fl" to show full process name
	}

	output, err := cmd.Output()
	if err != nil {
		return false
	}

	// Compile regex with case-insensitivity
	re := regexp.MustCompile(`(?i)` + processRegex)
	for _, line := range strings.Split(string(output), "\n") {
		if re.MatchString(line) {
			return true
		}
	}
	return false
}

func updateChecker(repoOwner string, repo string, repoName string, repodl string) (string, bool) {
	// uc := updatechecker.New("amarillier", "KrankyBearTimer", "Kranky Bear Timer", "", 1, false)
	uc := updatechecker.New(repoOwner, repo, repoName, repodl, 0, false)
	uc.CheckForUpdate(appVersion)
	// uc.PrintMessage()
	updtmsg := uc.Message
	return updtmsg, uc.UpdateAvailable
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
