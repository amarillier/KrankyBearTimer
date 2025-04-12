#! /bin/sh
# compile, then create a dmg package
# https://github.com/create-dmg/create-dmg

# go build .
# GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o bin/MacOSARM64/
GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -ldflags="-w -s" -o bin/MacOSARM64/

# set executable icon
./setIcon.sh TaniumTimer.png bin/MacOSARM64/TaniumTimer

# cp TaniumTimer TaniumTimer.app/Contents/MacOS
cp bin/MacOSARM64/TaniumTimer TaniumTimer.app/Contents/MacOS
test -f TaniumTimerARM.dmg && rm TaniumTimerARM.dmg
#   --volicon "TaniumTimer.icns" \
create-dmg \
  --volname "TaniumTimer" \
  --window-pos 200 120 \
  --window-size 800 400 \
  --icon-size 100 \
  --icon "TaniumTimer.app" 200 190 \
  --hide-extension "TaniumTimer.app" \
  --app-drop-link 600 185 \
  --eula license.txt \
  "TaniumTimerARM.dmg" \
  "TaniumTimer.app"
  # --add-file TaniumTimer.app ./TaniumTimer.app
  # "./"

# set dmg icon
./setIcon.sh TaniumTimer.png TaniumTimerARM.dmg
if [ ! -d installers ]
then
  mkdir installers
fi
cp TaniumTimerARM.dmg installers
cp TaniumTimerARM.dmg ~/OneDrive\ -\ Tanium\ Inc/Apps/

# "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942