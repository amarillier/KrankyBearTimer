#! /bin/sh
# compile, then create a dmg package
# https://github.com/create-dmg/create-dmg

# go build .
# GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o bin/MacOSARM64/
GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -ldflags="-w -s" -o bin/MacOSARM64/

# set executable icon
./setIcon.sh Resources/Images/KrankyBearFedoraRed.png bin/MacOSARM64/KrankyBearTimer

# cp KrankyBearTimer KrankyBearTimer.app/Contents/MacOS
cp bin/MacOSARM64/KrankyBearTimer KrankyBearTimer.app/Contents/MacOS
test -f KrankyBearTimerARM.dmg && rm KrankyBearTimerARM.dmg
#   --volicon "KrankyBearTimer.icns" \
create-dmg \
  --volname "KrankyBearTimer" \
  --window-pos 200 120 \
  --window-size 800 400 \
  --icon-size 100 \
  --icon "KrankyBearTimer.app" 200 190 \
  --hide-extension "KrankyBearTimer.app" \
  --app-drop-link 600 185 \
  --eula license.txt \
  "KrankyBearTimerARM.dmg" \
  "KrankyBearTimer.app"
  # --add-file KrankyBearTimer.app ./KrankyBearTimer.app
  # "./"

# set dmg icon
./setIcon.sh Resources/Images/KrankyBearFedoraRed.png KrankyBearTimerARM.dmg
if [ ! -d installers ]
then
  mkdir installers
fi
cp KrankyBearTimerARM.dmg installers

# "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942