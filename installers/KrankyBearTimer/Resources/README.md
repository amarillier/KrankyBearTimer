# KrankyBear Timer - personal fun / learning side project

preferences stored via fyne preferences API land in
~/Library/Preferences/fyne/com.github.amarillier.KrankyBearTimer/preferences.json
~\AppData\Roaming\fyne\com.github.amarillier.KrankyBearTimer\preferences.json
MacOS resource location (sounds and backgrounds): /Applications/KrankyBear Timer.app/Contents/Resources
Theme preferences are in ~/Library/Preferences/fyne/settings.json


## Features

Basic list - see help in the app for more detail
* Ad hoc time settable in 5 minute steps
* Bio break timer 10 minutes
* Lunch break timer 60 minutes
* Notifications when the timer is done
* Color highlight when time is running out
* Optional customizable desktop clock with seconds, UTC time, date, hourly chime
* System tray access

# To-do / known problems
- Allow optional always on top, save in prefs - may not be possible on Mac
https://www.google.com/search?q=fyne+golang+always+on+top&oq=fyne+golang+always+on+top&gs_lcrp=EgZjaHJvbWUyBggAEEUYOTIKCAEQABiABBiiBDIKCAIQABiABBiiBDIKCAMQABiABBiiBDIKCAQQABiABBiiBNIBCDg5MTBqMGoxqAIAsAIA&sourceid=chrome&ie=UTF-8

- Known problems - needs OpenGL drivers on some Windows
- Possibly add one or more clock alarms - one time, recurring etc.
- Allow a setting to disable hourly chime after hours when hourly chime is enabled

# Info for compiling / modifying

# modules
go mod init KrankyBearTimer

go mod tidy

go get fyne.io/fyne/v2@latest

go install fyne.io/fyne/v2/cmd/fyne@latest

go install fyne.io/fyne/v2/cmd/fyne_demo@latest // gets fyne_demo etc

go get -u github.com/gopxl/beep/v2

go get -u github.com/gopxl/beep/mp3

go get -u github.com/gopxl/beep/v2/mid

go get github.com/spiretechnology/go-autostart/v2@v2.0.0

go get -u github.com/itchyny/volume-go

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
mkdir KrankyBearTimer.app
cp KrankyBearTimer KrankyBearTimer.app
cp bg.tiff KrankyBearTimer/.bg.tiff
cp Icon* KrankyBearTimer.app
cp README.md KrankyBearTimer.app


# Audio: audio converter https://online-audio-convert.com/en/mpeg-to-mp3/


# dmg creation: https://github.com/create-dmg/create-dmg

manual below is difficult
MacOS extended / journaled, no encryption, no partition map
-partitionType none
-noaddpmap


hdiutil create -megabytes 80 -readwrite -volname "KrankyBearTimer" -srcfolder "KrankyBearTimer.app" -ov -format UDZO "KrankyBearTimer.dmg"
hdiutil attach -owners on ./KrankyBearTimer.dmg -shadow
cp "Applications alias" /Volumes/KrankyBearTimer
cp bg.tiff /Volumes/KrankyBearTimer/.bg.tiff
disk=$(diskutil list | grep KrankyBearTimer | awk '{ print $NF }')
hdiutil detach /dev/$disk
hdiutil convert KrankyBearTimer.dmg -format UDRO -o ./KrankyBearTimerRO.dmg



.app to .dmg installer
https://www.youtube.com/watch?v=FqW8Fwfed0U&t=342s
Use InvisibliX and image.tiff for icon


.app to .dmg installer
https://milanpanchal24.medium.com/a-guide-to-converting-app-to-dmg-on-macos-c19f9d81f871


# Generate the DMG file with debug option
hdiutil create -volname "KrankyBearTimer" -srcfolder "KrankyBearTimer.app" -ov -format UDZO "KrankyBearTimer.dmg" -debug

# Generate the DMG file with encryption [AES-128|AES-256]
hdiutil create -volname "KrankyBearTimer" -srcfolder "KrankyBearTimer.app" -ov -format UDZO "KrankyBearTimer.dmg" -encryption AES-128

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

