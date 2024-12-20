# Tanium Timer

# preferences stored via fyne preferences API land in
# ~/Library/Preferences/fyne/com.tanium.taniumtimer/preferences.json
# ~\AppData\Roaming\fyne\com.tanium.taniumtimer\preferences.json


## Features

* Ad hoc time settable in 5 minute steps
* Bio break timer 10 minutes
* Lunch break timer 60 minutes
* Notifications when the timer is done
* Color highlight when time is running out
* System tray access

# modules
go mod init TaniumTimer
go mod tidy
go get fyne.io/fyne/v2@latest
go install fyne.io/fyne/v2/cmd/fyne@latest
go install fyne.io/fyne/v2/cmd/fyne_demo@latest // gets fyne_demo etc
go get github.com/gen2brain/beeep
go get -u github.com/gopxl/beep/v2
go get -u github.com/gopxl/beep/mp3
go get -u github.com/gopxl/beep/v2/mid
go get github.com/spiretechnology/go-autostart/v2@v2.0.0

Occasionally go mod vendor to resolve problems
or for: build constraints exclude all Go files in ....
go clean -modcache
go mod tidy
https://stackoverflow.com/questions/55348458/build-constraints-exclude-all-go-files-in


# error logging
- https://rollbar.com/blog/golang-error-logging-guide/


# cross compile for Windows
https://stackoverflow.com/questions/36915134/go-golang-cross-compile-from-mac-to-windows-fatal-error-windows-h-file-not-f
brew install mingw-w64

# cross compile for Linux
?


# audio (mp3 / wav / midi) player
https://github.com/gopxl/beep

# beeep - prefer gopxl beep over this: https://pkg.go.dev/github.com/gen2brain/beeep#section-readme
https://pkg.go.dev/github.com/gen2brain/beeep#Alert



# png to svg online converter:
BEST: Use Inkscape (free)
- Open .png, .jpg etc, choose option (default) embed image
- Use selection tool arrow, click in image, verify selected
- click Path / Trace Bitmap / Pixel Art
- check image preview, make changes if needed, update preview
- Apply, wait a while ...
- File, Save As, ...svg

https://new.express.adobe.com/tools/convert-to-svg
https://convertio.co/
https://www.freeconvert.com/png-to-svg/download

# use https://www.aconvert.com/image/png-to-icns/ for png to icns conversion
mkdir TaniumTimer.app
cp TaniumTimer TaniumTimer.app
cp bg.tiff TaniumTimer/.bg.tiff
cp Icon* TaniumTimer.app
cp README.md TaniumTimer.app


# Audio: audio converter https://online-audio-convert.com/en/mpeg-to-mp3/


# dmg creation: https://github.com/create-dmg/create-dmg

manual below is difficult
MacOS extended / journaled, no encryption, no partition map
-partitionType none
-noaddpmap


hdiutil create -megabytes 80 -readwrite -volname "TaniumTimer" -srcfolder "TaniumTimer.app" -ov -format UDZO "TaniumTimer.dmg"
hdiutil attach -owners on ./TaniumTimer.dmg -shadow
cp "Applications alias" /Volumes/TaniumTimer
cp bg.tiff /Volumes/TaniumTimer/.bg.tiff
disk=$(diskutil list | grep TaniumTimer | awk '{ print $NF }')
hdiutil detach /dev/$disk
hdiutil convert TaniumTimer.dmg -format UDRO -o ./TaniumTimerRO.dmg



.app to .dmg installer
https://www.youtube.com/watch?v=FqW8Fwfed0U&t=342s
Use InvisibliX and image.tiff for icon


.app to .dmg installer
https://milanpanchal24.medium.com/a-guide-to-converting-app-to-dmg-on-macos-c19f9d81f871


# Generate the DMG file with debug option
hdiutil create -volname "TaniumTimer" -srcfolder "TaniumTimer.app" -ov -format UDZO "TaniumTimer.dmg" -debug

# Generate the DMG file with encryption [AES-128|AES-256]
hdiutil create -volname "TaniumTimer" -srcfolder "TaniumTimer.app" -ov -format UDZO "TaniumTimer.dmg" -encryption AES-128

https://stackoverflow.com/questions/37292756/how-to-create-a-dmg-file-for-a-app-for-mac

Copy your app to a new folder.
Open Disk Utility -> File -> New Image -> Image From Folder.
Select the folder where you have placed the App. Give a name for the DMG and save. This creates a distributable image for you.
If needed you can add a link to applications to DMG. It helps user in installing by drag and drop.

To create a disk image using the Terminal on a Mac, you can use the hdiutil command:
Open Terminal
Type hdiutil create -volname N -srcfolder P -ov N.dmg
Replace N with the name of the disk image file and P with the path of the source volume
Press Return

